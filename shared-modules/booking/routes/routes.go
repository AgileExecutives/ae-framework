package routes

import (
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/gin-gonic/gin"

	"github.com/AgileExecutives/ae-framwork/shared-modules/booking/handlers"
	bookingMiddleware "github.com/AgileExecutives/ae-framwork/shared-modules/booking/middleware"
)

// RouteProvider provides routing functionality for booking management
type RouteProvider struct {
	bookingHandler  *handlers.BookingHandler
	tokenMiddleware *bookingMiddleware.BookingTokenMiddleware
}

// NewRouteProvider creates a new route provider
func NewRouteProvider(bookingHandler *handlers.BookingHandler, tokenMiddleware *bookingMiddleware.BookingTokenMiddleware) *RouteProvider {
	return &RouteProvider{bookingHandler: bookingHandler, tokenMiddleware: tokenMiddleware}
}

// RegisterRoutes registers the booking management routes with the provided router group
func (rp *RouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	// Public endpoints (no authentication) - just token validation
	router.GET("/booking/freeslots/:token", rp.tokenMiddleware.ValidateBookingToken(), rp.bookingHandler.GetFreeSlots)
	router.GET("/client/:token", rp.tokenMiddleware.ValidateBookingToken(), rp.bookingHandler.GetClientByToken)

	// Authenticated endpoints - use ModuleContext auth
	// Booking templates/configurations CRUD endpoints (authenticated)
	templates := router.Group("/booking/templates")
	templates.Use(ctx.Auth.RequireAuth())
	{
		templates.POST("", rp.bookingHandler.CreateConfiguration)
		templates.GET("", rp.bookingHandler.GetAllConfigurations)
		templates.GET("/:id", rp.bookingHandler.GetConfiguration)
		templates.PUT("/:id", rp.bookingHandler.UpdateConfiguration)
		templates.DELETE("/:id", rp.bookingHandler.DeleteConfiguration)

		// Additional query endpoints
		templates.GET("/by-user", rp.bookingHandler.GetConfigurationsByUser)
		templates.GET("/by-calendar", rp.bookingHandler.GetConfigurationsByCalendar)
	}

	// Booking link generation endpoint (authenticated)
	router.POST("/booking/link", ctx.Auth.RequireAuth(), rp.bookingHandler.CreateBookingLink)
}

// GetPrefix returns the route prefix for booking management endpoints
func (rp *RouteProvider) GetPrefix() string { return "" }

// GetMiddleware returns middleware to apply to all routes
// Returns empty array since we apply auth selectively in RegisterRoutes
func (rp *RouteProvider) GetMiddleware() []gin.HandlerFunc { return []gin.HandlerFunc{} }

// GetSwaggerTags returns swagger tags for the routes
func (rp *RouteProvider) GetSwaggerTags() []string { return []string{"booking"} }
