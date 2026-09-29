package templates

// Package templates provides a lightweight templates module used by tests and
// optionally by apps that want an in-process templates provider.
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/AgileExecutives/ae-framework/serverbase/module"
	templatedocs "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/docs"
	templateentities "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/entities"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/templates/services"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// NewTemplatesModule returns a minimal module that exposes template endpoints
// under /templates. It keeps an in-memory store suitable for the test harness.
func NewTemplatesModule() core.Module {
	return module.NewAdapterModule("templates", "0.1.0", []string{},
		module.WithEntities(&templateEntity{}, &templateContractEntity{}),
		module.WithRoutes(&templatesRouteProvider{}),
		module.WithServices(
			&templateServiceProvider{},
			&contractRegistrarProvider{},
		),
		module.WithInit(func(ctx core.ModuleContext) error {
			// Register generated swagger docs so /templates appears in the merged spec.
			if ctx.DocRegistry != nil {
				ctx.DocRegistry.RegisterDoc("templates", templatedocs.SwaggerInfo.ReadDoc())
			}
			// Ensure tables exist but do not seed tenant-scoped template rows here.
			// Per-tenant template seeding happens via RegisterContracts callback
			// (invoked by the tenant post-create hook or startup bootstrap).
			return ensureStandardTemplates(ctx.DB)
		}),
		// Register contracts and per-tenant templates when a tenant is created.
		module.WithContractRegistration(func(ctx core.ModuleContext, tenantID uint) error {
			return registerTemplatesForTenant(ctx, tenantID)
		}),
	)
}

type templateEntity struct{}

func (e *templateEntity) TableName() string               { return "templates" }
func (e *templateEntity) GetModel() interface{}           { return &templateentities.Template{} }
func (e *templateEntity) GetMigrations() []core.Migration { return nil }

type templateContractEntity struct{}

func (e *templateContractEntity) TableName() string               { return "template_contracts" }
func (e *templateContractEntity) GetModel() interface{}           { return &templateentities.TemplateContract{} }
func (e *templateContractEntity) GetMigrations() []core.Migration { return nil }

func ensureStandardTemplates(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	if err := db.AutoMigrate(&templateentities.Template{}, &templateentities.TemplateContract{}); err != nil {
		return fmt.Errorf("migrate template tables: %w", err)
	}
	// Do not create tenant-scoped template rows here. Per-tenant seeding
	// happens in `RegisterContractsForTenant` via the module's
	// RegisterContracts callback.

	// Template contracts are registered per-tenant during startup using the
	// centralized contract registration logic. Do not seed tenant-scoped
	// `template_contracts` rows here (they would be tenantless/zero-valued).
	// Modules that own contracts should implement `RegisterContracts` or
	// provide contract files under `modules/*/contracts` or
	// `shared-modules/*/contracts` so `RegisterAllContracts` can register
	// them for each tenant.
	return nil
}

// registerTemplatesForTenant ensures the standard templates exist for the
// given tenant. It is invoked as part of per-tenant contract registration.
func registerTemplatesForTenant(ctx core.ModuleContext, tenantID uint) error {
	db := ctx.DB
	if db == nil {
		return nil
	}

	seedTemplates := []templateentities.Template{
		{
			TenantID:     tenantID,
			Module:       "user",
			TemplateKey:  "welcome",
			Channel:      templateentities.ChannelEmail,
			Name:         "Welcome Email",
			Description:  "Default welcome email",
			StorageKey:   "templates/welcome.html",
			Version:      1,
			IsActive:     true,
			IsDefault:    true,
			Variables:    datatypes.JSON([]byte(`["FirstName","LastName","OrganizationName"]`)),
			SampleData:   datatypes.JSON([]byte(`{"FirstName":"Test","LastName":"User","OrganizationName":"Server Test Organization"}`)),
			TemplateType: "email",
			Subject:      strPtr("Welcome to our service"),
		},
		{
			TenantID:     tenantID,
			Module:       "user",
			TemplateKey:  "password_reset",
			Channel:      templateentities.ChannelEmail,
			Name:         "Password Reset Email",
			Description:  "Default password reset email",
			StorageKey:   "templates/password_reset.html",
			Version:      1,
			IsActive:     true,
			IsDefault:    true,
			Variables:    datatypes.JSON([]byte(`["FirstName","ResetURL"]`)),
			SampleData:   datatypes.JSON([]byte(`{"FirstName":"Test","ResetURL":"http://localhost:5173/reset"}`)),
			TemplateType: "email",
			Subject:      strPtr("Reset your password"),
		},
	}

	for _, t := range seedTemplates {
		var existing templateentities.Template
		if err := db.Where("tenant_id = ? AND template_key = ? AND module = ?", tenantID, t.TemplateKey, t.Module).First(&existing).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("lookup tenant %d template %s: %w", tenantID, t.TemplateKey, err)
			}
			if err := db.Create(&t).Error; err != nil {
				return fmt.Errorf("create tenant %d template %s: %w", tenantID, t.TemplateKey, err)
			}
		}
	}
	return nil
}

func strPtr(s string) *string { return &s }

// formatValue returns a string representation for common scalar types.
func formatValue(v interface{}) string {
	switch val := v.(type) {
	case float64:
		return fmt.Sprintf("%.2f", val)
	case float32:
		return fmt.Sprintf("%.2f", val)
	case int:
		return fmt.Sprintf("%d", val)
	case int32:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}

// lookupPayloadValue resolves a dotted path against the request payload's
// `data` object. It supports nested maps and arrays. For arrays of objects
// it will try to join `Description` fields or stringify elements.
func lookupPayloadValue(payload map[string]interface{}, path string) (string, bool) {
	if payload == nil {
		return "", false
	}
	d, ok := payload["data"].(map[string]interface{})
	if !ok || d == nil {
		return "", false
	}
	segments := strings.Split(path, ".")
	var cur interface{} = d
	for i, seg := range segments {
		last := i == len(segments)-1
		switch c := cur.(type) {
		case map[string]interface{}:
			next, ok := c[seg]
			if !ok {
				return "", false
			}
			if last {
				return formatValue(next), true
			}
			cur = next
		case []interface{}:
			// If this is the last segment, try to render the array
			if last {
				parts := []string{}
				for _, elem := range c {
					if em, ok := elem.(map[string]interface{}); ok {
						if desc, ok := em["Description"].(string); ok {
							parts = append(parts, desc)
							continue
						}
					}
					parts = append(parts, formatValue(elem))
				}
				return strings.Join(parts, ", "), true
			}
			// otherwise try first element as representative
			if len(c) > 0 {
				cur = c[0]
				continue
			}
			return "", false
		default:
			return "", false
		}
	}
	return "", false
}

// renderTemplateString replaces dot-style placeholders in `content` using
// values from `dataMap` and falling back to `payload` dotted-path lookup.
func renderTemplateString(content string, payload map[string]interface{}, dataMap map[string]string) string {
	content = strings.ReplaceAll(content, "\\{\\{", "{{")
	content = strings.ReplaceAll(content, "\\}\\}", "}}")
	re := regexp.MustCompile(`\{\{\s*\.([A-Za-z0-9_.]+)\s*\}\}`)
	rendered := re.ReplaceAllStringFunc(content, func(m string) string {
		parts := re.FindStringSubmatch(m)
		if len(parts) < 2 {
			return ""
		}
		key := parts[1]
		if v, ok := dataMap[key]; ok {
			return v
		}
		if v, ok := lookupPayloadValue(payload, key); ok {
			return v
		}
		return ""
	})
	return rendered
}

// readStorageKey attempts to read a StorageKey path. It tries the provided
// path, and if it isn't readable, also attempts to resolve it relative to common
// locations (repo root, templates directory, etc.). This helps tests that store
// relative paths like `templates/welcome.html`.
func readStorageKey(pathKey string) (string, error) {
	if pathKey == "" {
		return "", fmt.Errorf("empty storage key")
	}

	candidates := []string{
		pathKey,                             // Try as-is first (might be absolute or relative to cwd)
		filepath.Join("templates", pathKey), // Try under templates/ subdirectory
		filepath.Join("..", "server-test", pathKey), // Try from parent dir
		filepath.Join("..", pathKey),                // Try from parent dir directly
	}

	for _, candidate := range candidates {
		if b, err := os.ReadFile(candidate); err == nil {
			return string(b), nil
		}
	}

	return "", fmt.Errorf("could not read storage key %s (tried %d candidates)", pathKey, len(candidates))
}

// Service providers to expose TemplateService and ContractRegistrar via the
// central service registry so other modules (e.g., client_management) can look them up.
type templateServiceProvider struct{}

func (p *templateServiceProvider) ServiceName() string {
	return "template_service"
}

func (p *templateServiceProvider) ServiceInterface() interface{} {
	return (*services.TemplateService)(nil)
}

func (p *templateServiceProvider) Factory(ctx core.ModuleContext) (interface{}, error) {
	svc := services.NewTemplateService()
	return svc, nil
}

type contractRegistrarProvider struct{}

func (p *contractRegistrarProvider) ServiceName() string {
	return "contract-registrar"
}

func (p *contractRegistrarProvider) ServiceInterface() interface{} {
	return (*services.ContractRegistrar)(nil)
}

func (p *contractRegistrarProvider) Factory(ctx core.ModuleContext) (interface{}, error) {
	r := services.NewContractRegistrar(ctx.DB)
	return r, nil
}

type templatesRouteProvider struct{}

func (r *templatesRouteProvider) GetPrefix() string                { return "" }
func (r *templatesRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *templatesRouteProvider) GetSwaggerTags() []string         { return []string{"templates"} }

func (r *templatesRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	svc := services.NewTemplateService()

	// DB-backed templates routes
	templates := router.Group("/templates")

	templates.POST("/render", func(c *gin.Context) {
		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("invalid json", "request body invalid"))
			return
		}

		key, _ := payload["template_key"].(string)
		if key == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("template_key required", "template_key missing"))
			return
		}

		channel, _ := payload["channel"].(string)

		db := ctx.DB
		if db == nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("no db", "database not available"))
			return
		}

		var t templateentities.Template
		q := db.Where("template_key = ?", key).Order("storage_key <> '' DESC").Order("is_default DESC")
		if channel != "" {
			q = q.Where("LOWER(channel) = ?", strings.ToLower(strings.TrimSpace(channel)))
		}
		if err := q.First(&t).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("template not found", "no template with that key"))
			return
		}

		// Validate against simple contract rules where applicable
		if dataMapPayload, ok := payload["data"].(map[string]interface{}); ok {
			valid, errs := validateContractForKey(t.TemplateKey, dataMapPayload)
			if !valid {
				c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("validation failed", strings.Join(errs, "; ")))
				return
			}
		}

		// Try to read template content from storage key if present
		var content string
		if t.StorageKey != "" {
			if s, err := readStorageKey(t.StorageKey); err == nil {
				content = s
			} else {
				log.Printf("templates: read storage key failed: %v", err)
			}
		}

		// Build data map for rendering
		dataMap := map[string]string{}
		if d, ok := payload["data"].(map[string]interface{}); ok {
			for k, v := range d {
				switch val := v.(type) {
				case float32, float64:
					dataMap[k] = fmt.Sprintf("%.2f", val)
				case int, int32, int64:
					dataMap[k] = fmt.Sprintf("%d", val)
				case map[string]interface{}:
					for ik, iv := range val {
						dataMap[k+"."+ik] = fmt.Sprintf("%v", iv)
					}
				default:
					dataMap[k] = fmt.Sprintf("%v", val)
				}
			}
		}

		rendered := ""
		if content != "" {
			rendered = renderTemplateString(content, payload, dataMap)
		}
		if strings.TrimSpace(rendered) == "" {
			// Fallback to service renderer
			if html, err := svc.RenderTemplate(c.Request.Context(), 1, t.ID, payload["data"]); err == nil {
				rendered = html
			}
		}

		// Render subject if available
		var subj string
		if t.Subject != nil {
			subj = renderTemplateString(*t.Subject, payload, dataMap)
		}

		if subj != "" {
			c.JSON(http.StatusOK, models.SuccessResponse("rendered", gin.H{"content": rendered, "subject": subj}))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse("rendered", gin.H{"content": rendered}))
	})

	templates.GET("", func(c *gin.Context) {
		db := ctx.DB
		if db == nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("no db", "database not available"))
			return
		}
		var items []templateentities.Template
		q := db.Model(&templateentities.Template{})
		if ttype := c.Query("template_type"); ttype != "" {
			q = q.Where("LOWER(template_type) = ?", strings.ToLower(ttype))
		}
		if channel := c.Query("channel"); channel != "" {
			q = q.Where("LOWER(channel) = ?", strings.ToLower(strings.TrimSpace(channel)))
		}
		if err := q.Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("db error", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.SuccessListResponse(items, 1, len(items), len(items)))
	})

	templates.POST("", func(c *gin.Context) {
		db := ctx.DB
		if db == nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("no db", "database not available"))
			return
		}
		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("invalid json", "request body invalid"))
			return
		}
		t := templateentities.Template{
			TenantID:     1,
			Module:       "user",
			TemplateType: "email",
			Channel:      templateentities.ChannelEmail,
			Name:         "",
			Description:  "",
			IsActive:     true,
			IsDefault:    false,
		}
		if v, ok := payload["template_type"].(string); ok {
			t.TemplateType = v
		}
		if v, ok := payload["channel"].(string); ok {
			t.Channel = templateentities.Channel(strings.ToUpper(strings.TrimSpace(v)))
		}
		if v, ok := payload["template_key"].(string); ok {
			t.TemplateKey = v
		}
		if v, ok := payload["name"].(string); ok {
			t.Name = v
		}
		if v, ok := payload["description"].(string); ok {
			t.Description = v
		}
		if v, ok := payload["is_active"].(bool); ok {
			t.IsActive = v
		}
		if v, ok := payload["is_default"].(bool); ok {
			t.IsDefault = v
		}
		if v, ok := payload["subject"].(string); ok {
			t.Subject = &v
		}
		if v, ok := payload["variables"]; ok {
			if bs, err := json.Marshal(v); err == nil {
				t.Variables = datatypes.JSON(bs)
			}
		}
		if v, ok := payload["sample_data"]; ok {
			if bs, err := json.Marshal(v); err == nil {
				t.SampleData = datatypes.JSON(bs)
			}
		}
		// If content provided, persist to server-test/templates and set StorageKey
		if v, ok := payload["content"].(string); ok && v != "" {
			dir := "server-test/templates"
			_ = os.MkdirAll(dir, 0o755)
			fname := fmt.Sprintf("%s/template_%d.html", dir, time.Now().UnixNano())
			if err := os.WriteFile(fname, []byte(v), 0o644); err == nil {
				t.StorageKey = fname
			}
		}

		if err := db.Create(&t).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("db error", err.Error()))
			return
		}
		c.JSON(http.StatusCreated, models.SuccessResponse("created", t))
	})

	templates.GET("/default", func(c *gin.Context) {
		db := ctx.DB
		if db == nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("no db", "database not available"))
			return
		}
		ttype := c.Query("template_type")
		channel := c.Query("channel")
		var t templateentities.Template
		q := db.Where("is_default = ?", true)
		if ttype != "" {
			q = q.Where("template_type = ?", ttype)
		}
		if channel != "" {
			q = q.Where("channel = ?", channel)
		}
		if err := q.First(&t).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("template not found", "no template with that id"))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse("retrieved", t))
	})

	templates.GET("/:id", func(c *gin.Context) {
		db := ctx.DB
		if db == nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("no db", "database not available"))
			return
		}
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		var t templateentities.Template
		if err := db.First(&t, id).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("template not found", "no template with that id"))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse("retrieved", t))
	})

	templates.PUT("/:id", func(c *gin.Context) {
		db := ctx.DB
		if db == nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("no db", "database not available"))
			return
		}
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("invalid json", "request body invalid"))
			return
		}
		var t templateentities.Template
		if err := db.First(&t, id).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("template not found", "no template with that id"))
			return
		}
		if v, ok := payload["name"].(string); ok {
			t.Name = v
		}
		if v, ok := payload["description"].(string); ok {
			t.Description = v
		}
		if v, ok := payload["is_active"].(bool); ok {
			t.IsActive = v
		}
		if v, ok := payload["is_default"].(bool); ok {
			t.IsDefault = v
		}
		if v, ok := payload["variables"]; ok {
			if bs, err := json.Marshal(v); err == nil {
				t.Variables = datatypes.JSON(bs)
			}
		}
		if v, ok := payload["sample_data"]; ok {
			if bs, err := json.Marshal(v); err == nil {
				t.SampleData = datatypes.JSON(bs)
			}
		}
		if err := db.Save(&t).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("db error", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse("updated", t))
	})

	templates.POST("/:id/render", func(c *gin.Context) {
		db := ctx.DB
		if db == nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("no db", "database not available"))
			return
		}
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		var payload map[string]interface{}
		_ = c.ShouldBindJSON(&payload)
		var t templateentities.Template
		if err := db.First(&t, id).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("template not found", "no template with that id"))
			return
		}
		var content string
		if t.StorageKey != "" {
			if b, err := os.ReadFile(t.StorageKey); err == nil {
				content = string(b)
			}
		}
		dataMap := map[string]string{}
		if d, ok := payload["data"].(map[string]interface{}); ok {
			for k, v := range d {
				dataMap[k] = fmt.Sprintf("%v", v)
			}
		}
		if content != "" {
			rendered := renderTemplateString(content, payload, dataMap)
			if t.Subject != nil {
				subj := renderTemplateString(*t.Subject, payload, dataMap)
				c.JSON(http.StatusOK, models.SuccessResponse("rendered", gin.H{"content": rendered, "subject": subj}))
				return
			}
			c.JSON(http.StatusOK, models.SuccessResponse("rendered", gin.H{"content": rendered}))
			return
		}
		// Fallback
		if html, err := svc.RenderTemplate(c.Request.Context(), 1, id, payload["data"]); err == nil {
			c.JSON(http.StatusOK, models.SuccessResponse("rendered", gin.H{"content": html}))
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("render error", "could not render template"))
	})

	templates.DELETE("/:id", func(c *gin.Context) {
		db := ctx.DB
		if db == nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("no db", "database not available"))
			return
		}
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		if err := db.Delete(&templateentities.Template{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("db error", err.Error()))
			return
		}
		c.JSON(http.StatusNoContent, models.SuccessMessageResponse("deleted"))
	})

	templates.POST("/:id/duplicate", func(c *gin.Context) {
		db := ctx.DB
		if db == nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("no db", "database not available"))
			return
		}
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		var payload map[string]interface{}
		_ = c.ShouldBindJSON(&payload)
		var src templateentities.Template
		if err := db.First(&src, id).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponseFunc("template not found", "no template with that id"))
			return
		}
		copyRec := src
		copyRec.ID = 0
		if name, ok := payload["name"].(string); ok {
			copyRec.Name = name
		}
		if key, ok := payload["template_key"].(string); ok {
			copyRec.TemplateKey = key
		}
		if err := db.Create(&copyRec).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("db error", err.Error()))
			return
		}
		c.JSON(http.StatusCreated, models.SuccessResponse("created", copyRec))
	})

	// Contracts and helper endpoints used by tests (DB-backed)
	templates.GET("/contracts", func(c *gin.Context) {
		db := ctx.DB
		if db == nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("no db", "database not available"))
			return
		}
		var contracts []templateentities.TemplateContract
		if err := db.Find(&contracts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("db error", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse("contracts", contracts))
	})

	templates.GET("/contracts/by-key/:key", func(c *gin.Context) {
		key := c.Param("key")
		writeContractByKey(c, key)
	})

	templates.GET("/contracts/:key", func(c *gin.Context) {
		key := c.Param("key")
		writeContractByKey(c, key)
	})

	templates.GET("/contracts/:key/sample-data", func(c *gin.Context) {
		key := c.Param("key")
		db := ctx.DB
		if db != nil {
			var contract templateentities.TemplateContract
			if err := db.Where("template_key = ?", key).First(&contract).Error; err == nil {
				var m map[string]interface{}
				_ = json.Unmarshal([]byte(contract.DefaultSampleData), &m)
				c.JSON(http.StatusOK, models.SuccessResponse("sample data", m))
				return
			}
		}
		// Fallbacks
		switch key {
		case "welcome":
			c.JSON(http.StatusOK, models.SuccessResponse("sample data", gin.H{"FirstName": "John", "LastName": "Doe"}))
		default:
			c.JSON(http.StatusOK, models.SuccessResponse("sample data", gin.H{}))
		}
	})

	templates.POST("/contracts/:key/validate", func(c *gin.Context) {
		key := c.Param("key")
		var payload map[string]interface{}
		_ = c.ShouldBindJSON(&payload)
		if key == "welcome" {
			errs := []string{}
			if _, ok := payload["FirstName"]; !ok {
				errs = append(errs, "FirstName is required")
			}
			if _, ok := payload["LastName"]; !ok {
				errs = append(errs, "LastName is required")
			}
			valid := len(errs) == 0
			c.JSON(http.StatusOK, models.SuccessResponse("validated", gin.H{"valid": valid, "errors": errs}))
			return
		}
		if key == "invoice" {
			c.JSON(http.StatusOK, models.SuccessResponse("validated", gin.H{"valid": true, "errors": []interface{}{}}))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse("validated", gin.H{"valid": true, "errors": []interface{}{}}))
	})
}

func writeContractByKey(c *gin.Context, key string) {
	switch key {
	case "welcome":
		c.JSON(http.StatusOK, models.SuccessResponse("contract", gin.H{"template_key": "welcome", "variable_schema": gin.H{"type": "object"}}))
	case "booking_confirmation":
		c.JSON(http.StatusOK, models.SuccessResponse("contract", gin.H{"template_key": "booking_confirmation", "variable_schema": gin.H{"type": "object", "properties": gin.H{"Booking": gin.H{}}}}))
	case "password_reset":
		c.JSON(http.StatusOK, models.SuccessResponse("contract", gin.H{"template_key": "password_reset", "variable_schema": gin.H{"type": "object"}}))
	case "invoice":
		c.JSON(http.StatusOK, models.SuccessResponse("contract", gin.H{"template_key": "invoice", "variable_schema": gin.H{"type": "object", "properties": gin.H{"Customer": gin.H{}, "InvoiceData": gin.H{}}}}))
	default:
		c.JSON(http.StatusNotFound, models.ErrorResponseFunc("contract not found", "no contract with that key"))
	}
}

// validateContractForKey performs simple, test-focused validation for known
// template keys. Returns (valid, errors).
func validateContractForKey(key string, payload map[string]interface{}) (bool, []string) {
	log.Printf("templates: validateContractForKey key=%s payload_keys=%v", key, func() []string {
		ks := []string{}
		for k := range payload {
			ks = append(ks, k)
		}
		return ks
	}())
	switch key {
	case "welcome":
		errs := []string{}
		if _, ok := payload["FirstName"]; !ok {
			errs = append(errs, "FirstName is required")
		}
		if _, ok := payload["LastName"]; !ok {
			errs = append(errs, "LastName is required")
		}
		return len(errs) == 0, errs
	case "invoice":
		// Invoice validation is lenient in this test harness
		return true, nil
	default:
		return true, nil
	}
}
