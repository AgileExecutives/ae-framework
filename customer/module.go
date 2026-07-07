package customer
// Package customer re-exports the saas-base customer module for the base-server.
// The implementation has been moved to github.com/ae/shared-modules/saas-base.
package customer

import (
    "context"

    handlers "github.com/ae/base-server/modules/customer/handlers"
    "github.com/ae/base-server/pkg/core"
    saasmodels "github.com/ae/shared-modules/saas-base/models"
    "github.com/gin-gonic/gin"
)

type customerModule struct{}

type customerModuleWithCtx struct{ ctx core.ModuleContext }

func (m *customerModule) Name() string           { return "customer" }
func (m *customerModule) Version() string        { return "0.0.0" }
func (m *customerModule) Dependencies() []string { return []string{} }
func (m *customerModule) Initialize(ctx core.ModuleContext) error { return nil }
func (m *customerModule) Start(ctx context.Context) error { return nil }
func (m *customerModule) Stop(ctx context.Context) error  { return nil }

type simpleEntity struct{ model interface{} }
func (e *simpleEntity) TableName() string               { return "" }
func (e *simpleEntity) GetModel() interface{}           { return e.model }
func (e *simpleEntity) GetMigrations() []core.Migration { return nil }

func (m *customerModule) Entities() []core.Entity {
    return []core.Entity{ &simpleEntity{model: &saasmodels.Customer{}}, &simpleEntity{model: &saasmodels.Plan{}} }
}

func (m *customerModule) Routes() []core.RouteProvider {
    planProvider := &planRouteProvider{}
    customerProvider := &customerRouteProvider{}
    return []core.RouteProvider{planProvider, customerProvider}
}

type planRouteProvider struct{}
func (r *planRouteProvider) GetPrefix() string                { return "/plans" }
func (r *planRouteProvider) GetMiddleware() []gin.HandlerFunc { return []gin.HandlerFunc{} }
func (r *planRouteProvider) GetSwaggerTags() []string         { return []string{"plans"} }
func (r *planRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
    ph := handlers.NewPlanHandlers(ctx.DB, ctx.Logger)
    pr := handlers.NewPlanRoutes(ph)
    pr.RegisterRoutes(router, ctx)
}

type customerRouteProvider struct{}
func (r *customerRouteProvider) GetPrefix() string                { return "/customers" }
func (r *customerRouteProvider) GetMiddleware() []gin.HandlerFunc { return []gin.HandlerFunc{} }
func (r *customerRouteProvider) GetSwaggerTags() []string         { return []string{"customers"} }
func (r *customerRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
    ch := handlers.NewCustomerHandlers(ctx.DB, ctx.Logger)
    cr := handlers.NewCustomerRoutes(ch)
    cr.RegisterRoutes(router, ctx)
}

func (m *customerModule) EventHandlers() []core.EventHandler    { return []core.EventHandler{} }
func (m *customerModule) Services() []core.ServiceProvider      { return []core.ServiceProvider{} }
func (m *customerModule) Middleware() []core.MiddlewareProvider { return []core.MiddlewareProvider{} }
func (m *customerModule) SwaggerPaths() []string                { return []string{} }

func NewCustomerModule() core.Module { return &customerModule{} }
