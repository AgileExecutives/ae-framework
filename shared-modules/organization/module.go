package organization

import (
	"context"

	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framwork/shared-modules/organization/entities"
)

// Minimal OrganizationModule that only provides the entity to avoid
// importing serverbase internal packages from shared-modules.
type OrganizationModule struct{}

func NewOrganizationModule() core.Module             { return &OrganizationModule{} }
func (m *OrganizationModule) Name() string           { return "organization" }
func (m *OrganizationModule) Version() string        { return "1.0.0" }
func (m *OrganizationModule) Dependencies() []string { return []string{"user"} }
func (m *OrganizationModule) Initialize(ctx core.ModuleContext) error {
	ctx.Logger.Info("Initializing organization module (minimal)")
	return nil
}
func (m *OrganizationModule) Start(ctx context.Context) error { return nil }
func (m *OrganizationModule) Stop(ctx context.Context) error  { return nil }
func (m *OrganizationModule) Entities() []core.Entity {
	return []core.Entity{entities.NewOrganizationEntity()}
}
func (m *OrganizationModule) Routes() []core.RouteProvider       { return []core.RouteProvider{} }
func (m *OrganizationModule) EventHandlers() []core.EventHandler { return []core.EventHandler{} }
func (m *OrganizationModule) Services() []core.ServiceProvider   { return []core.ServiceProvider{} }
func (m *OrganizationModule) Middleware() []core.MiddlewareProvider {
	return []core.MiddlewareProvider{}
}
func (m *OrganizationModule) SwaggerPaths() []string { return []string{} }
