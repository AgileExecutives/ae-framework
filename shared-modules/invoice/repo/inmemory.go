package repo

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/AgileExecutives/ae-framwork/shared-modules/invoice/entities"
)

type InMemoryInvoiceRepo struct {
	mu       sync.Mutex
	id       uint
	invoices map[uint]*entities.Invoice
}

func NewInMemoryInvoiceRepo() *InMemoryInvoiceRepo {
	return &InMemoryInvoiceRepo{invoices: make(map[uint]*entities.Invoice)}
}

func (r *InMemoryInvoiceRepo) CreateInvoice(ctx context.Context, inv *entities.Invoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.id++
	inv.ID = r.id
	now := time.Now()
	inv.CreatedAt = now
	inv.UpdatedAt = now
	for i := range inv.Items {
		inv.Items[i].InvoiceID = inv.ID
	}
	r.invoices[inv.ID] = inv
	return nil
}

func (r *InMemoryInvoiceRepo) GetInvoice(ctx context.Context, tenantID, invoiceID uint) (*entities.Invoice, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.invoices[invoiceID]
	if !ok || inv.TenantID != tenantID {
		return nil, errors.New("not found")
	}
	return inv, nil
}

func (r *InMemoryInvoiceRepo) ListInvoices(ctx context.Context, tenantID uint, organizationID *uint, status *entities.InvoiceStatus, page, pageSize int) ([]entities.Invoice, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.Invoice
	for _, inv := range r.invoices {
		if inv.TenantID != tenantID {
			continue
		}
		if organizationID != nil && inv.OrganizationID != *organizationID {
			continue
		}
		if status != nil && inv.Status != *status {
			continue
		}
		out = append(out, *inv)
	}
	total := int64(len(out))
	start := (page - 1) * pageSize
	if start >= len(out) {
		return []entities.Invoice{}, total, nil
	}
	end := start + pageSize
	if end > len(out) {
		end = len(out)
	}
	return out[start:end], total, nil
}

func (r *InMemoryInvoiceRepo) UpdateInvoice(ctx context.Context, inv *entities.Invoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ex, ok := r.invoices[inv.ID]
	if !ok || ex.TenantID != inv.TenantID {
		return errors.New("not found")
	}
	inv.UpdatedAt = time.Now()
	r.invoices[inv.ID] = inv
	return nil
}

func (r *InMemoryInvoiceRepo) DeleteInvoice(ctx context.Context, tenantID, invoiceID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ex, ok := r.invoices[invoiceID]
	if !ok || ex.TenantID != tenantID {
		return errors.New("not found")
	}
	delete(r.invoices, invoiceID)
	return nil
}

func (r *InMemoryInvoiceRepo) MarkAsPaid(ctx context.Context, tenantID, invoiceID uint, paymentDate time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ex, ok := r.invoices[invoiceID]
	if !ok || ex.TenantID != tenantID {
		return errors.New("not found")
	}
	ex.Status = entities.InvoiceStatusPaid
	ex.PaymentDate = &paymentDate
	ex.UpdatedAt = time.Now()
	return nil
}

func (r *InMemoryInvoiceRepo) LinkDocument(ctx context.Context, tenantID, invoiceID, documentID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ex, ok := r.invoices[invoiceID]
	if !ok || ex.TenantID != tenantID {
		return errors.New("not found")
	}
	ex.DocumentID = &documentID
	ex.UpdatedAt = time.Now()
	return nil
}
