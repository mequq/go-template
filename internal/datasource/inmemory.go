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

func NewInmemoryDB(logger *slog.Logger, controller app.Controller) (*InmemoryDB, error) {
	db, err := sql.Open("ramsql", "inmemory")
	if err != nil {
		return nil, err
	}

	mem := &InmemoryDB{
		DB:     db,
		logger: logger.With("layer", "InmemoryDB"),
	}

	controller.RegisterHealthz("inmemorydb", func(ctx context.Context) error {
		return db.PingContext(ctx)
	})
	controller.RegisterShutdown("inmemorydb", func(_ context.Context) error {
		return db.Close()
	})

	return mem, nil
}
