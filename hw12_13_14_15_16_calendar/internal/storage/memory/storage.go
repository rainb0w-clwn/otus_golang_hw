package memorystorage

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string]*entity.Event
}

func New() *Storage {
	return &Storage{}
}

func NewWithEvents(events map[string]*entity.Event) *Storage {
	return &Storage{data: events}
}

func (s *Storage) GetByID(id string) (*entity.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, has := s.data[id]
	if !has {
		return nil, entity.ErrEventNotFound
	}

	return event, nil
}

func (s *Storage) GetAll() (*entity.Events, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make(entity.Events, 0, len(s.data))

	for _, event := range s.data {
		events = append(events, event)
	}

	return &events, nil
}

func (s *Storage) Create(event entity.Event) (string, error) {
	event.ID = uuid.New().String()

	s.mu.Lock()
	s.data[event.ID] = &event
	s.mu.Unlock()

	return event.ID, nil
}

func (s *Storage) Update(event entity.Event) error {
	_, err := s.GetByID(event.ID)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.data[event.ID] = &event
	s.mu.Unlock()

	return nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	delete(s.data, id)
	s.mu.Unlock()

	return nil
}

func (s *Storage) GetForTime(t time.Time) (*entity.Event, error) {
	for _, event := range s.data {
		if event.DateTime == t {
			return event, nil
		}
	}

	return nil, entity.ErrEventNotFound
}

func (s *Storage) GetForPeriod(periodStart time.Time, periodEnd time.Time) (*entity.Events, error) {
	periodEvents := make(entity.Events, 0)

	for _, event := range s.data {
		if event.DateTime.After(periodStart) && event.DateTime.Before(periodEnd) {
			periodEvents = append(periodEvents, event)
		}
	}

	return &periodEvents, nil
}

func (s *Storage) GetForRemind() (*entity.Events, error) {
	remindEvents := make(entity.Events, 0)

	for _, event := range s.data {
		if event.RemindSentTime.IsZero() && slices.Contains(
			[]int{-1, 0},
			event.RemindTime.Compare(time.Now().UTC()),
		) {
			remindEvents = append(remindEvents, event)
		}
	}

	return &remindEvents, nil
}

func (s *Storage) MarkAsReminded(id string) error {
	event, err := s.GetByID(id)
	if err != nil {
		return nil //nolint:nilerr
	}
	s.mu.Lock()
	event.RemindSentTime = time.Now().UTC()
	s.mu.Unlock()

	return nil
}

func (s *Storage) DeleteOlderThan(t time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	for _, event := range s.data {
		if event.DateTime.Compare(t) == -1 {
			err = s.Delete(event.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Storage) Connect(_ context.Context) error {
	s.data = make(map[string]*entity.Event)

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = nil

	return nil
}
