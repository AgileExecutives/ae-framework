package services

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator { return &PDFGenerator{} }

// GeneratePDF is a lightweight stub that simulates PDF generation for tests.
// It returns the generated filename (fileName + ".pdf") and no error.
func (p *PDFGenerator) GeneratePDF(data map[string]interface{}, templateName, fileName string) (string, error) {
	// In production this would render a template and produce a PDF file.
	return fileName + ".pdf", nil
}
