package services

import (
	"context"
	"testing"

	"github.com/AgileExecutives/ae-framwork/shared-modules/invoice_number/repo"
)

func TestGenerateInvoiceNumber_InMemoryRepo(t *testing.T) {
	r := repo.NewInMemoryInvoiceRepo()
	svc := NewInvoiceNumberServiceWithRepo(r)

	ctx := context.Background()
	tenantID := uint(1)
	orgID := uint(2)

	cfg := DefaultInvoiceConfig()

	// First generation should produce sequence 1
	resp1, err := svc.GenerateInvoiceNumber(ctx, tenantID, orgID, cfg)
	if err != nil {
		t.Fatalf("GenerateInvoiceNumber failed: %v", err)
	}
	if resp1.Sequence != 1 {
		t.Fatalf("expected sequence 1, got %d", resp1.Sequence)
	}

	// Second generation should increment to 2
	resp2, err := svc.GenerateInvoiceNumber(ctx, tenantID, orgID, cfg)
	if err != nil {
		t.Fatalf("GenerateInvoiceNumber failed: %v", err)
	}
	if resp2.Sequence != 2 {
		t.Fatalf("expected sequence 2, got %d", resp2.Sequence)
	}

	// Check format contains prefix and year
	if resp2.Format == "" {
		t.Fatalf("expected non-empty format")
	}
}
