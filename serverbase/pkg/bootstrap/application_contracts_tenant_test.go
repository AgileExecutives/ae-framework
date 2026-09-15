package bootstrap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AgileExecutives/ae-framework/serverbase/modules/base"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/templates"
	pkgconfig "github.com/AgileExecutives/ae-framework/serverbase/pkg/config"
)

// Test that after startup/registering contracts the important contracts are
// associated with tenant ID 1 (not tenant 0 or omitted). This prevents the
// seeded template contracts from being global/tenantless.
func TestContractsAreRegisteredForTenant1(t *testing.T) {
	root := ensureRepoRoot(t)

	// Use in-memory DB
	os.Setenv("USE_IN_MEMORY_DB", "true")
	defer os.Unsetenv("USE_IN_MEMORY_DB")

	cfg := pkgconfig.Load()
	cfg.Server.Mode = "test"

	app := NewApplication(cfg)

	// Register base and templates modules (templates seeds standard contracts)
	if err := app.RegisterModule(base.NewBaseModule()); err != nil {
		t.Fatalf("register base module: %v", err)
	}
	if err := app.RegisterModule(templates.NewTemplatesModule()); err != nil {
		t.Fatalf("register templates module: %v", err)
	}

	if err := app.initializeCoreServices(); err != nil {
		t.Fatalf("initialize core services: %v", err)
	}
	if err := app.registry.InitializeAll(app.context); err != nil {
		t.Fatalf("initialize modules: %v", err)
	}
	if err := app.runMigrations(); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	// Ensure shared-modules contract files exist so fallback registration picks them up
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

	// Perform seeding which will create tenant 1 and call tenant-level registration
	if err := app.seedDatabase(); err != nil {
		t.Fatalf("seed database: %v", err)
	}

	// Now register all contracts (mirrors Application.Initialize ordering)
	if err := app.registerContracts(); err != nil {
		t.Fatalf("register contracts: %v", err)
	}

	db := app.DB()

	// Check that expected keys have at least one contract with tenant_id = 1
	keys := []struct{ module, key string }{
		{"user", "welcome"},
		{"user", "password_reset"},
		{"booking", "booking_confirmation"},
		{"invoice", "invoice"},
	}

	for _, k := range keys {
		var cnt int64
		if err := db.Table("template_contracts").Where("module = ? AND template_key = ? AND tenant_id = ?", k.module, k.key, 1).Count(&cnt).Error; err != nil {
			t.Fatalf("count contract %s/%s: %v", k.module, k.key, err)
		}
		if cnt == 0 {
			t.Fatalf("expected contract %s/%s to be registered for tenant 1, none found", k.module, k.key)
		}
	}
}
