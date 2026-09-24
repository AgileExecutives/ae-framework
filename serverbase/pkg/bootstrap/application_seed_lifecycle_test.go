package bootstrap

import (
	"context"
	"testing"

	"github.com/AgileExecutives/ae-framework/serverbase/modules/base"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/templates"
	pkgconfig "github.com/AgileExecutives/ae-framework/serverbase/pkg/config"
	minimalorg "github.com/AgileExecutives/ae-framework/shared-modules/organization"
)

func TestStartupSeedCreatesTenantScopedTemplatesViaTenantCreatedEvent(t *testing.T) {
	ensureRepoRoot(t)
	t.Setenv("USE_IN_MEMORY_DB", "true")
	t.Setenv("ADMIN_USER", "")
	t.Setenv("ADMIN_PASSWORD", "")

	cfg := pkgconfig.Load()
	cfg.Server.Mode = "test"
	app := NewApplication(cfg)

	for _, module := range []struct {
		name string
		add  func() error
	}{
		{"base", func() error { return app.RegisterModule(base.NewBaseModule()) }},
		{"organization", func() error { return app.RegisterModule(minimalorg.NewOrganizationModule()) }},
		{"templates", func() error { return app.RegisterModule(templates.NewTemplatesModule()) }},
	} {
		if err := module.add(); err != nil {
			t.Fatalf("register %s module: %v", module.name, err)
		}
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

	if err := app.startModulesAndSeed(context.Background()); err != nil {
		t.Fatalf("start modules and seed: %v", err)
	}
	t.Cleanup(func() { _ = app.Stop(context.Background()) })

	var tenantCount int64
	if err := app.DB().Table("tenants").Count(&tenantCount).Error; err != nil {
		t.Fatalf("count tenants: %v", err)
	}
	if tenantCount != 1 {
		t.Fatalf("expected one seeded tenant, got %d", tenantCount)
	}

	for _, templateKey := range []string{"welcome", "password_reset"} {
		var count int64
		if err := app.DB().Table("templates").
			Where("tenant_id = ? AND module = ? AND template_key = ?", 1, "user", templateKey).
			Count(&count).Error; err != nil {
			t.Fatalf("count tenant template %s: %v", templateKey, err)
		}
		if count != 1 {
			t.Fatalf("expected tenant-scoped %s template after tenant.created event, got %d", templateKey, count)
		}
	}

	var contractCount int64
	if err := app.DB().Table("template_contracts").
		Where("tenant_id = ? AND module = ? AND template_key = ?", 1, "user", "welcome").
		Count(&contractCount).Error; err != nil {
		t.Fatalf("count tenant contract: %v", err)
	}
	if contractCount != 1 {
		t.Fatalf("expected tenant-scoped user/welcome contract after tenant.created event, got %d", contractCount)
	}
}
