package storage

import (
	docsstorage "github.com/AgileExecutives/shared-modules/documents/services/storage"
)

// Re-export types and helpers from documents storage to satisfy other modules
type MinIOStorage = docsstorage.MinIOStorage
type StoreRequest = docsstorage.StoreRequest
type DocumentMeta = docsstorage.DocumentMeta

type MinIOConfig = docsstorage.MinIOConfig

// NewMinIOStorage delegates to the documents storage implementation
func NewMinIOStorage(cfg MinIOConfig) (*MinIOStorage, error) {
	return docsstorage.NewMinIOStorage(cfg)
}
