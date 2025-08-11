package biz

import (
	"context"
	"log/slog"

	"application/pkg/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type healthz struct {
	repo   RepositoryHealthzer
	logger *slog.Logger
	tracer trace.Tracer
}

// New Usecase
func NewHealthz(repo RepositoryHealthzer, logger *slog.Logger) *healthz {
	return &healthz{
		repo:   repo,
		logger: logger.With("layer", "Healthz"),
		tracer: otel.Tracer("HealthzUseCase"),
	}
}

func (uc *healthz) Readiness(ctx context.Context) error {
	ctx, span := uc.tracer.Start(ctx, "ReadinessUsecase", trace.WithAttributes(attribute.Bool("readiness", true)))
	logger := uc.logger.With("method", "Readiness", "ctx", utils.GetLoggerContext(ctx))
	logger.DebugContext(ctx, "Readiness")
	defer span.End()
	span.AddEvent("Readiness", trace.WithStackTrace(true))
	return uc.repo.Readiness(ctx)
}

func (uc *healthz) Liveness(ctx context.Context) error {
	logger := uc.logger.With("method", "LivenessUsecase", "ctx", utils.GetLoggerContext(ctx))
	ctx, sp := uc.tracer.Start(ctx, "Liveness", trace.WithAttributes(attribute.Bool("liveness", true)))
	defer sp.End()
	sp.AddEvent("Liveness", trace.WithStackTrace(true))
	logger.Debug("Liveness")
	return uc.repo.Liveness(ctx)
}
