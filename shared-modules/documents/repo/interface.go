package repo

import (
	"context"

	"github.com/AgileExecutives/ae-framework/shared-modules/documents/entities"
)

// DocumentRepo provides DB persistence operations for documents
type DocumentRepo interface {
	Create(ctx context.Context, doc *entities.Document) error
	GetByID(ctx context.Context, tenantID, documentID uint) (*entities.Document, error)
	List(ctx context.Context, tenantID uint, req entities.ListDocumentsRequest) ([]entities.Document, int64, error)
	Delete(ctx context.Context, tenantID, documentID uint) error
	ListByReference(ctx context.Context, tenantID uint, referenceType string, referenceID uint) ([]entities.Document, error)
}
