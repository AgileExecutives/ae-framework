package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestTemplateContractsEndpoints(t *testing.T) {
	_, cleanup := StartServer(t)
	t.Cleanup(cleanup)

	token, _, _, _ := RegisterUserAndToken(t)
	client := &http.Client{Timeout: 5 * time.Second}

	authReq := func(method, url string, body []byte) (*http.Response, []byte, error) {
		var req *http.Request
		var err error
		if body != nil {
			req, err = http.NewRequest(method, url, bytes.NewReader(body))
		} else {
			req, err = http.NewRequest(method, url, nil)
		}
		if err != nil {
			return nil, nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, nil, err
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return resp, b, nil
	}

	// List all contracts
	resp, body, err := authReq("GET", "http://localhost:8080/api/v1/templates/contracts", nil)
	if err != nil {
		t.Fatalf("list contracts failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("list contracts status=%d body=%s", resp.StatusCode, string(body))
	}

	// Get welcome contract
	resp, body, err = authReq("GET", "http://localhost:8080/api/v1/templates/contracts/welcome", nil)
	if err != nil {
		t.Fatalf("get welcome contract failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("welcome contract status=%d body=%s", resp.StatusCode, string(body))
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("invalid welcome response: %v", err)
	}
	data := parsed["data"].(map[string]interface{})
	if data["template_key"] != "welcome" {
		t.Fatalf("welcome key mismatch: %v", data["template_key"])
	}

	// Get booking_confirmation contract
	resp, body, err = authReq("GET", "http://localhost:8080/api/v1/templates/contracts/booking_confirmation", nil)
	if err != nil {
		t.Fatalf("get booking contract failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("booking contract status=%d body=%s", resp.StatusCode, string(body))
	}

	// Get password_reset contract
	resp, body, err = authReq("GET", "http://localhost:8080/api/v1/templates/contracts/password_reset", nil)
	if err != nil {
		t.Fatalf("get password reset contract failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("password reset contract status=%d body=%s", resp.StatusCode, string(body))
	}

	// Get invoice contract
	resp, body, err = authReq("GET", "http://localhost:8080/api/v1/templates/contracts/invoice", nil)
	if err != nil {
		t.Fatalf("get invoice contract failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("invoice contract status=%d body=%s", resp.StatusCode, string(body))
	}

	// Get sample data for welcome
	resp, body, err = authReq("GET", "http://localhost:8080/api/v1/templates/contracts/welcome/sample-data", nil)
	if err != nil {
		t.Fatalf("get welcome sample failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("welcome sample status=%d body=%s", resp.StatusCode, string(body))
	}

	// Validate valid data against welcome
	valid := map[string]interface{}{"FirstName": "John", "LastName": "Doe", "OrganizationName": "Test Company"}
	vb, _ := json.Marshal(valid)
	resp, body, err = authReq("POST", "http://localhost:8080/api/v1/templates/contracts/welcome/validate", vb)
	if err != nil {
		t.Fatalf("validate welcome failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("validate welcome status=%d body=%s", resp.StatusCode, string(body))
	}
	var vresp map[string]interface{}
	if err := json.Unmarshal(body, &vresp); err != nil {
		t.Fatalf("invalid validate response: %v", err)
	}
	vdata := vresp["data"].(map[string]interface{})
	if validRes, ok := vdata["valid"].(bool); !ok || !validRes {
		t.Fatalf("expected valid=true, got %v", vdata["valid"])
	}

	// Validate invalid data (missing LastName)
	invalid := map[string]interface{}{"FirstName": "John"}
	ib, _ := json.Marshal(invalid)
	resp, body, err = authReq("POST", "http://localhost:8080/api/v1/templates/contracts/welcome/validate", ib)
	if err != nil {
		t.Fatalf("validate welcome invalid failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("validate welcome invalid status=%d body=%s", resp.StatusCode, string(body))
	}
	if err := json.Unmarshal(body, &vresp); err != nil {
		t.Fatalf("invalid validate response 2: %v", err)
	}
	vdata = vresp["data"].(map[string]interface{})
	if validRes, ok := vdata["valid"].(bool); !ok || validRes {
		t.Fatalf("expected valid=false, got %v", vdata["valid"])
	}

	// Complex invoice validation
	invoice := map[string]interface{}{
		"Customer":         map[string]interface{}{"Name": "Test Customer Ltd.", "Email": "customer@test.com", "Address": map[string]interface{}{"Street": "123 Test St", "City": "Test City", "PostalCode": "12345", "Country": "Germany"}, "TaxId": "DE123456789"},
		"InvoiceData":      map[string]interface{}{"InvoiceNumber": "INV-2024-001", "IssueDate": "2024-01-15T00:00:00Z", "DueDate": "2024-02-14T00:00:00Z", "Items": []interface{}{map[string]interface{}{"Description": "Consulting Service", "Quantity": 10, "UnitPrice": 150.0, "Total": 1500.0, "VATRate": 19.0, "VATAmount": 285.0}}, "Subtotal": 1500.0, "TotalVAT": 285.0, "Total": 1785.0, "Currency": "EUR", "Language": "en"},
		"OrganizationData": map[string]interface{}{"Name": "Test Organization GmbH", "Address": map[string]interface{}{"Street": "Organization St 1", "City": "Berlin", "PostalCode": "10115", "Country": "Germany"}, "TaxId": "DE987654321", "Email": "billing@testorg.com", "Phone": "+49 30 12345678"},
	}
	ibb, _ := json.Marshal(invoice)
	resp, body, err = authReq("POST", "http://localhost:8080/api/v1/templates/contracts/invoice/validate", ibb)
	if err != nil {
		t.Fatalf("invoice validate failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("invoice validate status=%d body=%s", resp.StatusCode, string(body))
	}
	if err := json.Unmarshal(body, &vresp); err != nil {
		t.Fatalf("invalid invoice validate response: %v", err)
	}
	vdata = vresp["data"].(map[string]interface{})
	if validRes, ok := vdata["valid"].(bool); !ok || !validRes {
		t.Fatalf("expected invoice valid=true, got %v", vdata["valid"])
	}

	// Non-existent contract should return 404
	resp, body, err = authReq("GET", "http://localhost:8080/api/v1/templates/contracts/non_existent_contract", nil)
	if err != nil {
		t.Fatalf("get non-existent contract failed: %v", err)
	}
	if resp.StatusCode != 404 {
		t.Fatalf("expected 404 for non-existent contract, got %d body=%s", resp.StatusCode, string(body))
	}
}
