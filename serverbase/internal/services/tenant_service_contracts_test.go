package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/AgileExecutives/ae-framework/serverbase/internal/models"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/startup"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// fakeTenantRepo is a minimal in-memory TenantRepo used for tests.
type fakeTenantRepo struct {
	nextID uint
}

func (f *fakeTenantRepo) FindByID(ctx context.Context, id uint) (*models.Tenant, error) {
	return &models.Tenant{ID: id}, nil
}
func (f *fakeTenantRepo) FindBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	return nil, nil
}
func (f *fakeTenantRepo) Save(ctx context.Context, t *models.Tenant) error {
	f.nextID++
	t.ID = f.nextID
	return nil
}
func (f *fakeTenantRepo) List(ctx context.Context) ([]models.Tenant, error) {
	return []models.Tenant{}, nil
}

func TestTenantCreateTriggersContractRegistration(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	// repo root is three levels up
	root = filepath.Clean(filepath.Join(root, "..", "..", ".."))
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Exec(`CREATE TABLE template_contracts (id integer primary key, tenant_id integer, module text, template_key text, variable_schema json, default_sample_data json, created_at datetime, updated_at datetime)`).Error; err != nil {
		t.Fatalf("create template_contracts table: %v", err)
	}

	// create a contract file to be discovered
	bookingDir := filepath.Join(root, "modules", "booking", "contracts")
	if err := os.MkdirAll(bookingDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	bookingFile := filepath.Join(bookingDir, "booking_confirmation-contract.json")
	if err := os.WriteFile(bookingFile, []byte(`{"type":"object","properties":{"Booking":{"type":"object"}},"template_key":"booking_confirmation"}`), 0o644); err != nil {
		t.Fatalf("write contract: %v", err)
	}

	repo := &fakeTenantRepo{nextID: 0}
	tenantSvc := NewTenantService(repo, nil)
	tenantSvc.SetPostCreateHook(func(tid uint) error {
		return startup.RegisterContractsForTenant(db, tid)
	})

	tnt, err := tenantSvc.CreateTenant(context.Background(), models.TenantCreateRequest{Name: "Acme"})
	if err != nil {
		t.Fatalf("CreateTenant: %v", err)
	}

	var cnt int64
	if err := db.Table("template_contracts").Where("tenant_id = ?", tnt.ID).Count(&cnt).Error; err != nil {
		t.Fatalf("count contracts: %v", err)
	}
	if cnt == 0 {
		t.Fatalf("expected contracts for tenant %d, got 0", tnt.ID)
	}
}
