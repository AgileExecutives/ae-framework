package handlers

// APIResponse is a minimal API response type used for Swagger generation
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"ok"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse used for Swagger error responses
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"error message"`
	Error   string `json:"error,omitempty" example:"internal error details"`
}

// ListResponse is a minimal list wrapper for Swagger
type ListResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"ok"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int         `json:"total"`
}

// PaginationResponse placeholder for swagger
type PaginationResponse struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}
