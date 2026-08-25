package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"context"

	templateentities "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/entities"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Template represents a minimal template metadata struct
type Template struct {
	ID   uint
	Name string
}

// TemplateService provides template rendering and retrieval.
type TemplateService struct{}

func NewTemplateService() *TemplateService { return &TemplateService{} }

func (s *TemplateService) RenderTemplate(ctx context.Context, tenantID uint, templateID uint, data interface{}) (string, error) {
	// minimal renderer: return placeholder HTML
	return "<html><body>Rendered</body></html>", nil
}

func (s *TemplateService) GetTemplate(ctx context.Context, tenantID uint, templateID uint) (*Template, error) {
	return &Template{ID: templateID, Name: "default"}, nil
}

// ListTemplates returns available templates for a tenant.
func (s *TemplateService) ListTemplates(ctx context.Context, tenantID uint) ([]*Template, error) {
	return []*Template{}, nil
}

// CreateTemplateRequest represents the payload used to create a new template
// in migration/bootstrapping code. Kept lightweight for compatibility.
type CreateTemplateRequest struct {
	TenantID       uint
	OrganizationID *uint
	Module         string
	TemplateKey    string
	Channel        string
	TemplateType   string
	Name           string
	Description    string
	Subject        *string
	Content        string
	IsActive       bool
	IsDefault      bool
}

// CreateTemplate persists a new template. For tests we only need it to be a noop
// that returns a simple Template with an ID.
func (s *TemplateService) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*Template, error) {
	return &Template{ID: 1, Name: req.Name}, nil
}

// CopyTemplatesFromTenant2Org2 copies default templates from tenant 2/org 2 to a new org.
// This is a lightweight stub used by the server-test harness.
func (s *TemplateService) CopyTemplatesFromTenant2Org2(ctx context.Context, tenantID uint, orgID uint) error {
	// In the real implementation this would duplicate templates and storage entries.
	return nil
}

// ContractRegistrar is a helper for registering template contracts.
type ContractRegistrar struct {
	db *gorm.DB
}

// NewContractRegistrar creates a new ContractRegistrar bound to a DB.
func NewContractRegistrar(db *gorm.DB) *ContractRegistrar { return &ContractRegistrar{db: db} }

// RegisterContractFromFile reads a JSON contract file and inserts a
// TemplateContract record if one does not already exist for the module/key.
func (r *ContractRegistrar) RegisterContractFromFile(tenantID uint, module string, path string) error {
	if r.db == nil {
		return fmt.Errorf("contract registrar has no db")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read contract file %s: %w", path, err)
	}
	return r.RegisterContractFromBytes(tenantID, module, path, b)
}

// RegisterContractFromBytes parses the provided JSON bytes and inserts a
// TemplateContract record if one does not already exist for the tenant/module/key.
func (r *ContractRegistrar) RegisterContractFromBytes(tenantID uint, module string, pathOrHint string, b []byte) error {
	if r.db == nil {
		return fmt.Errorf("contract registrar has no db")
	}

	// Try to parse the file as generic JSON
	var payload map[string]interface{}
	if err := json.Unmarshal(b, &payload); err != nil {
		// If not an object, treat entire file as the variable schema
		payload = map[string]interface{}{"_raw": nil}
	}
	// Derive template key: prefer explicit fields, otherwise filename minus suffix
	var templateKey string
	if v, ok := payload["template_key"].(string); ok && v != "" {
		templateKey = v
	} else if v, ok := payload["templateKey"].(string); ok && v != "" {
		templateKey = v
	} else {
		base := filepath.Base(pathOrHint)
		templateKey = strings.TrimSuffix(base, "-contract.json")
		templateKey = strings.TrimSuffix(templateKey, ".json")
	}

	// Extract variable schema: look for common keys, otherwise assume whole file is schema
	var variableSchema datatypes.JSON
	if v, ok := payload["variable_schema"]; ok {
		if bs, err := json.Marshal(v); err == nil {
			variableSchema = datatypes.JSON(bs)
		}
	} else if v, ok := payload["variableSchema"]; ok {
		if bs, err := json.Marshal(v); err == nil {
			variableSchema = datatypes.JSON(bs)
		}
	} else if _, hasType := payload["type"]; hasType {
		// treat the whole document as the schema
		variableSchema = datatypes.JSON(b)
	} else {
		// fallback empty object
		variableSchema = datatypes.JSON([]byte(`{}`))
	}

	// Extract default sample data if present
	var defaultSample datatypes.JSON
	if v, ok := payload["default_sample_data"]; ok {
		if bs, err := json.Marshal(v); err == nil {
			defaultSample = datatypes.JSON(bs)
		}
	} else if v, ok := payload["defaultSampleData"]; ok {
		if bs, err := json.Marshal(v); err == nil {
			defaultSample = datatypes.JSON(bs)
		}
	} else if v, ok := payload["example"]; ok {
		if bs, err := json.Marshal(v); err == nil {
			defaultSample = datatypes.JSON(bs)
		}
	} else {
		defaultSample = datatypes.JSON([]byte(`{}`))
	}

	// Check existing contract for this tenant
	var existingCount int64
	if err := r.db.Table("template_contracts").Where("tenant_id = ? AND module = ? AND template_key = ?", tenantID, module, templateKey).Count(&existingCount).Error; err != nil {
		return fmt.Errorf("count existing contracts: %w", err)
	}
	if existingCount > 0 {
		return nil
	}

	// Insert contract
	contract := templateentities.TemplateContract{
		TenantID:          tenantID,
		Module:            module,
		TemplateKey:       templateKey,
		VariableSchema:    variableSchema,
		DefaultSampleData: defaultSample,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := r.db.Create(&contract).Error; err != nil {
		return fmt.Errorf("create contract %s/%s: %w", module, templateKey, err)
	}
	return nil
}
