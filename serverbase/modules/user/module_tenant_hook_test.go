package user_test

import (
	"testing"

	templateentities "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/entities"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/user"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/testutils"
)

func TestUserModuleSubscribesToTenantCreated(t *testing.T) {
	db := testutils.SetupTestDB(t)
	ctx := core.ModuleContext{
		DB:       db,
		EventBus: core.NewEventBus(),
		Services: core.NewServiceRegistry(),
		Logger:   core.NewLogger(),
	}

	m := user.NewUserModule()
	if err := m.Initialize(ctx); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	for _, handler := range m.EventHandlers() {
		if handler.EventType() == "tenant.created" {
			return
		}
	}
	t.Fatal("expected user module to subscribe to tenant.created")
}

func TestUserModuleRegistersEmbeddedContractsForTenant(t *testing.T) {
	db := testutils.SetupTestDB(t)
	if err := db.AutoMigrate(&templateentities.TemplateContract{}); err != nil {
		t.Fatalf("migrate template contracts: %v", err)
	}

	module := user.NewUserModule()
	registrar, ok := module.(interface {
		RegisterContracts(core.ModuleContext, uint) error
	})
	if !ok {
		t.Fatal("user module does not implement contract registration")
	}

	const tenantID = uint(42)
	if err := registrar.RegisterContracts(core.ModuleContext{DB: db}, tenantID); err != nil {
		t.Fatalf("register contracts: %v", err)
	}

	var contractCount int64
	if err := db.Model(&templateentities.TemplateContract{}).
		Where("tenant_id = ? AND module = ?", tenantID, "user").
		Count(&contractCount).Error; err != nil {
		t.Fatalf("count tenant contracts: %v", err)
	}
	if contractCount != 2 {
		t.Fatalf("expected two contracts for tenant %d, got %d", tenantID, contractCount)
	}
}
