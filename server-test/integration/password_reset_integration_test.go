package integration

import (
    "encoding/json"
    "regexp"
    "testing"
)

func TestPasswordResetFullFlow(t *testing.T) {
    _, cleanup := StartServer(t)
    t.Cleanup(cleanup)

    // Register a fresh user and trigger forgot-password for that email
    _, email, _, _ := RegisterUserAndToken(t)

    payload := map[string]string{"email": email}
    pb, _ := json.Marshal(payload)
    resp, _ := AuthRequest(t, "", "POST", "http://localhost:8080/api/v1/auth/forgot-password", pb)
    if resp.StatusCode != 200 {
        t.Fatalf("forgot-password failed status=%d", resp.StatusCode)
    }

    // Retrieve latest mock emails
    resp, respBody := AuthRequest(t, "", "GET", "http://localhost:8080/api/v1/emails/latest-emails", nil)
    if resp.StatusCode != 200 {
        t.Fatalf("latest-emails failed status=%d body=%s", resp.StatusCode, string(respBody))
    }
    var parsed map[string]interface{}
    if err := json.Unmarshal(respBody, &parsed); err != nil {
        t.Fatalf("invalid latest-emails response: %v", err)
    }
    data, _ := parsed["data"].([]interface{})
    if len(data) == 0 {
        t.Fatalf("no mock emails found")
    }
    // Use last email HTML to extract token
    last := data[len(data)-1].(map[string]interface{})
    html, _ := last["html"].(string)
    re := regexp.MustCompile(`new-password\?token=([A-Za-z0-9\-\._~%]+)`) 
    m := re.FindStringSubmatch(html)
    if len(m) < 2 {
        t.Fatalf("reset token not found in email html")
    }
    token := m[1]

    // Reset password using token
    newPw := map[string]string{"new_password": "NewPass!234"}
    nb, _ := json.Marshal(newPw)
    resetURL := "http://localhost:8080/api/v1/auth/new-password/" + token
    resp, _ = AuthRequest(t, "", "POST", resetURL, nb)
    if resp.StatusCode != 200 {
        t.Fatalf("reset password failed status=%d", resp.StatusCode)
    }

    // Login with new password
    login := map[string]string{"email": email, "password": "NewPass!234"}
    lb, _ := json.Marshal(login)
    resp, _ = AuthRequest(t, "", "POST", "http://localhost:8080/api/v1/auth/login", lb)
    if resp.StatusCode != 200 {
        t.Fatalf("login failed status=%d", resp.StatusCode)
    }
}
