# 10_calendar.hurl - Calendar Module API Tests
# Testing Calendar and RecurringEvent CRUD operations

# Test variables from auth flow
# Assuming JWT token is available from previous test or can be obtained

# First, let's login to get a JWT token
POST {{base_url}}/api/v1/auth/login
Content-Type: application/json

{
  "email": "testuser@unburdy.de",
  "password": "newpass123"
}

HTTP/1.1 200
[Captures]
jwt_token: jsonpath "$.data.token"

# ==============================
# CALENDAR CRUD TESTS
# ==============================

# 1. Create a new calendar
POST {{base_url}}/api/v1/calendars
Authorization: Bearer {{jwt_token}}
Content-Type: application/json

{
  "title": "My Work Calendar",
  "schedule_days_ahead": 30,
  "schedule_days_notice": 7,
  "schedule_unit": "days"
}

HTTP/1.1 201
[Captures]
calendar_id: jsonpath "$.data.id"
[Asserts]
jsonpath "$.data.title" == "My Work Calendar"

# 2. Get all calendars
GET {{base_url}}/api/v1/calendars
Authorization: Bearer {{jwt_token}}

HTTP/1.1 200
[Asserts]
jsonpath "$.data" isCollection
jsonpath "$.data[0].title" exists
jsonpath "$.data[0].id" exists

# 3. Get specific calendar
GET {{base_url}}/api/v1/calendars/{{calendar_id}}
Authorization: Bearer {{jwt_token}}

HTTP/1.1 200
[Asserts]
jsonpath "$.data.id" == {{calendar_id}}
jsonpath "$.data.title" == "My Work Calendar"

# 4. Update calendar
PUT {{base_url}}/api/v1/calendars/{{calendar_id}}
Authorization: Bearer {{jwt_token}}
Content-Type: application/json

{
  "title": "Updated Work Calendar",
  "schedule_days_ahead": 45,
  "schedule_days_notice": 14,
  "schedule_unit": "days"
}

HTTP/1.1 200
[Asserts]
jsonpath "$.data.title" == "Updated Work Calendar"

# ==============================
# RECURRING EVENT CRUD TESTS
# ==============================

# 5. Create a recurring event
POST {{base_url}}/api/v1/calendar-series
Authorization: Bearer {{jwt_token}}
Content-Type: application/json

{
  "calendar_id": {{calendar_id}},
  "title": "Weekly Team Meeting",
  "participants": ["alice@company.com", "bob@company.com"],
  "interval_type": "weekly",
  "interval_value": 1,
  "start_time": "2025-11-02T09:00:00Z",
  "end_time": "2025-11-02T10:00:00Z",
  "description": "Regular team sync meeting",
  "location": "Conference Room A"
}

HTTP/1.1 201
[Captures]
recurring_event_id: jsonpath "$.data.series.id"
[Asserts]
jsonpath "$.data.series.title" == "Weekly Team Meeting"
jsonpath "$.data.series.interval_type" == "weekly"
jsonpath "$.data.series.interval_value" == 1
jsonpath "$.data.series.start_time" == "2025-11-02T09:00:00Z"
jsonpath "$.data.series.end_time" == "2025-11-02T10:00:00Z"
jsonpath "$.data.series.calendar_id" == {{calendar_id}}
# 6. Get all recurring events
GET {{base_url}}/api/v1/calendar-series
Authorization: Bearer {{jwt_token}}

HTTP/1.1 200
[Asserts]
jsonpath "$.items" isCollection
jsonpath "$.items[0].title" exists
jsonpath "$.items[0].id" exists

# 7. Get specific recurring event with calendar preloaded
GET {{base_url}}/api/v1/calendar-series/{{recurring_event_id}}
Authorization: Bearer {{jwt_token}}

HTTP/1.1 200
[Asserts]
jsonpath "$.data.id" == {{recurring_event_id}}
jsonpath "$.data.title" == "Weekly Team Meeting"
# 8. Update recurring event
PUT {{base_url}}/api/v1/calendar-series/{{recurring_event_id}}
Authorization: Bearer {{jwt_token}}
Content-Type: application/json

{
  "calendar_id": {{calendar_id}},
  "title": "Updated Team Meeting",
  "participants": ["alice@company.com", "bob@company.com", "charlie@company.com"],
  "interval_type": "weekly",
  "interval_value": 2,
  "start_time": "2025-11-03T10:00:00Z",
  "end_time": "2025-11-03T11:00:00Z",
  "description": "Updated team sync meeting",
  "location": "Conference Room B"
}

HTTP/1.1 200
[Asserts]
jsonpath "$.data.title" == "Updated Team Meeting"
jsonpath "$.data.interval_type" == "weekly"
jsonpath "$.data.interval_value" == 2
jsonpath "$.data.start_time" == "2025-11-03T10:00:00Z"
jsonpath "$.data.end_time" == "2025-11-03T11:00:00Z"
# ==============================
# CLEANUP TESTS
# ==============================

# 9. Delete recurring event
DELETE {{base_url}}/api/v1/calendar-series/{{recurring_event_id}}
Authorization: Bearer {{jwt_token}}
Content-Type: application/json

{ "delete_mode": "all" }

HTTP/1.1 200
[Asserts]
jsonpath "$.message" exists

# 10. Verify recurring event is deleted
GET {{base_url}}/api/v1/calendar-series/{{recurring_event_id}}
Authorization: Bearer {{jwt_token}}

HTTP/1.1 404
[Asserts]
jsonpath "$.error" exists

# 11. Delete calendar
DELETE {{base_url}}/api/v1/calendars/{{calendar_id}}
Authorization: Bearer {{jwt_token}}

HTTP/1.1 200
[Asserts]
jsonpath "$.message" exists

# 12. Verify calendar is deleted
GET {{base_url}}/api/v1/calendars/{{calendar_id}}
Authorization: Bearer {{jwt_token}}

HTTP/1.1 404
[Asserts]
jsonpath "$.error" exists

# ==============================
# ERROR HANDLING TESTS
# ==============================

# 13. Try to create recurring event without calendar_id
POST {{base_url}}/api/v1/calendar-series
Authorization: Bearer {{jwt_token}}
Content-Type: application/json

{
  "title": "Invalid Meeting",
  "weekday": "monday",
  "interval": "weekly",
  "start_time": "09:00",
  "end_time": "10:00"
}

HTTP/1.1 400
[Asserts]
jsonpath "$.error" exists

# 14. Try to access non-existent calendar
GET {{base_url}}/api/v1/calendars/99999
Authorization: Bearer {{jwt_token}}

HTTP/1.1 404
[Asserts]
jsonpath "$.error" exists

# 15. Try to access without authentication
GET {{base_url}}/api/v1/calendars

HTTP/1.1 401
[Asserts]
jsonpath "$.error" exists