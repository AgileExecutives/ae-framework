package templates

// Local swagger-friendly types for templates documentation.
// They mirror the structures used in shared-modules but are kept local
// so swag can resolve them when generating docs for this module.

type CreateTemplateRequest struct {
	OrganizationID *uint                  `json:"organization_id,omitempty"`
	TemplateType   string                 `json:"template_type"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description,omitempty"`
	Content        string                 `json:"content"`
	Variables      []string               `json:"variables,omitempty"`
	SampleData     map[string]interface{} `json:"sample_data,omitempty"`
	IsActive       bool                   `json:"is_active,omitempty"`
	IsDefault      bool                   `json:"is_default,omitempty"`
}

type UpdateTemplateRequest struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Content     *string                `json:"content,omitempty"`
	Variables   []string               `json:"variables,omitempty"`
	SampleData  map[string]interface{} `json:"sample_data,omitempty"`
	IsActive    *bool                  `json:"is_active,omitempty"`
	IsDefault   *bool                  `json:"is_default,omitempty"`
}

type TemplateResponse struct {
	ID             uint                   `json:"id"`
	TenantID       uint                   `json:"tenant_id"`
	OrganizationID *uint                  `json:"organization_id,omitempty"`
	TemplateType   string                 `json:"template_type"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description,omitempty"`
	Version        int                    `json:"version,omitempty"`
	IsActive       bool                   `json:"is_active,omitempty"`
	IsDefault      bool                   `json:"is_default,omitempty"`
	StorageKey     string                 `json:"storage_key,omitempty"`
	Variables      []string               `json:"variables,omitempty"`
	SampleData     map[string]interface{} `json:"sample_data,omitempty"`
}

type TemplateWithContentResponse struct {
	Template TemplateResponse `json:"template"`
	Content  string           `json:"content"`
}

type ListTemplatesResponse struct {
	Success    bool               `json:"success"`
	Data       []TemplateResponse `json:"data"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

type RenderTemplateRequest struct {
	Data map[string]interface{} `json:"data"`
}

type DuplicateTemplateRequest struct {
	Name string `json:"name"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
