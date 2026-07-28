package repo

import (
	"context"

	"github.com/AgileExecutives/ae-framwork/shared-modules/documents/entities"
	"gorm.io/gorm"
)

type GormDocumentRepo struct{ db *gorm.DB }

func NewGormDocumentRepo(db *gorm.DB) *GormDocumentRepo { return &GormDocumentRepo{db: db} }

func (r *GormDocumentRepo) Create(ctx context.Context, doc *entities.Document) error {
	return r.db.Create(doc).Error
}

func (r *GormDocumentRepo) GetByID(ctx context.Context, tenantID, documentID uint) (*entities.Document, error) {
	var d entities.Document
	if err := r.db.Where("id = ? AND tenant_id = ?", documentID, tenantID).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *GormDocumentRepo) List(ctx context.Context, tenantID uint, req entities.ListDocumentsRequest) ([]entities.Document, int64, error) {
	var docs []entities.Document
	var total int64
	q := r.db.Model(&entities.Document{}).Where("tenant_id = ?", tenantID)
	if req.DocumentType != "" {
		q = q.Where("document_type = ?", req.DocumentType)
	}
	if req.ReferenceType != "" {
		q = q.Where("reference_type = ?", req.ReferenceType)
	}
	if req.ReferenceID > 0 {
		q = q.Where("reference_id = ?", req.ReferenceID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	offset := (req.Page - 1) * req.Limit
	if err := q.Order("created_at DESC").Limit(req.Limit).Offset(offset).Find(&docs).Error; err != nil {
		return nil, 0, err
	}
	return docs, total, nil
}

func (r *GormDocumentRepo) Delete(ctx context.Context, tenantID, documentID uint) error {
	return r.db.Where("id = ? AND tenant_id = ?", documentID, tenantID).Delete(&entities.Document{}).Error
}

func (r *GormDocumentRepo) ListByReference(ctx context.Context, tenantID uint, referenceType string, referenceID uint) ([]entities.Document, error) {
	var docs []entities.Document
	if err := r.db.Where("tenant_id = ? AND reference_type = ? AND reference_id = ?", tenantID, referenceType, referenceID).Order("created_at DESC").Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}
