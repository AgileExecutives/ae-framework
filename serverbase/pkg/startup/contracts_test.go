package startup

import (
	"os"
	"path/filepath"
	"testing"

	templateServices "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/services"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ensureRepoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	// repo root is three levels up from serverbase/pkg/startup
	root := filepath.Clean(filepath.Join(cwd, "..", "..", ".."))
	return root
}

func TestRegisterAllContractsFromFiles(t *testing.T) {
	root := ensureRepoRoot(t)
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir repo root: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	// create tables
	if err := db.Exec(`CREATE TABLE tenants (id integer primary key, deleted_at datetime)`).Error; err != nil {
		t.Fatalf("create tenants table: %v", err)
	}
	if err := db.Exec(`CREATE TABLE template_contracts (id integer primary key, tenant_id integer, module text, template_key text, variable_schema json, default_sample_data json, created_at datetime, updated_at datetime)`).Error; err != nil {
		t.Fatalf("create template_contracts table: %v", err)
	}

	// create two tenants
	if err := db.Exec(`INSERT INTO tenants (id) VALUES (1), (2)`).Error; err != nil {
		t.Fatalf("insert tenants: %v", err)
	}

	// Create contract file for booking shared-module
	bookingDir := filepath.Join(root, "shared-modules", "booking", "contracts")
	if err := os.MkdirAll(bookingDir, 0o755); err != nil {
		t.Fatalf("mkdir booking contracts: %v", err)
	}
	bookingFile := filepath.Join(bookingDir, "booking_confirmation-contract.json")
	bookingContent := `{"type":"object","properties":{"Booking":{"type":"object"}},"template_key":"booking_confirmation"}`
	if err := os.WriteFile(bookingFile, []byte(bookingContent), 0o644); err != nil {
		t.Fatalf("write booking contract: %v", err)
	}

	// Create contract file for invoice shared-module
	invoiceDir := filepath.Join(root, "shared-modules", "invoice", "contracts")
	if err := os.MkdirAll(invoiceDir, 0o755); err != nil {
		t.Fatalf("mkdir invoice contracts: %v", err)
	}
	invoiceFile := filepath.Join(invoiceDir, "std_invoice-contract.json")
	invoiceContent := `{"type":"object","properties":{"Invoice":{"type":"object"}},"template_key":"invoice"}`
	if err := os.WriteFile(invoiceFile, []byte(invoiceContent), 0o644); err != nil {
		t.Fatalf("write invoice contract: %v", err)
	}

	// Create contract files for user/auth module
	userDir := filepath.Join(root, "modules", "user", "contracts")
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		t.Fatalf("mkdir user contracts: %v", err)
	}
	welcomeFile := filepath.Join(userDir, "welcome-contract.json")
	welcomeContent := `{"type":"object","properties":{"User":{"type":"object"}},"template_key":"welcome"}`
	if err := os.WriteFile(welcomeFile, []byte(welcomeContent), 0o644); err != nil {
		t.Fatalf("write welcome contract: %v", err)
	}
	pwdFile := filepath.Join(userDir, "password_reset-contract.json")
	pwdContent := `{"type":"object","properties":{"Reset":{"type":"object"}},"template_key":"password_reset"}`
	if err := os.WriteFile(pwdFile, []byte(pwdContent), 0o644); err != nil {
		t.Fatalf("write password reset contract: %v", err)
	}

	// Call RegisterAllContracts - our modules map currently includes booking as placeholder,
	// but RegisterAllContracts should pick up contracts from module registration functions.
	if err := RegisterAllContracts(db); err != nil {
		t.Fatalf("RegisterAllContracts: %v", err)
	}

	var cnt int64
	if err := db.Table("template_contracts").Count(&cnt).Error; err != nil {
		t.Fatalf("count template_contracts: %v", err)
	}
	if cnt == 0 {
		t.Fatalf("expected some template_contracts to be registered, got 0")
	}
}

func TestRegisterContractFromBytesCreatesTenantScopedContract(t *testing.T) {
	root := ensureRepoRoot(t)
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir repo root: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Exec(`CREATE TABLE template_contracts (id integer primary key, tenant_id integer, module text, template_key text, variable_schema json, default_sample_data json, created_at datetime, updated_at datetime)`).Error; err != nil {
		t.Fatalf("create template_contracts table: %v", err)
	}

	r := templateServices.NewContractRegistrar(db)
	// Use byte payload directly (simulates MinIO-read)
	b := []byte(`{"type":"object","properties":{"X":{"type":"string"}},"template_key":"from_bytes_test"}`)
	if err := r.RegisterContractFromBytes(42, "testing", "hint.json", b); err != nil {
		t.Fatalf("RegisterContractFromBytes: %v", err)
	}
	var cnt int64
	if err := db.Table("template_contracts").Where("tenant_id = ?", 42).Count(&cnt).Error; err != nil {
		t.Fatalf("count tenant contracts: %v", err)
	}
	if cnt != 1 {
		t.Fatalf("expected 1 contract for tenant 42, got %d", cnt)
	}
}
