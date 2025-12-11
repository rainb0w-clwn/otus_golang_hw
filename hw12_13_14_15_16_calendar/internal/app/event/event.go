package event

import (
	"errors"
	"time"

	common "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app/_common"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
)

type App struct {
	*common.Deps
}

var (
	ErrNotFound      = errors.New("event not found")
	ErrDateBusy      = errors.New("time not available")
	ErrEventIsActive = errors.New("can't modify active event")
)

// CreateEvent create event if requested time is not busy.
func (a App) CreateEvent(event entity.Event) (string, error) {
	existingEvent, getErr := a.Storage.GetForTime(event.DateTime)
	if existingEvent != nil {
		return "", ErrDateBusy
	}

	if getErr != nil && !errors.Is(getErr, entity.ErrEventNotFound) {
		a.Logger.Error(getErr.Error())

		return "", getErr
	}

	id, createErr := a.Storage.Create(event)
	if createErr != nil {
		a.Logger.Error(createErr.Error())

		return "", createErr
	}

	return id, nil
}

// UpdateEvent updates event if it is not active and requested time is not busy.
func (a App) UpdateEvent(id string, event entity.Event) error {
	// check has event
	existingEvent, readErr := a.Storage.GetByID(id)
	if readErr != nil {
		a.Logger.Error(readErr.Error())

		if errors.Is(readErr, entity.ErrEventNotFound) {
			return ErrNotFound
		}
	}

	// check not active
	if existingEvent.DateTime.Round(time.Minute) == time.Now().Round(time.Minute) {
		return ErrEventIsActive
	}

	// check new time not busy
	eventWithRequestedTime, getErr := a.Storage.GetForTime(event.DateTime)
	if eventWithRequestedTime != nil {
		return ErrDateBusy
	}

	if getErr != nil && !errors.Is(getErr, entity.ErrEventNotFound) {
		a.Logger.Error(getErr.Error())

		return getErr
	}

	// update
	updateErr := a.Storage.Update(event)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

// DeleteEvent deletes event if it is not active.
func (a App) DeleteEvent(id string) error {
	event, readErr := a.Storage.GetByID(id)
	if readErr != nil {
		a.Logger.Error(readErr.Error())

		return readErr
	}

	if event.DateTime.Round(time.Minute) == time.Now().Round(time.Minute) {
		return ErrEventIsActive
	}

	deleteErr := a.Storage.Delete(id)
	if deleteErr != nil {
		a.Logger.Error(deleteErr.Error())

		return deleteErr
	}

	return nil
}

// GetDayEvents returns events for passed day. Use UTC time format.
func (a App) GetDayEvents(day time.Time) (*entity.Events, error) {
	events, err := a.Storage.GetForPeriod(
		day.Truncate(time.Hour*24),
		day.Round(time.Hour*24),
	)
	if err != nil {
		a.Logger.Error(err.Error())

		return nil, err
	}

	return events, nil
}

// GetWeekEvents returns events for week starts with weekStart. Use UTC time format.
func (a App) GetWeekEvents(weekStart time.Time) (*entity.Events, error) {
	weekStart = weekStart.Truncate(time.Hour * 24)
	weekEnd := weekStart.Add(time.Hour * 24 * 7)

	events, err := a.Storage.GetForPeriod(weekStart, weekEnd)
	if err != nil {
		a.Logger.Error(err.Error())

		return nil, err
	}

	return events, nil
}

// GetMonthEvents returns events for month starts with monthStart. Use UTC time format.
func (a App) GetMonthEvents(monthStart time.Time) (*entity.Events, error) {
	monthStart = monthStart.Truncate(time.Hour * 24)
	monthEnd := time.Date(monthStart.Year(), monthStart.Month()+1, monthStart.Day(), 0, 0, 0, 0, monthStart.Location())

	events, err := a.Storage.GetForPeriod(monthStart, monthEnd)
	if err != nil {
		a.Logger.Error(err.Error())

		return nil, err
	}

	return events, nil
}

func (a App) DeleteEventsOlderThan(t time.Time) error {
	return a.Storage.DeleteOlderThan(t)
}

func (a App) GetEventsForRemind() (*entity.Events, error) {
	events, err := a.Storage.GetForRemind()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (a App) MarkEventAsReminded(id string) error {
	err := a.Storage.MarkAsReminded(id)
	if err != nil {
		return err
	}
	return nil
}
