package routes

import (
	baseAPI "github.com/AgileExecutives/serverbase/api"
	templateServices "github.com/AgileExecutives/serverbase/modules/templates/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/shared-modules/documents/handlers"
	"github.com/AgileExecutives/shared-modules/documents/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PDFRoutes implements RouteProvider for PDF generation endpoints
type PDFRoutes struct {
	pdfService      *services.PDFService
	templateService *templateServices.TemplateService
	db              *gorm.DB
}

// NewPDFRoutes creates a new PDFRoutes instance
func NewPDFRoutes(pdfService *services.PDFService, templateService *templateServices.TemplateService, db *gorm.DB) *PDFRoutes {
	return &PDFRoutes{
		pdfService:      pdfService,
		templateService: templateService,
		db:              db,
	}
}

// RegisterRoutes registers all PDF generation routes
func (r *PDFRoutes) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	handler := handlers.NewPDFHandler(r.pdfService, r.templateService, ctx.DB)

	// PDF routes
	pdfs := router.Group("/pdfs")
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
func (r *PDFRoutes) GetPrefix() string {
	return ""
}

// GetMiddleware returns middleware to apply to all PDF routes
func (r *PDFRoutes) GetMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		baseAPI.AuthMiddleware(r.db), // Require authentication for tenant ID extraction
	}
}

// GetSwaggerTags returns Swagger tags for documentation
func (r *PDFRoutes) GetSwaggerTags() []string {
	return []string{"PDFs"}
}
