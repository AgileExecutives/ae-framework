package pdf

import (
	"context"

	"github.com/AgileExecutives/shared-modules/pdf/events"
	"github.com/AgileExecutives/shared-modules/pdf/handlers"
	"github.com/AgileExecutives/shared-modules/pdf/services"
	"github.com/AgileExecutives/serverbase/pkg/core"
)

type PDFModule struct {
	pdfHandler    *handlers.PDFHandler
	pdfService    *services.PDFGenerator
	eventHandlers []core.EventHandler
}

func NewPDFModule() *PDFModule              { return &PDFModule{} }
func (m *PDFModule) Name() string           { return "pdf" }
func (m *PDFModule) Version() string        { return "1.0.0" }
func (m *PDFModule) Description() string    { return "PDF generation and document management system" }
func (m *PDFModule) Dependencies() []string { return []string{} }
func (m *PDFModule) Initialize(ctx core.ModuleContext) error {
	ctx.Logger.Info("Initializing PDF module...")
	m.pdfService = services.NewPDFGenerator()
	m.pdfHandler = handlers.NewPDFHandler(m.pdfService, ctx.DB)
	m.eventHandlers = []core.EventHandler{events.NewPDFGeneratedHandler(ctx.Logger), events.NewPDFFailedHandler(ctx.Logger)}
	ctx.Logger.Info("PDF module initialized successfully")
	return nil
}
func (m *PDFModule) Start(ctx context.Context) error    { return nil }
func (m *PDFModule) Stop(ctx context.Context) error     { return nil }
func (m *PDFModule) Entities() []core.Entity            { return []core.Entity{} }
func (m *PDFModule) Routes() []core.RouteProvider       { return []core.RouteProvider{m.pdfHandler} }
func (m *PDFModule) EventHandlers() []core.EventHandler { return m.eventHandlers }
func (m *PDFModule) Services() []core.ServiceProvider {
	return []core.ServiceProvider{&PDFServiceProvider{pdfService: m.pdfService}}
}
func (m *PDFModule) Middleware() []core.MiddlewareProvider { return []core.MiddlewareProvider{} }
func (m *PDFModule) SwaggerPaths() []string                { return []string{} }

type PDFServiceProvider struct{ pdfService *services.PDFGenerator }

func (p *PDFServiceProvider) ServiceName() string           { return "pdf-generator" }
func (p *PDFServiceProvider) ServiceInterface() interface{} { return p.pdfService }
func (p *PDFServiceProvider) Factory(ctx core.ModuleContext) (interface{}, error) {
	return p.pdfService, nil
}
