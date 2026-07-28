package audit

import (
	"github.com/AgileExecutives/ae-framwork/shared-modules/audit/entities"
	"github.com/AgileExecutives/ae-framwork/shared-modules/audit/handlers"
	"github.com/AgileExecutives/ae-framwork/shared-modules/audit/routes"
	"github.com/AgileExecutives/ae-framwork/shared-modules/audit/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewModule remains for direct construction when embedding the module in other code.
type Module struct {
	db            *gorm.DB
	service       *services.AuditService
	handler       *handlers.AuditHandler
	routeProvider *routes.RouteProvider
}

// Prefer registering the module via `NewCoreModule()` adapters used by the
// bootstrap system; direct constructors may still be used for embedding.

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
