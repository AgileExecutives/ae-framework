# Template Rendering Tests (renamed)

POST {{host}}/api/v1/templates/render
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "template_key": "welcome",
  "channel": "EMAIL",
  "data": {"FirstName": "Sarah", "LastName": "Johnson", "OrganizationName": "Hurl Testing Corp", "ActivationLink": "https://app.example.com/activate?token=abc123"}
}

HTTP 200
[Asserts]
jsonpath "$.data.content" contains "Sarah Johnson"
jsonpath "$.data.content" contains "Hurl Testing Corp"
jsonpath "$.data.content" contains "activate?token=abc123"

# Additional rendering tests preserved

POST {{host}}/api/v1/templates/render
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "template_key": "password_reset",
  "channel": "EMAIL",
  "data": {"FirstName": "Alice", "LastName": "Wilson", "ResetLink": "https://app.example.com/reset-password?token=xyz789", "ExpirationTime": "24 hours"}
}

HTTP 200
[Asserts]
jsonpath "$.data.content" contains "Alice Wilson"
