package repo

import (
	"context"
	"time"

	"github.com/AgileExecutives/shared-modules/audit/entities"
	"gorm.io/gorm"
)

type GormAuditRepo struct{ db *gorm.DB }

func NewGormAuditRepo(db *gorm.DB) *GormAuditRepo { return &GormAuditRepo{db: db} }

func (r *GormAuditRepo) CreateLog(ctx context.Context, log *entities.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *GormAuditRepo) GetLogs(ctx context.Context, filter entities.AuditLogFilter) ([]entities.AuditLog, int64, error) {
	var logs []entities.AuditLog
	var total int64
	query := r.db.Model(&entities.AuditLog{}).Where("tenant_id = ?", filter.TenantID)
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.EntityType != nil {
		query = query.Where("entity_type = ?", *filter.EntityType)
	}
	if filter.EntityID != nil {
		query = query.Where("entity_id = ?", *filter.EntityID)
	}
	if filter.Action != nil {
		query = query.Where("action = ?", *filter.Action)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *GormAuditRepo) GetLogsByEntity(ctx context.Context, tenantID, entityID uint, entityType entities.EntityType) ([]entities.AuditLog, error) {
	var logs []entities.AuditLog
	if err := r.db.Where("tenant_id = ? AND entity_type = ? AND entity_id = ?", tenantID, entityType, entityID).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *GormAuditRepo) GetLogsForExport(ctx context.Context, filter entities.AuditLogFilter) ([]entities.AuditLog, error) {
	filter.Limit = 0
	filter.Page = 0
	var logs []entities.AuditLog
	query := r.db.Model(&entities.AuditLog{}).Where("tenant_id = ?", filter.TenantID)
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.EntityType != nil {
		query = query.Where("entity_type = ?", *filter.EntityType)
	}
	if filter.EntityID != nil {
		query = query.Where("entity_id = ?", *filter.EntityID)
	}
	if filter.Action != nil {
		query = query.Where("action = ?", *filter.Action)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	if err := query.Order("created_at ASC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *GormAuditRepo) GetStatistics(ctx context.Context, tenantID uint, startDate, endDate *time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	query := r.db.Model(&entities.AuditLog{}).Where("tenant_id = ?", tenantID)
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total_logs"] = total

	var actionCounts []struct {
		Action entities.AuditAction
		Count  int64
	}
	if err := query.Select("action, COUNT(*) as count").Group("action").Scan(&actionCounts).Error; err != nil {
		return nil, err
	}
	actionStats := map[string]int64{}
	for _, ac := range actionCounts {
		actionStats[string(ac.Action)] = ac.Count
	}
	stats["by_action"] = actionStats

	var entityCounts []struct {
		EntityType entities.EntityType
		Count      int64
	}
	if err := query.Select("entity_type, COUNT(*) as count").Group("entity_type").Scan(&entityCounts).Error; err != nil {
		return nil, err
	}
	entityStats := map[string]int64{}
	for _, ec := range entityCounts {
		entityStats[string(ec.EntityType)] = ec.Count
	}
	stats["by_entity_type"] = entityStats

	return stats, nil
}
