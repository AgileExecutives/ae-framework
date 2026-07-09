package repo

import (
	"context"
	"time"

	"github.com/AgileExecutives/shared-modules/invoice/entities"
)

// InvoiceRepo defines persistence operations for invoices
type InvoiceRepo interface {
	CreateInvoice(ctx context.Context, inv *entities.Invoice) error
	GetInvoice(ctx context.Context, tenantID, invoiceID uint) (*entities.Invoice, error)
	ListInvoices(ctx context.Context, tenantID uint, organizationID *uint, status *entities.InvoiceStatus, page, pageSize int) ([]entities.Invoice, int64, error)
	UpdateInvoice(ctx context.Context, inv *entities.Invoice) error
	DeleteInvoice(ctx context.Context, tenantID, invoiceID uint) error
	MarkAsPaid(ctx context.Context, tenantID, invoiceID uint, paymentDate time.Time) error
	LinkDocument(ctx context.Context, tenantID, invoiceID, documentID uint) error
}
