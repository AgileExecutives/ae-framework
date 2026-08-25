# Full password reset flow: trigger, extract token, reset, login

# Trigger forgot-password
POST http://localhost:8080/api/v1/auth/forgot-password
Content-Type: application/json
{
  "email": "test_1787651213_104b3454@example.com"
}

HTTP 200
[Asserts]
jsonpath "$.success" == true

# Retrieve latest mock emails
GET http://localhost:8080/api/v1/emails/latest-emails

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data" isCollection

[Captures]
# Extract the reset token from the email HTML link
reset_token: regex /new-password\?token=([A-Za-z0-9\-\._~%]+)/

# Use the captured token to reset the password (path param)
POST http://localhost:8080/api/v1/auth/new-password/{{reset_token}}
Content-Type: application/json
{
  "new_password": "HurlFlow1!"
}

HTTP 200
[Asserts]
jsonpath "$.success" == true

# Login with the new password
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json
{
  "email": "test_1787651213_104b3454@example.com",
  "password": "HurlFlow1!"
}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.message" == "Login successful"
