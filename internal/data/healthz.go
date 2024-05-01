package data

import (
	"application/internal/biz"
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type HealthzRepo struct {
	logger *slog.Logger
	ds     *DataSource
}

func NewHealthzRepo(logger *slog.Logger, ds *DataSource) biz.HealthzRepoInterface {
	return &HealthzRepo{
		logger: logger.With("layer", "HealthzRepo"),
		ds:     ds,
	}
}

func (r *HealthzRepo) Readiness(ctx context.Context) error {
	// logger := r.logger.With("method", "Readiness", "ctx", biz.LogContext(ctx))
	// logger.Debug("repo Readiness")
	_, span := otel.Tracer("repo", trace.WithInstrumentationVersion("12"), trace.WithInstrumentationAttributes(attribute.String("a", "d"))).Start(ctx, "Readiness")
	defer span.End()
	return nil
}

func (r *HealthzRepo) Liveness(ctx context.Context) error {
	logger := r.logger.With("method", "Liveness", "ctx", ctx)
	logger.Debug("repo Liveness")
	return nil
}
