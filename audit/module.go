package audit

import (
	"github.com/AgileExecutives/shared-modules/audit/entities"
	"github.com/AgileExecutives/shared-modules/audit/handlers"
	repo "github.com/AgileExecutives/shared-modules/audit/repo"
	"github.com/AgileExecutives/shared-modules/audit/routes"
	"github.com/AgileExecutives/shared-modules/audit/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	db            *gorm.DB
	service       *services.AuditService
	handler       *handlers.AuditHandler
	routeProvider *routes.RouteProvider
}

func NewModule(db *gorm.DB) *Module {
	gormRepo := repo.NewGormAuditRepo(db)
	service := services.NewAuditServiceWithRepo(gormRepo)
	handler := handlers.NewAuditHandler(service)
	routeProvider := routes.NewRouteProvider(handler)

	return &Module{
		db:            db,
		service:       service,
		handler:       handler,
		routeProvider: routeProvider,
	}
}

func (m *Module) GetService() *services.AuditService {
	return m.service
}

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.routeProvider.RegisterRoutes(router)
}

func (m *Module) AutoMigrate() error {
	return m.db.AutoMigrate(&entities.AuditLog{})
}

func (m *Module) GetName() string {
	return "audit"
}
