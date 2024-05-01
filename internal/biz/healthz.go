package biz

import (
	"context"

	"log/slog"

	// otelhttp "go.opentelemetry.io/exporter/otlp/otlphttp"
	"go.opentelemetry.io/otel"
)

type HealthzRepoInterface interface {
	Readiness(ctx context.Context) error
	Liveness(ctx context.Context) error
}

type HealthzUseCaseInterface interface {
	Readiness(ctx context.Context) error
	Liveness(ctx context.Context) error
}

type HealthzUseCase struct {
	repo   HealthzRepoInterface
	logger *slog.Logger
}

// New Usecase
func NewHealthzUseCase(repo HealthzRepoInterface, logger *slog.Logger) HealthzUseCaseInterface {
	return &HealthzUseCase{
		repo:   repo,
		logger: logger.With("layer", "HealthzUseCase"),
	}
}

func (uc *HealthzUseCase) Readiness(ctx context.Context) error {

	// logger := uc.logger.With("method", "Readiness", "ctx", LogContext(ctx))

	ctx, span := otel.Tracer("usecase").Start(ctx, "rediness")
	defer span.End()

	// client := otelhttp.New

	// zclient := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	// req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:9090", nil)
	// if err != nil {
	// 	logger.Error("error creating request", "error", err)
	// 	return err
	// }

	// _, err = client.Do(req)
	// if err != nil {
	// 	logger.Error("error doing request", "error", err)
	// 	return err
	// }

	// logger.Debug("Readiness")
	return uc.repo.Readiness(ctx)
}

func (uc *HealthzUseCase) Liveness(ctx context.Context) error {
	logger := uc.logger.With("method", "Liveness", "ctx", ctx)
	logger.Debug("Liveness")
	return uc.repo.Liveness(ctx)
}
