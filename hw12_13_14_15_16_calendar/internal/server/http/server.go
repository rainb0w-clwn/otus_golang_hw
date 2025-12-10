package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	proto "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/api"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/http/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Options struct {
	Host, Port                string
	ReadTimeout, WriteTimeout time.Duration
}

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type server struct {
	*http.Server
	logger logger.Logger
}

func New(options Options, logger logger.Logger) Server {
	serverHTTP := &http.Server{
		Addr:         net.JoinHostPort(options.Host, options.Port),
		ReadTimeout:  options.ReadTimeout,
		WriteTimeout: options.WriteTimeout,
	}
	return &server{
		serverHTTP,
		logger,
	}
}

func (s *server) Start(ctx context.Context) error {
	cfg := config.GetFromContext(ctx)
	if cfg == nil {
		return config.ErrNoConfigInContext
	}

	conn, err := grpc.NewClient(
		net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	mux := runtime.NewServeMux()
	for _, f := range []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error{
		proto.RegisterEventServiceHandler,
	} {
		if err = f(ctx, mux, conn); err != nil {
			return err
		}
	}
	s.Handler = log.NewHandler(s.logger, mux)

	err = s.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
