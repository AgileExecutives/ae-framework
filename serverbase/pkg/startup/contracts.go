package startup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	templateServices "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/services"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
)

// RegisterAllContracts registers all module contracts with the template system.
// Variant 2: modules may implement an optional method `RegisterContracts(ctx core.ModuleContext, tenantID uint) error`.
func RegisterAllContracts(ctx core.ModuleContext) error {
	db := ctx.DB
	fmt.Println("🔧 Registering template contracts for all tenants...")

	// Get all tenants from database
	var tenants []struct{ ID uint }
	if err := db.Table("tenants").Select("id").Where("deleted_at IS NULL").Find(&tenants).Error; err != nil {
		return fmt.Errorf("failed to fetch tenants: %w", err)
	}

	fmt.Printf("Found %d tenants for contract registration\n", len(tenants))

	// Prepare fallback shared-module contract files map (scan once)
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

	// Iterate tenants and invoke per-module registration when implemented.
	for _, tenant := range tenants {
		fmt.Printf("📝 Registering contracts for tenant ID: %d\n", tenant.ID)

		// Iterate all registered modules and call RegisterContracts if available
		modules := ctx.ModuleRegistry.GetAll()
		for _, mod := range modules {
			if rc, ok := mod.(interface {
				RegisterContracts(core.ModuleContext, uint) error
			}); ok {
				if err := rc.RegisterContracts(ctx, tenant.ID); err != nil {
					return fmt.Errorf("module %s RegisterContracts: %w", mod.Name(), err)
				}
			}
		}

		// Fallback: register any shared-modules contract files found on disk
		registrar := templateServices.NewContractRegistrar(db)
		for moduleName, files := range contractFiles {
			for _, f := range files {
				if err := registrar.RegisterContractFromFile(tenant.ID, moduleName, f); err != nil {
					return fmt.Errorf("register contract %s for module %s tenant %d: %w", f, moduleName, tenant.ID, err)
				}
			}
		}
	}

	// Verify contracts were created
	var contractCount int64
	db.Table("template_contracts").Count(&contractCount)
	fmt.Printf("✅ All contracts registered successfully - Total contracts in database: %d\n", contractCount)
	return nil
}

// RegisterContractsForTenant registers contracts for a single tenant. This
// centralizes the per-tenant registration logic so it can be reused from the
// tenant post-create hook or admin endpoints.
func RegisterContractsForTenant(ctx core.ModuleContext, tenantID uint) error {
	db := ctx.DB

	// Iterate all registered modules and call RegisterContracts if available
	modules := ctx.ModuleRegistry.GetAll()
	for _, mod := range modules {
		if rc, ok := mod.(interface {
			RegisterContracts(core.ModuleContext, uint) error
		}); ok {
			if err := rc.RegisterContracts(ctx, tenantID); err != nil {
				return fmt.Errorf("module %s RegisterContracts: %w", mod.Name(), err)
			}
		}
	}

	// Fallback: scan shared-modules/<mod>/contracts for JSON files
	registrar := templateServices.NewContractRegistrar(db)
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
				if err := registrar.RegisterContractFromFile(tenantID, modName, filepath.Join(contractsDir, f.Name())); err != nil {
					return fmt.Errorf("register contract %s for module %s tenant %d: %w", f.Name(), modName, tenantID, err)
				}
			}
		}
	}

	return nil
}
