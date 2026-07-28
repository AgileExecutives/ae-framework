package services

import (
	"context"
	"testing"

	"github.com/AgileExecutives/ae-framework/serverbase/modules/customers/repo"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framwork/shared-modules/saas-base/models"
)

type testLogger struct{}

func (l *testLogger) Debug(args ...interface{})                      {}
func (l *testLogger) Info(args ...interface{})                       {}
func (l *testLogger) Warn(args ...interface{})                       {}
func (l *testLogger) Error(args ...interface{})                      {}
func (l *testLogger) Fatal(args ...interface{})                      {}
func (l *testLogger) With(key string, value interface{}) core.Logger { return l }

func TestGetByTenant_ReturnsCustomers(t *testing.T) {
	r := repo.NewInMemoryCustomerRepo()
	// create two customers for tenant 42
	c1 := &models.Customer{Name: "Alpha", Email: "alpha@example.com", PlanID: 1, TenantID: 42}
	c2 := &models.Customer{Name: "Beta", Email: "beta@example.com", PlanID: 1, TenantID: 42}
	if err := r.Save(context.Background(), c1); err != nil {
		t.Fatalf("save c1: %v", err)
	}
	if err := r.Save(context.Background(), c2); err != nil {
		t.Fatalf("save c2: %v", err)
	}

	svc := NewCustomerService(r, &testLogger{})
	res, err := svc.GetByTenant(42)
	if err != nil {
		t.Fatalf("GetByTenant err: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 customers, got %d", len(res))
	}
}

func TestGetByTenant_IsolatedByTenantID(t *testing.T) {
	r := repo.NewInMemoryCustomerRepo()
	// create customers for tenant 1 and 2
	c1 := &models.Customer{Name: "OneA", Email: "onea@example.com", PlanID: 1, TenantID: 1}
	c2 := &models.Customer{Name: "TwoA", Email: "twoa@example.com", PlanID: 1, TenantID: 2}
	if err := r.Save(context.Background(), c1); err != nil {
		t.Fatalf("save c1: %v", err)
	}
	if err := r.Save(context.Background(), c2); err != nil {
		t.Fatalf("save c2: %v", err)
	}

	svc := NewCustomerService(r, &testLogger{})

	res1, err := svc.GetByTenant(1)
	if err != nil {
		t.Fatalf("GetByTenant err: %v", err)
	}
	if len(res1) != 1 || res1[0].TenantID != 1 {
		t.Fatalf("expected tenant 1 to have its customer, got %+v", res1)
	}

	res2, err := svc.GetByTenant(2)
	if err != nil {
		t.Fatalf("GetByTenant err: %v", err)
	}
	if len(res2) != 1 || res2[0].TenantID != 2 {
		t.Fatalf("expected tenant 2 to have its customer, got %+v", res2)
	}
}
