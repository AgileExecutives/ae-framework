# Password reset full flow (mock email)

# Trigger forgot-password which should send a mock email
POST http://localhost:8080/api/v1/auth/forgot-password
Content-Type: application/json
{
  "email": "test_1787651213_104b3454@example.com"
}

HTTP 200
[Asserts]
jsonpath "$.success" == true

# Retrieve latest mock emails and verify the reset link is present
GET http://localhost:8080/api/v1/emails/latest-emails

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data" isCollection
jsonpath "$.data[0].html" contains "new-password?token="
