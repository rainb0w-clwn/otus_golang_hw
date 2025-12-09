package common

import (
	common "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app/_common"
	"time"
)

type Logger interface {
	common.Logger
}

type Options struct {
	Host, Port                string
	ReadTimeout, WriteTimeout time.Duration
}
