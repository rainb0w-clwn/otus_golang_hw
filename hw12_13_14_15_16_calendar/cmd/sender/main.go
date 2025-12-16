package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	common "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app/_common"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/service/sender"
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

	qManager := queue.NewRabbitManager(logg)
	err = qManager.Connect(ctx)
	if err != nil {
		logg.Error("Error connecting to RabbitMQ server: %v", err)
		return 1
	}

	service := sender.New(logg, qManager)

	err = service.Run(ctx)
	if err != nil {
		logg.Error("Error starting sender: %v", err)
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
