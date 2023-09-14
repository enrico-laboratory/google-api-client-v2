package googleapiclient

import (
	"context"
	"errors"
	"github.com/enrico-laboratory/go-validator"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"time"
)

const (
	TzEuropeAmsterdam = "Europe/Amsterdam"
)

type Calendar struct {
	GCalendars
	GEvent
}

type GCalendars struct {
	calendar *calendar.Service
}

type GEvent struct {
	event *calendar.Service
}

type CalendarModel struct {
	Id          string
	Description string
	Location    string
	Summary     string
}

type GEventModel struct {
	EventID       string
	Description   string
	EndDateTime   GEventDateTime
	Location      string
	StartDateTime GEventDateTime
	Summary       string
}

type GEventDateTime struct {
	Date     time.Time
	DateTime time.Time
}

func NewCalendarService(keypath, projectId string, scopes ...string) (*Calendar, error) {
	c, err := getCredentials(keypath, projectId, scopes...)
	if err != nil {
		return nil, err
	}
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(c))
	if err != nil {
		return nil, err
	}

	cal := &Calendar{
		GCalendars: GCalendars{
			calendar: srv,
		},
		GEvent: GEvent{
			event: srv,
		},
	}

	return cal, nil
}

func (c *GCalendars) Insert(summary string) (string, error) {

	resp, err := c.calendar.Calendars.Insert(&calendar.Calendar{
		Summary: summary,
	}).Do()
	if err != nil {
		return "", err
	}

	return resp.Id, nil
}
func (c *GCalendars) Get(calendarId string) (string, error) {
	resp, err := c.calendar.Calendars.Get(calendarId).Do()
	if err != nil {
		return "", err
	}
	return resp.Summary, nil
}

// Patch adds only default reminders (method: popup, minutes: 90) the argument "cal" is not considered
func (c *GCalendars) Patch(calendarID string, cal *CalendarModel) (string, error) {
	defaultReminder := &calendar.EventReminder{
		Method:  "popup",
		Minutes: 90,
	}
	var reminderList []*calendar.EventReminder
	reminderList = append(reminderList, defaultReminder)

	calendarList := &calendar.CalendarListEntry{
		//ForegroundColor:  cal.ColorId,
		DefaultReminders: reminderList,
		Description:      cal.Description,
		Id:               calendarID,
		Location:         cal.Location,
		//SummaryOverride: cal.Summary,
		//TimeZone: c.config.timeZone,
	}

	resp, err := c.calendar.CalendarList.Patch(calendarID, calendarList).Do()

	if err != nil {
		return "", err
	}

	return resp.Summary, err
}

func (c *GCalendars) Delete(calendarID string) error {
	err := c.calendar.Calendars.Delete(calendarID).Do()
	if err != nil {
		return err
	}
	return nil
}

func (c *GCalendars) List() ([]CalendarModel, error) {
	list, err := c.calendar.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}
	var gCalendarList []CalendarModel
	for _, cal := range list.Items {
		var gCalendar CalendarModel
		gCalendar.Description = cal.Description
		gCalendar.Location = cal.Location
		gCalendar.Summary = cal.Summary
		gCalendar.Id = cal.Id

		gCalendarList = append(gCalendarList, gCalendar)
	}

	return gCalendarList, nil

}

func (c *GEvent) Insert(calendarID string, event *GEventModel) (string, error) {
	v := validator.New()

	if validateDates(v, event); !v.Valid() {
		return "", errors.New(v.ErrorsToString())
	}

	var endDateTime calendar.EventDateTime
	if event.EndDateTime.Date.IsZero() && !event.EndDateTime.DateTime.IsZero() {
		endDateTime.DateTime = event.EndDateTime.DateTime.Format(time.RFC3339)
	} else if !event.EndDateTime.Date.IsZero() && event.EndDateTime.DateTime.IsZero() {
		endDateTime.Date = event.EndDateTime.Date.Format("2006-01-02")
	}
	endDateTime.TimeZone = TzEuropeAmsterdam

	var startDateTime calendar.EventDateTime
	if event.StartDateTime.Date.IsZero() && !event.StartDateTime.DateTime.IsZero() {
		startDateTime.DateTime = event.StartDateTime.DateTime.Format(time.RFC3339)
	} else if !event.StartDateTime.Date.IsZero() && event.StartDateTime.DateTime.IsZero() {
		startDateTime.Date = event.StartDateTime.Date.Format("2006-01-02")
	}
	startDateTime.TimeZone = TzEuropeAmsterdam

	var overrides []*calendar.EventReminder
	var override calendar.EventReminder
	override.Method = "popup"
	override.Minutes = 90
	overrides = append(overrides, &override)
	//reminders := &calendar.EventReminders{
	//	Overrides:  overrides,
	//	UseDefault: false,
	//}

	eventParsed := &calendar.Event{
		Description: event.Description,
		End:         &endDateTime,
		Location:    event.Location,
		Start:       &startDateTime,
		//Reminders:   reminders,
		Status:  "confirmed",
		Summary: event.Summary,
	}

	resp, err := c.event.Events.Insert(calendarID, eventParsed).Do()
	if err != nil {
		return "", err
	}

	return resp.Id, nil

}

func (c *GEvent) List(calendarID string) ([]GEventModel, error) {

	resp, err := c.ListByTimeMin(calendarID, time.Date(1900, 01, 01, 00, 00, 00, 00, &time.Location{}))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *GEvent) ListByTimeMin(calendarID string, timeMax time.Time) ([]GEventModel, error) {

	resp, err := c.event.Events.List(calendarID).TimeMin(timeMax.Format(time.RFC3339)).Do()
	if err != nil {
		return nil, err
	}

	var idList []string
	for _, event := range resp.Items {
		idList = append(idList, event.Id)
	}

	var gEvents []GEventModel

	for _, event := range resp.Items {

		var startDateTimeObject GEventDateTime
		var endDateTimeObject GEventDateTime

		if event.Start.Date != "" {
			layoutDate := "2006-01-02"
			startDate, err := time.Parse(layoutDate, event.Start.Date)
			if err != nil {
				return nil, err
			}
			endDate, err := time.Parse(layoutDate, event.Start.Date)
			if err != nil {
				return nil, err
			}
			startDateTimeObject.Date = startDate
			endDateTimeObject.Date = endDate
		} else {
			startDateTime, err := time.Parse(time.RFC3339, event.Start.DateTime)
			if err != nil {
				return nil, err
			}
			endDateTime, err := time.Parse(time.RFC3339, event.End.DateTime)
			if err != nil {
				return nil, err
			}
			startDateTimeObject.DateTime = startDateTime
			endDateTimeObject.DateTime = endDateTime
		}

		gEvent := GEventModel{
			EventID:       event.Id,
			Description:   event.Description,
			EndDateTime:   endDateTimeObject,
			Location:      event.Location,
			StartDateTime: startDateTimeObject,
			Summary:       event.Summary,
		}
		gEvents = append(gEvents, gEvent)
	}

	return gEvents, nil
}

func (c *GEvent) Delete(calendarID string, eventId string) error {
	err := c.event.Events.Delete(calendarID, eventId).Do()
	if err != nil {
		return err
	}
	return nil
}

func validateDates(v *validator.Validator, event *GEventModel) {

	v.Check(bothDateNotFull(&event.StartDateTime), "Start-date", "chose Date or DateTime, cannot input both")
	v.Check(bothDateNotEmpty(&event.StartDateTime), "Start-date", "at least one date must be present")
	v.Check(bothDateNotFull(&event.EndDateTime), "end-date", "chose Date or DateTime, cannot input both")
	v.Check(bothDateNotEmpty(&event.EndDateTime), "end-date", "at least one date must be present")
}

func bothDateNotEmpty(date *GEventDateTime) bool {
	if date.Date.IsZero() && date.DateTime.IsZero() {
		return false
	}
	return true
}

func bothDateNotFull(date *GEventDateTime) bool {
	if !date.Date.IsZero() && !date.DateTime.IsZero() {
		return false
	}
	return true
}
