package repo

import "context"

// DocumentRepo defines an abstraction for storing generated documents
type DocumentRepo interface {
	// SaveDocument persists the document bytes and returns an ID
	SaveDocument(ctx context.Context, filename string, data []byte) (uint, error)
}
