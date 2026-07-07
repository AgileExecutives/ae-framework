package handlers

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type staticRouteProvider struct {
	h  *StaticHandlers
	db *gorm.DB
}

func NewStaticRoutes(h *StaticHandlers, db *gorm.DB) core.RouteProvider {
	return &staticRouteProvider{h: h, db: db}
}

func (r *staticRouteProvider) GetPrefix() string { return "/static" }
func (r *staticRouteProvider) GetMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{middleware.AuthMiddleware(r.db)}
}
func (r *staticRouteProvider) GetSwaggerTags() []string { return []string{"static"} }

func (r *staticRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	// Support both `/static` and `/static/`
	router.GET("", r.h.ListStaticJSON)
	router.GET("/", r.h.ListStaticJSON)
	router.GET("/:filename", r.h.ServeStaticJSON)
}
