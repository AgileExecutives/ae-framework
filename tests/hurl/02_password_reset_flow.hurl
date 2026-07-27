# Password reset full flow (mock email)

# Trigger forgot-password which should send a mock email
POST {{host}}/api/v1/auth/forgot-password
Content-Type: application/json
{
  "email": "testuser@unburdy.de"
}

HTTP 200
[Asserts]
jsonpath "$.success" == true

# Retrieve latest mock emails and verify the reset link is present
GET {{host}}/api/v1/emails/latest-emails

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data" isCollection
jsonpath "$.data[0].html" contains "new-password?token="
