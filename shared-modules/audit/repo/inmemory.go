package repo

import (
	"context"
	"sync"
	"time"

	"github.com/AgileExecutives/ae-framework/shared-modules/audit/entities"
)

type InMemoryAuditRepo struct {
	mu   sync.Mutex
	id   uint
	logs map[uint]*entities.AuditLog
}

func NewInMemoryAuditRepo() *InMemoryAuditRepo {
	return &InMemoryAuditRepo{logs: make(map[uint]*entities.AuditLog)}
}

func (r *InMemoryAuditRepo) CreateLog(ctx context.Context, log *entities.AuditLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.id++
	log.ID = r.id
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	r.logs[log.ID] = log
	return nil
}

func (r *InMemoryAuditRepo) GetLogs(ctx context.Context, filter entities.AuditLogFilter) ([]entities.AuditLog, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.AuditLog
	for _, l := range r.logs {
		if l.TenantID != filter.TenantID {
			continue
		}
		if filter.UserID != nil && l.UserID != *filter.UserID {
			continue
		}
		if filter.EntityType != nil && l.EntityType != *filter.EntityType {
			continue
		}
		if filter.EntityID != nil && l.EntityID != *filter.EntityID {
			continue
		}
		if filter.Action != nil && l.Action != *filter.Action {
			continue
		}
		if filter.StartDate != nil && l.CreatedAt.Before(*filter.StartDate) {
			continue
		}
		if filter.EndDate != nil && l.CreatedAt.After(*filter.EndDate) {
			continue
		}
		out = append(out, *l)
	}
	total := int64(len(out))
	// pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 50
	}
	start := (page - 1) * limit
	if start >= len(out) {
		return []entities.AuditLog{}, total, nil
	}
	end := start + limit
	if end > len(out) {
		end = len(out)
	}
	return out[start:end], total, nil
}

func (r *InMemoryAuditRepo) GetLogsByEntity(ctx context.Context, tenantID, entityID uint, entityType entities.EntityType) ([]entities.AuditLog, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.AuditLog
	for _, l := range r.logs {
		if l.TenantID == tenantID && l.EntityID == entityID && l.EntityType == entityType {
			out = append(out, *l)
		}
	}
	return out, nil
}

func (r *InMemoryAuditRepo) GetLogsForExport(ctx context.Context, filter entities.AuditLogFilter) ([]entities.AuditLog, error) {
	// ignore pagination and return all matching
	filter.Limit = 0
	filter.Page = 0
	logs, _, _ := r.GetLogs(ctx, filter)
	return logs, nil
}

func (r *InMemoryAuditRepo) GetStatistics(ctx context.Context, tenantID uint, startDate, endDate *time.Time) (map[string]interface{}, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	stats := make(map[string]interface{})
	var total int64
	actionCounts := map[string]int64{}
	entityCounts := map[string]int64{}
	for _, l := range r.logs {
		if l.TenantID != tenantID {
			continue
		}
		if startDate != nil && l.CreatedAt.Before(*startDate) {
			continue
		}
		if endDate != nil && l.CreatedAt.After(*endDate) {
			continue
		}
		total++
		actionCounts[string(l.Action)]++
		entityCounts[string(l.EntityType)]++
	}
	stats["total_logs"] = total
	stats["by_action"] = actionCounts
	stats["by_entity_type"] = entityCounts
	return stats, nil
}
