package connection

import (
	"context"
	"database/sql"

	"github.com/SashaMelva/car_catalog/internal/config"
	"go.uber.org/zap"
)

type StorageConnection struct {
	StorageDb *sql.DB
}

func New(confDB *config.ConfigDB, log *zap.SugaredLogger) *StorageConnection {
	dsn := "postgres://" + confDB.User + ":" + confDB.Password + "@" + confDB.Host + ":" + confDB.Port + "/" + confDB.NameDB
	storage, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Fatal("Cannot open pgx driver: %w", err)
	}

	log.Debug("DSN connection database " + dsn)

	return &StorageConnection{storage}
}

func (s *StorageConnection) Connect(ctx context.Context) error {
	err := s.StorageDb.PingContext(ctx)
	return err
}

func (s *StorageConnection) Close(ctx context.Context) error {
	s.StorageDb.Close()
	return nil
}
