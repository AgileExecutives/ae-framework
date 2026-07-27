package repo

import (
	"context"

	"gorm.io/gorm"
)

type GormDocumentRepo struct {
	db *gorm.DB
}

func NewGormDocumentRepo(db *gorm.DB) *GormDocumentRepo {
	return &GormDocumentRepo{db: db}
}

// SaveDocument is a no-op adapter placeholder; real implementations may store
// metadata or upload bytes to blob storage. Returns a synthetic ID.
func (r *GormDocumentRepo) SaveDocument(ctx context.Context, filename string, data []byte) (uint, error) {
	// Placeholder: do not persist to DB here; return id 0 and nil error
	return 0, nil
}
