package routes

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/shared-modules/documents/handlers"
	"github.com/AgileExecutives/shared-modules/documents/middleware"
	"github.com/AgileExecutives/shared-modules/documents/services"
	"github.com/gin-gonic/gin"
)

// DocumentRoutes implements RouteProvider for document management endpoints
type DocumentRoutes struct{ documentService *services.DocumentService }

// NewDocumentRoutes creates a new DocumentRoutes instance
func NewDocumentRoutes(documentService *services.DocumentService) *DocumentRoutes {
	return &DocumentRoutes{documentService: documentService}
}

// RegisterRoutes registers all document management routes
func (r *DocumentRoutes) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	documentHandler := handlers.NewDocumentHandler(r.documentService)

	// Document routes with tenant isolation
	documents := router.Group("/documents")
	// Apply auth middleware at registration using ModuleContext
	documents.Use(ctx.Auth.RequireAuth())
	{
		// Upload a document
		documents.POST("", documentHandler.UploadDocument)

		// List documents with pagination and filters
		documents.GET("", documentHandler.ListDocuments)

		// Get document metadata (with tenant access check)
		documents.GET("/:id",
			middleware.EnsureTenantAccess(ctx.DB),
			documentHandler.GetDocument,
		)

		// Get document download URL (with tenant access check)
		documents.GET("/:id/download",
			middleware.EnsureTenantAccess(ctx.DB),
			documentHandler.DownloadDocument,
		)

		// Delete document (with tenant access check)
		documents.DELETE("/:id",
			middleware.EnsureTenantAccess(ctx.DB),
			documentHandler.DeleteDocument,
		)
	}
}

// GetPrefix returns the base path for document routes
func (r *DocumentRoutes) GetPrefix() string { return "" }

// GetMiddleware returns middleware to apply to all document routes
func (r *DocumentRoutes) GetMiddleware() []gin.HandlerFunc { return nil }

// GetSwaggerTags returns Swagger tags for documentation
func (r *DocumentRoutes) GetSwaggerTags() []string { return []string{"Documents"} }
