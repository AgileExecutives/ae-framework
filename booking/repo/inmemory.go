package repo

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/AgileExecutives/shared-modules/booking/entities"
)

type InMemoryBookingRepo struct {
	mu      sync.Mutex
	id      uint
	configs map[uint]*entities.BookingTemplate
	// simple blacklist store: tokenID -> expiry
	blacklist map[string]time.Time
	// in-memory calendar availability and entries for tests
	calendars map[uint][]byte
	entries   []struct {
		ID               uint
		CalendarID       uint
		TenantID         uint
		StartTime        *time.Time
		EndTime          *time.Time
		SeriesID         *uint
		PositionInSeries *int
	}
}

func NewInMemoryBookingRepo() *InMemoryBookingRepo {
	return &InMemoryBookingRepo{configs: make(map[uint]*entities.BookingTemplate), blacklist: make(map[string]time.Time), calendars: make(map[uint][]byte), entries: make([]struct {
		ID               uint
		CalendarID       uint
		TenantID         uint
		StartTime        *time.Time
		EndTime          *time.Time
		SeriesID         *uint
		PositionInSeries *int
	}, 0)}
}

func (r *InMemoryBookingRepo) CreateConfiguration(ctx context.Context, cfg *entities.BookingTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.id++
	cfg.ID = r.id
	now := time.Now()
	cfg.CreatedAt = now
	cfg.UpdatedAt = now
	r.configs[cfg.ID] = cfg
	return nil
}

func (r *InMemoryBookingRepo) FindConfiguration(ctx context.Context, id, tenantID uint) (*entities.BookingTemplate, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.configs[id]
	if !ok || c.TenantID != tenantID {
		return nil, errors.New("not found")
	}
	return c, nil
}

func (r *InMemoryBookingRepo) ListConfigurations(ctx context.Context, tenantID uint, page, limit int) ([]entities.BookingTemplate, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.BookingTemplate
	for _, c := range r.configs {
		if c.TenantID == tenantID {
			out = append(out, *c)
		}
	}
	total := int64(len(out))
	start := (page - 1) * limit
	if start >= len(out) {
		return []entities.BookingTemplate{}, total, nil
	}
	end := start + limit
	if end > len(out) {
		end = len(out)
	}
	return out[start:end], total, nil
}

func (r *InMemoryBookingRepo) FindConfigurationsByUser(ctx context.Context, userID, tenantID uint) ([]entities.BookingTemplate, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.BookingTemplate
	for _, c := range r.configs {
		if c.UserID == userID && c.TenantID == tenantID {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (r *InMemoryBookingRepo) FindConfigurationsByCalendar(ctx context.Context, calendarID, tenantID uint) ([]entities.BookingTemplate, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.BookingTemplate
	for _, c := range r.configs {
		if c.CalendarID == calendarID && c.TenantID == tenantID {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (r *InMemoryBookingRepo) UpdateConfiguration(ctx context.Context, cfg *entities.BookingTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ex, ok := r.configs[cfg.ID]
	if !ok || ex.TenantID != cfg.TenantID {
		return errors.New("not found")
	}
	cfg.UpdatedAt = time.Now()
	r.configs[cfg.ID] = cfg
	return nil
}

func (r *InMemoryBookingRepo) DeleteConfiguration(ctx context.Context, id, tenantID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ex, ok := r.configs[id]
	if !ok || ex.TenantID != tenantID {
		return errors.New("not found")
	}
	delete(r.configs, id)
	return nil
}

func (r *InMemoryBookingRepo) BlacklistToken(ctx context.Context, tokenID string, userID uint, expiresAt time.Time, reason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blacklist[tokenID] = expiresAt
	return nil
}

func (r *InMemoryBookingRepo) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	exp, ok := r.blacklist[tokenID]
	if !ok {
		return false, nil
	}
	if time.Now().After(exp) {
		delete(r.blacklist, tokenID)
		return false, nil
	}
	return true, nil
}

func (r *InMemoryBookingRepo) GetCalendarWeeklyAvailability(ctx context.Context, calendarID, tenantID uint) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if b, ok := r.calendars[calendarID]; ok {
		return b, nil
	}
	return nil, nil
}

func (r *InMemoryBookingRepo) GetCalendarEntries(ctx context.Context, calendarID, tenantID uint, start, end time.Time) ([]CalendarEntryRow, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []CalendarEntryRow
	for _, e := range r.entries {
		if e.CalendarID != calendarID || e.TenantID != tenantID {
			continue
		}
		// simple time range check if start/end present
		if e.StartTime != nil && e.EndTime != nil {
			if e.StartTime.Before(start) || e.StartTime.After(end) {
				continue
			}
		}
		out = append(out, CalendarEntryRow{ID: e.ID, CalendarID: e.CalendarID, TenantID: e.TenantID, StartTime: e.StartTime, EndTime: e.EndTime, SeriesID: e.SeriesID, PositionInSeries: e.PositionInSeries})
	}
	return out, nil
}
