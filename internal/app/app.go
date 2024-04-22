package app

import (
	"github.com/SashaMelva/car_catalog/internal/storage/memory"
	"go.uber.org/zap"
)

type App struct {
	storage *memory.Storage
	Logger  *zap.SugaredLogger
	JwtKey  string
}

func New(logger *zap.SugaredLogger, storage *memory.Storage) *App {
	return &App{
		storage: storage,
		Logger:  logger,
	}
}
