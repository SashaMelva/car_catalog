package memory

import (
	"context"
	"database/sql"
	"sync"
)

type Storage struct {
	Ctx          context.Context
	ConnectionDB *sql.DB
	sync.RWMutex
}

func New(connection *sql.DB) *Storage {
	return &Storage{
		ConnectionDB: connection,
	}
}
