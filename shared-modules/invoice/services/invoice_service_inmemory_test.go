package services

import (
	"context"
	"testing"
	"time"

	"github.com/AgileExecutives/ae-framwork/shared-modules/invoice/entities"
	repo "github.com/AgileExecutives/ae-framwork/shared-modules/invoice/repo"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvoice_WithInMemoryRepo(t *testing.T) {
	r := repo.NewInMemoryInvoiceRepo()
	svc := NewInvoiceServiceWithRepo(r)

	now := time.Now()
	req := &entities.CreateInvoiceRequest{
		OrganizationID: 1,
		InvoiceNumber:  "INV-100",
		InvoiceDate:    now,
		CustomerName:   "ACME Corp",
		Items:          []entities.InvoiceItemData{{Description: "Test", Quantity: 1, UnitPrice: 10, TaxRate: 0}},
	}

	inv, err := svc.CreateInvoice(context.Background(), 1, 2, req)
	assert.NoError(t, err)
	assert.Equal(t, "INV-100", inv.InvoiceNumber)
	assert.Equal(t, float64(10), inv.SubtotalAmount)
}
