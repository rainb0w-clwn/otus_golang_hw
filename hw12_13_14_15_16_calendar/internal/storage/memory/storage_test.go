package memorystorage

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/stretchr/testify/require"
)

var (
	now   = time.Now()
	event = entity.Event{
		Title:       "some event",
		DateTime:    now,
		Description: "this is some event",
		Duration:    "60",
		RemindTime:  "15",
		CreatedAt:   now,
		UpdatedAt:   now,
		UserID:      1,
	}
)
var initialDate = time.Date(2025, 12, 5, 12, 0o0, 0, 0, time.UTC)

func TestStorageModify(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		memStorage := New()
		connectErr := memStorage.Connect(context.Background())
		require.NoError(t, connectErr)
		_, createErr := memStorage.Create(event)
		storageEvents, _ := memStorage.GetAll()
		require.NoError(t, createErr)
		require.Len(t, *storageEvents, 1)
	})
	t.Run("update", func(t *testing.T) {
		newTitle := "new title"
		memStorage := New()
		connectErr := memStorage.Connect(context.Background())
		require.NoError(t, connectErr)
		// create to init storage
		id, createErr := memStorage.Create(event)
		require.NoError(t, createErr)
		// read & update
		event1, readErr := memStorage.GetByID(id)
		require.NoError(t, readErr)
		event1.Title = newTitle
		updateErr := memStorage.Update(*event1)
		require.NoError(t, updateErr)
		// assert
		event2, readErr := memStorage.GetByID(id)
		require.NoError(t, readErr)
		require.Equal(t, newTitle, event2.Title)
	})
	t.Run("update unknown", func(t *testing.T) {
		memStorage := New()
		connectErr := memStorage.Connect(context.Background())
		require.NoError(t, connectErr)
		// create to init storage
		_, createErr := memStorage.Create(event)
		require.NoError(t, createErr)
		// generate random ID & try to update
		event.ID = uuid.New().String()
		updateErr := memStorage.Update(event)
		// assert
		require.ErrorIs(t, updateErr, entity.ErrEventNotFound)
	})
	t.Run("delete", func(t *testing.T) {
		memStorage := New()
		connectErr := memStorage.Connect(context.Background())
		require.NoError(t, connectErr)
		// create to init storage
		id, createErr := memStorage.Create(event)
		require.NoError(t, createErr)
		// delete
		deleteErr := memStorage.Delete(id)
		// assert
		events, _ := memStorage.GetAll()
		require.NoError(t, deleteErr)
		require.Len(t, *events, 0)
	})
}

func TestStorageRead(t *testing.T) {
	t.Run("read unknown", func(t *testing.T) {
		memStorage := New()
		connectErr := memStorage.Connect(context.Background())
		require.NoError(t, connectErr)
		// create to init storage
		_, createErr := memStorage.Create(event)
		require.NoError(t, createErr)
		// generate random ID & try to read
		id := uuid.New().String()
		_, readErr := memStorage.GetByID(id)
		// assert
		require.ErrorIs(t, readErr, entity.ErrEventNotFound)
	})

	t.Run("read all", func(t *testing.T) {
		n := 3
		memStorage := New()
		connectErr := memStorage.Connect(context.Background())
		require.NoError(t, connectErr)
		// create to init storage
		for i := 0; i < n; i++ {
			_, createErr := memStorage.Create(event)
			require.NoError(t, createErr)
		}
		events, _ := memStorage.GetAll()
		require.Len(t, *events, n)
	})

	t.Run("read for day", func(t *testing.T) {
		strg := NewWithEvents(map[string]entity.Event{
			"1": {ID: "1", Title: "1", DateTime: initialDate, UserID: 1},
			"2": {ID: "2", Title: "2", DateTime: initialDate.Add(time.Hour * 2), UserID: 1},
			"3": {ID: "3", Title: "3", DateTime: initialDate.Add(-time.Hour * 2), UserID: 1},
			"4": {ID: "4", Title: "4", DateTime: initialDate.Add(time.Hour * 24), UserID: 1},
			"5": {ID: "5", Title: "5", DateTime: initialDate.Add(-time.Hour * 24), UserID: 1},
		})
		dayBeginning := time.Date(
			initialDate.Year(),
			initialDate.Month(),
			initialDate.Day(),
			0, 0, 0, 0,
			initialDate.Location(),
		)
		events, err := strg.GetForPeriod(
			dayBeginning,
			dayBeginning.Add(time.Hour*24-time.Second),
		)
		require.NoError(t, err)
		require.Len(t, *events, 3)
		require.Equal(t, []string{"1", "2", "3"}, getKeys(t, events))
	})

	t.Run("read for week", func(t *testing.T) {
		str := NewWithEvents(map[string]entity.Event{
			"1": {ID: "1", Title: "1", DateTime: initialDate, UserID: 1},
			"2": {ID: "2", Title: "2", DateTime: initialDate.Add(time.Hour * 24 * 10), UserID: 1},
			"3": {ID: "3", Title: "3", DateTime: initialDate.Add(-time.Hour * 24 * 10), UserID: 1},
			"4": {ID: "4", Title: "4", DateTime: initialDate.Add(time.Hour * 24 * 2), UserID: 1},
		})
		weekBeginning := weekStartDate(initialDate)
		events, err := str.GetForPeriod(
			weekBeginning,
			weekBeginning.AddDate(0, 0, 7).Add(-time.Second),
		)
		require.NoError(t, err)
		require.Len(t, *events, 2)
		require.Equal(t, []string{"1", "4"}, getKeys(t, events))
	})

	t.Run("read for month", func(t *testing.T) {
		str := NewWithEvents(map[string]entity.Event{
			"1": {ID: "1", Title: "1", DateTime: initialDate, UserID: 1},
			"2": {ID: "2", Title: "2", DateTime: initialDate.Add(time.Hour * 24 * 30), UserID: 1},
			"3": {ID: "3", Title: "3", DateTime: initialDate.Add(time.Hour * 24 * 10), UserID: 1},
			"4": {ID: "4", Title: "4", DateTime: initialDate.Add(time.Hour * 24 * 2), UserID: 1},
			"5": {ID: "5", Title: "5", DateTime: initialDate.Add(-time.Hour * 24 * 30), UserID: 1},
		})
		monthBeginning := time.Date(initialDate.Year(), initialDate.Month(), 1, 0, 0, 0, 0, initialDate.Location())
		events, err := str.GetForPeriod(
			monthBeginning,
			monthBeginning.AddDate(0, 1, 0).Add(-time.Second),
		)
		require.NoError(t, err)
		require.Len(t, *events, 3)
		require.Equal(t, []string{"1", "3", "4"}, getKeys(t, events))
	})
}

func weekStartDate(date time.Time) time.Time {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	result := time.
		Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, initialDate.Location()).
		Add(time.Duration(offset*24) * time.Hour)
	return result
}

func getKeys(t *testing.T, events *entity.Events) []string {
	t.Helper()
	keys := make([]string, 0, len(*events))

	for _, event := range *events {
		keys = append(keys, event.ID)
	}

	sort.Strings(keys)

	return keys
}
