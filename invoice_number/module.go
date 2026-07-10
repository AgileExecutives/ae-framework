package invoicenumber

import (
	"context"

	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/shared-modules/invoice_number/entities"
	repo "github.com/AgileExecutives/shared-modules/invoice_number/repo"
	"github.com/AgileExecutives/shared-modules/invoice_number/routes"
	"github.com/AgileExecutives/shared-modules/invoice_number/services"
	sbsettingsrepo "github.com/AgileExecutives/serverbase/pkg/settings/repository"
	sbsettings "github.com/AgileExecutives/serverbase/pkg/settings/services"
)

// InvoiceNumberModule represents the invoice number generation module
type InvoiceNumberModule struct {
	invoiceNumberService *services.InvoiceNumberService
	invoiceNumberRoutes  *routes.InvoiceNumberRoutes
}

// NewInvoiceNumberModule creates a new invoice number module instance
func NewInvoiceNumberModule() core.Module {
	return &InvoiceNumberModule{}
}

func (m *InvoiceNumberModule) Name() string {
	return "invoice_number"
}

func (m *InvoiceNumberModule) Version() string {
	return "1.0.0"
}

func (m *InvoiceNumberModule) Dependencies() []string {
	return []string{"base"}
}

func (m *InvoiceNumberModule) Initialize(ctx core.ModuleContext) error {
	ctx.Logger.Info("Initializing invoice number module...")

	// Initialize service (prefer repo-backed)
	gormRepo := repo.NewGormInvoiceNumberRepo(ctx.DB)
	m.invoiceNumberService = services.NewInvoiceNumberServiceWithRepo(gormRepo)
	// Wire settings service so handlers and consumers get per-tenant formats
	settingsRepo := sbsettingsrepo.NewSettingsRepository(ctx.DB)
	settingsSvc := sbsettings.NewSettingsService(settingsRepo)
	m.invoiceNumberService.SetSettingsService(settingsSvc)

	// Initialize routes
	m.invoiceNumberRoutes = routes.NewInvoiceNumberRoutes(m.invoiceNumberService)

	ctx.Logger.Info("Invoice number module initialized successfully")
	return nil
}

func (m *InvoiceNumberModule) Start(ctx context.Context) error {
	return nil
}

func (m *InvoiceNumberModule) Stop(ctx context.Context) error {
	return nil
}

func (m *InvoiceNumberModule) Entities() []core.Entity {
	return []core.Entity{
		entities.NewInvoiceNumberEntity(),
		entities.NewInvoiceNumberLogEntity(),
	}
}

func (m *InvoiceNumberModule) Routes() []core.RouteProvider {
	return []core.RouteProvider{
		m.invoiceNumberRoutes,
	}
}

func (m *InvoiceNumberModule) EventHandlers() []core.EventHandler {
	return []core.EventHandler{}
}

func (m *InvoiceNumberModule) Services() []core.ServiceProvider {
	return []core.ServiceProvider{}
}

func (m *InvoiceNumberModule) Middleware() []core.MiddlewareProvider {
	return []core.MiddlewareProvider{}
}

func (m *InvoiceNumberModule) SwaggerPaths() []string {
	return []string{
		"./modules/invoice_number/handlers",
		"./modules/invoice_number/entities",
	}
}
