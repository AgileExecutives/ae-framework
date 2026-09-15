package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

// TestDatabaseSeeding verifies that when the server is started with an empty
// database the bootstrap process creates the initial tenant, organization
// and seeds the initial admin user as configured by the test harness.
func TestDatabaseSeeding(t *testing.T) {
	_, cleanup := StartServer(t)
	t.Cleanup(cleanup)

	client := &http.Client{Timeout: 5 * time.Second}

	// The test harness StartServer sets deterministic admin credentials
	// (see StartServer in testutil.go). Use those to login and inspect
	// the returned user payload for tenant/org ids.
	creds := map[string]string{"email": "testuser@unburdy.de", "password": "newpass123"}
	b, _ := json.Marshal(creds)

	resp, err := client.Post("http://localhost:8080/api/v1/auth/login", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("admin login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("admin login failed status=%d body=%s", resp.StatusCode, string(body))
	}

	var parsed map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		t.Fatalf("invalid admin login response: %v", err)
	}

	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing data in login response")
	}
	user, ok := data["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing user in login response")
	}

	// JSON numbers decode as float64
	tenantF, _ := user["tenant_id"].(float64)
	orgF, _ := user["organization_id"].(float64)

	if int(tenantF) != 1 {
		t.Fatalf("expected seeded tenant id 1, got %v", int(tenantF))
	}
	if int(orgF) != 1 {
		t.Fatalf("expected seeded organization id 1, got %v", int(orgF))
	}

	// basic sanity: ensure the admin token can be used to access a protected endpoint
	token, _ := data["token"].(string)
	if token == "" {
		t.Fatalf("login response missing token")
	}

	// call /auth/me to confirm token is accepted
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("auth/me request failed: %v", err)
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		b, _ := io.ReadAll(r.Body)
		t.Fatalf("auth/me failed status=%d body=%s", r.StatusCode, string(b))
	}

	// allow some breathing room for async background initialization
	time.Sleep(50 * time.Millisecond)
}
