package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CreateEventRequestData struct {
	UserID      string `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DateTime    string `json:"dateTime"`
	Duration    string `json:"duration"`
	RemindTime  string `json:"remindTime"`
}
type CreateEventRequest struct {
	EventData CreateEventRequestData `json:"eventData"`
}

type EventID struct {
	ID string `json:"id"`
}

type CreateEventResponse struct {
	EventID EventID `json:"eventId"`
}

type Event struct {
	EventID EventID `json:"eventId"`
	Data    struct {
		UserID         string    `json:"userId"`
		Title          string    `json:"title"`
		DateTime       time.Time `json:"dateTime"`
		RemindTime     time.Time `json:"remindTime"`
		RemindSentTime time.Time `json:"remindSentTime"`
	} `json:"eventData"`
}

type Response struct {
	Events []Event `json:"events"`
}

func TestCreateEvent_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now().UTC()
	datetime := now.Add(1 * time.Hour)

	req := CreateEventRequest{
		EventData: CreateEventRequestData{
			UserID:      "1",
			Title:       "Meeting",
			Description: "Team meeting",
			DateTime:    datetime.Format(time.RFC3339),
			Duration:    "02:00:00",
			RemindTime:  datetime.Add(2 * time.Hour).Format(time.RFC3339),
		},
	}
	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequestWithContext(
		ctx, "POST", calendarBaseURL+"/event.EventService/CreateEvent", bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result CreateEventResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	require.NotEmpty(t, result.EventID.ID)
}

func TestCreateEvent_Error(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now().UTC()
	datetime := now.Add(1 * time.Hour)

	req := CreateEventRequest{
		EventData: CreateEventRequestData{
			Title:       "Meeting",
			Description: "Team meeting",
			DateTime:    datetime.Format(time.RFC3339),
			Duration:    "02:00:00",
			RemindTime:  datetime.Add(2 * time.Hour).Format(time.RFC3339),
		},
	}
	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequestWithContext(
		ctx, "POST", calendarBaseURL+"/event.EventService/CreateEvent", bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NotEqual(t, http.StatusOK, resp.StatusCode)
}

type GetDateEventRequest struct {
	StartDate string `json:"startDate"`
}

func createTestEvent(ctx context.Context, t *testing.T, userID, title string, start, end int64) {
	t.Helper()
	req := CreateEventRequest{
		EventData: CreateEventRequestData{
			UserID:   userID,
			Title:    title,
			DateTime: time.Unix(start, 0).Format(time.RFC3339),
			Duration: time.Time{}.Add(time.Duration(end) * time.Second).Format(time.TimeOnly),
		},
	}
	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequestWithContext(
		ctx, "POST", calendarBaseURL+"/event.EventService/CreateEvent", bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()
}

func TestListEvents_Day(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userID := "1"
	now := time.Now().UTC()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day()+3, 0, 0, 0, 0, time.UTC)

	createTestEvent(ctx, t, userID, "Today event", now.Unix(), now.Add(time.Hour).Unix())

	req := GetDateEventRequest{
		StartDate: startOfDay.Format(time.RFC3339),
	}
	body, _ := json.Marshal(req)

	url := fmt.Sprintf("%s/event.EventService/GetDayEvents", calendarBaseURL)

	doRangeRequestAndVerify(ctx, t, url, body, func(events []Event) {
		assert.Empty(t, events)
	})
}

func TestListEvents_Week(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userID := "1"
	now := time.Now().UTC()

	var weekStart time.Time
	if now.Weekday() == time.Sunday {
		weekStart = now.AddDate(0, 0, -6)
	} else {
		weekStart = now.AddDate(0, 0, -int(now.Weekday())+1)
	}
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC)
	weekEnd := weekStart.Add(7 * 24 * time.Hour).Add(-time.Nanosecond)

	eventTime := weekStart.Add(12 * time.Hour)
	createTestEvent(ctx, t, userID, "Weekly event", eventTime.Unix(), eventTime.Add(time.Hour).Unix())

	req := GetDateEventRequest{
		StartDate: weekStart.Format(time.RFC3339),
	}
	body, _ := json.Marshal(req)

	url := fmt.Sprintf("%s/event.EventService/GetWeekEvents", calendarBaseURL)

	doRangeRequestAndVerify(ctx, t, url, body, func(events []Event) {
		assert.NotEmpty(t, events)
		for _, e := range events {
			assert.Equal(t, userID, e.Data.UserID)
			assert.True(t, e.Data.DateTime.Compare(weekStart) >= 0 && e.Data.DateTime.Compare(weekEnd) <= 0)
		}
	})
}

func TestListEvents_Month(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userID := "1"
	now := time.Now().UTC()

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)

	eventTime := monthStart.Add(24 * time.Hour)
	createTestEvent(ctx, t, userID, "Monthly event", eventTime.Unix(), eventTime.Add(time.Hour).Unix())

	req := GetDateEventRequest{
		StartDate: monthStart.Format(time.RFC3339),
	}
	body, _ := json.Marshal(req)

	url := fmt.Sprintf("%s/event.EventService/GetMonthEvents", calendarBaseURL)

	doRangeRequestAndVerify(ctx, t, url, body, func(events []Event) {
		assert.NotEmpty(t, events)
		for _, e := range events {
			assert.Equal(t, userID, e.Data.UserID)
			assert.True(t, e.Data.DateTime.Compare(monthStart) >= 0 && e.Data.DateTime.Compare(monthEnd) <= 0)
		}
	})
}

func doRangeRequestAndVerify(ctx context.Context, t *testing.T, url string, body []byte, verifyFunc func([]Event)) {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	t.Logf("Found %d events in range: %s", len(result.Events), url)

	verifyFunc(result.Events)
}
