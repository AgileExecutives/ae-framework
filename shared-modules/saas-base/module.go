package saasbase

// Package saasbase provides shared SaaS building blocks: customer management,
// subscription plans, and newsletter subscriptions.
//
// Deprecated: This package is now library-only. The serverbase module
// `serverbase/modules/saas` provides the `core.Module` integration and
// ServiceProviders that should be used by the application. Prefer
// `serverbase/modules/saas.NewSaaSModule()` for runtime registration.

import (
	"github.com/AgileExecutives/ae-framework/shared-modules/saas-base/entities"
	"github.com/AgileExecutives/ae-framework/shared-modules/saas-base/handlers"
	"github.com/AgileExecutives/ae-framework/shared-modules/saas-base/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module is the top-level entry point for direct usage.
type Module struct {
	db                 *gorm.DB
	customerHandlers   *handlers.CustomerHandlers
	planHandlers       *handlers.PlanHandlers
	newsletterHandlers *handlers.NewsletterHandlers
}

// NewModule creates a new saas-base Module.
// Prefer runtime registration via `NewCoreModule()` adapters used by the
// bootstrap system; direct construction is deprecated for runtime use.

// AutoMigrate runs the database migrations for all saas-base entities.
func (m *Module) AutoMigrate() error {
	return m.db.AutoMigrate(
		&models.Plan{},
		&models.Customer{},
		&models.Newsletter{},
	)
}

// GetName returns the module name.
func (m *Module) GetName() string {
	return "saas-base"
}

// GetEntitiesForMigration returns the GORM models for auto-migration.
func (m *Module) GetEntitiesForMigration() []interface{} {
	return []interface{}{
		&models.Plan{},
		&models.Customer{},
		&models.Newsletter{},
	}
}

// RegisterRoutes registers all saas-base routes on the provided router group.
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	customers := router.Group("/customers")
	{
		customers.GET("", m.customerHandlers.GetCustomers)
		customers.POST("", m.customerHandlers.CreateCustomer)
		customers.GET("/:id", m.customerHandlers.GetCustomer)
		customers.PUT("/:id", m.customerHandlers.UpdateCustomer)
		customers.DELETE("/:id", m.customerHandlers.DeleteCustomer)
	}

	plans := router.Group("/plans")
	{
		plans.GET("", m.planHandlers.GetPlans)
		plans.GET("/:id", m.planHandlers.GetPlan)
		plans.POST("", m.planHandlers.CreatePlan)
		plans.PUT("/:id", m.planHandlers.UpdatePlan)
		plans.DELETE("/:id", m.planHandlers.DeletePlan)
	}

	newsletter := router.Group("/newsletter")
	{
		newsletter.GET("", m.newsletterHandlers.GetSubscribers)
		newsletter.POST("/subscribe", m.newsletterHandlers.Subscribe)
		newsletter.POST("/unsubscribe", m.newsletterHandlers.Unsubscribe)
		newsletter.DELETE("/:id", m.newsletterHandlers.DeleteSubscriber)
	}
}

// Entities returns the core.Entity implementations for use with the base-server module system.
func Entities() []interface{} {
	return []interface{}{
		entities.NewPlanEntity(),
		entities.NewCustomerEntity(),
		entities.NewNewsletterEntity(),
	}
}
