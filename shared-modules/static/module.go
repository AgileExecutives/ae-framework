package static

import (
	"github.com/AgileExecutives/serverbase/module"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/shared-modules/static/handlers"
	"github.com/gin-gonic/gin"
)

// NewStaticModule returns a lightweight adapter-based module
func NewStaticModule() core.Module {
	// route provider constructs handlers at RegisterRoutes time using ctx
	rp := &staticRouteProvider{}
	return module.NewAdapterModule("static", "1.0.0", []string{},
		module.WithRoutes(rp),
		module.WithSwaggerPaths("./modules/static/handlers"),
	)
}

type staticRouteProvider struct{}

func (r *staticRouteProvider) GetPrefix() string                { return "/static" }
func (r *staticRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *staticRouteProvider) GetSwaggerTags() []string         { return []string{"static"} }
func (r *staticRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	repo := handlers.NewFSStaticRepo("./statics/json")
	h := handlers.NewStaticHandlers(ctx.Logger, repo)
	auth := router.Group("")
	auth.Use(ctx.Auth.RequireAuth())
	auth.GET("", h.ListStaticJSON)
	auth.GET("/", h.ListStaticJSON)
	auth.GET("/:filename", h.ServeStaticJSON)
}
