package invoice

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	templateServices "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/services"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framework/shared-modules/invoice/entities"
	"github.com/AgileExecutives/ae-framework/shared-modules/invoice/handlers"
	repo "github.com/AgileExecutives/ae-framework/shared-modules/invoice/repo"
	"github.com/AgileExecutives/ae-framework/shared-modules/invoice/routes"
	"github.com/AgileExecutives/ae-framework/shared-modules/invoice/services"
)

// CoreModule implements the core.Module interface for the invoice module
type CoreModule struct {
	invoiceService *services.InvoiceService
	invoiceHandler *handlers.InvoiceHandler
	invoiceRoutes  *routes.InvoiceRoutes
}

// NewCoreModule creates a new invoice module instance
func NewCoreModule() *CoreModule {
	return &CoreModule{}
}

// Name returns the module name
func (m *CoreModule) Name() string {
	return "invoice"
}

// Version returns the module version
func (m *CoreModule) Version() string {
	return "1.0.0"
}

// Dependencies returns module dependencies
func (m *CoreModule) Dependencies() []string {
	return []string{"user"} // Depends on user (base) module for auth and database
}

// Initialize sets up the module with dependencies
func (m *CoreModule) Initialize(ctx core.ModuleContext) error {
	// Initialize service
	gormRepo := repo.NewGormInvoiceRepo(ctx.DB)
	m.invoiceService = services.NewInvoiceServiceWithRepo(gormRepo)

	// Initialize handler
	m.invoiceHandler = handlers.NewInvoiceHandler(m.invoiceService)

	// Initialize routes
	m.invoiceRoutes = routes.NewInvoiceRoutes(m.invoiceHandler)

	return nil
}

// Start starts the module
func (m *CoreModule) Start(ctx context.Context) error {
	return nil
}

// Stop stops the module
func (m *CoreModule) Stop(ctx context.Context) error {
	return nil
}

// Entities returns the list of entities for database migration
func (m *CoreModule) Entities() []core.Entity {
	return []core.Entity{
		entities.NewInvoiceEntity(),
		entities.NewInvoiceItemEntity(),
	}
}

// Routes returns the list of route handlers for this module
func (m *CoreModule) Routes() []core.RouteProvider {
	providers := []core.RouteProvider{}
	if m.invoiceRoutes != nil {
		providers = append(providers, m.invoiceRoutes)
	}
	return providers
}

// EventHandlers returns event handlers for the module
func (m *CoreModule) EventHandlers() []core.EventHandler {
	return []core.EventHandler{}
}

// Middleware returns middleware providers
func (m *CoreModule) Middleware() []core.MiddlewareProvider {
	return []core.MiddlewareProvider{}
}

// Services returns service providers
func (m *CoreModule) Services() []core.ServiceProvider {
	return []core.ServiceProvider{}
}

// SwaggerPaths returns Swagger documentation paths
func (m *CoreModule) SwaggerPaths() []string {
	return []string{}
}

// RegisterContracts registers invoice module contracts from shared-modules/invoice/contracts
func (m *CoreModule) RegisterContracts(ctx core.ModuleContext, tenantID uint) error {
	registrar := templateServices.NewContractRegistrar(ctx.DB)
	contractsDir := filepath.Join("shared-modules", "invoice", "contracts")
	files, err := os.ReadDir(contractsDir)
	if err != nil {
		return nil
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".json") {
			if err := registrar.RegisterContractFromFile(tenantID, "invoice", filepath.Join(contractsDir, f.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetInvoiceService returns the invoice service instance
func (m *CoreModule) GetInvoiceService() *services.InvoiceService {
	return m.invoiceService
}
