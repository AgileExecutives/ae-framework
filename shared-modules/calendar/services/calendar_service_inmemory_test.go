package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/AgileExecutives/ae-framwork/shared-modules/calendar/entities"
	repo "github.com/AgileExecutives/ae-framwork/shared-modules/calendar/repo"
	"github.com/stretchr/testify/require"
)

func TestCalendarService_WithInMemoryRepo_WeekYearViews(t *testing.T) {
	r := repo.NewInMemoryCalendarRepo()
	svc := NewCalendarServiceWithRepo(r)

	// create calendar
	req := makeCalendarReqForTest()
	cal, err := svc.CreateCalendar(req, 1, 1)
	require.NoError(t, err)

	// create an entry within a specific week and year
	start := time.Date(2026, 7, 6, 9, 0, 0, 0, time.UTC) // Tuesday
	end := time.Date(2026, 7, 6, 10, 0, 0, 0, time.UTC)
	entryReq := makeEntryReq(cal.ID)
	entryReq.StartTime = &start
	entryReq.EndTime = &end

	_, err = svc.CreateCalendarEntry(entryReq, 1, 1)
	require.NoError(t, err)

	// Week view for that date should include the entry
	entriesWeek, err := svc.GetCalendarWeekView(start, 1, 1)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(entriesWeek), 1)

	// Year view for 2026 should include the entry
	entriesYear, err := svc.GetCalendarYearView(2026, 1, 1)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(entriesYear), 1)
}

// helper mimics test fixture used elsewhere
func makeCalendarReqForTest() entities.CreateCalendarRequest {
	wa, _ := json.Marshal(map[string]interface{}{"monday": []string{"09:00-17:00"}})
	return entities.CreateCalendarRequest{
		Title:              "InMemory Cal",
		Color:              "#00FF00",
		WeeklyAvailability: wa,
		Timezone:           "UTC",
	}
}
