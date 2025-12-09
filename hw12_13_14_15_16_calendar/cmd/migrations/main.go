package main

import (
	"context"
	"flag"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
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

	args := os.Args[1:]
	if len(args) < 2 {
		flag.Usage()
		return 0
	}
	command := args[1]

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

	db, err := goose.OpenDBWithDriver("postgres", cfg.DB.Dsn)
	if err != nil {
		log.Printf("goose: failed to open DB: %v\n\n", err)

		return 1
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	var arguments []string
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.RunContext(context.Background(), command, db, cfg.DB.MigrationsDir, arguments...); err != nil {
		log.Printf("goose %v: %v\n", command, err)
		return 1
	}

	return 0
}
