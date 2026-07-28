// Package settings provides a core.Module that wires the serverbase settings
// system into the module registry so it appears as a first-class module.
package settings

import (
	"github.com/AgileExecutives/ae-framework/serverbase/module"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	sbsettings "github.com/AgileExecutives/ae-framework/serverbase/pkg/settings"
	sbentities "github.com/AgileExecutives/ae-framework/serverbase/pkg/settings/entities"
	"github.com/gin-gonic/gin"
)

// NewSettingsModule returns a core.Module that initializes the settings system
// and registers its routes under /api/v1/settings.
func NewSettingsModule() core.Module {
	return module.NewAdapterModule("settings", "1.0.0", []string{},
		module.WithEntities(sbentities.NewSettingEntity(), sbentities.NewSettingDefinitionEntity()),
		module.WithInit(func(ctx core.ModuleContext) error {
			// Register pre-rendered docs if present (settings doesn't currently
			// include generated docs). Keep hook for future documentation.
			return nil
		}),
		module.WithRoutes(&settingsRouteProvider{}),
	)
}

type settingsRouteProvider struct{}

func (r *settingsRouteProvider) GetPrefix() string                { return "" }
func (r *settingsRouteProvider) GetMiddleware() []gin.HandlerFunc { return nil }
func (r *settingsRouteProvider) GetSwaggerTags() []string         { return []string{"settings"} }

func (r *settingsRouteProvider) RegisterRoutes(router *gin.RouterGroup, ctx core.ModuleContext) {
	// Construct a fresh settings system for this module using the ModuleContext DB
	sys, err := sbsettings.NewSettingsSystem(ctx.DB)
	if err != nil {
		if ctx.Logger != nil {
			ctx.Logger.Error("failed to initialize settings system", err)
		}
		return
	}

	// Register routes under /settings so the final path becomes /api/v1/settings
	settingsGroup := router.Group("/settings")
	{
		settingsGroup.GET("/health", sys.Handler.HealthCheck)
		settingsGroup.GET("/modules", sys.Handler.GetRegisteredModules)
		settingsGroup.GET("/version", sys.Handler.GetVersion)

		orgGroup := settingsGroup.Group("/organizations/:organization_id")
		{
			orgGroup.GET("", sys.Handler.GetOrganizationSettings)
			orgGroup.POST("", sys.Handler.SetOrganizationSetting)
			orgGroup.PUT("/:domain/:key", sys.Handler.UpdateOrganizationSetting)
			orgGroup.DELETE("/:domain/:key", sys.Handler.DeleteOrganizationSetting)
			orgGroup.POST("/bulk", sys.Handler.BulkSetOrganizationSettings)
			orgGroup.GET("/domains", sys.Handler.GetOrganizationDomains)
			orgGroup.GET("/domains/:domain", sys.Handler.GetOrganizationDomainSettings)
			orgGroup.POST("/domains/:domain", sys.Handler.SetOrganizationDomainSettings)
			orgGroup.DELETE("/domains/:domain", sys.Handler.DeleteOrganizationDomainSettings)
			orgGroup.POST("/validate", sys.Handler.ValidateOrganizationSettings)
			orgGroup.GET("/export", sys.Handler.ExportOrganizationSettings)
			orgGroup.POST("/import", sys.Handler.ImportOrganizationSettings)
		}
	}
}
