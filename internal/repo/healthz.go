package repo

import (
	"context"
	"log/slog"

	healthzusecase "application/internal/biz"
	"application/internal/datasource"
	"application/pkg/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type HealthzDS struct {
	logger *slog.Logger
	memDB  *datasource.InmemoryDB
}

func NewHealthzDS(logger *slog.Logger, memDB *datasource.InmemoryDB) healthzusecase.HealthzRepoInterface {
	return &HealthzDS{
		logger: logger.With("layer", "HealthzDS"),
		memDB:  memDB,
	}
}

func (r *HealthzDS) Readiness(ctx context.Context) error {
	_, span := otel.Tracer("repo", trace.WithInstrumentationVersion("12"),
		trace.WithInstrumentationAttributes(attribute.String("a", "d"))).Start(ctx, "Readiness")
	defer span.End()
	if err := r.memDB.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (r *HealthzDS) Liveness(ctx context.Context) error {
	logger := r.logger.With("method", "Liveness", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("repo Liveness")
	return nil
}
