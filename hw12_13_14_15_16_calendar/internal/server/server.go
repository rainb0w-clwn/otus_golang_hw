package server

import (
	"time"

	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	serverGRPC "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/grpc"
	serverHTTP "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/http"
)

type Options struct {
	HTTP serverHTTP.Options
	GRPC serverGRPC.Options
}

type Server struct {
	GRPC    *serverGRPC.Server
	HTTP    *serverHTTP.Server
	logger  logger.Logger
	options Options
}

type ApplicationEvent interface {
	GetEvent(id string) (*entity.Event, error)
	CreateEvent(event entity.Event) (string, error)
	UpdateEvent(id string, event entity.Event) error
	DeleteEvent(id string) error
	GetDayEvents(day time.Time) (*entity.Events, error)
	GetWeekEvents(weekStart time.Time) (*entity.Events, error)
	GetMonthEvents(monthStart time.Time) (*entity.Events, error)
}

type Application interface {
	ApplicationEvent
}

func New(options Options, logger logger.Logger, app Application) *Server {
	grpcServer := serverGRPC.New(
		options.GRPC,
		logger,
		app,
	)
	httpServer := serverHTTP.New(
		options.HTTP,
		logger,
	)
	return &Server{
		GRPC:    &grpcServer,
		HTTP:    &httpServer,
		logger:  logger,
		options: options,
	}
}
