package server

import (
	"context"
	"net"
	"time"

	proto "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/api"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/grpc/log"
	"google.golang.org/grpc"
)

type Options struct {
	Host, Port     string
	ConnectTimeout time.Duration
}

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type server struct {
	*grpc.Server
	logger logger.Logger
}

type Application interface {
	CreateEvent(event entity.Event) (string, error)
	UpdateEvent(id string, event entity.Event) error
	DeleteEvent(id string) error
	GetDayEvents(day time.Time) (*entity.Events, error)
	GetWeekEvents(weekStart time.Time) (*entity.Events, error)
	GetMonthEvents(monthStart time.Time) (*entity.Events, error)
}

func New(options Options, logger logger.Logger, app Application) Server {
	serverGRPC := grpc.NewServer(
		grpc.ConnectionTimeout(options.ConnectTimeout),
		grpc.ChainUnaryInterceptor(
			log.New(logger),
		),
	)
	proto.RegisterEventServiceServer(serverGRPC, NewService(app, logger))
	return &server{serverGRPC, logger}
}

func (s *server) Start(ctx context.Context) error {
	cfg := config.GetFromContext(ctx)
	if cfg == nil {
		return config.ErrNoConfigInContext
	}

	var lc net.ListenConfig

	listener, err := lc.Listen(
		ctx,
		"tcp",
		net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port),
	)
	if err != nil {
		return err
	}

	err = s.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) Stop(_ context.Context) error {
	s.GracefulStop()
	return nil
}
