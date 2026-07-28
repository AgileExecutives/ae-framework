package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/AgileExecutives/ae-framework/serverbase/pkg/settings/repository"
	sbsettings "github.com/AgileExecutives/ae-framework/serverbase/pkg/settings/services"
	"github.com/AgileExecutives/ae-framework/shared-modules/invoice_number/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupInMemoryDBForFormatTests(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&entities.InvoiceNumber{}, &entities.InvoiceNumberLog{}); err != nil {
		t.Fatalf("automigrate invoice entities failed: %v", err)
	}
	repo := repository.NewSettingsRepository(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("automigrate settings failed: %v", err)
	}
	return db
}

func TestInvoiceNumber_Formats_TableDriven(t *testing.T) {
	db := setupInMemoryDBForFormatTests(t)
	settingsRepo := repository.NewSettingsRepository(db)
	settingsSvc := sbsettings.NewSettingsService(settingsRepo)

	svc := NewInvoiceNumberServiceWithSettings(db, settingsSvc)

	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	tests := []struct {
		name     string
		settings map[string]interface{}
	}{
		{"default", map[string]interface{}{}},
		{"custom_prefix", map[string]interface{}{"invoice_prefix": "ACME"}},
		{"yy_year_m_single_month_sep", map[string]interface{}{"year_format": "YY", "month_format": "M", "separator": "_"}},
		{"no_prefix", map[string]interface{}{"invoice_prefix": ""}},
	}

	tenantID := uint(1)

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			orgID := uint(100 + i)
			orgStr := strconv.FormatUint(uint64(orgID), 10)
			// Clear domain settings for invoice first
			_ = settingsSvc.DeleteDomainSettings(tenantID, orgStr, "invoice")

			for k, v := range tc.settings {
				if err := settingsSvc.SetSetting(tenantID, orgStr, "invoice", k, v, "json"); err != nil {
					t.Fatalf("SetSetting failed: %v", err)
				}
			}

			resp, err := svc.GenerateNextInvoiceNumber(context.Background(), tenantID, orgID)
			if err != nil {
				t.Fatalf("GenerateNextInvoiceNumber failed: %v", err)
			}

			// Build expected string dynamically according to settings used
			// Determine config as service would
			cfg := DefaultInvoiceConfig()
			if v, ok := tc.settings["invoice_prefix"].(string); ok {
				cfg.Prefix = v
			}
			if v, ok := tc.settings["year_format"].(string); ok {
				cfg.YearFormat = v
			}
			if v, ok := tc.settings["month_format"].(string); ok {
				cfg.MonthFormat = v
			}
			// Note: avoid numeric settings in SQLite-backed tests (JSON scanning issues).
			if v, ok := tc.settings["separator"].(string); ok {
				cfg.Separator = v
			}

			parts := []string{}
			if cfg.Prefix != "" {
				parts = append(parts, cfg.Prefix)
			}
			if cfg.YearFormat == "YYYY" {
				parts = append(parts, strconv.Itoa(year))
			} else if cfg.YearFormat == "YY" {
				parts = append(parts, strconv.Itoa(year%100))
			}
			// Build month
			if cfg.MonthFormat == "MM" {
				parts = append(parts, sprintfTwoDigit(month))
			} else if cfg.MonthFormat == "M" {
				parts = append(parts, strconv.Itoa(month))
			}

			// Sequence for first call is 1 with padding
			seq := fmtSequence(cfg.Padding, 1)
			parts = append(parts, seq)

			expected := strings.Join(parts, cfg.Separator)

			if resp.InvoiceNumber != expected {
				t.Fatalf("expected invoice number %s got %s", expected, resp.InvoiceNumber)
			}
		})
	}

	// Direct-config test for padding behavior (avoids storing numeric setting)
	t.Run("direct_config_padding", func(t *testing.T) {
		cfg := DefaultInvoiceConfig()
		cfg.Padding = 2
		cfg.Prefix = "PAD"
		// use a fresh org id
		orgID := uint(200)

		resp, err := svc.GenerateInvoiceNumber(context.Background(), tenantID, orgID, cfg)
		if err != nil {
			t.Fatalf("GenerateInvoiceNumber failed: %v", err)
		}

		// expected: PAD-YYYY-MM-01 with two-digit padding
		expected := strings.Join([]string{cfg.Prefix, strconv.Itoa(year), sprintfTwoDigit(month), fmtSequence(cfg.Padding, 1)}, cfg.Separator)
		if resp.InvoiceNumber != expected {
			t.Fatalf("expected invoice number %s got %s", expected, resp.InvoiceNumber)
		}
	})
}

// helper to format two-digit month
func sprintfTwoDigit(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

// helper to format sequence with padding
func fmtSequence(padding int, seq int) string {
	return fmt.Sprintf("%0*d", padding, seq)
}
