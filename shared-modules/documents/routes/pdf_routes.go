package routes

import (
	templateServices "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/services"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framwork/shared-modules/documents/handlers"
	"github.com/AgileExecutives/ae-framwork/shared-modules/documents/services"
	"github.com/gin-gonic/gin"
)

// PDFRoutes implements RouteProvider for PDF generation endpoints
type PDFRoutes struct {
	pdfService      *services.PDFService
	templateService *templateServices.TemplateService
}

// NewPDFRoutes creates a new PDFRoutes instance
func NewPDFRoutes(pdfService *services.PDFService, templateService *templateServices.TemplateService) *PDFRoutes {
	return &PDFRoutes{pdfService: pdfService, templateService: templateService}
}

// RegisterRoutes registers all PDF generation routes
func (r *PDFRoutes) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	handler := handlers.NewPDFHandler(r.pdfService, r.templateService)

	// PDF routes
	pdfs := router.Group("/pdfs")
	pdfs.Use(ctx.Auth.RequireAuth())
	{
		// Generate PDF from HTML
		pdfs.POST("/generate", handler.GeneratePDFFromHTML)

		// Generate PDF from template
		pdfs.POST("/from-template", handler.GeneratePDFFromTemplate)

		// Generate invoice PDF
		pdfs.POST("/invoice/:invoice_id", handler.GenerateInvoicePDF)
	}
}

// GetPrefix returns the base path for PDF routes
func (r *PDFRoutes) GetPrefix() string { return "" }

// GetMiddleware returns middleware to apply to all PDF routes
func (r *PDFRoutes) GetMiddleware() []gin.HandlerFunc { return nil }

// GetSwaggerTags returns Swagger tags for documentation
func (r *PDFRoutes) GetSwaggerTags() []string { return []string{"PDFs"} }
