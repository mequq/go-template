package datasource

import (
	"application/app"
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	*sql.DB

	logger     *slog.Logger
	controller app.Controller
}

const (
	// PostgresDriver is the driver name for PostgreSQL.
	PostgresDriver  string        = "pgx"
	MaxOpenConns    int           = 25
	MaxIdleConns    int           = 5
	ConnMaxLifetime time.Duration = 1 * time.Hour
	ConnMaxIdleTime time.Duration = 1 * time.Minute
)

func NewPostgresDB(ctx context.Context, logger *slog.Logger, controller app.Controller) (*PostgresDB, error) {
	dsn := "postgresql://user:password@localhost:5432/db?sslmode=disable"

	db, err := otelsql.Open(
		PostgresDriver,
		dsn,
		otelsql.WithAttributes(
			otelsql.AttributesFromDSN(dsn)...,
		),
	)
	if err != nil {
		return nil, err
	}

	otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		otelsql.AttributesFromDSN(dsn)...,
	))

	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)
	db.SetConnMaxIdleTime(ConnMaxIdleTime)

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	pg := &PostgresDB{
		DB:         db,
		logger:     logger.With("layer", "PostgresDB"),
		controller: controller,
	}

	controller.RegisterHealthz("pgx", pg.healthz)
	controller.RegisterShutdown("pgx", pg.shutdown)

	return pg, nil
}

// healthz implements DatasourceHealthzer.
func (p *PostgresDB) healthz(ctx context.Context) error {
	if err := p.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

// shutdown implements app.Shutdowner.
func (p *PostgresDB) shutdown(ctx context.Context) error {
	p.logger.Info("shutting down PostgresDB")

	if err := p.Close(); err != nil {
		return err
	}

	return nil
}
