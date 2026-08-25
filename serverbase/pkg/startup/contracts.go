package startup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	baseServices "github.com/AgileExecutives/ae-framework/serverbase/modules/base/services"
	emailServices "github.com/AgileExecutives/ae-framework/serverbase/modules/email/services"
	templateServices "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/services"
	userServices "github.com/AgileExecutives/ae-framework/serverbase/modules/user/services"
	"gorm.io/gorm"
)

// RegisterAllContracts registers all module contracts with the template system
func RegisterAllContracts(db *gorm.DB) error {
	fmt.Println("🔧 Registering template contracts for all tenants...")

	// Get all tenants from database
	var tenants []struct{ ID uint }
	if err := db.Table("tenants").Select("id").Where("deleted_at IS NULL").Find(&tenants).Error; err != nil {
		return fmt.Errorf("failed to fetch tenants: %w", err)
	}

	fmt.Printf("Found %d tenants for contract registration\n", len(tenants))

	// Register for each tenant using the single-tenant helper
	for _, tenant := range tenants {
		fmt.Printf("📝 Registering contracts for tenant ID: %d\n", tenant.ID)
		if err := RegisterContractsForTenant(db, tenant.ID); err != nil {
			return fmt.Errorf("failed to register contracts for tenant %d: %w", tenant.ID, err)
		}
	}

	// Verify contracts were created
	var contractCount int64
	db.Table("template_contracts").Count(&contractCount)
	fmt.Printf("✅ All contracts registered successfully - Total contracts in database: %d\n", contractCount)
	return nil
}

// RegisterContractsForTenant registers all discovered JSON contract files for a single tenant.
// It scans `modules/*/contracts/*.json` and `shared-modules/*/contracts/*.json` and
// calls the ContractRegistrar for each file.
func RegisterContractsForTenant(db *gorm.DB, tenantID uint) error {
	registrar := templateServices.NewContractRegistrar(db)

	// Prefer module-level registration functions. Each module should register
	// its own contracts using the provided registrar. Fall back to no-ops.
	if err := baseServices.RegisterBaseContracts(registrar, tenantID); err != nil {
		return fmt.Errorf("base module register contracts: %w", err)
	}
	if err := emailServices.RegisterEmailContracts(registrar, tenantID); err != nil {
		return fmt.Errorf("email module register contracts: %w", err)
	}
	if err := userServices.RegisterUserContracts(registrar, tenantID); err != nil {
		return fmt.Errorf("user module register contracts: %w", err)
	}

	// For shared-modules, fall back to scanning their contracts directories.
	contractFiles := map[string][]string{}
	entries, _ := os.ReadDir("shared-modules")
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		modName := e.Name()
		contractsDir := filepath.Join("shared-modules", modName, "contracts")
		files, err := os.ReadDir(contractsDir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if strings.HasSuffix(f.Name(), ".json") {
				contractFiles[modName] = append(contractFiles[modName], filepath.Join(contractsDir, f.Name()))
			}
		}
	}

	for moduleName, files := range contractFiles {
		for _, f := range files {
			if err := registrar.RegisterContractFromFile(tenantID, moduleName, f); err != nil {
				return fmt.Errorf("register contract %s for module %s tenant %d: %w", f, moduleName, tenantID, err)
			}
		}
	}

	return nil
}
