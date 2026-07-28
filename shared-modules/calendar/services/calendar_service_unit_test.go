package services

import (
	"testing"
	"time"

	"github.com/AgileExecutives/ae-framework/shared-modules/calendar/entities"
	"github.com/AgileExecutives/ae-framework/shared-modules/calendar/repo"
)

func TestCalendarService_InMemoryRepo_CreateAndEntry(t *testing.T) {
	r := repo.NewInMemoryCalendarRepo()
	svc := NewCalendarServiceWithRepo(r)

	tenantID := uint(10)
	userID := uint(20)

	// Create calendar
	req := entities.CreateCalendarRequest{Title: "Test Calendar", Timezone: "UTC"}
	cal, err := svc.CreateCalendar(req, tenantID, userID)
	if err != nil {
		t.Fatalf("CreateCalendar failed: %v", err)
	}
	if cal.ID == 0 {
		t.Fatalf("expected calendar ID to be set")
	}

	// Create calendar entry
	start := time.Now().UTC()
	end := start.Add(1 * time.Hour)
	entryReq := entities.CreateCalendarEntryRequest{
		CalendarID: cal.ID,
		Title:      "Meeting",
		StartTime:  &start,
		EndTime:    &end,
	}
	entry, err := svc.CreateCalendarEntry(entryReq, tenantID, userID)
	if err != nil {
		t.Fatalf("CreateCalendarEntry failed: %v", err)
	}
	if entry.ID == 0 {
		t.Fatalf("expected entry ID to be set")
	}

	// Retrieve calendar by ID
	got, err := svc.GetCalendarByID(cal.ID, tenantID, userID)
	if err != nil {
		t.Fatalf("GetCalendarByID failed: %v", err)
	}
	if got.ID != cal.ID {
		t.Fatalf("expected calendar ID %d got %d", cal.ID, got.ID)
	}
}
