package integration

import (
    "encoding/json"
    "fmt"
    "testing"
    "time"
)

func TestPlansBasicCRUD(t *testing.T) {
    _, cleanup := StartServer(t)
    t.Cleanup(cleanup)

    // Register a fresh admin user via helper and use its token
    token, _, _, _ := RegisterUserAndToken(t)

    // GET /api/v1/plans
    resp, _ := AuthRequest(t, token, "GET", "http://localhost:8080/api/v1/plans", nil)
    if resp.StatusCode != 200 { t.Fatalf("GET plans status=%d", resp.StatusCode) }

    // GET /api/v1/plans/1
    resp, _ = AuthRequest(t, token, "GET", "http://localhost:8080/api/v1/plans/1", nil)
    if resp.StatusCode != 200 { t.Fatalf("GET plan/1 status=%d", resp.StatusCode) }

    // Create a new plan as admin
    plan := map[string]interface{}{
        "name": "Integration Test Plan",
        "slug": fmt.Sprintf("int-plan-%d", time.Now().UnixNano()),
        "description": "created by integration test",
        "price": 9.99,
        "currency": "USD",
        "invoice_period": "monthly",
        "max_users": 2,
        "active": true,
    }
    pb, _ := json.Marshal(plan)
    resp, body := AuthRequest(t, token, "POST", "http://localhost:8080/api/v1/plans", pb)
    if resp.StatusCode != 201 && resp.StatusCode != 200 {
        t.Fatalf("create plan failed status=%d body=%s", resp.StatusCode, string(body))
    }
    var created map[string]interface{}
    if err := json.Unmarshal(body, &created); err != nil {
        t.Fatalf("invalid create plan response: %v", err)
    }
    dataResp, _ := created["data"].(map[string]interface{})
    if dataResp == nil || dataResp["name"] != "Integration Test Plan" {
        t.Fatalf("unexpected create response: %v", created)
    }
}
