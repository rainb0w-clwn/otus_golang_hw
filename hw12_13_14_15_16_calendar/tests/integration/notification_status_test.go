package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNotificationIsSent(t *testing.T) {
	t.Run("WaitForNotificationStatus", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		userID := "2"
		now := time.Now().UTC()
		req := CreateEventRequest{
			EventData: CreateEventRequestData{
				UserID:     userID,
				Title:      "Test",
				DateTime:   now.Format(time.RFC3339),
				Duration:   now.Add(time.Hour).Format(time.TimeOnly),
				RemindTime: now.Add(-6 * time.Hour).Format(time.RFC3339),
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
		var result CreateEventResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		resp.Body.Close()
		id := result.EventID.ID
		require.NotEmpty(t, id)

		<-time.After(10 * time.Second)
		reqG := EventID{
			ID: id,
		}
		body, _ = json.Marshal(reqG)
		httpReq, err = http.NewRequestWithContext(
			ctx, "POST", calendarBaseURL+"/event.EventService/GetEvent", bytes.NewBuffer(body),
		)
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client = &http.Client{}
		resp, err = client.Do(httpReq)
		require.NoError(t, err)
		var resultG Event
		err = json.NewDecoder(resp.Body).Decode(&resultG)
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, id, resultG.EventID.ID)
		t.Logf("Event: %v\n", resultG.Data)
		require.NotZero(t, resultG.Data.RemindSentTime)
	})
}
