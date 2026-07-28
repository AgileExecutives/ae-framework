package services

import (
	repo "github.com/AgileExecutives/ae-framwork/shared-modules/pdf/repo"
)

type PDFGenerator struct {
	repo repo.DocumentRepo
}

func NewPDFGenerator() *PDFGenerator { return &PDFGenerator{} }

func NewPDFGeneratorWithRepo(r repo.DocumentRepo) *PDFGenerator { return &PDFGenerator{repo: r} }

// GeneratePDF is a lightweight stub that simulates PDF generation for tests.
// It returns the generated filename (fileName + ".pdf") and no error.
func (p *PDFGenerator) GeneratePDF(data map[string]interface{}, templateName, fileName string) (string, error) {
	// In production this would render a template and produce a PDF file.
	// If a repo is provided we could persist bytes; for now this is a no-op.
	return fileName + ".pdf", nil
}
