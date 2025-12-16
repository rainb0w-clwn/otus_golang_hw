package calendar

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server"
)

type Calendar struct {
	server *server.Server
	logger logger.Logger
}

func New(
	server *server.Server,
	logger logger.Logger,
) *Calendar {
	return &Calendar{
		server: server,
		logger: logger,
	}
}

func (s *Calendar) Run(ctx context.Context) error {
	cfg := config.GetFromContext(ctx)
	if cfg == nil {
		return config.ErrNoConfigInContext
	}
	var err error
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()

		s.logger.Info("GRPC server starting...")
		err = (*s.server.GRPC).Start(ctx)
		if err != nil {
			s.logger.Error("Failed to start GRPC server: %v", err)
			// cancel()
		}
	}()
	go func() {
		defer wg.Done()

		s.logger.Info("HTTP server starting...")
		err = (*s.server.HTTP).Start(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Failed to start HTTP server: %v", err)
			// cancel()
		}
	}()
	go func() {
		defer wg.Done()
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		s.logger.Info("GRPC server stopping...")
		if err = (*s.server.GRPC).Stop(ctx); err != nil {
			s.logger.Error("Failed to stop GRPC server: %v", err)
		}

		s.logger.Info("HTTP server stopping...")
		if err = (*s.server.HTTP).Stop(ctx); err != nil {
			s.logger.Error("Failed to stop HTTP server: %v", err)
		}

		s.logger.Info("Calendar stopped")
	}()
	wg.Wait()
	return err
}
