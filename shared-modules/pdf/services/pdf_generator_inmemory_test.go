package services

import (
	"testing"

	repo "github.com/AgileExecutives/ae-framework/shared-modules/pdf/repo"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePDF_InMemoryRepo(t *testing.T) {
	r := repo.NewInMemoryDocumentRepo()
	svc := NewPDFGeneratorWithRepo(r)

	filename, err := svc.GeneratePDF(map[string]interface{}{"x": 1}, "tmpl", "testfile")
	assert.NoError(t, err)
	assert.Equal(t, "testfile.pdf", filename)
}
