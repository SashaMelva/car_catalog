package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SashaMelva/car_catalog/internal/app"
	"github.com/SashaMelva/car_catalog/internal/config"
	"github.com/SashaMelva/car_catalog/internal/logger"
	"github.com/SashaMelva/car_catalog/internal/storage/connection"
	"github.com/SashaMelva/car_catalog/internal/storage/memory"
	"github.com/SashaMelva/car_catalog/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../configFiles/", "Path to configuration file")
}

func main() {

	config := config.New(configFile)
	log := logger.New(config.Logger, "../logFiles/")

	connectionDB := connection.New(config.DataBase, log)

	memstorage := memory.New(connectionDB.StorageDb)
	app := app.New(log, memstorage)

	httpServer := http.NewServer(log, app, config.HttpServer)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}
	}()

	log.Info("calendar is running...")

	if err := httpServer.Start(ctx); err != nil {
		log.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
