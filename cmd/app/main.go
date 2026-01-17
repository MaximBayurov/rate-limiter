package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaximBayurov/rate-limiter/internal/app"
	"github.com/MaximBayurov/rate-limiter/internal/configuration"
	"github.com/MaximBayurov/rate-limiter/internal/database"
	"github.com/MaximBayurov/rate-limiter/internal/logger"
	"github.com/MaximBayurov/rate-limiter/internal/server"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/config/app/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config, err := configuration.New(configFile)
	if err != nil {
		panic(err)
	}

	logg, err := logger.New(config.Logger)
	if err != nil {
		panic(err)
	}

	db, err := database.New(ctx, config.Database)
	if err != nil {
		panic(err)
	}

	application := app.New(logg, db, config.App)
	serv := server.New(logg, application, config.Server)
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serv.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	if err := serv.Start(ctx); err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
