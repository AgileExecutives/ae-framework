package repo

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/AgileExecutives/shared-modules/calendar/entities"
)

type InMemoryCalendarRepo struct {
	mu        sync.Mutex
	calID     uint
	entryID   uint
	seriesID  uint
	extID     uint
	calendars map[uint]*entities.Calendar
	entries   map[uint]*entities.CalendarEntry
	series    map[uint]*entities.CalendarSeries
	externals map[uint]*entities.ExternalCalendar
}

func NewInMemoryCalendarRepo() *InMemoryCalendarRepo {
	return &InMemoryCalendarRepo{
		calendars: make(map[uint]*entities.Calendar),
		entries:   make(map[uint]*entities.CalendarEntry),
		series:    make(map[uint]*entities.CalendarSeries),
		externals: make(map[uint]*entities.ExternalCalendar),
		calID:     0,
		entryID:   0,
		seriesID:  0,
		extID:     0,
	}
}

// Calendars
func (r *InMemoryCalendarRepo) CreateCalendar(ctx context.Context, c *entities.Calendar) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calID++
	c.ID = r.calID
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	r.calendars[c.ID] = c
	return nil
}

func (r *InMemoryCalendarRepo) FindCalendarByID(ctx context.Context, id, tenantID, userID uint) (*entities.Calendar, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.calendars[id]
	if !ok || c.TenantID != tenantID || c.UserID != userID {
		return nil, errors.New("not found")
	}
	return c, nil
}

func (r *InMemoryCalendarRepo) ListCalendars(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.Calendar, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.Calendar
	for _, c := range r.calendars {
		if c.TenantID == tenantID && c.UserID == userID {
			out = append(out, *c)
		}
	}
	total := len(out)
	// apply pagination
	start := (page - 1) * limit
	if start >= total {
		return []entities.Calendar{}, total, nil
	}
	end := start + limit
	if end > total {
		end = total
	}
	return out[start:end], total, nil
}

func (r *InMemoryCalendarRepo) ListCalendarsNoPaging(ctx context.Context, tenantID, userID uint) ([]entities.Calendar, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.Calendar
	for _, c := range r.calendars {
		if c.TenantID == tenantID && c.UserID == userID {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (r *InMemoryCalendarRepo) UpdateCalendar(ctx context.Context, c *entities.Calendar) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.calendars[c.ID]
	if !ok || existing.TenantID != c.TenantID || existing.UserID != c.UserID {
		return errors.New("not found")
	}
	c.UpdatedAt = time.Now()
	r.calendars[c.ID] = c
	return nil
}

func (r *InMemoryCalendarRepo) DeleteCalendar(ctx context.Context, id, tenantID, userID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.calendars[id]
	if !ok || c.TenantID != tenantID || c.UserID != userID {
		return errors.New("not found")
	}
	delete(r.calendars, id)
	return nil
}

// Entries
func (r *InMemoryCalendarRepo) CreateCalendarEntry(ctx context.Context, e *entities.CalendarEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Ensure calendar exists
	cal, ok := r.calendars[e.CalendarID]
	if !ok || cal.TenantID != e.TenantID || cal.UserID != e.UserID {
		return errors.New("calendar not found")
	}
	r.entryID++
	e.ID = r.entryID
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	r.entries[e.ID] = e
	return nil
}

func (r *InMemoryCalendarRepo) FindCalendarEntryByID(ctx context.Context, id, tenantID, userID uint) (*entities.CalendarEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[id]
	if !ok || e.TenantID != tenantID || e.UserID != userID {
		return nil, errors.New("not found")
	}
	return e, nil
}

func (r *InMemoryCalendarRepo) ListCalendarEntries(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.CalendarEntry, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.CalendarEntry
	for _, e := range r.entries {
		if e.TenantID == tenantID && e.UserID == userID {
			out = append(out, *e)
		}
	}
	total := len(out)
	start := (page - 1) * limit
	if start >= total {
		return []entities.CalendarEntry{}, total, nil
	}
	end := start + limit
	if end > total {
		end = total
	}
	return out[start:end], total, nil
}

func (r *InMemoryCalendarRepo) UpdateCalendarEntry(ctx context.Context, e *entities.CalendarEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ex, ok := r.entries[e.ID]
	if !ok || ex.TenantID != e.TenantID || ex.UserID != e.UserID {
		return errors.New("not found")
	}
	e.UpdatedAt = time.Now()
	r.entries[e.ID] = e
	return nil
}

func (r *InMemoryCalendarRepo) DeleteCalendarEntry(ctx context.Context, id, tenantID, userID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[id]
	if !ok || e.TenantID != tenantID || e.UserID != userID {
		return errors.New("not found")
	}
	delete(r.entries, id)
	return nil
}

// Series
func (r *InMemoryCalendarRepo) CreateCalendarSeries(ctx context.Context, s *entities.CalendarSeries) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// ensure calendar exists
	cal, ok := r.calendars[s.CalendarID]
	if !ok || cal.TenantID != s.TenantID || cal.UserID != s.UserID {
		return errors.New("calendar not found")
	}
	r.seriesID++
	s.ID = r.seriesID
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	r.series[s.ID] = s
	return nil
}

func (r *InMemoryCalendarRepo) FindCalendarSeriesByID(ctx context.Context, id, tenantID, userID uint) (*entities.CalendarSeries, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.series[id]
	if !ok || s.TenantID != tenantID || s.UserID != userID {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (r *InMemoryCalendarRepo) ListCalendarSeries(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.CalendarSeries, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.CalendarSeries
	for _, s := range r.series {
		if s.TenantID == tenantID && s.UserID == userID {
			out = append(out, *s)
		}
	}
	total := len(out)
	start := (page - 1) * limit
	if start >= total {
		return []entities.CalendarSeries{}, total, nil
	}
	end := start + limit
	if end > total {
		end = total
	}
	return out[start:end], total, nil
}

func (r *InMemoryCalendarRepo) UpdateCalendarSeries(ctx context.Context, s *entities.CalendarSeries) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ex, ok := r.series[s.ID]
	if !ok || ex.TenantID != s.TenantID || ex.UserID != s.UserID {
		return errors.New("not found")
	}
	s.UpdatedAt = time.Now()
	r.series[s.ID] = s
	return nil
}

func (r *InMemoryCalendarRepo) DeleteCalendarSeries(ctx context.Context, id, tenantID, userID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.series[id]
	if !ok || s.TenantID != tenantID || s.UserID != userID {
		return errors.New("not found")
	}
	delete(r.series, id)
	return nil
}

// External calendars
func (r *InMemoryCalendarRepo) CreateExternalCalendar(ctx context.Context, e *entities.ExternalCalendar) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// ensure calendar exists
	cal, ok := r.calendars[e.CalendarID]
	if !ok || cal.TenantID != e.TenantID || cal.UserID != e.UserID {
		return errors.New("calendar not found")
	}
	r.extID++
	e.ID = r.extID
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	r.externals[e.ID] = e
	return nil
}

func (r *InMemoryCalendarRepo) FindExternalCalendarByID(ctx context.Context, id, tenantID, userID uint) (*entities.ExternalCalendar, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.externals[id]
	if !ok || e.TenantID != tenantID || e.UserID != userID {
		return nil, errors.New("not found")
	}
	return e, nil
}

func (r *InMemoryCalendarRepo) ListExternalCalendars(ctx context.Context, page, limit int, tenantID, userID uint) ([]entities.ExternalCalendar, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.ExternalCalendar
	for _, e := range r.externals {
		if e.TenantID == tenantID && e.UserID == userID {
			out = append(out, *e)
		}
	}
	total := len(out)
	start := (page - 1) * limit
	if start >= total {
		return []entities.ExternalCalendar{}, total, nil
	}
	end := start + limit
	if end > total {
		end = total
	}
	return out[start:end], total, nil
}

func (r *InMemoryCalendarRepo) UpdateExternalCalendar(ctx context.Context, e *entities.ExternalCalendar) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ex, ok := r.externals[e.ID]
	if !ok || ex.TenantID != e.TenantID || ex.UserID != e.UserID {
		return errors.New("not found")
	}
	e.UpdatedAt = time.Now()
	r.externals[e.ID] = e
	return nil
}

func (r *InMemoryCalendarRepo) DeleteExternalCalendar(ctx context.Context, id, tenantID, userID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.externals[id]
	if !ok || e.TenantID != tenantID || e.UserID != userID {
		return errors.New("not found")
	}
	delete(r.externals, id)
	return nil
}

// ListCalendarEntriesInRange returns entries overlapping the provided date range
func (r *InMemoryCalendarRepo) ListCalendarEntriesInRange(ctx context.Context, tenantID, userID uint, dateFrom, dateTo time.Time) ([]entities.CalendarEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.CalendarEntry
	for _, e := range r.entries {
		if e.TenantID != tenantID || e.UserID != userID {
			continue
		}
		// Consider entries overlapping the range: entry.StartTime <= dateTo && entry.EndTime >= dateFrom
		if e.StartTime != nil && e.EndTime != nil {
			if !e.StartTime.After(dateTo) && !e.EndTime.Before(dateFrom) {
				out = append(out, *e)
			}
		}
	}
	return out, nil
}
