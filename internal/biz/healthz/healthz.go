package healthzusecase

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
func NewHealthzBiz(repo HealthzRepoInterface, logger *slog.Logger) HealthzUseCaseInterface {
	return &HealthzBiz{
		repo:   repo,
		logger: logger.With("layer", "HealthzBiz"),
		tracer: otel.Tracer("usecase"),
	}
}

func (uc *HealthzBiz) Readiness(ctx context.Context) error {
	ctx, span := otel.Tracer("usecase").Start(ctx, "rediness")
	defer span.End()
	span.AddEvent("Readiness", trace.WithStackTrace(true))
	return uc.repo.Readiness(ctx)
}

func (uc *HealthzBiz) Liveness(ctx context.Context) error {
	logger := uc.logger.With("method", "Liveness", "ctx", utils.GetLoggerContext(ctx))
	ctx, sp := uc.tracer.Start(ctx, "Liveness", trace.WithAttributes(attribute.Bool("liveness", true)))
	defer sp.End()
	logger.Debug("Liveness")
	return uc.repo.Liveness(ctx)
}
