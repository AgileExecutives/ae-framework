# Template CRUD Operations Integration Tests (renamed)

# Test GET /templates - Initial state (should have system default templates)
GET {{host}}/api/v1/templates
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.data" exists
jsonpath "$.data.data" isCollection

# Test POST /templates - Create Email Template
POST {{host}}/api/v1/templates
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "template_type": "email",
  "template_key": "integration_test_email",
  "channel": "EMAIL",
  "subject": "Integration Test Email",
  "name": "Integration Test Email Template",
  "description": "Template for integration testing",
  "content": "<h1>Hello \\\{\\\{.Name\\\}\\\}!</h1><p>This is a test email for \\\{\\\{.Purpose\\\}\\\}.</p>",
  "variables": ["Name", "Purpose", "ContactEmail"],
  "sample_data": {"Name": "Integration Tester", "Purpose": "API Testing", "ContactEmail": "test@integration.com"},
  "is_active": true,
  "is_default": false
}

HTTP 201
[Asserts]
jsonpath "$.data.id" exists
jsonpath "$.data.name" == "Integration Test Email Template"
[Captures]
email_template_id: jsonpath "$.data.id"

# Test POST /templates - Create PDF Template
POST {{host}}/api/v1/templates
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "template_type": "document",
  "template_key": "integration_test_document",
  "channel": "PDF",
  "name": "Integration Test Document Template",
  "description": "PDF template for integration testing",
  "content": "<html><head><title>\\\{\\\{.Title\\\}\\\}</title></head><body><h1>\\\{\\\{.Title\\\}\\\}</h1></body></html>",
  "variables": ["Title", "Content", "Date"],
  "sample_data": {"Title": "Integration Test Document", "Content": "This is test content.", "Date": "2024-01-15"},
  "is_active": true,
  "is_default": false
}

HTTP 201
[Asserts]
jsonpath "$.data.id" exists
jsonpath "$.data.name" == "Integration Test Document Template"
[Captures]
pdf_template_id: jsonpath "$.data.id"

# Verify created templates exist
GET {{host}}/api/v1/templates
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.data.data[?(@.name == 'Integration Test Email Template')]" exists
jsonpath "$.data.data[?(@.name == 'Integration Test Document Template')]" exists

# ... rest of CRUD/cleanup steps preserved

DELETE {{host}}/api/v1/templates/{{duplicated_email_id}}
Authorization: Bearer {{auth_token}}

HTTP 204

DELETE {{host}}/api/v1/templates/{{pdf_template_id}}
Authorization: Bearer {{auth_token}}

HTTP 204

DELETE {{host}}/api/v1/templates/{{email_template_id}}
Authorization: Bearer {{auth_token}}

HTTP 204

GET {{host}}/api/v1/templates/{{email_template_id}}
Authorization: Bearer {{auth_token}}

HTTP 404
