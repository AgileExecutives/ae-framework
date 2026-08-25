package integration

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "testing"
    "time"
    "fmt"
)

// RegisterUserAndToken registers a fresh user and returns an auth token.
func RegisterUserAndToken(t *testing.T) (string, string, string, string) {
    t.Helper()
    uniq := fmt.Sprintf("h%d", time.Now().UnixNano())
    reg := map[string]interface{}{
        "email": fmt.Sprintf("helper+%s@example.com", uniq),
        "username": "helper_" + uniq,
        "password": "helperPass123",
        "first_name": "Helper",
        "last_name": "User",
        "company_name": "Helper Co " + uniq,
        "tenant_name": "tenant_" + uniq,
        "accept_terms": true,
    }
    b, _ := json.Marshal(reg)
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Post("http://localhost:8080/api/v1/auth/register", "application/json", bytes.NewReader(b))
    if err != nil {
        t.Fatalf("helper register request failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        body, _ := io.ReadAll(resp.Body)
        t.Fatalf("helper register failed status=%d body=%s", resp.StatusCode, string(body))
    }
    var parsed map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
        t.Fatalf("helper invalid register response: %v", err)
    }
    data := parsed["data"].(map[string]interface{})
    token, _ := data["token"].(string)
    // use the original registration payload for email/username/password to be deterministic
    email, _ := reg["email"].(string)
    username, _ := reg["username"].(string)
    password, _ := reg["password"].(string)
    if token == "" {
        t.Fatalf("helper token missing in register response")
    }
    return token, email, username, password
}

// AuthRequest performs an authenticated HTTP request and returns the response and body.
func AuthRequest(t *testing.T, token, method, url string, body []byte) (*http.Response, []byte) {
    t.Helper()
    var req *http.Request
    var err error
    if body != nil {
        req, err = http.NewRequest(method, url, bytes.NewReader(body))
    } else {
        req, err = http.NewRequest(method, url, nil)
    }
    if err != nil {
        t.Fatalf("failed to create request: %v", err)
    }
    req.Header.Set("Authorization", "Bearer "+token)
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    b, _ := io.ReadAll(resp.Body)
    resp.Body.Close()
    return resp, b
}

// CreateTemplate posts a template payload and returns created template ID.
func CreateTemplate(t *testing.T, token string, payload map[string]interface{}) int {
    t.Helper()
    b, _ := json.Marshal(payload)
    resp, body := AuthRequest(t, token, "POST", "http://localhost:8080/api/v1/templates", b)
    if resp.StatusCode != 201 {
        t.Fatalf("create template failed status=%d body=%s", resp.StatusCode, string(body))
    }
    var parsed map[string]interface{}
    if err := json.Unmarshal(body, &parsed); err != nil {
        t.Fatalf("invalid create template response: %v", err)
    }
    idFloat := parsed["data"].(map[string]interface{})["id"].(float64)
    return int(idFloat)
}

// RenderTemplate calls the /templates/render endpoint with provided payload and returns status and body.
func RenderTemplate(t *testing.T, token string, payload map[string]interface{}) (int, []byte) {
    t.Helper()
    b, _ := json.Marshal(payload)
    resp, body := AuthRequest(t, token, "POST", "http://localhost:8080/api/v1/templates/render", b)
    return resp.StatusCode, body
}

// DeleteTemplate deletes a template by id and returns the status code.
func DeleteTemplate(t *testing.T, token string, id int) int {
    t.Helper()
    url := fmt.Sprintf("http://localhost:8080/api/v1/templates/%d", id)
    resp, body := AuthRequest(t, token, "DELETE", url, nil)
    _ = body
    return resp.StatusCode
}
