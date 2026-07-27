package storage

import (
	"context"
	"time"
)

// Lightweight local storage types used during migration.
type MinIOStorage struct{}

type StoreRequest struct {
	Bucket      string
	Key         string
	Data        []byte
	ContentType string
	Metadata    map[string]string
}

type DocumentMeta struct {
	Key string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Region          string
}

func NewMinIOStorage(cfg MinIOConfig) (*MinIOStorage, error) {
	return &MinIOStorage{}, nil
}

func (m *MinIOStorage) Store(ctx context.Context, req StoreRequest) (DocumentMeta, error) {
	return DocumentMeta{Key: req.Key}, nil
}

func (m *MinIOStorage) Retrieve(ctx context.Context, bucket, key string) ([]byte, error) {
	return []byte{}, nil
}

func (m *MinIOStorage) Delete(ctx context.Context, bucket, key string) error {
	return nil
}

func (m *MinIOStorage) GetURL(ctx context.Context, bucket, key string, expiresIn time.Duration) (string, error) {
	return "", nil
}

func (m *MinIOStorage) Exists(ctx context.Context, bucket, key string) (bool, error) {
	return false, nil
}
