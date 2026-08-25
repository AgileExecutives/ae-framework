# Forgot password (mock) test
POST http://localhost:8080/api/v1/auth/forgot-password
Content-Type: application/json
{
  "email": "test_1787651213_104b3454@example.com"
}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.message" exists
