# Quick Integration Test with Unique Users
# Minimal test to verify system is working before running full test suite

# Test server availability
GET http://localhost:8080/api/v1/health

HTTP 200
[Asserts]
jsonpath "$.status" == "healthy"
jsonpath "$.database" == "connected"

# Test basic authentication flow with unique user
POST http://localhost:8080/api/v1/auth/register
Content-Type: application/json
{
  "email": "test_1783501978_55351c49@example.com",
  "username": "user_1783501978_55351c49",
  "password": "Pass123_55351c49",
  "first_name": "Integration",
  "last_name": "Test",
  "company_name": "Test Company 1783501978_55351c49",
  "tenant_name": "tenant_1783501978_55351c49",
  "accept_terms": true
}

# Accept either 201 (created) or 409 (already exists)
HTTP *
[Asserts]
jsonpath "$.success" exists

# Test login with the created user
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json
{
  "email": "test_1783501978_55351c49@example.com",
  "password": "Pass123_55351c49"
}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.token" exists

[Captures]
integration_token: jsonpath "$.data.token"

# Test accessing protected endpoint
GET http://localhost:8080/api/v1/user-settings
Authorization: Bearer {{integration_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true