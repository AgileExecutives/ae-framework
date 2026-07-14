package handlers

// ErrorResponse represents a generic error response for Swagger documentation
// This mirrors the external ErrorResponse shape so swag can resolve types when
// generating docs for this module.
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Unauthorized"`
	Error   string `json:"error,omitempty" example:"detailed error message"`
}
