package repo

import (
	"application/internal/datasource"
	"application/pkg/utils"
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type healthz struct {
	logger *slog.Logger
	memDB  *datasource.InmemoryDB
	tracer trace.Tracer
}

func NewHealthzDS(logger *slog.Logger, memDB *datasource.InmemoryDB) *healthz {
	return &healthz{
		logger: logger.With("layer", "healthzRepo"),
		memDB:  memDB,
		tracer: otel.Tracer("healthzRepo"),
	}
}

func (r *healthz) Readiness(ctx context.Context) error {
	logger := r.logger.With("method", "readiness")

	ctx, span := r.tracer.Start(
		ctx,
		"readiness",
		trace.WithAttributes(attribute.Bool("readiness", true)),
	)
	defer span.End()

	logger.DebugContext(ctx, "repo readiness")

	if err := r.memDB.DB.PingContext(ctx); err != nil {
		span.SetStatus(codes.Error, "PingContext failed")
		span.RecordError(err)
		r.logger.ErrorContext(ctx, "PingContext failed", "error", err)

		return err
	}

	return nil
}

func (r *healthz) Liveness(ctx context.Context) error {
	logger := r.logger.With("method", "Liveness", "ctx", utils.GetLoggerContext(ctx))

	_, span := r.tracer.Start(
		ctx,
		"Liveness",
		trace.WithAttributes(attribute.Bool("liveness", true)),
	)
	defer span.End()

	logger.Debug("repo Liveness")

	return nil
}
