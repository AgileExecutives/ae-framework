# Template Contract API Tests (renamed)

# Test GET /templates/contracts - List all contracts
GET {{host}}/api/v1/templates/contracts
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists
jsonpath "$.data" isCollection

# Test GET /templates/contracts/welcome - Get specific contract
GET {{host}}/api/v1/templates/contracts/welcome
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.data.template_key" == "welcome"
jsonpath "$.data.variable_schema" exists

# Test sample-data and validation
GET {{host}}/api/v1/templates/contracts/welcome/sample-data
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists
jsonpath "$.data.FirstName" exists

POST {{host}}/api/v1/templates/contracts/welcome/validate
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "FirstName": "John",
  "LastName": "Doe",
  "OrganizationName": "Test Company"
}

HTTP 200
[Asserts]
jsonpath "$.data.valid" == true

POST {{host}}/api/v1/templates/contracts/welcome/validate
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "FirstName": "John"
}

HTTP 200
[Asserts]
jsonpath "$.data.valid" == false
jsonpath "$.data.errors" isCollection
