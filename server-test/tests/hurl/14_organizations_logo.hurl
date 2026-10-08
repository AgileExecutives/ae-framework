# Test Organization Logo Upload and Retrieval

# Create an organization first
POST http://localhost:8080/api/v1/organizations
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Test Org for Logo",
  "legal_notice_url": "https://example.com/legal",
  "privacy_policy_url": "https://example.com/privacy"
}

HTTP 201
[Captures]
org_id: jsonpath("$.data.id")

###

# Upload a logo to the organization
POST http://localhost:8080/api/v1/organizations/{{org_id}}/logo
Authorization: Bearer {{access_token}}

[MultipartFormData]
logo: file,/tmp/test_logo.png;type=image/png

HTTP 200
[Asserts]
jsonpath("$.success") == true
jsonpath("$.data.logo_content_type") == "image/png"

###

# Retrieve the logo publicly (no auth required)
GET http://localhost:8080/api/v1/organization/logo?organization_id={{org_id}}

HTTP 200
[Asserts]
header "Content-Type" == "image/png"

###

# Upload a print logo
POST http://localhost:8080/api/v1/organizations/{{org_id}}/printlogo
Authorization: Bearer {{access_token}}

[MultipartFormData]
logo_print: file,/tmp/test_logo.png;type=image/png

HTTP 200
[Asserts]
jsonpath("$.success") == true
jsonpath("$.data.logo_print_content_type") == "image/png"

###

# Retrieve the print logo publicly
GET http://localhost:8080/api/v1/organization/printlogo?organization_id={{org_id}}

HTTP 200
[Asserts]
header "Content-Type" == "image/png"
