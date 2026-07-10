package services

import (
	"testing"

	"github.com/AgileExecutives/serverbase/pkg/settings/repository"
)

func Test_SetAndGetSetting_InMemory(t *testing.T) {
	repo := repository.NewInMemorySettingsRepository()
	svc := NewSettingsService(repo)

	tenantID := uint(1)
	domain := "billing"
	key := "tax"
	value := map[string]interface{}{"rate": 19}

	if err := svc.SetSetting(tenantID, "", domain, key, value, "json"); err != nil {
		t.Fatalf("SetSetting failed: %v", err)
	}

	v, err := svc.GetSetting(tenantID, "", domain, key)
	if err != nil {
		t.Fatalf("GetSetting error: %v", err)
	}
	if v == nil {
		t.Fatalf("expected non-nil value")
	}
}
