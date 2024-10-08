package healthzusecase

import (
	"context"
	"log/slog"

	"application/pkg/utils"

	"go.opentelemetry.io/otel"
)

type HealthzBiz struct {
	repo   HealthzRepoInterface
	logger *slog.Logger
}

// New Usecase
func NewHealthzBiz(repo HealthzRepoInterface, logger *slog.Logger) HealthzUseCaseInterface {
	return &HealthzBiz{
		repo:   repo,
		logger: logger.With("layer", "HealthzBiz"),
	}
}

func (uc *HealthzBiz) Readiness(ctx context.Context) error {
	ctx, span := otel.Tracer("usecase").Start(ctx, "rediness")
	defer span.End()
	return uc.repo.Readiness(ctx)
}

func (uc *HealthzBiz) Liveness(ctx context.Context) error {
	logger := uc.logger.With("method", "Liveness", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("Liveness")
	return uc.repo.Liveness(ctx)
}
