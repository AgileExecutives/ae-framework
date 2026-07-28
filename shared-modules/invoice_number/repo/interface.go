package repo

import (
	"context"

	"github.com/AgileExecutives/ae-framework/shared-modules/invoice_number/entities"
)

// InvoiceNumberRepo defines persistence operations for invoice numbers
type InvoiceNumberRepo interface {
	Find(ctx context.Context, tenantID, organizationID uint, year, month int) (*entities.InvoiceNumber, error)
	Create(ctx context.Context, rec *entities.InvoiceNumber) error
	Update(ctx context.Context, rec *entities.InvoiceNumber) error
	CreateLog(ctx context.Context, log *entities.InvoiceNumberLog) error
	GetLogs(ctx context.Context, tenantID, organizationID uint, year, month int, page, pageSize int) ([]entities.InvoiceNumberLog, int64, error)
}
