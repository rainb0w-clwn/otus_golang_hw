package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	common "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app/_common"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server"
	serverGRPC "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/grpc"
	serverHTTP "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/service/calendar"
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

	application := app.New(logg, st)

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
		application,
	)

	service := calendar.New(srv, logg)

	err = service.Run(ctx)
	if err != nil {
		logg.Error("Error starting calendar: %v", err)
		return 1
	}

	return 0
}
