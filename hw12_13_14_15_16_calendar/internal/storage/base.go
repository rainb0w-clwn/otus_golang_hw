package storage

import (
	"context"
	"errors"
	"time"

	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type Type string

const (
	Memory Type = "memory"
	DB     Type = "db"
)

var ErrInvalidStorageValue = errors.New("invalid storage value in config")

type ConnectionStorage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

type Storage interface {
	ConnectionStorage

	Create(entity.Event) (string, error)
	Update(entity.Event) error
	Delete(string) error
	GetAll() (*entity.Events, error)
	GetByID(string) (*entity.Event, error)
	GetForPeriod(time.Time, time.Time) (*entity.Events, error)
	GetForTime(time.Time) (*entity.Event, error)
}

func Get(storageType string) (Storage, error) {
	if Type(storageType) != Memory && Type(storageType) != DB {
		return nil, ErrInvalidStorageValue
	}

	if Type(storageType) == Memory {
		return memorystorage.New(), nil
	}

	return sqlstorage.New(), nil
}
