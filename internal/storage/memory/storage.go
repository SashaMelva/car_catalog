package memory

import (
	"context"
	"database/sql"
	"sync"

	"go.uber.org/zap"
)

type Storage struct {
	Logger       *zap.SugaredLogger
	Ctx          context.Context
	ConnectionDB *sql.DB
	sync.RWMutex
}

func New(connection *sql.DB, log *zap.SugaredLogger) *Storage {
	return &Storage{
		Logger:       log,
		ConnectionDB: connection,
	}
}
