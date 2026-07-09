package calendar

import (
	baseAPI "github.com/AgileExecutives/serverbase/api"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/AgileExecutives/shared-modules/calendar/handlers"
	repo "github.com/AgileExecutives/shared-modules/calendar/repo"
	"github.com/AgileExecutives/shared-modules/calendar/routes"
	"github.com/AgileExecutives/shared-modules/calendar/services"
)

// BasicModule implements the baseAPI.ModuleRouteProvider interface (legacy compatibility)
type BasicModule struct {
	routeProvider *routes.RouteProvider
	db            *gorm.DB
}

// NewBasicModule creates a new basic calendar module (legacy compatibility)
func NewBasicModule(db *gorm.DB) baseAPI.ModuleRouteProvider {

	// Initialize GORM-backed repo and service
	gormRepo := repo.NewGormCalendarRepo(db)
	calendarService := services.NewCalendarServiceWithRepo(gormRepo)

	// Initialize handlers
	calendarHandler := handlers.NewCalendarHandler(calendarService)

	// Initialize route provider
	routeProvider := routes.NewRouteProvider(calendarHandler)

	return &BasicModule{
		routeProvider: routeProvider,
		db:            db,
	}
}

// RegisterRoutes implements baseAPI.ModuleRouteProvider
func (m *BasicModule) RegisterRoutes(router *gin.RouterGroup) {
	// Legacy compatibility: call route provider with minimal ModuleContext
	m.routeProvider.RegisterRoutes(router, core.ModuleContext{DB: m.db})
}

// GetPrefix implements baseAPI.ModuleRouteProvider
func (m *BasicModule) GetPrefix() string {
	return m.routeProvider.GetPrefix()
}

// GetMiddleware forwards middleware from the underlying route provider
func (m *BasicModule) GetMiddleware() []gin.HandlerFunc {
	return m.routeProvider.GetMiddleware()
}

// GetSwaggerTags forwards Swagger tags from the underlying route provider
func (m *BasicModule) GetSwaggerTags() []string {
	return m.routeProvider.GetSwaggerTags()
}

// NewModule creates a new calendar module with auto-migration support (now the default)
// For basic module without auto-migration, use NewBasicModule instead
func NewModule(db *gorm.DB) baseAPI.ModuleRouteProvider {
	return NewModuleWithAutoMigration(db)
}
