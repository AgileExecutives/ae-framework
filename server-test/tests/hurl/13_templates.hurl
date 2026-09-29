// Combined templates smoke tests (renamed from templates.hurl)

# Test GET /templates - List all templates
GET {{host}}/api/v1/templates
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.data" exists
jsonpath "$.data.data" isCollection
jsonpath "$.data.pagination" exists

// Additional create/get/put/render/duplicate/default/contracts/delete steps
// (reused content from previous templates.hurl)

POST {{host}}/api/v1/templates
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "template_type": "email",
  "template_key": "test_welcome",
  "channel": "EMAIL",
  "subject": "Welcome to Test Platform",
  "name": "Test Welcome Template",
  "description": "A test welcome email template",
  "content": "<h1>Welcome \\\{\\\{.FirstName\\\}\\\} \\\{\\\{.LastName\\\}\\\}!</h1><p>Thank you for joining \\\{\\\{.OrganizationName\\\}\\\}.</p>",
  "variables": ["FirstName", "LastName", "OrganizationName"],
  "sample_data": {
    "FirstName": "John",
    "LastName": "Doe",
    "OrganizationName": "Test Company"
  },
  "is_active": true,
  "is_default": false
}

HTTP 201
[Asserts]
jsonpath "$.data.id" exists
jsonpath "$.data.name" == "Test Welcome Template"
[Captures]
template_id: jsonpath "$.data.id"

GET {{host}}/api/v1/templates/{{template_id}}
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.data.id" == {{template_id}}
jsonpath "$.data.name" == "Test Welcome Template"

PUT {{host}}/api/v1/templates/{{template_id}}
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "name": "Updated Welcome Template"
}

HTTP 200
[Asserts]
jsonpath "$.data.name" == "Updated Welcome Template"

POST {{host}}/api/v1/templates/{{template_id}}/render
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{
  "data": {"FirstName": "Jane", "LastName": "Smith", "OrganizationName": "Hurl Test Company"}
}

HTTP 200
[Asserts]
jsonpath "$.data.content" contains "Jane Smith"

POST {{host}}/api/v1/templates/{{template_id}}/duplicate
Authorization: Bearer {{auth_token}}
Content-Type: application/json

{"name":"Duplicated Welcome Template"}

HTTP 201
[Asserts]
jsonpath "$.data.name" == "Duplicated Welcome Template"
[Captures]
duplicated_template_id: jsonpath "$.data.id"

GET {{host}}/api/v1/templates/default?template_type=email&channel=EMAIL
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists
jsonpath "$.data.is_default" == true

GET {{host}}/api/v1/templates/contracts
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists
jsonpath "$.data" isCollection

GET {{host}}/api/v1/templates/contracts/by-key/welcome
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.data.template_key" == "welcome"

DELETE {{host}}/api/v1/templates/{{template_id}}
Authorization: Bearer {{auth_token}}

HTTP 204

DELETE {{host}}/api/v1/templates/{{duplicated_template_id}}
Authorization: Bearer {{auth_token}}

HTTP 204

GET {{host}}/api/v1/templates/{{template_id}}
Authorization: Bearer {{auth_token}}

HTTP 404
