package handlers

import (
	"context"
	"os"
	"path/filepath"
)

type FSStaticRepo struct{ basePath string }

func NewFSStaticRepo(basePath string) *FSStaticRepo { return &FSStaticRepo{basePath: basePath} }
func (r *FSStaticRepo) ListFiles(ctx context.Context) ([]string, error) {
	var out []string
	entries, err := os.ReadDir(r.basePath)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out, nil
}
func (r *FSStaticRepo) ReadFile(ctx context.Context, name string) ([]byte, error) {
	full := filepath.Join(r.basePath, name)
	return os.ReadFile(full)
}
