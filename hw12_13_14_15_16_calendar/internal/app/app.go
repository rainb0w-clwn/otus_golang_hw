package app

import (
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app/_common"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app/event"
)

type App struct {
	*common.Deps
	Event event.App
}

func New(logger common.Logger, storage common.Storage) *App {
	deps := &common.Deps{
		Logger:  logger,
		Storage: storage,
	}
	return &App{
		Deps:  deps,
		Event: event.App{Deps: deps},
	}
}
