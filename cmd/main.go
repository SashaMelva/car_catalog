package main

import (
	"context"
	"flag"
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

// @title           Car Catalog API
// @version         1.0
// @description     This is a sample server car catalog.

// @host      localhost:8080
// @BasePath  /

func main() {

	config := config.New(configFile)
	log := logger.New(config.Logger, "../logFiles/")

	connectionDB := connection.New(config.DataBase, log)

	memstorage := memory.New(connectionDB.StorageDb, log)
	app := app.New(log, memstorage, config.HostClientApi)

	httpServer := http.NewServer(log, app, config.HttpServer)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		httpServer.Stop(ctx)
	}()
}
