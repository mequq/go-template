package datasource

import (
	"database/sql"
	"log/slog"

	_ "github.com/proullon/ramsql/driver"
)

type InmemoryDB struct {
	DB     *sql.DB
	logger *slog.Logger
}

func NewInmemoryDB(logger *slog.Logger) *InmemoryDB {
	db, err := sql.Open("ramsql", "TestLoadUserAddresses")
	if err != nil {
		panic(err)
	}

	return &InmemoryDB{
		DB:     db,
		logger: logger.With("layer", "InmemoryDB"),
	}
}
