package routes

import (
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framwork/shared-modules/invoice_number/handlers"
	"github.com/AgileExecutives/ae-framwork/shared-modules/invoice_number/services"
	"github.com/gin-gonic/gin"
)

// InvoiceNumberRoutes implements RouteProvider for invoice number endpoints
type InvoiceNumberRoutes struct {
	invoiceNumberService *services.InvoiceNumberService
}

// NewInvoiceNumberRoutes creates a new InvoiceNumberRoutes instance
func NewInvoiceNumberRoutes(invoiceNumberService *services.InvoiceNumberService) *InvoiceNumberRoutes {
	return &InvoiceNumberRoutes{invoiceNumberService: invoiceNumberService}
}

// RegisterRoutes registers all invoice number routes
func (r *InvoiceNumberRoutes) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	handler := handlers.NewInvoiceNumberHandler(r.invoiceNumberService)

	// Apply auth middleware at registration time using ModuleContext
	grp := router.Group("/invoice-numbers")
	grp.Use(ctx.Auth.RequireAuth())

	// Invoice number routes
	{
		// Generate next invoice number
		grp.POST("/generate", handler.GenerateInvoiceNumber)

		// Get current sequence without incrementing
		grp.GET("/current", handler.GetCurrentSequence)

		// Get invoice number history
		grp.GET("/history", handler.GetInvoiceNumberHistory)

		// Void an invoice number
		grp.POST("/void", handler.VoidInvoiceNumber)
	}
}

// GetPrefix returns the base path for invoice number routes
func (r *InvoiceNumberRoutes) GetPrefix() string { return "" }

// GetMiddleware returns middleware to apply to all invoice number routes
func (r *InvoiceNumberRoutes) GetMiddleware() []gin.HandlerFunc { return nil }

// GetSwaggerTags returns Swagger tags for documentation
func (r *InvoiceNumberRoutes) GetSwaggerTags() []string { return []string{"Invoice Numbers"} }
