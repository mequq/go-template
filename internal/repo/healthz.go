package repo

import (
	"context"
	"log/slog"

	healthzusecase "application/internal/biz"
	"application/internal/datasource"
	"application/pkg/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type HealthzDS struct {
	logger *slog.Logger
	memDB  *datasource.InmemoryDB
	tracer trace.Tracer
}

func NewHealthzDS(logger *slog.Logger, memDB *datasource.InmemoryDB) healthzusecase.HealthzRepoInterface {
	return &HealthzDS{
		logger: logger.With("layer", "healthzRepo"),
		memDB:  memDB,
		tracer: otel.Tracer("healthzRepo"),
	}
}

func (r *HealthzDS) Readiness(ctx context.Context) error {
	ctx, span := r.tracer.Start(ctx, "readiness", trace.WithAttributes(attribute.Bool("readiness", true)))
	defer span.End()
	logger := r.logger.With("method", "readiness")
	logger.DebugContext(ctx, "repo readiness")
	if err := r.memDB.DB.PingContext(ctx); err != nil {
		span.SetStatus(codes.Error, "PingContext failed")
		span.RecordError(err)
		r.logger.ErrorContext(ctx, "PingContext failed", "error", err)
		return err
	}
	return nil
}

func (r *HealthzDS) Liveness(ctx context.Context) error {
	ctx, span := r.tracer.Start(ctx, "Liveness", trace.WithAttributes(attribute.Bool("liveness", true)))
	defer span.End()

	logger := r.logger.With("method", "Liveness", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("repo Liveness")
	return nil
}
