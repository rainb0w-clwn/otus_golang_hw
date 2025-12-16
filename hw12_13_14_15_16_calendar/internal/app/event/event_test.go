package event

import (
	"context"
	"io"
	"log"
	"testing"
	"time"

	common "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app/_common"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func createApp(t *testing.T) *App {
	t.Helper()

	st, err := storage.Get(string(storage.Memory))
	if err != nil {
		log.Fatal(err)
	}

	storageInitErr := st.Connect(context.Background())
	if storageInitErr != nil {
		log.Fatal(storageInitErr)
	}

	return &App{&common.Deps{
		Logger:  logger.New(logger.Debug, io.Discard),
		Storage: st,
	}}
}

func TestCreateEvent(t *testing.T) {
	app := createApp(t)
	dateTime := time.Now()
	event := entity.Event{
		Title:       "event 1",
		Description: "this is event 1",
		DateTime:    dateTime,
		Duration:    "02:00:00",
		RemindTime:  dateTime.Add(-time.Minute * 15),
		UserID:      1,
	}
	id1, err := app.CreateEvent(event)
	require.NoError(t, err)
	require.NotEmpty(t, id1)

	eventFromStorage, err := app.Storage.GetByID(id1)
	require.NoError(t, err)
	require.NotNil(t, eventFromStorage)

	_, err = app.CreateEvent(
		entity.Event{
			Title:       "event 2",
			Description: "this is event 2",
			DateTime:    dateTime,
			Duration:    "03:00:00",
			RemindTime:  dateTime.Add(-time.Minute * 15),
			UserID:      1,
		},
	)
	require.ErrorIs(t, err, ErrDateBusy)
}

func TestUpdateEvent(t *testing.T) {
	app := createApp(t)
	dateTime := time.Now()
	event := entity.Event{
		Title:       "event 1",
		Description: "this is event 1",
		DateTime:    dateTime.AddDate(0, 0, -1),
		Duration:    "02:00:00",
		RemindTime:  dateTime.Add(-time.Minute * 15),
		UserID:      1,
	}

	// update unknown
	err := app.UpdateEvent("random id", event)
	require.ErrorIs(t, err, ErrNotFound)

	// fill storage
	id1, err := app.CreateEvent(event)
	require.NoError(t, err)
	event2 := event
	event2.DateTime = dateTime.AddDate(0, 0, 1)
	_, err = app.CreateEvent(event2)
	require.NoError(t, err)

	// update to busy date
	event.DateTime = event2.DateTime
	err = app.UpdateEvent(id1, event)
	require.ErrorIs(t, err, ErrDateBusy)

	// successful update
	event.ID = id1
	event.DateTime = dateTime
	err = app.UpdateEvent(id1, event)
	require.NoError(t, err)

	// update active
	event.Title = "event 2"
	err = app.UpdateEvent(id1, event)
	require.ErrorIs(t, err, ErrEventIsActive)
}

func TestDeleteEvent(t *testing.T) {
	app := createApp(t)
	dateTime := time.Now()

	id1, err := app.CreateEvent(
		entity.Event{
			Title:       "event 1",
			Description: "this is event 1",
			DateTime:    dateTime.Add(-time.Hour * 24),
			Duration:    "02:00:00",
			RemindTime:  dateTime.Add(-time.Minute * 15),
			UserID:      1,
		},
	)
	require.NoError(t, err)
	err = app.DeleteEvent(id1)
	require.NoError(t, err)

	id2, err := app.CreateEvent(
		entity.Event{
			Title:       "event 2",
			Description: "this is event 2",
			DateTime:    dateTime,
			Duration:    "03:00:00",
			RemindTime:  dateTime.Add(-time.Minute * 15),
			UserID:      1,
		},
	)
	require.NoError(t, err)
	deleteErr2 := app.DeleteEvent(id2)
	require.ErrorIs(t, deleteErr2, ErrEventIsActive)
}
