package handlers
package handlers

import (
    "github.com/ae/base-server/modules/pdf/services"
    "github.com/ae/base-server/pkg/core"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type PDFHandler struct{ service *services.PDFGenerator; db *gorm.DB }
func NewPDFHandler(service *services.PDFGenerator, db *gorm.DB) *PDFHandler { return &PDFHandler{service: service, db: db} }
func (h *PDFHandler) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) { /* register routes as needed */ }
func (h *PDFHandler) GetPrefix() string { return "/pdf" }
func (h *PDFHandler) GetMiddleware() []gin.HandlerFunc { return []gin.HandlerFunc{} }
func (h *PDFHandler) GetSwaggerTags() []string { return []string{"pdf"} }
