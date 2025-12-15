package integration

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)

var calendarBaseURL = "http://calendar-app:" + os.Getenv("HTTP_PORT")

func TestCalendarHealth(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", calendarBaseURL+"/health", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}
