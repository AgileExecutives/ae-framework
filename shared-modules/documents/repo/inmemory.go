package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/AgileExecutives/ae-framework/shared-modules/documents/entities"
)

type InMemoryDocumentRepo struct {
	mu   sync.Mutex
	id   uint
	docs map[uint]*entities.Document
}

func NewInMemoryDocumentRepo() *InMemoryDocumentRepo {
	return &InMemoryDocumentRepo{docs: make(map[uint]*entities.Document)}
}

func (r *InMemoryDocumentRepo) Create(ctx context.Context, doc *entities.Document) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.id++
	doc.ID = r.id
	r.docs[doc.ID] = doc
	return nil
}

func (r *InMemoryDocumentRepo) GetByID(ctx context.Context, tenantID, documentID uint) (*entities.Document, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, ok := r.docs[documentID]
	if !ok || d.TenantID != tenantID {
		return nil, errors.New("not found")
	}
	return d, nil
}

func (r *InMemoryDocumentRepo) List(ctx context.Context, tenantID uint, req entities.ListDocumentsRequest) ([]entities.Document, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.Document
	for _, d := range r.docs {
		if d.TenantID != tenantID {
			continue
		}
		if req.DocumentType != "" && d.DocumentType != req.DocumentType {
			continue
		}
		if req.ReferenceType != "" && d.ReferenceType != req.ReferenceType {
			continue
		}
		if req.ReferenceID > 0 && d.ReferenceID != nil && *d.ReferenceID != req.ReferenceID {
			continue
		}
		out = append(out, *d)
	}
	total := int64(len(out))
	// Simple pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	start := (req.Page - 1) * req.Limit
	if start >= len(out) {
		return []entities.Document{}, total, nil
	}
	end := start + req.Limit
	if end > len(out) {
		end = len(out)
	}
	return out[start:end], total, nil
}

func (r *InMemoryDocumentRepo) Delete(ctx context.Context, tenantID, documentID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, ok := r.docs[documentID]
	if !ok || d.TenantID != tenantID {
		return errors.New("not found")
	}
	delete(r.docs, documentID)
	return nil
}

func (r *InMemoryDocumentRepo) ListByReference(ctx context.Context, tenantID uint, referenceType string, referenceID uint) ([]entities.Document, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.Document
	for _, d := range r.docs {
		if d.TenantID != tenantID {
			continue
		}
		if d.ReferenceType == referenceType && d.ReferenceID != nil && *d.ReferenceID == referenceID {
			out = append(out, *d)
		}
	}
	return out, nil
}
