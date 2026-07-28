package repo

import (
	"context"
	"time"

	"github.com/AgileExecutives/ae-framework/shared-modules/audit/entities"
)

// AuditRepo defines persistence operations for audit logs
type AuditRepo interface {
	CreateLog(ctx context.Context, log *entities.AuditLog) error
	GetLogs(ctx context.Context, filter entities.AuditLogFilter) ([]entities.AuditLog, int64, error)
	GetLogsByEntity(ctx context.Context, tenantID, entityID uint, entityType entities.EntityType) ([]entities.AuditLog, error)
	GetLogsForExport(ctx context.Context, filter entities.AuditLogFilter) ([]entities.AuditLog, error)
	GetStatistics(ctx context.Context, tenantID uint, startDate, endDate *time.Time) (map[string]interface{}, error)
}
