package services

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/AgileExecutives/serverbase/pkg/settings/repository"
	sbsettings "github.com/AgileExecutives/serverbase/pkg/settings/services"
	"github.com/AgileExecutives/shared-modules/invoice_number/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupInMemoryDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	// Automigrate required entities
	if err := db.AutoMigrate(&entities.InvoiceNumber{}, &entities.InvoiceNumberLog{}); err != nil {
		t.Fatalf("automigrate invoice entities failed: %v", err)
	}
	// Automigrate settings via repository helper
	repo := repository.NewSettingsRepository(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("automigrate settings failed: %v", err)
	}
	return db
}

func TestGenerateNextInvoiceNumber_Default(t *testing.T) {
	db := setupInMemoryDB(t)
	svc := NewInvoiceNumberServiceWithSettings(db, nil)

	ctx := context.Background()
	tenantID := uint(1)
	orgID := uint(10)

	resp, err := svc.GenerateNextInvoiceNumber(ctx, tenantID, orgID)
	if err != nil {
		t.Fatalf("GenerateNextInvoiceNumber failed: %v", err)
	}

	expectedPrefix := "INV"
	expectedSeq := 1

	if resp.Sequence != expectedSeq {
		t.Fatalf("expected sequence %d got %d", expectedSeq, resp.Sequence)
	}
	if resp.TenantID != tenantID {
		t.Fatalf("tenant id mismatch")
	}
	if resp.InvoiceNumber == "" {
		t.Fatalf("invoice number empty")
	}
	// basic check: contains prefix and year
	if !strings.HasPrefix(resp.InvoiceNumber, expectedPrefix) {
		t.Fatalf("expected prefix %s in %s", expectedPrefix, resp.InvoiceNumber)
	}

	// second call increments
	resp2, err := svc.GenerateNextInvoiceNumber(ctx, tenantID, orgID)
	if err != nil {
		t.Fatalf("second GenerateNextInvoiceNumber failed: %v", err)
	}
	if resp2.Sequence != 2 {
		t.Fatalf("expected sequence 2 got %d", resp2.Sequence)
	}
}

func TestGenerateNextInvoiceNumber_WithSettings(t *testing.T) {
	db := setupInMemoryDB(t)
	// create settings repo/service
	settingsRepo := repository.NewSettingsRepository(db)
	settingsSvc := sbsettings.NewSettingsService(settingsRepo)

	svc := NewInvoiceNumberServiceWithSettings(db, settingsSvc)

	// set settings for invoice domain
	tenantID := uint(1)
	orgID := uint(20)

	// set a custom prefix and padding
	orgStr := strconv.FormatUint(uint64(orgID), 10)
	if err := settingsSvc.SetSetting(tenantID, orgStr, "invoice", "invoice_prefix", "ACME", "json"); err != nil {
		t.Fatalf("SetSetting failed: %v", err)
	}
	// avoid storing numeric types in SQLite JSON column for this test

	resp, err := svc.GenerateNextInvoiceNumber(context.Background(), tenantID, orgID)
	if err != nil {
		t.Fatalf("GenerateNextInvoiceNumber with settings failed: %v", err)
	}

	if resp.Sequence != 1 {
		t.Fatalf("expected sequence 1 got %d", resp.Sequence)
	}
	if resp.TenantID != tenantID {
		t.Fatalf("tenant id mismatch")
	}
	if len(resp.InvoiceNumber) == 0 {
		t.Fatalf("invoice number empty")
	}
	if resp.InvoiceNumber[:4] != "ACME" {
		t.Fatalf("expected prefix ACME, got %s", resp.InvoiceNumber)
	}
}
