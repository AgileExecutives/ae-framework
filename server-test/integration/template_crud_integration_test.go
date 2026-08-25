package integration

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestTemplateCRUD(t *testing.T) {
	_, cleanup := StartServer(t)
	t.Cleanup(cleanup)

	token, _, _, _ := RegisterUserAndToken(t)
	uniq := fmt.Sprintf("tc%d", time.Now().UnixNano())

	// Create Email Template
	emailTpl := map[string]interface{}{
		"template_type": "email",
		"template_key":  "integration_test_email_" + uniq,
		"channel":       "EMAIL",
		"subject":       "Integration Test Email",
		"name":          "Integration Test Email Template",
		"description":   "Template for integration testing",
		"content":       "<h1>Hello {{.Name}}!</h1><p>This is a test email for {{.Purpose}}.</p><p>Contact: {{.ContactEmail}}</p>",
		"variables":     []string{"Name", "Purpose", "ContactEmail"},
		"sample_data":   map[string]interface{}{"Name": "Integration Tester", "Purpose": "API Testing", "ContactEmail": "test@integration.com"},
		"is_active":     true,
		"is_default":    false,
	}
	emailID := CreateTemplate(t, token, emailTpl)

	// Create PDF template
	pdfTpl := map[string]interface{}{
		"template_type": "document",
		"template_key":  "integration_test_document_" + uniq,
		"channel":       "PDF",
		"name":          "Integration Test Document Template",
		"description":   "PDF template for integration testing",
		"content":       "<html><head><title>{{.Title}}</title></head><body><h1>{{.Title}}</h1><p>{{.Content}}</p></body></html>",
		"variables":     []string{"Title", "Content", "Date"},
		"sample_data":   map[string]interface{}{"Title": "Integration Test Document", "Content": "This is test content", "Date": "2024-01-15"},
		"is_active":     true,
	}
	pdfID := CreateTemplate(t, token, pdfTpl)

	// Verify GET /templates contains both
	resp, body := AuthRequest(t, token, "GET", "http://localhost:8080/api/v1/templates", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("get templates status=%d body=%s", resp.StatusCode, string(body))
	}
	if !strings.Contains(string(body), "Integration Test Email Template") || !strings.Contains(string(body), "Integration Test Document Template") {
		t.Fatalf("templates listing missing created templates: %s", string(body))
	}

	// Get specific email template
	resp, _ = AuthRequest(t, token, "GET", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d", emailID), nil)
	if resp.StatusCode != 200 {
		t.Fatalf("get email tpl status=%d", resp.StatusCode)
	}

	// Update email template
	update := map[string]interface{}{"name": "Updated Integration Test Email", "description": "Updated description", "content": "<h1>Hi {{.Name}}!</h1><p>Updated content {{.Purpose}}</p>", "variables": []string{"Name", "Purpose"}, "sample_data": map[string]interface{}{"Name": "John", "Purpose": "CRUD"}}
	ub, _ := json.Marshal(update)
	resp, _ = AuthRequest(t, token, "PUT", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d", emailID), ub)
	if resp.StatusCode != 200 {
		t.Fatalf("update tpl status=%d", resp.StatusCode)
	}

	// Render updated template
	renderPayload := map[string]interface{}{"data": map[string]interface{}{"Name": "John Doe", "Purpose": "CRUD Testing"}}
	rb2, _ := json.Marshal(renderPayload)
	resp, body = AuthRequest(t, token, "POST", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d/render", emailID), rb2)
	if resp.StatusCode != 200 {
		t.Fatalf("render updated tpl status=%d body=%s", resp.StatusCode, string(body))
	}
	if !strings.Contains(string(body), "John Doe") {
		t.Fatalf("render content missing: %s", string(body))
	}

	// Duplicate template
	dup := map[string]interface{}{"name": "Duplicated Integration Email", "template_key": "duplicated_integration_email_" + uniq}
	db, _ := json.Marshal(dup)
	resp, body = AuthRequest(t, token, "POST", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d/duplicate", emailID), db)
	if resp.StatusCode != 201 {
		t.Fatalf("duplicate tpl status=%d body=%s", resp.StatusCode, string(body))
	}
	var created map[string]interface{}
	if err := json.Unmarshal(body, &created); err != nil {
		t.Fatalf("invalid duplicate response: %v", err)
	}
	dupID := int(created["data"].(map[string]interface{})["id"].(float64))

	// Deactivate then reactivate
	status := map[string]interface{}{"is_active": false}
	sb, _ := json.Marshal(status)
	resp, _ = AuthRequest(t, token, "PUT", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d", emailID), sb)
	if resp.StatusCode != 200 {
		t.Fatalf("deactivate status=%d", resp.StatusCode)
	}

	status["is_active"] = true
	sb, _ = json.Marshal(status)
	resp, _ = AuthRequest(t, token, "PUT", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d", emailID), sb)
	if resp.StatusCode != 200 {
		t.Fatalf("reactivate status=%d", resp.StatusCode)
	}

	// Filter by template_type and channel
	resp, _ = AuthRequest(t, token, "GET", "http://localhost:8080/api/v1/templates?template_type=email", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("filter by type status=%d", resp.StatusCode)
	}

	resp, _ = AuthRequest(t, token, "GET", "http://localhost:8080/api/v1/templates?channel=EMAIL", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("filter by channel status=%d", resp.StatusCode)
	}

	// Cleanup deletes
	resp, _ = AuthRequest(t, token, "DELETE", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d", dupID), nil)
	if resp.StatusCode != 204 {
		t.Fatalf("delete dup status=%d", resp.StatusCode)
	}

	resp, _ = AuthRequest(t, token, "DELETE", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d", pdfID), nil)
	if resp.StatusCode != 204 {
		t.Fatalf("delete pdf status=%d", resp.StatusCode)
	}

	resp, _ = AuthRequest(t, token, "DELETE", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d", emailID), nil)
	if resp.StatusCode != 204 {
		t.Fatalf("delete email status=%d", resp.StatusCode)
	}

	// Verify deleted templates return 404
	resp, _ = AuthRequest(t, token, "GET", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d", emailID), nil)
	if resp.StatusCode != 404 {
		t.Fatalf("expected 404 for deleted email, got %d", resp.StatusCode)
	}

	resp, _ = AuthRequest(t, token, "GET", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d", pdfID), nil)
	if resp.StatusCode != 404 {
		t.Fatalf("expected 404 for deleted pdf, got %d", resp.StatusCode)
	}

	resp, _ = AuthRequest(t, token, "GET", fmt.Sprintf("http://localhost:8080/api/v1/templates/%d", dupID), nil)
	if resp.StatusCode != 404 {
		t.Fatalf("expected 404 for deleted dup, got %d", resp.StatusCode)
	}
}
