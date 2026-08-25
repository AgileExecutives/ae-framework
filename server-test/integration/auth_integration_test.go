package integration

import (
    "encoding/json"
    "io"
    "net/http"
    "testing"
    "time"
)

func TestAuthRegisterAndLogin(t *testing.T) {
    serverDir, cleanup := StartServer(t)
    _ = serverDir
    t.Cleanup(cleanup)

    token, email, username, _ := RegisterUserAndToken(t)

    // Verify login by calling /auth/me
    req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/auth/me", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("me request failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        body, _ := io.ReadAll(resp.Body)
        t.Fatalf("me failed status=%d body=%s", resp.StatusCode, string(body))
    }
    var parsed map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
        t.Fatalf("invalid me response: %v", err)
    }
    data, _ := parsed["data"].(map[string]interface{})
    if data == nil {
        t.Fatalf("me response missing data")
    }
    if data["username"] != username {
        t.Fatalf("unexpected username: got %v want %v", data["username"], username)
    }
    if data["email"] != email {
        t.Fatalf("unexpected email: got %v want %v", data["email"], email)
    }
}
