package googleapiclient

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var scopesCalendar = os.Getenv("CALENDAR_SCOPE")

func TestCalendar(t *testing.T) {

	c, err := NewCalendarService(keyPath, projectId, scopesCalendar)
	if err != nil {
		log.Fatal(err)
	}

	testEventDate := &GEventModel{
		Description: "Short description",
		EndDateTime: GEventDateTime{
			Date: time.Now().Add(24 * time.Hour),
		},
		Location: "Test Location",
		StartDateTime: GEventDateTime{
			Date: time.Now().Add(24 * time.Hour), // Starts Tomorrow
		},
		Summary: "TEST GEvent Date",
	}

	testEventDateTime := &GEventModel{
		Description: "Short description",
		EndDateTime: GEventDateTime{
			DateTime: time.Now().Add(48 * time.Hour).Add(2 * time.Hour),
		},
		Location: "Test Location",
		StartDateTime: GEventDateTime{
			DateTime: time.Now().Add(48 * time.Hour), // Starts in 2 days
		},
		Summary: "TEST GEvent DateTime",
	}

	var eventIDList []GEventModel
	var calendarID string

	t.Run("INSERT Calendar", func(t *testing.T) {
		summary := "Test GEvent Calendar"
		result, err := c.GCalendars.Insert(summary)
		if err != nil {
			log.Fatal(err)
		}
		t.Log(result)
		assert.Empty(t, err)
		assert.NotEmpty(t, result)

		calendarID = result
	})

	t.Run("GET Calendar", func(t *testing.T) {
		result, err := c.GCalendars.Get(calendarID)
		expected := "Test GEvent Calendar"
		actual := result
		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("UPDATE Calendar", func(t *testing.T) {
		calendar := &CalendarModel{
			Description: "Test description",
			Location:    "Unknown location",
			Summary:     "Test Calendar Override",
		}
		result, err := c.GCalendars.Patch(calendarID, calendar)
		expected := "Test GEvent Calendar"
		actual := result
		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("INSERT GEvent with Date", func(t *testing.T) {
		result, err := c.GEvent.Insert(calendarID, testEventDate)

		assert.Empty(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("INSERT GEvent with DateTime", func(t *testing.T) {
		result, err := c.GEvent.Insert(calendarID, testEventDateTime)
		t.Log(result)
		assert.Empty(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("LIST all events TimeMax, finds event in 2 days", func(t *testing.T) {
		result, err := c.GEvent.ListByTimeMin(calendarID, time.Now().Add(46*time.Hour))
		t.Log(result[0])
		assert.Empty(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("LIST all events", func(t *testing.T) {
		result, err := c.GEvent.List(calendarID)
		assert.Empty(t, err)
		assert.Equal(t, 2, len(result))

		eventIDList = append(eventIDList, result...)
	})

	t.Run("DELETE all events", func(t *testing.T) {
		for _, event := range eventIDList {
			err := c.GEvent.Delete(calendarID, event.EventID)
			assert.Empty(t, err)
		}
		result, err := c.GEvent.List(calendarID)
		assert.Empty(t, err)
		assert.True(t, len(result) == 0)
	})

	eventDateValidation := []GEventDateTime{
		{
			Date:     time.Now(),
			DateTime: time.Now(),
		},
		{
			Date:     time.Time{},
			DateTime: time.Time{},
		},
	}

	for _, date := range eventDateValidation {
		t.Run("INSERT test Start date validation", func(t *testing.T) {
			testEventValidation := &GEventModel{
				Description: "Short description",
				EndDateTime: GEventDateTime{
					DateTime: time.Now().Add(48 * time.Hour).Add(2 * time.Hour),
				},
				Location:      "Test Location",
				StartDateTime: date,
				Summary:       "TEST GEvent DateTime",
			}
			result, err := c.GEvent.Insert(calendarID, testEventValidation)
			t.Log(err)
			assert.NotEmpty(t, err)
			assert.Empty(t, result)
		})

		for _, date := range eventDateValidation {
			t.Run("INSERT test End Date validation", func(t *testing.T) {
				testEventValidation := &GEventModel{
					Description: "Short description",
					EndDateTime: date,
					Location:    "Test Location",
					StartDateTime: GEventDateTime{
						DateTime: time.Now().Add(48 * time.Hour).Add(2 * time.Hour),
					},
					Summary: "TEST GEvent DateTime",
				}
				result, err := c.GEvent.Insert(calendarID, testEventValidation)
				t.Log(err)
				assert.NotEmpty(t, err)
				assert.Empty(t, result)
			})

		}
	}

	t.Run("DELETE Calendar", func(t *testing.T) {
		err := c.GCalendars.Delete(calendarID)
		assert.Empty(t, err)
	})

	t.Run("LIST Calendars", func(t *testing.T) {
		list, err := c.GCalendars.List()
		assert.Empty(t, err)

		assert.Equal(t, 1, len(list))
		for _, calendar := range list {
			log.Println("Summary:", calendar.Summary)
		}
	})
}
