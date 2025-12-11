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
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/service/scheduler"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yml", "Path to configuration file")
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	configR, err := os.Open(configFile)
	if err != nil {
		log.Printf("Error opening config file: %v", err)
		return 1
	}

	cfg, err := config.New(configR)
	if err != nil {
		log.Printf("Error parsing config file: %v", err)
		return 1
	}
	ctx = cfg.WithContext(ctx)

	logg := logger.New(common.LevelMap[cfg.Logger.Level], os.Stdout)

	st, err := storage.Get(cfg.Storage)
	if err != nil {
		logg.Error("Error getting storage: %v", err)
		return 1
	}

	err = st.Connect(ctx)
	if err != nil {
		logg.Error("Error init storage: %v", err)
		return 1
	}

	qManager := queue.NewRabbitManager(logg)
	err = qManager.Connect(ctx)
	if err != nil {
		logg.Error("Error connecting to RabbitMQ server: %v", err)
		return 1
	}

	application := app.New(logg, st)

	service := scheduler.New(
		application,
		logg,
		qManager,
	)
	err = service.Run(ctx)
	if err != nil {
		logg.Error("Error starting scheduler: %v", err)
		return 1
	}

	err = qManager.Close()
	if err != nil {
		logg.Error("Error closing RabbitMQ connection: %v", err)

		return 1
	}

	logg.Info("Connection closed.")

	return 0
}
