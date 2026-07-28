package routes

import (
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framwork/shared-modules/settings/handlers"
	"github.com/gin-gonic/gin"
)

type RouteProvider struct{ handler *handlers.SettingsHandler }

func NewRouteProvider(handler *handlers.SettingsHandler) *RouteProvider {
	return &RouteProvider{handler: handler}
}
func (rp *RouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	settingsGroup := router.Group("/settings")
	{
		tenantGroup := settingsGroup.Group("/organizations/:tenant_id")
		{
			tenantGroup.GET("/domains/:domain", rp.handler.GetDomainSettings)
			tenantGroup.POST("/domains/:domain", rp.handler.UpdateDomainSettings)
			tenantGroup.PUT("/domains/:domain", rp.handler.UpdateDomainSettings)
			tenantGroup.GET("/domains/:domain/:key", rp.handler.GetSetting)
			tenantGroup.POST("/domains/:domain/:key", rp.handler.UpdateSetting)
			tenantGroup.PUT("/domains/:domain/:key", rp.handler.UpdateSetting)
		}
	}
}
func (rp *RouteProvider) GetPrefix() string                { return "" }
func (rp *RouteProvider) GetMiddleware() []gin.HandlerFunc { return []gin.HandlerFunc{} }
func (rp *RouteProvider) GetSwaggerTags() []string         { return []string{"settings"} }
