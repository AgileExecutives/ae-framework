package static

import (
	"context"

	"github.com/AgileExecutives/shared-modules/static/handlers"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"gorm.io/gorm"
)

type StaticModule struct {
	staticHandlers *handlers.StaticHandlers
	db             *gorm.DB
}

func NewStaticModule() core.Module             { return &StaticModule{} }
func (m *StaticModule) Name() string           { return "static" }
func (m *StaticModule) Version() string        { return "1.0.0" }
func (m *StaticModule) Dependencies() []string { return []string{} }
func (m *StaticModule) Initialize(ctx core.ModuleContext) error {
	ctx.Logger.Info("Initializing static module...")
	m.db = ctx.DB
	repo := handlers.NewFSStaticRepo("./statics/json")
	m.staticHandlers = handlers.NewStaticHandlers(ctx.Logger, repo)
	ctx.Logger.Info("Static module initialized successfully")
	return nil
}
func (m *StaticModule) Start(ctx context.Context) error { return nil }
func (m *StaticModule) Stop(ctx context.Context) error  { return nil }
func (m *StaticModule) Entities() []core.Entity         { return []core.Entity{} }
func (m *StaticModule) Routes() []core.RouteProvider {
	return []core.RouteProvider{handlers.NewStaticRoutes(m.staticHandlers, m.db)}
}
func (m *StaticModule) EventHandlers() []core.EventHandler    { return []core.EventHandler{} }
func (m *StaticModule) Services() []core.ServiceProvider      { return []core.ServiceProvider{} }
func (m *StaticModule) Middleware() []core.MiddlewareProvider { return []core.MiddlewareProvider{} }
func (m *StaticModule) SwaggerPaths() []string                { return []string{"./modules/static/handlers"} }
