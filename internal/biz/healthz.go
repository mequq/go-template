package biz

import (
	"context"
	"log/slog"

	"application/pkg/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type HealthzBiz struct {
	repo   HealthzRepoInterface
	logger *slog.Logger
	tracer trace.Tracer
}

// New Usecase
func NewHealthzBiz(repo HealthzRepoInterface, logger *slog.Logger) *HealthzBiz {
	return &HealthzBiz{
		repo:   repo,
		logger: logger.With("layer", "HealthzBiz"),
		tracer: otel.Tracer("HealthzUseCase"),
	}
}

func (uc *HealthzBiz) Readiness(ctx context.Context) error {
	ctx, span := uc.tracer.Start(ctx, "ReadinessUsecase", trace.WithAttributes(attribute.Bool("readiness", true)))
	logger := uc.logger.With("method", "Readiness", "ctx", utils.GetLoggerContext(ctx))
	logger.DebugContext(ctx, "Readiness")
	defer span.End()
	span.AddEvent("Readiness", trace.WithStackTrace(true))
	return uc.repo.Readiness(ctx)
}

func (uc *HealthzBiz) Liveness(ctx context.Context) error {
	logger := uc.logger.With("method", "LivenessUsecase", "ctx", utils.GetLoggerContext(ctx))
	ctx, sp := uc.tracer.Start(ctx, "Liveness", trace.WithAttributes(attribute.Bool("liveness", true)))
	defer sp.End()
	sp.AddEvent("Liveness", trace.WithStackTrace(true))
	logger.Debug("Liveness")
	return uc.repo.Liveness(ctx)
}
