package services

import (
	"context"
	"testing"

	"github.com/AgileExecutives/ae-framwork/shared-modules/documents/entities"
	repo "github.com/AgileExecutives/ae-framwork/shared-modules/documents/repo"
	"github.com/AgileExecutives/ae-framwork/shared-modules/documents/services/storage"
	"github.com/stretchr/testify/assert"
)

func TestStoreDocument_WithInMemoryRepo(t *testing.T) {
	r := repo.NewInMemoryDocumentRepo()
	minio := storage.NewInMemoryStorage()
	svc := NewDocumentServiceWithRepo(r, minio)

	req := entities.StoreDocumentRequest{
		DocumentType: "invoice",
		FileName:     "test.pdf",
		Bucket:       "docs",
		Path:         "test.pdf",
		Content:      []byte("hello"),
	}

	doc, err := svc.StoreDocument(context.Background(), 1, 2, req)
	assert.NoError(t, err)
	assert.Equal(t, "test.pdf", doc.FileName)
}
