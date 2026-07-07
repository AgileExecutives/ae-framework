package events
package events

import "github.com/ae/base-server/pkg/core"

type PDFGeneratedHandler struct{ logger core.Logger }
func NewPDFGeneratedHandler(logger core.Logger) *PDFGeneratedHandler { return &PDFGeneratedHandler{logger: logger} }
func (h *PDFGeneratedHandler) EventType() string { return "pdf.generated" }
func (h *PDFGeneratedHandler) Handle(event interface{}) error { return nil }
func (h *PDFGeneratedHandler) Priority() int { return 100 }

type PDFFailedHandler struct{ logger core.Logger }
func NewPDFFailedHandler(logger core.Logger) *PDFFailedHandler { return &PDFFailedHandler{logger: logger} }
func (h *PDFFailedHandler) EventType() string { return "pdf.failed" }
func (h *PDFFailedHandler) Handle(event interface{}) error { return nil }
func (h *PDFFailedHandler) Priority() int { return 100 }
