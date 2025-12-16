package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	common "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app/_common"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	servercommon "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/_common"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/server/http"
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
		log.Println("Error opening config file.")
		return 1
	}

	cfg, err := config.New(configR)
	if err != nil {
		log.Println("Error parsing config file.")
		return 1
	}
	ctx = cfg.WithContext(ctx)

	logg := logger.New(common.LevelMap[cfg.Logger.Level], os.Stdout)

	st, err := storage.Get(cfg.Storage)
	if err != nil {
		log.Println("Error getting storage.")
		return 1
	}

	err = st.Connect(ctx)
	if err != nil {
		log.Println("Error init storage.")
		return 1
	}

	srv := server.NewServer(
		servercommon.Options{
			Host:         cfg.HTTP.Host,
			Port:         cfg.HTTP.Port,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
		},
		logg,
		app.New(logg, st).Event,
	)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := srv.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		logg.Info("calendar stopped.")
	}()

	logg.Info("calendar is running...")

	if err := srv.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		return 1
	}

	wg.Wait()

	return 0
}
