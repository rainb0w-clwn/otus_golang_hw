package entity

import (
	"errors"
	"time"
)

var ErrEventNotFound = errors.New("event not found")

type Events []*Event

type Event struct {
	ID             string
	UserID         int
	Title          string
	DateTime       time.Time
	Description    string
	Duration       string
	RemindTime     time.Time
	RemindSentTime time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type EventMsg struct {
	ID       string
	UserID   int
	Title    string
	DateTime time.Time
}

func (e Event) ToMsg() EventMsg {
	return EventMsg{
		ID:       e.ID,
		UserID:   e.UserID,
		Title:    e.Title,
		DateTime: e.DateTime,
	}
}
