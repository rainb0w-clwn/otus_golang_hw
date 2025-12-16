package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
	common "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/_common"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/http/log"
)

type Server struct {
	server http.Server
	logger common.Logger
}

type Application interface {
	CreateEvent(event entity.Event) (string, error)
	UpdateEvent(id string, event entity.Event) error
	DeleteEvent(id string) error
	GetDayEvents(day time.Time) (*entity.Events, error)
	GetWeekEvents(weekStart time.Time) (*entity.Events, error)
	GetMonthEvents(monthStart time.Time) (*entity.Events, error)
}

func NewServer(options common.Options, logger common.Logger, app Application) *Server {
	return &Server{
		server: http.Server{
			Addr:         net.JoinHostPort(options.Host, options.Port),
			Handler:      log.NewHandler(logger, NewHandler(app, logger)),
			ReadTimeout:  options.ReadTimeout,
			WriteTimeout: options.WriteTimeout,
		},
		logger: logger,
	}
}

func (s *Server) Start(_ context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewHandler(_ Application, logger common.Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		_, err := writer.Write([]byte("Hello World"))
		if err != nil {
			logger.Error(err.Error())
		}
	}
}
