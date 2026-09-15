package base

// To regenerate swagger docs for this module run from serverbase/modules/base:
//
//go:generate swag init -g doc.go --dir .,handlers,../user/entities,../user/models,../../internal/models,../../pkg/utils --output docs --instanceName base

import (
	"context"

	internalTenantSvc "github.com/AgileExecutives/ae-framework/serverbase/internal/services"
	basedocs "github.com/AgileExecutives/ae-framework/serverbase/modules/base/docs"
	baseRepo "github.com/AgileExecutives/ae-framework/serverbase/modules/base/repo"
	baseServices "github.com/AgileExecutives/ae-framework/serverbase/modules/base/services"
	templateServices "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/services"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/user/entities"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/user/events"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/user/handlers"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/user/middleware"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/user/services"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/repos"
	settingsentities "github.com/AgileExecutives/ae-framework/serverbase/pkg/settings/entities"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/startup"
)

// BaseModule provides core authentication, user management, and contact functionality
type BaseModule struct {
	authHandlers         *handlers.AuthHandlers
	contactHandlers      *handlers.ContactHandlers
	healthHandlers       *handlers.HealthHandlers
	userSettingsHandlers *handlers.UserSettingsHandlers
	userSettingsService  *baseServices.UserSettingsService
	authService          *services.AuthService
	contactService       *services.ContactService
	eventHandlers        *events.BaseEventHandlers
	authMiddleware       *middleware.AuthMiddleware
	moduleContext        core.ModuleContext
}

// NewBaseModule creates a new base module instance
func NewBaseModule() core.Module {
	return &BaseModule{}
}

func (m *BaseModule) Name() string {
	return "user"
}

func (m *BaseModule) Version() string {
	return "1.0.0"
}

func (m *BaseModule) Dependencies() []string {
	return []string{} // No dependencies - this is the base module
}

func (m *BaseModule) Initialize(ctx core.ModuleContext) error {
	ctx.Logger.Info("Initializing base module...")

	// Store context for later use
	m.moduleContext = ctx

	// Initialize services
	m.authService = services.NewAuthService(ctx.DB, ctx.Logger)

	// Create tenant service and wire into auth service so tenant creation
	// side-effects (buckets, etc.) are centralized.
	rf := repos.NewGormRepoFactory(ctx.DB)
	tenantSvc := internalTenantSvc.NewTenantService(rf.TenantRepo(), nil)
	// Wire event bus into tenant service so it can publish tenant lifecycle events
	if ctx.EventBus != nil {
		tenantSvc.SetEventBus(ctx.EventBus)
	}
	// Register post-create hook to register template contracts for the new tenant.
	// Delegate to the startup package so registration logic is centralized and
	// can be reused by admin endpoints or CLI commands.
	tenantSvc.SetPostCreateHook(func(tid uint) error {
		// Note: Post-create hook must not fail tenant creation, so log/return
		// errors to the caller but do not convert them into creation failures.
		if err := startup.RegisterContractsForTenant(m.moduleContext, tid); err != nil {
			m.moduleContext.Logger.Warn("post-create contract registration failed", "tenant", tid, "err", err)
			// continue to publish event even if contract registration failed
		}

		// Publish a standard tenant created event on the module event bus so
		// other modules (e.g., documents) can react (create buckets, etc.).
		if m.moduleContext.EventBus != nil {
			// publish tenant ID as numeric payload
			if err := m.moduleContext.EventBus.Publish("tenant.created", tid); err != nil {
				m.moduleContext.Logger.Warn("failed to publish tenant.created event", "tenant", tid, "err", err)
			}
		}

		return nil
	})
	m.authService.SetTenantService(tenantSvc)

	// Initialize handlers (pass authService for newer handler constructors)
	m.authHandlers = handlers.NewAuthHandlers(ctx, m.authService, ctx.Logger)
	contactRepo := rf.ContactRepo()
	m.contactService = services.NewContactServiceWithRepo(contactRepo, ctx.Logger)
	m.contactHandlers = handlers.NewContactHandlers(ctx, m.contactService, ctx.Logger)
	m.healthHandlers = handlers.NewHealthHandlers(ctx, ctx.Logger)
	// Wire user settings service + handlers
	// Construct UserSettingsService with a GORM-backed repo
	m.userSettingsService = baseServices.NewUserSettingsService(baseRepo.NewGormUserSettingsRepo(ctx.DB))
	m.userSettingsHandlers = handlers.NewUserSettingsHandlers(ctx, ctx.Logger)

	// Initialize event handlers
	m.eventHandlers = events.NewBaseEventHandlers(ctx.EventBus, ctx.Logger)

	// Initialize middleware
	m.authMiddleware = middleware.NewAuthMiddleware(ctx, ctx.Logger)

	// Register pre-generated swagger docs so they appear in the combined spec.
	if ctx.DocRegistry != nil {
		ctx.DocRegistry.RegisterDoc(m.Name(), basedocs.SwaggerInfobase.ReadDoc())
	}

	ctx.Logger.Info("Base module initialized successfully")
	return nil
}

func (m *BaseModule) Start(ctx context.Context) error {
	// Set module registry in auth handler now that all modules are initialized
	// This needs to be done in Start() because it happens after all Initialize() calls
	if m.authHandlers != nil && m.moduleContext.ModuleRegistry != nil {
		m.authHandlers.SetModuleRegistry(m.moduleContext.ModuleRegistry)
		m.moduleContext.Logger.Info("Base module started - auth handlers configured with module registry for template seeding")
	} else {
		m.moduleContext.Logger.Warn("Base module started - module registry not available for template seeding")
	}
	return nil
}

func (m *BaseModule) Stop(ctx context.Context) error {
	// Stop any background services if needed
	return nil
}

func (m *BaseModule) Entities() []core.Entity {
	return []core.Entity{
		entities.NewUserEntity(),
		entities.NewTenantEntity(),
		entities.NewContactEntity(),
		entities.NewNewsletterEntity(),
		entities.NewTokenBlacklistEntity(),
		entities.NewUserSettingsEntity(),
		settingsentities.NewSettingDefinitionEntity(),
		settingsentities.NewSettingEntity(),
	}
}

func (m *BaseModule) Routes() []core.RouteProvider {
	return []core.RouteProvider{
		handlers.NewAuthRoutes(m.authHandlers),
		handlers.NewContactRoutes(m.contactHandlers),
		handlers.NewUserSettingsRoutes(m.userSettingsHandlers),
		handlers.NewHealthRoutes(m.healthHandlers),
	}
}

func (m *BaseModule) EventHandlers() []core.EventHandler {
	return []core.EventHandler{
		events.NewUserCreatedHandler(m.eventHandlers),
		events.NewUserLoginHandler(m.eventHandlers),
		events.NewContactFormSubmittedHandler(m.eventHandlers),
	}
}

func (m *BaseModule) Services() []core.ServiceProvider {
	return []core.ServiceProvider{
		services.NewAuthServiceProvider(m.authService),
		services.NewContactServiceProvider(m.contactService),
	}
}

func (m *BaseModule) Middleware() []core.MiddlewareProvider {
	return []core.MiddlewareProvider{
		middleware.NewAuthMiddlewareProvider(m.authMiddleware),
	}
}

func (m *BaseModule) SwaggerPaths() []string {
	return []string{
		"./modules/user/handlers",
		"./modules/user/entities",
	}
}

// RegisterContracts implements optional module-level contract registration for
// the base module. Currently there are no base contracts; this exists to keep
// the contract bootstrap pathway consistent.
func (m *BaseModule) RegisterContracts(ctx core.ModuleContext, tenantID uint) error {
	registrar := templateServices.NewContractRegistrar(ctx.DB)
	return baseServices.RegisterBaseContracts(registrar, tenantID)
}
