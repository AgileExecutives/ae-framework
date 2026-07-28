package handlers

import (
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framwork/shared-modules/pdf/services"
	"github.com/gin-gonic/gin"
)

// PdfHandlerAdapter provides a ModuleContext-aware constructor for the
// shared pdf handler service usage.
type PdfHandlerAdapter struct{ pdfService *services.PDFGenerator }

func NewPDFHandlerWithCtx(ctx core.ModuleContext) *PdfHandlerAdapter {
	if ctx.Services != nil {
		if s, ok := ctx.Services.Get("pdf-generator"); ok {
			if gen, ok := s.(*services.PDFGenerator); ok {
				return &PdfHandlerAdapter{pdfService: gen}
			}
		}
	}
	return &PdfHandlerAdapter{pdfService: services.NewPDFGenerator()}
}

func (h *PdfHandlerAdapter) GeneratePDFFromTemplate(c *gin.Context) {
	type PDFGenerateRequest struct {
		Data         map[string]interface{} `json:"data"`
		TemplateName string                 `json:"templateName"`
		FileName     string                 `json:"fileName"`
	}
	var req PDFGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	if req.Data == nil || req.TemplateName == "" || req.FileName == "" {
		c.JSON(400, gin.H{"error": "Missing required fields"})
		return
	}
	name, err := h.pdfService.GeneratePDF(req.Data, req.TemplateName, req.FileName)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate PDF", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "filename": name})
}

// Note: do not declare NewPDFHandler to avoid colliding with the primary
// constructor in pdf_handler.go. Use NewPDFHandlerWithCtx or the existing
// NewPDFHandler(service) defined in the package.
