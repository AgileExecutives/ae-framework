package services

import (
	"context"
	"testing"
	"time"

	"github.com/AgileExecutives/ae-framework/shared-modules/invoice/entities"
	repo "github.com/AgileExecutives/ae-framework/shared-modules/invoice/repo"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceService_CreateInvoice_TableDriven(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		items        []entities.InvoiceItemData
		currency     string
		wantSubtotal float64
		wantTax      float64
		wantTotal    float64
	}{
		{
			name:         "single item no tax",
			items:        []entities.InvoiceItemData{{Description: "Test", Quantity: 1, UnitPrice: 10, TaxRate: 0}},
			currency:     "",
			wantSubtotal: 10,
			wantTax:      0,
			wantTotal:    10,
		},
		{
			name: "multiple items with tax",
			items: []entities.InvoiceItemData{
				{Description: "A", Quantity: 2, UnitPrice: 5, TaxRate: 10}, // amount 10, tax 1
				{Description: "B", Quantity: 3, UnitPrice: 2, TaxRate: 5},  // amount 6, tax 0.3
			},
			currency:     "USD",
			wantSubtotal: 16,
			wantTax:      1.3,
			wantTotal:    17.3,
		},
		{
			name:         "no items",
			items:        []entities.InvoiceItemData{},
			currency:     "EUR",
			wantSubtotal: 0,
			wantTax:      0,
			wantTotal:    0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := repo.NewInMemoryInvoiceRepo()
			svc := NewInvoiceServiceWithRepo(r)

			req := &entities.CreateInvoiceRequest{
				OrganizationID: 1,
				InvoiceNumber:  "INV-001",
				InvoiceDate:    now,
				CustomerName:   "ACME",
				Items:          tc.items,
				Currency:       tc.currency,
			}

			inv, err := svc.CreateInvoice(context.Background(), 1, 2, req)
			assert.NoError(t, err)
			// amounts may be floats; allow small delta
			assert.InDelta(t, tc.wantSubtotal, inv.SubtotalAmount, 0.0001)
			assert.InDelta(t, tc.wantTax, inv.TaxAmount, 0.0001)
			assert.InDelta(t, tc.wantTotal, inv.TotalAmount, 0.0001)

			// default currency when empty should be EUR
			if tc.currency == "" {
				assert.Equal(t, "EUR", inv.Currency)
			} else {
				assert.Equal(t, tc.currency, inv.Currency)
			}
		})
	}
}
