package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	common "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app/_common"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server"
	serverGRPC "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/grpc"
	serverHTTP "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/configs/config.yml", "Path to configuration file")
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return 0
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	configR, err := os.Open(configFile)
	if err != nil {
		log.Printf("%s", "Error opening config file: "+err.Error())
		return 1
	}

	cfg, err := config.New(configR)
	if err != nil {
		log.Printf("%s", "Error parsing config file: "+err.Error())
		return 1
	}
	ctx = cfg.WithContext(ctx)

	logg := logger.New(common.LevelMap[cfg.Logger.Level], os.Stdout)

	st, err := storage.Get(cfg.Storage)
	if err != nil {
		logg.Error("Error getting storage: " + err.Error())
		return 1
	}

	err = st.Connect(ctx)
	if err != nil {
		logg.Error("Error init storage: " + err.Error())
		return 1
	}

	srv := server.New(
		server.Options{
			GRPC: serverGRPC.Options{
				Host:           cfg.GRPC.Host,
				Port:           cfg.GRPC.Port,
				ConnectTimeout: cfg.GRPC.ConnectTimeout,
			},
			HTTP: serverHTTP.Options{
				Host:         cfg.HTTP.Host,
				Port:         cfg.HTTP.Port,
				ReadTimeout:  cfg.HTTP.ReadTimeout,
				WriteTimeout: cfg.HTTP.WriteTimeout,
			},
		},
		logg,
		app.New(logg, st),
	)

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()

		logg.Info("GRPC server starting...")
		err := (*srv.GRPC).Start(ctx)
		if err != nil {
			logg.Error("Failed to start GRPC server: " + err.Error())
			cancel()
		}
	}()
	go func() {
		defer wg.Done()

		logg.Info("HTTP server starting...")
		err := (*srv.HTTP).Start(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Error("Failed to start HTTP server: " + err.Error())
			cancel()
		}
	}()
	go func() {
		defer wg.Done()
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		logg.Info("GRPC server stopping...")
		if err := (*srv.GRPC).Stop(ctx); err != nil {
			logg.Error("Failed to stop GRPC server: " + err.Error())
		}

		logg.Info("HTTP server stopping...")
		if err := (*srv.HTTP).Stop(ctx); err != nil {
			logg.Error("Failed to stop HTTP server: " + err.Error())
		}

		logg.Info("Calendar stopped")
	}()

	wg.Wait()

	return 0
}
