package storage

import (
	"context"
	"errors"
	"sync"
	"time"
)

type InMemoryStorage struct {
	mu    sync.Mutex
	store map[string][]byte // key = bucket+"/"+key
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{store: make(map[string][]byte)}
}

func (s *InMemoryStorage) Store(ctx context.Context, req StoreRequest) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := req.Key
	s.store[req.Bucket+"/"+key] = req.Data
	return key, nil
}

func (s *InMemoryStorage) Retrieve(ctx context.Context, bucket, key string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.store[bucket+"/"+key]
	if !ok {
		return nil, errors.New("not found")
	}
	return data, nil
}

func (s *InMemoryStorage) GetURL(ctx context.Context, bucket, key string, expiresIn time.Duration) (string, error) {
	// Return a fake URL for tests
	return "https://inmemory.local/" + bucket + "/" + key, nil
}

func (s *InMemoryStorage) Delete(ctx context.Context, bucket, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, bucket+"/"+key)
	return nil
}

func (s *InMemoryStorage) List(ctx context.Context, bucket, prefix string) ([]DocumentMeta, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []DocumentMeta
	for k := range s.store {
		// k format: bucket/key
		// simple prefix check
		if len(k) >= len(bucket)+1 && k[:len(bucket)+1] == bucket+"/" {
			out = append(out, DocumentMeta{Key: k[len(bucket)+1:], Size: int64(len(s.store[k])), LastModified: time.Now()})
		}
	}
	return out, nil
}

func (s *InMemoryStorage) Exists(ctx context.Context, bucket, key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.store[bucket+"/"+key]
	return ok, nil
}
