package handlers

// APIResponse represents a generic API response for Swagger docs
// This mirrors the common API response structure used across modules.
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents a generic error response for Swagger docs
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Error occurred"`
	Error   string `json:"error,omitempty" example:"detailed error message"`
}
