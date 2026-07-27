package handlers

import "context"

type StaticRepo interface {
	ListFiles(ctx context.Context) ([]string, error)
	ReadFile(ctx context.Context, name string) ([]byte, error)
}
