package repo

import (
	"context"

	"github.com/AgileExecutives/ae-framwork/shared-modules/invoice_number/entities"
	"gorm.io/gorm"
)

type GormInvoiceNumberRepo struct{ db *gorm.DB }

func NewGormInvoiceNumberRepo(db *gorm.DB) *GormInvoiceNumberRepo {
	return &GormInvoiceNumberRepo{db: db}
}

func (r *GormInvoiceNumberRepo) Find(ctx context.Context, tenantID, organizationID uint, year, month int) (*entities.InvoiceNumber, error) {
	var rec entities.InvoiceNumber
	err := r.db.Where("tenant_id = ? AND organization_id = ? AND year = ? AND month = ?", tenantID, organizationID, year, month).First(&rec).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *GormInvoiceNumberRepo) Create(ctx context.Context, rec *entities.InvoiceNumber) error {
	return r.db.Create(rec).Error
}

func (r *GormInvoiceNumberRepo) Update(ctx context.Context, rec *entities.InvoiceNumber) error {
	return r.db.Save(rec).Error
}

func (r *GormInvoiceNumberRepo) CreateLog(ctx context.Context, log *entities.InvoiceNumberLog) error {
	return r.db.Create(log).Error
}

func (r *GormInvoiceNumberRepo) GetLogs(ctx context.Context, tenantID, organizationID uint, year, month int, page, pageSize int) ([]entities.InvoiceNumberLog, int64, error) {
	var logs []entities.InvoiceNumberLog
	var total int64
	query := r.db.Model(&entities.InvoiceNumberLog{}).Where("tenant_id = ?", tenantID)
	if organizationID > 0 {
		query = query.Where("organization_id = ?", organizationID)
	}
	if year > 0 {
		query = query.Where("year = ?", year)
	}
	if month > 0 {
		query = query.Where("month = ?", month)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := query.Order("generated_at DESC").Limit(pageSize).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
