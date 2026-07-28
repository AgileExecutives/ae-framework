package repo

import (
	"context"
	"time"

	"github.com/AgileExecutives/ae-framwork/shared-modules/booking/entities"
)

// BookingRepo defines persistence operations for booking templates
type BookingRepo interface {
	CreateConfiguration(ctx context.Context, cfg *entities.BookingTemplate) error
	FindConfiguration(ctx context.Context, id, tenantID uint) (*entities.BookingTemplate, error)
	ListConfigurations(ctx context.Context, tenantID uint, page, limit int) ([]entities.BookingTemplate, int64, error)
	FindConfigurationsByUser(ctx context.Context, userID, tenantID uint) ([]entities.BookingTemplate, error)
	FindConfigurationsByCalendar(ctx context.Context, calendarID, tenantID uint) ([]entities.BookingTemplate, error)
	UpdateConfiguration(ctx context.Context, cfg *entities.BookingTemplate) error
	DeleteConfiguration(ctx context.Context, id, tenantID uint) error
	// Token blacklist operations
	BlacklistToken(ctx context.Context, tokenID string, userID uint, expiresAt time.Time, reason string) error
	IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error)
	// Calendar helpers
	GetCalendarWeeklyAvailability(ctx context.Context, calendarID, tenantID uint) ([]byte, error)
	GetCalendarEntries(ctx context.Context, calendarID, tenantID uint, start, end time.Time) ([]CalendarEntryRow, error)
}

// CalendarEntryRow is a lightweight row representation for calendar entries used by services
type CalendarEntryRow struct {
	ID               uint
	CalendarID       uint
	TenantID         uint
	StartTime        *time.Time
	EndTime          *time.Time
	SeriesID         *uint
	PositionInSeries *int
}
