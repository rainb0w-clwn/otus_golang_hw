package entity

import (
	"errors"
	"time"
)

var ErrEventNotFound = errors.New("event not found")

type Events []Event

type Event struct {
	ID          string
	Title       string
	DateTime    time.Time
	Description string
	Duration    string
	RemindTime  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	UserID int
}
