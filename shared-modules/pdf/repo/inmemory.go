package repo

import (
	"context"
	"sync"

	"errors"
)

type InMemoryDocumentRepo struct {
	mu    sync.Mutex
	id    uint
	store map[uint][]byte
}

func NewInMemoryDocumentRepo() *InMemoryDocumentRepo {
	return &InMemoryDocumentRepo{store: make(map[uint][]byte)}
}

func (r *InMemoryDocumentRepo) SaveDocument(ctx context.Context, filename string, data []byte) (uint, error) {
	if data == nil {
		return 0, errors.New("no data")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.id++
	r.store[r.id] = data
	return r.id, nil
}
