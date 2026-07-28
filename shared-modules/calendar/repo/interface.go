package repo

import (
	"context"
	"time"

	"github.com/AgileExecutives/ae-framwork/shared-modules/calendar/entities"
)

// CalendarRepo defines persistence operations for calendar domain
type CalendarRepo interface {
	// Calendars
	CreateCalendar(ctx context.Context, c *entities.Calendar) error
	FindCalendarByID(ctx context.Context, id, tenantID, userID uint) (*entities.Calendar, error)
	ListCalendars(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.Calendar, int, error)
	ListCalendarsNoPaging(ctx context.Context, tenantID, userID uint) ([]entities.Calendar, error)
	UpdateCalendar(ctx context.Context, c *entities.Calendar) error
	DeleteCalendar(ctx context.Context, id, tenantID, userID uint) error

	// Entries
	CreateCalendarEntry(ctx context.Context, e *entities.CalendarEntry) error
	FindCalendarEntryByID(ctx context.Context, id, tenantID, userID uint) (*entities.CalendarEntry, error)
	ListCalendarEntries(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.CalendarEntry, int, error)
	UpdateCalendarEntry(ctx context.Context, e *entities.CalendarEntry) error
	DeleteCalendarEntry(ctx context.Context, id, tenantID, userID uint) error

	// Series
	CreateCalendarSeries(ctx context.Context, s *entities.CalendarSeries) error
	FindCalendarSeriesByID(ctx context.Context, id, tenantID, userID uint) (*entities.CalendarSeries, error)
	ListCalendarSeries(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.CalendarSeries, int, error)
	UpdateCalendarSeries(ctx context.Context, s *entities.CalendarSeries) error
	DeleteCalendarSeries(ctx context.Context, id, tenantID, userID uint) error

	// External calendars
	CreateExternalCalendar(ctx context.Context, e *entities.ExternalCalendar) error
	FindExternalCalendarByID(ctx context.Context, id, tenantID, userID uint) (*entities.ExternalCalendar, error)
	ListExternalCalendars(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.ExternalCalendar, int, error)
	UpdateExternalCalendar(ctx context.Context, e *entities.ExternalCalendar) error
	DeleteExternalCalendar(ctx context.Context, id, tenantID, userID uint) error

	// Views / Queries
	ListCalendarEntriesInRange(ctx context.Context, tenantID, userID uint, dateFrom, dateTo time.Time) ([]entities.CalendarEntry, error)
}
