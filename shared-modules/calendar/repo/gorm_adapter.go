package repo

import (
	"context"
	"time"

	"github.com/AgileExecutives/ae-framwork/shared-modules/calendar/entities"
	"gorm.io/gorm"
)

// GormCalendarRepo is a GORM-backed implementation of CalendarRepo
type GormCalendarRepo struct {
	db *gorm.DB
}

func NewGormCalendarRepo(db *gorm.DB) *GormCalendarRepo {
	return &GormCalendarRepo{db: db}
}

// Calendars
func (r *GormCalendarRepo) CreateCalendar(ctx context.Context, c *entities.Calendar) error {
	return r.db.Create(c).Error
}

func (r *GormCalendarRepo) FindCalendarByID(ctx context.Context, id, tenantID, userID uint) (*entities.Calendar, error) {
	var c entities.Calendar
	if err := r.db.Preload("CalendarSeries").Preload("CalendarEntries").Preload("ExternalCalendars").
		Where("id = ? AND tenant_id = ? AND user_id = ?", id, tenantID, userID).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *GormCalendarRepo) ListCalendars(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.Calendar, int, error) {
	var calendars []entities.Calendar
	var total int64
	offset := (page - 1) * limit
	if err := r.db.Model(&entities.Calendar{}).Where("tenant_id = ? AND user_id = ?", tenantID, userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Preload("CalendarEntries").Preload("CalendarEntries.Series").Preload("CalendarSeries").
		Preload("CalendarSeries.CalendarEntries").Preload("ExternalCalendars").
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).Offset(offset).Limit(limit).Find(&calendars).Error; err != nil {
		return nil, 0, err
	}
	return calendars, int(total), nil
}

func (r *GormCalendarRepo) ListCalendarsNoPaging(ctx context.Context, tenantID, userID uint) ([]entities.Calendar, error) {
	var calendars []entities.Calendar
	if err := r.db.Preload("CalendarEntries").Preload("CalendarEntries.Series").Preload("CalendarSeries").
		Preload("CalendarSeries.CalendarEntries").Preload("ExternalCalendars").
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).Find(&calendars).Error; err != nil {
		return nil, err
	}
	return calendars, nil
}

func (r *GormCalendarRepo) UpdateCalendar(ctx context.Context, c *entities.Calendar) error {
	return r.db.Save(c).Error
}

func (r *GormCalendarRepo) DeleteCalendar(ctx context.Context, id, tenantID, userID uint) error {
	var c entities.Calendar
	if err := r.db.Where("id = ? AND tenant_id = ? AND user_id = ?", id, tenantID, userID).First(&c).Error; err != nil {
		return err
	}
	return r.db.Delete(&c).Error
}

// Entries
func (r *GormCalendarRepo) CreateCalendarEntry(ctx context.Context, e *entities.CalendarEntry) error {
	return r.db.Create(e).Error
}

func (r *GormCalendarRepo) FindCalendarEntryByID(ctx context.Context, id, tenantID, userID uint) (*entities.CalendarEntry, error) {
	var e entities.CalendarEntry
	if err := r.db.Preload("Calendar").Preload("Series").
		Where("id = ? AND tenant_id = ? AND user_id = ?", id, tenantID, userID).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *GormCalendarRepo) ListCalendarEntries(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.CalendarEntry, int, error) {
	var entries []entities.CalendarEntry
	var total int64
	offset := (page - 1) * limit
	if err := r.db.Model(&entities.CalendarEntry{}).Where("tenant_id = ? AND user_id = ?", tenantID, userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Preload("Calendar").Preload("Series").Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Offset(offset).Limit(limit).Find(&entries).Error; err != nil {
		return nil, 0, err
	}
	return entries, int(total), nil
}

func (r *GormCalendarRepo) UpdateCalendarEntry(ctx context.Context, e *entities.CalendarEntry) error {
	return r.db.Save(e).Error
}

func (r *GormCalendarRepo) DeleteCalendarEntry(ctx context.Context, id, tenantID, userID uint) error {
	var e entities.CalendarEntry
	if err := r.db.Where("id = ? AND tenant_id = ? AND user_id = ?", id, tenantID, userID).First(&e).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&e).Error; err != nil {
		return err
	}
	return nil
}

// Series
func (r *GormCalendarRepo) CreateCalendarSeries(ctx context.Context, s *entities.CalendarSeries) error {
	return r.db.Create(s).Error
}

func (r *GormCalendarRepo) FindCalendarSeriesByID(ctx context.Context, id, tenantID, userID uint) (*entities.CalendarSeries, error) {
	var s entities.CalendarSeries
	if err := r.db.Where("id = ? AND tenant_id = ? AND user_id = ?", id, tenantID, userID).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *GormCalendarRepo) ListCalendarSeries(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.CalendarSeries, int, error) {
	var list []entities.CalendarSeries
	var total int64
	offset := (page - 1) * limit
	if err := r.db.Model(&entities.CalendarSeries{}).Where("tenant_id = ? AND user_id = ?", tenantID, userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Where("tenant_id = ? AND user_id = ?", tenantID, userID).Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, int(total), nil
}

func (r *GormCalendarRepo) UpdateCalendarSeries(ctx context.Context, s *entities.CalendarSeries) error {
	return r.db.Save(s).Error
}

func (r *GormCalendarRepo) DeleteCalendarSeries(ctx context.Context, id, tenantID, userID uint) error {
	var s entities.CalendarSeries
	if err := r.db.Where("id = ? AND tenant_id = ? AND user_id = ?", id, tenantID, userID).First(&s).Error; err != nil {
		return err
	}
	return r.db.Delete(&s).Error
}

// External calendars
func (r *GormCalendarRepo) CreateExternalCalendar(ctx context.Context, e *entities.ExternalCalendar) error {
	return r.db.Create(e).Error
}

func (r *GormCalendarRepo) FindExternalCalendarByID(ctx context.Context, id, tenantID, userID uint) (*entities.ExternalCalendar, error) {
	var e entities.ExternalCalendar
	if err := r.db.Where("id = ? AND tenant_id = ? AND user_id = ?", id, tenantID, userID).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *GormCalendarRepo) ListExternalCalendars(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.ExternalCalendar, int, error) {
	var list []entities.ExternalCalendar
	var total int64
	offset := (page - 1) * limit
	if err := r.db.Model(&entities.ExternalCalendar{}).Where("tenant_id = ? AND user_id = ?", tenantID, userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Where("tenant_id = ? AND user_id = ?", tenantID, userID).Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, int(total), nil
}

func (r *GormCalendarRepo) UpdateExternalCalendar(ctx context.Context, e *entities.ExternalCalendar) error {
	return r.db.Save(e).Error
}

func (r *GormCalendarRepo) DeleteExternalCalendar(ctx context.Context, id, tenantID, userID uint) error {
	var e entities.ExternalCalendar
	if err := r.db.Where("id = ? AND tenant_id = ? AND user_id = ?", id, tenantID, userID).First(&e).Error; err != nil {
		return err
	}
	return r.db.Delete(&e).Error
}

// ListCalendarEntriesInRange returns entries overlapping the provided date range
func (r *GormCalendarRepo) ListCalendarEntriesInRange(ctx context.Context, tenantID, userID uint, dateFrom, dateTo time.Time) ([]entities.CalendarEntry, error) {
	var entries []entities.CalendarEntry
	// Use StartTime/EndTime overlap logic
	if err := r.db.Preload("Calendar").Preload("Series").
		Where("tenant_id = ? AND user_id = ? AND start_time <= ? AND end_time >= ?", tenantID, userID, dateTo, dateFrom).
		Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}
