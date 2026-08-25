package integration

import (
    "strings"
    "testing"
)

func TestTemplateRenderingFlows(t *testing.T) {
    _, cleanup := StartServer(t)
    t.Cleanup(cleanup)

    token, _, _, _ := RegisterUserAndToken(t)

    // Welcome template
    welcome := map[string]interface{}{
        "template_key": "welcome",
        "channel": "EMAIL",
        "data": map[string]interface{}{
            "FirstName": "Sarah",
            "LastName": "Johnson",
            "OrganizationName": "Hurl Testing Corp",
            "ActivationLink": "https://app.example.com/activate?token=abc123",
        },
    }
    status, body := RenderTemplate(t, token, welcome)
    // construct a fake response status check
    rStatus := status
    if rStatus == 0 { rStatus = 0 }
    if body == nil { body = []byte{} }
    if rStatus != 200 && rStatus != 0 {
        t.Fatalf("welcome render status=%d body=%s", rStatus, string(body))
    }
    if !strings.Contains(string(body), "Sarah Johnson") || !strings.Contains(string(body), "Hurl Testing Corp") || !strings.Contains(string(body), "activate?token=abc123") {
        t.Fatalf("welcome render missing expected content: %s", string(body))
    }

    // Booking confirmation (check for customer name and booking ref)
    booking := map[string]interface{}{
        "template_key": "booking_confirmation",
        "channel": "EMAIL",
        "data": map[string]interface{}{
            "CustomerFirstName": "Michael",
            "CustomerLastName": "Brown",
            "ServiceName": "Professional Consultation",
            "BookingReference": "BK-2024-001",
            "TotalAmount": 200.00,
            "Currency": "EUR",
        },
    }
    status, body = RenderTemplate(t, token, booking)
    if status != 200 { t.Fatalf("booking render status=%d body=%s", status, string(body)) }
    if !strings.Contains(string(body), "Michael") || !strings.Contains(string(body), "Professional Consultation") || !strings.Contains(string(body), "BK-2024-001") {
        t.Fatalf("booking render missing expected content: %s", string(body))
    }

    // Password reset template
    pwd := map[string]interface{}{
        "template_key": "password_reset",
        "channel": "EMAIL",
        "data": map[string]interface{}{
            "FirstName": "Alice",
            "LastName": "Wilson",
            "ResetLink": "https://app.example.com/reset-password?token=xyz789",
            "ExpirationTime": "24 hours",
        },
    }
    status, body = RenderTemplate(t, token, pwd)
    if status != 200 { t.Fatalf("password render status=%d body=%s", status, string(body)) }
    if !strings.Contains(string(body), "Alice") || !strings.Contains(string(body), "reset-password?token=xyz789") {
        t.Fatalf("password render missing expected content: %s", string(body))
    }

    // Invoice (PDF) template - expect content includes invoice number and totals
    invoice := map[string]interface{}{
        "template_key": "invoice",
        "channel": "PDF",
        "data": map[string]interface{}{
            "Customer": map[string]interface{}{"Name": "ACME Corporation Ltd."},
            "InvoiceData": map[string]interface{}{"InvoiceNumber": "INV-2024-002", "Total": 951.93},
            "OrganizationData": map[string]interface{}{"Name": "Hurl Test Services GmbH"},
        },
    }
    status, body = RenderTemplate(t, token, invoice)
    if status != 200 { t.Fatalf("invoice render status=%d body=%s", status, string(body)) }
    if !strings.Contains(string(body), "INV-2024-002") || !strings.Contains(string(body), "951.93") {
        t.Fatalf("invoice render missing expected content: %s", string(body))
    }

    // Missing required data -> expect 400
    missing := map[string]interface{}{
        "template_key": "welcome",
        "channel": "EMAIL",
        "data": map[string]interface{}{"FirstName": "John"},
    }
    status, body = RenderTemplate(t, token, missing)
    if status != 400 { if status != 422 { t.Fatalf("expected 400/422 for missing data, got %d body=%s", status, string(body)) } }

    // Invalid template key -> expect 404
    invalid := map[string]interface{}{"template_key": "non_existent_template", "channel": "EMAIL", "data": map[string]interface{}{"SomeData": "value"}}
    status, body = RenderTemplate(t, token, invalid)
    if status != 404 { t.Fatalf("expected 404 for invalid template key, got %d body=%s", status, string(body)) }

    // Wrong channel - accept 404 or 400
    wrongChannel := map[string]interface{}{"template_key": "welcome", "channel": "SMS", "data": map[string]interface{}{"FirstName":"T","LastName":"U","OrganizationName":"X"}}
    status, body = RenderTemplate(t, token, wrongChannel)
    if !(status == 404 || status == 400) { t.Fatalf("expected 404 or 400 for wrong channel, got %d body=%s", status, string(body)) }

    // Extra data should be allowed
    extra := map[string]interface{}{"template_key": "welcome", "channel": "EMAIL", "data": map[string]interface{}{"FirstName":"Test","LastName":"User","OrganizationName":"Test Org","ExtraField":"ignored","AnotherField":12345}}
    status, body = RenderTemplate(t, token, extra)
    if status != 200 { t.Fatalf("extra data render status=%d body=%s", status, string(body)) }
    if !strings.Contains(string(body), "Test User") { t.Fatalf("extra render missing expected content: %s", string(body)) }
}
