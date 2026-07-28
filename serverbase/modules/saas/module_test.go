package saas

import "testing"

func TestNewSaaSModuleExposesMigrationEntities(t *testing.T) {
	mod := NewSaaSModule()
	entities := mod.Entities()

	if len(entities) != 2 {
		t.Fatalf("expected 2 migration entities, got %d", len(entities))
	}

	tables := make(map[string]bool, len(entities))
	for _, entity := range entities {
		tables[entity.TableName()] = true
	}

	for _, table := range []string{"plans", "customers"} {
		if !tables[table] {
			t.Fatalf("expected %s to be registered for migration", table)
		}
	}

	if tables["newsletters"] {
		t.Fatal("did not expect newsletters to be registered by the saas adapter")
	}
}
