package repo

import (
	"context"
	"time"

	"github.com/AgileExecutives/ae-framework/shared-modules/invoice/entities"
	"gorm.io/gorm"
)

type GormInvoiceRepo struct {
	db *gorm.DB
}

func NewGormInvoiceRepo(db *gorm.DB) *GormInvoiceRepo {
	return &GormInvoiceRepo{db: db}
}

func (r *GormInvoiceRepo) CreateInvoice(ctx context.Context, inv *entities.Invoice) error {
	return r.db.WithContext(ctx).Create(inv).Error
}

func (r *GormInvoiceRepo) GetInvoice(ctx context.Context, tenantID, invoiceID uint) (*entities.Invoice, error) {
	var inv entities.Invoice
	if err := r.db.WithContext(ctx).Preload("Items").Where("id = ? AND tenant_id = ?", invoiceID, tenantID).First(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *GormInvoiceRepo) ListInvoices(ctx context.Context, tenantID uint, organizationID *uint, status *entities.InvoiceStatus, page, pageSize int) ([]entities.Invoice, int64, error) {
	var invoices []entities.Invoice
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Invoice{}).Where("tenant_id = ?", tenantID)
	if organizationID != nil {
		query = query.Where("organization_id = ?", *organizationID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Preload("Items").Order("invoice_date DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

func (r *GormInvoiceRepo) UpdateInvoice(ctx context.Context, inv *entities.Invoice) error {
	return r.db.WithContext(ctx).Save(inv).Error
}

func (r *GormInvoiceRepo) DeleteInvoice(ctx context.Context, tenantID, invoiceID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", invoiceID, tenantID).Delete(&entities.Invoice{}).Error
}

func (r *GormInvoiceRepo) MarkAsPaid(ctx context.Context, tenantID, invoiceID uint, paymentDate time.Time) error {
	return r.db.WithContext(ctx).Model(&entities.Invoice{}).Where("id = ? AND tenant_id = ?", invoiceID, tenantID).Updates(map[string]interface{}{
		"status":       "paid",
		"payment_date": paymentDate,
	}).Error
}

func (r *GormInvoiceRepo) LinkDocument(ctx context.Context, tenantID, invoiceID, documentID uint) error {
	return r.db.WithContext(ctx).Model(&entities.Invoice{}).Where("id = ? AND tenant_id = ?", invoiceID, tenantID).Update("document_id", documentID).Error
}
