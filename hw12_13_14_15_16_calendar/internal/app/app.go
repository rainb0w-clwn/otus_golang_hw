package app

import (
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app/_common"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app/event"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	*common.Deps
	*event.App
}

func New(logger logger.Logger, storage storage.Storage) *App {
	deps := &common.Deps{
		Logger:  logger,
		Storage: storage,
	}
	return &App{
		Deps: deps,
		App:  &event.App{Deps: deps},
	}
}
