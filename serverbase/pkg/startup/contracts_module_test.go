package startup

import (
    "testing"

    "github.com/AgileExecutives/ae-framework/serverbase/module"
    templateServices "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/services"
    "github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func TestRegisterAllContractsCallsModuleRegisterContracts(t *testing.T) {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("open db: %v", err)
    }
    if err := db.Exec(`CREATE TABLE tenants (id integer primary key, deleted_at datetime)`).Error; err != nil {
        t.Fatalf("create tenants table: %v", err)
    }
    if err := db.Exec(`CREATE TABLE template_contracts (id integer primary key, tenant_id integer, module text, template_key text, variable_schema json, default_sample_data json, created_at datetime, updated_at datetime)`).Error; err != nil {
        t.Fatalf("create template_contracts table: %v", err)
    }
    if err := db.Exec(`INSERT INTO tenants (id) VALUES (1), (2)`).Error; err != nil {
        t.Fatalf("insert tenants: %v", err)
    }

    // Create a registry and register an AdapterModule that uses WithContractRegistration
    reg := core.NewModuleRegistry()
    mod := module.NewAdapterModule("fake", "0.1.0", []string{}, module.WithContractRegistration(func(ctx core.ModuleContext, tenantID uint) error {
        r := templateServices.NewContractRegistrar(ctx.DB)
        b := []byte(`{"type":"object","properties":{"X":{"type":"string"}},"template_key":"from_module_test"}`)
        return r.RegisterContractFromBytes(tenantID, "fake", "from_module_test.json", b)
    }))
    if err := reg.Register(mod); err != nil {
        t.Fatalf("register module: %v", err)
    }

    ctx := core.ModuleContext{DB: db, ModuleRegistry: reg}
    if err := RegisterAllContracts(ctx); err != nil {
        t.Fatalf("RegisterAllContracts: %v", err)
    }

    var cnt int64
    if err := db.Table("template_contracts").Count(&cnt).Error; err != nil {
        t.Fatalf("count template_contracts: %v", err)
    }
    if cnt != 2 { // two tenants each should have one contract
        t.Fatalf("expected 2 template_contracts, got %d", cnt)
    }
}
