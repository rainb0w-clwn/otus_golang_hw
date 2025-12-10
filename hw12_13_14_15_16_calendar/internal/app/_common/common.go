package common

import (
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type Deps struct {
	Storage storage.Storage
	Logger  logger.Logger
}

var LevelMap = map[string]logger.Level{
	"debug":   logger.Debug,
	"info":    logger.Info,
	"warning": logger.Warning,
	"error":   logger.Error,
}
