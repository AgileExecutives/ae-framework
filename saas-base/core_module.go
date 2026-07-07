package saasbase

import (
	"context"

	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/shared-modules/saas-base/entities"
	"github.com/AgileExecutives/shared-modules/saas-base/handlers"
	"github.com/AgileExecutives/shared-modules/saas-base/services"
	"github.com/gin-gonic/gin"
)

// CoreModule implements core.Module for integration with the base-server module system.
type CoreModule struct {
	module          *Module
	customerService *services.CustomerService
}

// NewCoreModule creates a CoreModule that implements core.Module.
func NewCoreModule() core.Module {
	return &CoreModule{}
}

func (m *CoreModule) Name() string {
	return "saas-base"
}

func (m *CoreModule) Version() string {
	return "1.0.0"
}

func (m *CoreModule) Dependencies() []string {
	return []string{"base"}
}

func (m *CoreModule) Initialize(ctx core.ModuleContext) error {
	ctx.Logger.Info("Initializing saas-base module...")

	m.module = NewModule(ctx.DB)
	m.customerService = services.NewCustomerService(ctx.DB, ctx.Logger)

	if err := m.module.AutoMigrate(); err != nil {
		return err
	}

	ctx.Services.Register("saas-base-customer", m.customerService)
	ctx.Logger.Info("saas-base module initialized successfully")
	return nil
}

func (m *CoreModule) Start(_ context.Context) error { return nil }
func (m *CoreModule) Stop(_ context.Context) error  { return nil }

func (m *CoreModule) Entities() []core.Entity {
	return []core.Entity{
		entities.NewPlanEntity(),
		entities.NewCustomerEntity(),
		entities.NewNewsletterEntity(),
	}
}

func (m *CoreModule) Routes() []core.RouteProvider {
	return []core.RouteProvider{
		&customerRouteProvider{handlers: handlers.NewCustomerHandlers(m.module.db)},
		&planRouteProvider{handlers: handlers.NewPlanHandlers(m.module.db)},
		&newsletterRouteProvider{handlers: handlers.NewNewsletterHandlers(m.module.db)},
	}
}

func (m *CoreModule) EventHandlers() []core.EventHandler    { return []core.EventHandler{} }
func (m *CoreModule) Middleware() []core.MiddlewareProvider { return []core.MiddlewareProvider{} }
func (m *CoreModule) Services() []core.ServiceProvider      { return []core.ServiceProvider{} }

func (m *CoreModule) SwaggerPaths() []string {
	return []string{}
}

// ─── Route providers ─────────────────────────────────────────────────────────

type customerRouteProvider struct {
	handlers *handlers.CustomerHandlers
}

func (r *customerRouteProvider) GetPrefix() string                { return "/customers" }
func (r *customerRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *customerRouteProvider) GetSwaggerTags() []string         { return []string{"customers"} }
func (r *customerRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	auth := router.Use(ctx.Auth.RequireAuth())
	auth.GET("", r.handlers.GetCustomers)
	auth.POST("", r.handlers.CreateCustomer)
	auth.GET("/:id", r.handlers.GetCustomer)
	auth.PUT("/:id", r.handlers.UpdateCustomer)
	auth.DELETE("/:id", r.handlers.DeleteCustomer)
}

type planRouteProvider struct {
	handlers *handlers.PlanHandlers
}

func (r *planRouteProvider) GetPrefix() string                { return "/plans" }
func (r *planRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *planRouteProvider) GetSwaggerTags() []string         { return []string{"plans"} }
func (r *planRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	// Public read routes
	router.GET("", r.handlers.GetPlans)
	router.GET("/:id", r.handlers.GetPlan)

	// Admin write routes
	admin := router.Use(ctx.Auth.RequireAuth(), ctx.Auth.RequireRole("admin", "super-admin"))
	admin.POST("", r.handlers.CreatePlan)
	admin.PUT("/:id", r.handlers.UpdatePlan)
	admin.DELETE("/:id", r.handlers.DeletePlan)
}

type newsletterRouteProvider struct {
	handlers *handlers.NewsletterHandlers
}

func (r *newsletterRouteProvider) GetPrefix() string                { return "/newsletter" }
func (r *newsletterRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *newsletterRouteProvider) GetSwaggerTags() []string         { return []string{"newsletter"} }
func (r *newsletterRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	// Public subscription routes
	router.POST("/subscribe", r.handlers.Subscribe)
	router.POST("/unsubscribe", r.handlers.Unsubscribe)

	// Admin routes
	admin := router.Use(ctx.Auth.RequireAuth(), ctx.Auth.RequireRole("admin", "super-admin"))
	admin.GET("", r.handlers.GetSubscribers)
	admin.DELETE("/:id", r.handlers.DeleteSubscriber)
}
