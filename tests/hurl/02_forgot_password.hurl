# Forgot password (mock) test
POST {{host}}/api/v1/auth/forgot-password
Content-Type: application/json
{
  "email": "testuser@unburdy.de"
}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.message" contains "Password reset"
