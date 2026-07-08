package handlers

import (
	"net/http"

	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/middleware"
	"github.com/AgileExecutives/shared-modules/pdf/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PDFHandler struct {
	service *services.PDFGenerator
	db      *gorm.DB
}

func NewPDFHandler(service *services.PDFGenerator, db *gorm.DB) *PDFHandler {
	return &PDFHandler{service: service, db: db}
}

// RegisterRoutes registers the PDF endpoints used by tests.
func (h *PDFHandler) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	// router is already mounted under /api/v1/<prefix> where prefix is /pdf
	router.POST("/create", func(c *gin.Context) {
		var req struct {
			Data         map[string]interface{} `json:"data"`
			TemplateName string                 `json:"templateName"`
			FileName     string                 `json:"fileName"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}
		if req.Data == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Data is required"})
			return
		}
		if req.TemplateName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Template name is required"})
			return
		}
		if req.FileName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File name is required"})
			return
		}

		filename, err := h.service.GeneratePDF(req.Data, req.TemplateName, req.FileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "PDF generated successfully", "filename": filename})
	})
}

func (h *PDFHandler) GetPrefix() string { return "/pdf" }
func (h *PDFHandler) GetMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{middleware.AuthMiddleware(h.db)}
}
func (h *PDFHandler) GetSwaggerTags() []string { return []string{"pdf"} }
