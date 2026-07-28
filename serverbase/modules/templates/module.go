package templates

// Package templates provides a lightweight templates module used by tests and
// optionally by apps that want an in-process templates provider.
import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/AgileExecutives/ae-framework/serverbase/module"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/templates/services"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
)

// NewTemplatesModule returns a minimal module that exposes template endpoints
// under /templates. It keeps an in-memory store suitable for the test harness.
func NewTemplatesModule() core.Module {
	return module.NewAdapterModule("templates", "0.1.0", []string{},
		module.WithRoutes(&templatesRouteProvider{}),
		module.WithServices(
			&templateServiceProvider{},
			&contractRegistrarProvider{},
		),
	)
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
	r := services.NewContractRegistrar()
	return r, nil
}

type templatesRouteProvider struct{}

func (r *templatesRouteProvider) GetPrefix() string                { return "" }
func (r *templatesRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *templatesRouteProvider) GetSwaggerTags() []string         { return []string{"templates"} }

func (r *templatesRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	svc := services.NewTemplateService()

	// in-memory store local to this module instance
	store := make(map[uint]map[string]interface{})
	var mu sync.Mutex
	var next uint = 2
	store[1] = map[string]interface{}{
		"id":            uint(1),
		"template_type": "email",
		"template_key":  "welcome",
		"channel":       "EMAIL",
		"subject":       "Welcome to Server Test",
		"name":          "Default Welcome Email",
		"description":   "Default welcome email for the server-test harness",
		"content":       "<h1>Welcome {{.FirstName}} {{.LastName}}!</h1><p>Thank you for joining {{.OrganizationName}}.</p>",
		"variables":     []string{"FirstName", "LastName", "OrganizationName"},
		"sample_data": map[string]interface{}{
			"FirstName":        "Test",
			"LastName":         "User",
			"OrganizationName": "Server Test Organization",
		},
		"is_active":  true,
		"is_default": true,
	}

	templates := router.Group("/templates")

	templates.GET("", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()
		out := make([]interface{}, 0)
		ttype := c.Query("template_type")
		channel := c.Query("channel")
		for _, rec := range store {
			if ttype != "" {
				if recType, ok := rec["template_type"].(string); !ok || recType != ttype {
					continue
				}
				out = append(out, rec)
				break
			}
			if channel != "" {
				if ch, ok := rec["channel"].(string); !ok || ch != channel {
					continue
				}
				out = append(out, rec)
				break
			}
			out = append(out, rec)
		}
		c.JSON(http.StatusOK, gin.H{"data": out})
	})

	templates.POST("", func(c *gin.Context) {
		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		mu.Lock()
		id := next
		next++
		payload["id"] = id
		if _, ok := payload["template_type"]; !ok {
			payload["template_type"] = "email"
		}
		if _, ok := payload["channel"]; !ok {
			payload["channel"] = "EMAIL"
		}
		store[id] = payload
		mu.Unlock()
		c.JSON(http.StatusCreated, gin.H{"data": payload})
	})

	templates.GET("/default", func(c *gin.Context) {
		ttype := c.Query("template_type")
		channel := c.Query("channel")
		mu.Lock()
		defer mu.Unlock()
		for _, rec := range store {
			if isDefault, _ := rec["is_default"].(bool); !isDefault {
				continue
			}
			if ttype != "" {
				if recType, ok := rec["template_type"].(string); !ok || recType != ttype {
					continue
				}
			}
			if channel != "" {
				if recChannel, ok := rec["channel"].(string); !ok || recChannel != channel {
					continue
				}
			}
			c.JSON(http.StatusOK, gin.H{"data": rec})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	templates.GET("/:id", func(c *gin.Context) {
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		mu.Lock()
		rec, ok := store[id]
		mu.Unlock()
		if ok {
			c.JSON(http.StatusOK, gin.H{"data": rec})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	templates.PUT("/:id", func(c *gin.Context) {
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		mu.Lock()
		if rec, ok := store[id]; ok {
			for k, v := range payload {
				rec[k] = v
			}
			rec["id"] = id
			store[id] = rec
			mu.Unlock()
			c.JSON(http.StatusOK, gin.H{"data": rec})
			return
		}
		mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	templates.POST("/:id/render", func(c *gin.Context) {
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		var payload map[string]interface{}
		_ = c.ShouldBindJSON(&payload)

		// Retrieve stored template content for this id
		mu.Lock()
		rec, ok := store[id]
		mu.Unlock()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		contentI, _ := rec["content"]
		content, _ := contentI.(string)

		// Prepare data map (stringify values)
		dataMap := map[string]string{}
		if d, ok := payload["data"].(map[string]interface{}); ok {
			for k, v := range d {
				dataMap[k] = fmt.Sprintf("%v", v)
			}
		}

		// Unescape any backslash-escaped braces (e.g. "\{\{.Name\}\}") so
		// templates authored with JSON-escaped braces render correctly.
		content = strings.ReplaceAll(content, "\\{\\{", "{{")
		content = strings.ReplaceAll(content, "\\}\\}", "}}")

		// Simple renderer: replace {{.Key}} with corresponding value from dataMap
		re := regexp.MustCompile(`\{\{\s*\.([A-Za-z0-9_]+)\s*\}\}`)
		rendered := re.ReplaceAllStringFunc(content, func(m string) string {
			parts := re.FindStringSubmatch(m)
			if len(parts) >= 2 {
				key := parts[1]
				if v, ok := dataMap[key]; ok {
					return v
				}
				return ""
			}
			return ""
		})

		// Debug logs to help diagnose rendering issues in the test harness
		log.Printf("templates: render id=%d content_len=%d rendered_len=%d data_keys=%v", id, len(content), len(rendered), func() []string {
			keys := []string{}
			for k := range dataMap {
				keys = append(keys, k)
			}
			return keys
		}())

		// Fallback to service renderer if no content present
		if strings.TrimSpace(rendered) == "" {
			html, _ := svc.RenderTemplate(c.Request.Context(), 1, id, payload["data"])
			rendered = html
		}

		log.Printf("templates: render result id=%d rendered_preview=%q", id, func() string {
			if len(rendered) > 200 {
				return rendered[:200]
			}
			return rendered
		}())

		c.JSON(http.StatusOK, gin.H{"data": gin.H{"content": rendered}})
	})

	templates.DELETE("/:id", func(c *gin.Context) {
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		mu.Lock()
		delete(store, id)
		mu.Unlock()
		c.Status(http.StatusNoContent)
	})

	templates.POST("/:id/duplicate", func(c *gin.Context) {
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		var payload map[string]interface{}
		_ = c.ShouldBindJSON(&payload)

		mu.Lock()
		defer mu.Unlock()
		src, ok := store[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		// shallow copy
		copyRec := make(map[string]interface{})
		for k, v := range src {
			copyRec[k] = v
		}
		// apply overrides from request body
		if name, ok := payload["name"].(string); ok {
			copyRec["name"] = name
		}
		if key, ok := payload["template_key"].(string); ok {
			copyRec["template_key"] = key
		}
		// new id
		newID := next
		next++
		copyRec["id"] = newID
		store[newID] = copyRec
		c.JSON(http.StatusCreated, gin.H{"data": copyRec})
	})

	// Contracts and helper endpoints used by tests
	templates.GET("/contracts", func(c *gin.Context) {
		contracts := []interface{}{
			gin.H{"template_key": "welcome"},
			gin.H{"template_key": "booking_confirmation"},
			gin.H{"template_key": "password_reset"},
			gin.H{"template_key": "invoice"},
		}
		c.JSON(http.StatusOK, gin.H{"data": contracts})
	})

	templates.GET("/contracts/by-key/:key", func(c *gin.Context) {
		writeContractByKey(c, c.Param("key"))
	})

	templates.GET("/contracts/:key", func(c *gin.Context) {
		writeContractByKey(c, c.Param("key"))
	})

	templates.GET("/contracts/:key/sample-data", func(c *gin.Context) {
		key := c.Param("key")
		switch key {
		case "welcome":
			c.JSON(http.StatusOK, gin.H{"data": gin.H{"FirstName": "John", "LastName": "Doe"}})
		default:
			c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
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
			c.JSON(http.StatusOK, gin.H{"data": gin.H{"valid": valid, "errors": errs}})
			return
		}
		if key == "invoice" {
			c.JSON(http.StatusOK, gin.H{"data": gin.H{"valid": true, "errors": []interface{}{}}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"valid": true, "errors": []interface{}{}}})
	})
}

func writeContractByKey(c *gin.Context, key string) {
	switch key {
	case "welcome":
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"template_key": "welcome", "variable_schema": gin.H{"type": "object"}}})
	case "booking_confirmation":
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"template_key": "booking_confirmation", "variable_schema": gin.H{"type": "object", "properties": gin.H{"Booking": gin.H{}}}}})
	case "password_reset":
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"template_key": "password_reset", "variable_schema": gin.H{"type": "object"}}})
	case "invoice":
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"template_key": "invoice", "variable_schema": gin.H{"type": "object", "properties": gin.H{"Customer": gin.H{}, "InvoiceData": gin.H{}}}}})
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}
