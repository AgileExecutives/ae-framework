package bootstrap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AgileExecutives/ae-framework/serverbase/modules/base"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/templates"
	pkgconfig "github.com/AgileExecutives/ae-framework/serverbase/pkg/config"
)

// ensureRepoRoot moves to repository root so relative module paths resolve.
func ensureRepoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(cwd, "..", "..", ".."))
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir repo root: %v", err)
	}
	return root
}

func TestApplication_RegisterContracts_EndToEnd(t *testing.T) {
	root := ensureRepoRoot(t)

	// Use in-memory DB
	os.Setenv("USE_IN_MEMORY_DB", "true")
	defer os.Unsetenv("USE_IN_MEMORY_DB")

	cfg := pkgconfig.Load()
	cfg.Server.Mode = "test"

	app := NewApplication(cfg)

	// Register modules required for migrations and template seeding
	if err := app.RegisterModule(base.NewBaseModule()); err != nil {
		t.Fatalf("register base module: %v", err)
	}
	if err := app.RegisterModule(templates.NewTemplatesModule()); err != nil {
		t.Fatalf("register templates module: %v", err)
	}

	// Initialize core services (connects DB and builds ModuleContext)
	if err := app.initializeCoreServices(); err != nil {
		t.Fatalf("initialize core services: %v", err)
	}

	// Initialize modules (runs their Init hooks which may auto-migrate template tables)
	if err := app.registry.InitializeAll(app.context); err != nil {
		t.Fatalf("initialize modules: %v", err)
	}

	// Run migrations to create tenant table and other entities
	if err := app.runMigrations(); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	// Insert tenants so contract registration has targets
	db := app.DB()
	if err := db.Exec(`INSERT INTO tenants (id, customer_id, name, slug) VALUES (1,1,'t1','t1'), (2,1,'t2','t2')`).Error; err != nil {
		t.Fatalf("insert tenants: %v", err)
	}

	// Create shared-modules contract files that RegisterAllContracts will discover
	bookingDir := filepath.Join(root, "shared-modules", "booking", "contracts")
	if err := os.MkdirAll(bookingDir, 0o755); err != nil {
		t.Fatalf("mkdir booking: %v", err)
	}
	bookingFile := filepath.Join(bookingDir, "booking_confirmation-contract.json")
	bookingContent := `{"type":"object","properties":{"Booking":{"type":"object"}},"template_key":"booking_confirmation"}`
	if err := os.WriteFile(bookingFile, []byte(bookingContent), 0o644); err != nil {
		t.Fatalf("write booking contract: %v", err)
	}

	invoiceDir := filepath.Join(root, "shared-modules", "invoice", "contracts")
	if err := os.MkdirAll(invoiceDir, 0o755); err != nil {
		t.Fatalf("mkdir invoice: %v", err)
	}
	invoiceFile := filepath.Join(invoiceDir, "std_invoice-contract.json")
	invoiceContent := `{"type":"object","properties":{"Invoice":{"type":"object"}},"template_key":"invoice"}`
	if err := os.WriteFile(invoiceFile, []byte(invoiceContent), 0o644); err != nil {
		t.Fatalf("write invoice contract: %v", err)
	}

	// Finally register contracts (this mirrors Application.Initialize ordering)
	if err := app.registerContracts(); err != nil {
		t.Fatalf("register contracts: %v", err)
	}

	// Assert template_contracts were created for tenants
	var cnt int64
	if err := db.Table("template_contracts").Count(&cnt).Error; err != nil {
		t.Fatalf("count template_contracts: %v", err)
	}
	if cnt == 0 {
		t.Fatalf("expected template_contracts to be registered, got 0")
	}

	// sanity: ensure at least one contract per tenant was created
	var cntT1 int64
	if err := db.Table("template_contracts").Where("tenant_id = ?", 1).Count(&cntT1).Error; err != nil {
		t.Fatalf("count tenant1: %v", err)
	}
	if cntT1 == 0 {
		t.Fatalf("expected contracts for tenant 1, got 0")
	}
}
