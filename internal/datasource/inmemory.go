package datasource

import (
	"application/app"
	"context"
	"database/sql"
	"log/slog"

	_ "github.com/proullon/ramsql/driver"
)

type InmemoryDB struct {
	*sql.DB

	logger *slog.Logger
}

func NewInmemoryDB(logger *slog.Logger, controller app.Controller) *InmemoryDB {
	db, err := sql.Open("ramsql", "TestLoadUserAddresses")
	if err != nil {
		panic(err)
	}

	controller.RegisterHealthz(
		"inmemorydb",
		func(ctx context.Context) error {
			return db.PingContext(ctx)
		},
	)

	return &InmemoryDB{
		DB:     db,
		logger: logger.With("layer", "InmemoryDB"),
	}
}
