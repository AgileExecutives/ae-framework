package booking

import (
	"context"

	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/AgileExecutives/shared-modules/booking/docs"
	"github.com/AgileExecutives/shared-modules/booking/entities"
	"github.com/AgileExecutives/shared-modules/booking/handlers"
	"github.com/AgileExecutives/shared-modules/booking/middleware"
	repo "github.com/AgileExecutives/shared-modules/booking/repo"
	"github.com/AgileExecutives/shared-modules/booking/routes"
	"github.com/AgileExecutives/shared-modules/booking/services"
)

// Module implements the complete core.Module interface for auto-migration support
type Module struct {
	db                 *gorm.DB
	routeProvider      *routes.RouteProvider
	bookingService     *services.BookingService
	bookingLinkService *services.BookingLinkService
	freeSlotsSvc       *services.FreeSlotsService
	tokenMiddleware    *middleware.BookingTokenMiddleware
	bookingHandler     *handlers.BookingHandler
}

// NewCoreModule creates a new booking module for the bootstrap system
// Initialization happens during the Initialize() lifecycle method
func NewCoreModule() *Module {
	return &Module{}
}

// NewModuleWithAutoMigration creates a new booking module with auto-migration support
func NewModuleWithAutoMigration(db *gorm.DB, jwtSecret string) *Module {
	// Initialize services
	gormRepo := repo.NewGormBookingRepo(db)
	bookingService := services.NewBookingServiceWithRepoAndDB(gormRepo, db)
	bookingLinkService := services.NewBookingLinkService(db, jwtSecret)
	freeSlotsSvc := services.NewFreeSlotsServiceWithRepo(gormRepo)

	// Initialize middleware (use repo for blacklist checks)
	tokenMiddleware := middleware.NewBookingTokenMiddleware(bookingLinkService, gormRepo)

	// Initialize handlers
	bookingHandler := handlers.NewBookingHandler(bookingService, bookingLinkService, freeSlotsSvc)

	// Initialize route provider
	routeProvider := routes.NewRouteProvider(bookingHandler, tokenMiddleware)

	return &Module{
		db:                 db,
		routeProvider:      routeProvider,
		bookingService:     bookingService,
		bookingLinkService: bookingLinkService,
		freeSlotsSvc:       freeSlotsSvc,
		tokenMiddleware:    tokenMiddleware,
		bookingHandler:     bookingHandler,
	}
}

// Name returns the module name
func (m *Module) Name() string {
	return "booking"
}

// Version returns the module version
func (m *Module) Version() string {
	return "1.0.0"
}

// Dependencies returns module dependencies
func (m *Module) Dependencies() []string {
	return []string{"user", "calendar"} // Depends on user (base) module for users/tenants and calendar for calendar entities
}

// Initialize initializes the module
func (m *Module) Initialize(ctx core.ModuleContext) error {
	ctx.Logger.Info("Initializing booking module...")

	// Store database reference
	m.db = ctx.DB

	// Initialize services (prefer repo-backed service)
	gormRepo := repo.NewGormBookingRepo(ctx.DB)
	m.bookingService = services.NewBookingServiceWithRepoAndDB(gormRepo, ctx.DB)

	// Try to get BookingLinkService from service registry (created by Factory)
	// If not available, create it directly
	if bookingLinkSvcRaw, ok := ctx.Services.Get("booking-link-service"); ok {
		ctx.Logger.Info("✅ Initialize: Found BookingLinkService in service registry")
		if bookingLinkSvc, ok := bookingLinkSvcRaw.(*services.BookingLinkService); ok {
			m.bookingLinkService = bookingLinkSvc
			ctx.Logger.Info("✅ Initialize: Using BookingLinkService from registry")
		} else {
			ctx.Logger.Error("❌ Initialize: BookingLinkService type assertion failed")
		}
	}

	// If service not in registry, create it (shouldn't happen in normal flow)
	if m.bookingLinkService == nil {
		ctx.Logger.Warn("⚠️ Initialize: BookingLinkService not in registry, creating directly")
		if ctx.TokenService != nil {
			ctx.Logger.Info("Using unified TokenService for booking link generation")
			tokenServiceAdapter := &tokenServiceAdapter{service: ctx.TokenService}
			m.bookingLinkService = services.NewBookingLinkServiceWithTokenService(ctx.DB, tokenServiceAdapter)
		} else {
			jwtSecret := "booking-link-secret-key-change-in-production"
			ctx.Logger.Warn("TokenService not available - using legacy booking token implementation")
			m.bookingLinkService = services.NewBookingLinkService(ctx.DB, jwtSecret)
		}
	}

	m.freeSlotsSvc = services.NewFreeSlotsServiceWithRepo(gormRepo)

	// Initialize middleware
	m.tokenMiddleware = middleware.NewBookingTokenMiddleware(m.bookingLinkService, gormRepo)

	// Initialize handlers
	m.bookingHandler = handlers.NewBookingHandler(m.bookingService, m.bookingLinkService, m.freeSlotsSvc)

	// Initialize route provider
	m.routeProvider = routes.NewRouteProvider(m.bookingHandler, m.tokenMiddleware)

	// Register pre-generated swagger docs with the server's doc registry so they
	// are merged into the combined spec served at /swagger/index.html.
	if ctx.DocRegistry != nil {
		ctx.DocRegistry.RegisterDoc(m.Name(), docs.SwaggerInfo.ReadDoc())
	}

	ctx.Logger.Info("Booking module initialized successfully")
	return nil
}

// Start starts the module
func (m *Module) Start(ctx context.Context) error {
	// Any startup logic here
	return nil
}

// Stop stops the module
func (m *Module) Stop(ctx context.Context) error {
	// Any cleanup logic here
	return nil
}

// Entities returns database entities for auto-migration
func (m *Module) Entities() []core.Entity {
	return []core.Entity{
		entities.NewBookingTemplateEntity(),
	}
}

// GetEntitiesForMigration returns GORM models for auto-migration (implements ModuleWithEntities interface)
func (m *Module) GetEntitiesForMigration() []interface{} {
	return []interface{}{
		&entities.BookingTemplate{},
	}
}

// Routes returns route providers
func (m *Module) Routes() []core.RouteProvider {
	if m.routeProvider == nil {
		return []core.RouteProvider{}
	}
	return []core.RouteProvider{
		&bookingRouteAdapter{
			provider: m.routeProvider,
		},
	}
}

// EventHandlers returns event handlers
func (m *Module) EventHandlers() []core.EventHandler {
	return []core.EventHandler{}
}

// Middleware returns middleware providers
func (m *Module) Middleware() []core.MiddlewareProvider {
	return []core.MiddlewareProvider{}
}

// Services returns service providers
func (m *Module) Services() []core.ServiceProvider {
	return []core.ServiceProvider{
		&bookingLinkServiceProvider{module: m},
	}
}

// bookingLinkServiceProvider implements core.ServiceProvider for BookingLinkService
type bookingLinkServiceProvider struct {
	module *Module
}

func (p *bookingLinkServiceProvider) ServiceName() string {
	return "booking-link-service"
}

func (p *bookingLinkServiceProvider) ServiceInterface() interface{} {
	return p.module.bookingLinkService
}

func (p *bookingLinkServiceProvider) Factory(ctx core.ModuleContext) (interface{}, error) {
	// Create service with unified TokenService if available
	if ctx.TokenService != nil {
		ctx.Logger.Info("✅ Factory: Creating BookingLinkService with unified TokenService")
		tokenServiceAdapter := &tokenServiceAdapter{service: ctx.TokenService}
		p.module.bookingLinkService = services.NewBookingLinkServiceWithTokenService(ctx.DB, tokenServiceAdapter)
	} else {
		// Fallback to legacy implementation (should not happen in modern setup)
		ctx.Logger.Warn("⚠️ Factory: TokenService not available - using legacy booking token implementation")
		jwtSecret := "booking-link-secret-key-change-in-production"
		p.module.bookingLinkService = services.NewBookingLinkService(ctx.DB, jwtSecret)
	}
	ctx.Logger.Info("✅ Factory: BookingLinkService created and stored in module")
	return p.module.bookingLinkService, nil
}

// SwaggerPaths returns Swagger documentation paths
func (m *Module) SwaggerPaths() []string {
	return []string{
		"/booking/templates",
		"/booking/templates/{id}",
		"/booking/link",
	}
}

// Legacy compatibility methods removed (ModuleRouteProvider compatibility no longer required)

// bookingRouteAdapter adapts the booking routes.RouteProvider to core.RouteProvider
type bookingRouteAdapter struct {
	provider *routes.RouteProvider
}

func (a *bookingRouteAdapter) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	a.provider.RegisterRoutes(router, ctx)
}

func (a *bookingRouteAdapter) GetPrefix() string {
	return a.provider.GetPrefix()
}

func (a *bookingRouteAdapter) GetMiddleware() []gin.HandlerFunc {
	// Middleware is handled by the route provider itself
	return a.provider.GetMiddleware()
}

func (a *bookingRouteAdapter) GetSwaggerTags() []string {
	return a.provider.GetSwaggerTags()
}

// tokenServiceAdapter adapts the base TokenService to BookingLinkService interface
type tokenServiceAdapter struct {
	service core.TokenService
}

func (a *tokenServiceAdapter) GenerateToken(claims jwt.Claims) (string, error) {
	return a.service.GenerateToken(interface{}(claims))
}

func (a *tokenServiceAdapter) ValidateToken(tokenString string, claims jwt.Claims) error {
	return a.service.ValidateToken(tokenString, interface{}(claims))
}
