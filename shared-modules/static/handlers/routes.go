package handlers

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
)

type staticRouteProvider struct {
	h *StaticHandlers
}

func NewStaticRoutes(h *StaticHandlers) core.RouteProvider {
	return &staticRouteProvider{h: h}
}

func (r *staticRouteProvider) GetPrefix() string                { return "/static" }
func (r *staticRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *staticRouteProvider) GetSwaggerTags() []string         { return []string{"static"} }

func (r *staticRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	// Apply auth middleware at registration time using ModuleContext
	auth := router.Use(ctx.Auth.RequireAuth())
	// Support both `/static` and `/static/`
	auth.GET("", r.h.ListStaticJSON)
	auth.GET("/", r.h.ListStaticJSON)
	auth.GET("/:filename", r.h.ServeStaticJSON)
}
