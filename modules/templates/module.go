package templates

// Package templates provides a lightweight templates module used by tests and
// optionally by apps that want an in-process templates provider.
import (
	"fmt"
	"net/http"
	"sync"

	"github.com/AgileExecutives/serverbase/module"
	"github.com/AgileExecutives/serverbase/modules/templates/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
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
	var next uint = 1

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
		// try to render via service; allow service to return placeholder
		html, _ := svc.RenderTemplate(c.Request.Context(), 1, id, payload["data"])
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"content": html}})
	})

	templates.DELETE("/:id", func(c *gin.Context) {
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		mu.Lock()
		delete(store, id)
		mu.Unlock()
		c.Status(http.StatusNoContent)
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

	templates.GET("/contracts/:key", func(c *gin.Context) {
		key := c.Param("key")
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
