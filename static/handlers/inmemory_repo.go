package handlers
package handlers

import "context"

type InMemoryStaticRepo struct{ files map[string][]byte }
func NewInMemoryStaticRepo(files map[string][]byte) *InMemoryStaticRepo { return &InMemoryStaticRepo{files: files} }
func (r *InMemoryStaticRepo) ListFiles(ctx context.Context) ([]string, error) { var out []string; for k := range r.files { out = append(out, k) }; return out, nil }
func (r *InMemoryStaticRepo) ReadFile(ctx context.Context, name string) ([]byte, error) { if d, ok := r.files[name]; ok { return d, nil }; return nil, nil }
