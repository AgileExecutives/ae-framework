package repo

import (
	"context"
	"time"

	"github.com/AgileExecutives/ae-framwork/shared-modules/booking/entities"
	"gorm.io/gorm"
)

type GormBookingRepo struct {
	db *gorm.DB
}

func NewGormBookingRepo(db *gorm.DB) *GormBookingRepo {
	return &GormBookingRepo{db: db}
}

func (r *GormBookingRepo) CreateConfiguration(ctx context.Context, cfg *entities.BookingTemplate) error {
	return r.db.Create(cfg).Error
}

func (r *GormBookingRepo) FindConfiguration(ctx context.Context, id, tenantID uint) (*entities.BookingTemplate, error) {
	var c entities.BookingTemplate
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *GormBookingRepo) ListConfigurations(ctx context.Context, tenantID uint, page, limit int) ([]entities.BookingTemplate, int64, error) {
	var configs []entities.BookingTemplate
	var total int64
	query := r.db.Model(&entities.BookingTemplate{}).Where("tenant_id = ?", tenantID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&configs).Error; err != nil {
		return nil, 0, err
	}
	return configs, total, nil
}

func (r *GormBookingRepo) FindConfigurationsByUser(ctx context.Context, userID, tenantID uint) ([]entities.BookingTemplate, error) {
	var configs []entities.BookingTemplate
	if err := r.db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *GormBookingRepo) FindConfigurationsByCalendar(ctx context.Context, calendarID, tenantID uint) ([]entities.BookingTemplate, error) {
	var configs []entities.BookingTemplate
	if err := r.db.Where("calendar_id = ? AND tenant_id = ?", calendarID, tenantID).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *GormBookingRepo) UpdateConfiguration(ctx context.Context, cfg *entities.BookingTemplate) error {
	return r.db.Save(cfg).Error
}

func (r *GormBookingRepo) DeleteConfiguration(ctx context.Context, id, tenantID uint) error {
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&entities.BookingTemplate{}).Error
}

func (r *GormBookingRepo) BlacklistToken(ctx context.Context, tokenID string, userID uint, expiresAt time.Time, reason string) error {
	entry := map[string]interface{}{
		"token_id":   tokenID,
		"user_id":    userID,
		"expires_at": expiresAt,
		"reason":     reason,
	}
	return r.db.Table("token_blacklist").Create(&entry).Error
}

func (r *GormBookingRepo) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	var count int64
	if err := r.db.Table("token_blacklist").Where("token_id = ? AND expires_at > ? AND deleted_at IS NULL", tokenID, time.Now()).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormBookingRepo) GetCalendarWeeklyAvailability(ctx context.Context, calendarID, tenantID uint) ([]byte, error) {
	var row struct {
		WeeklyAvailability []byte `gorm:"column:weekly_availability"`
	}
	if err := r.db.Table("calendars").Select("weekly_availability").Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", calendarID, tenantID).Take(&row).Error; err != nil {
		return nil, err
	}
	return row.WeeklyAvailability, nil
}

func (r *GormBookingRepo) GetCalendarEntries(ctx context.Context, calendarID, tenantID uint, start, end time.Time) ([]CalendarEntryRow, error) {
	var rows []CalendarEntryRow
	query := r.db.Table("calendar_entries").Select("id, calendar_id, tenant_id, start_time, end_time, series_id, position_in_series").Where("calendar_id = ? AND tenant_id = ? AND start_time >= ? AND start_time <= ? AND deleted_at IS NULL", calendarID, tenantID, start, end)
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
