package memory

import (
	biz "application/internal/biz/healthz"
	"application/pkg/utils"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

type HealthzDS struct {
	logger *slog.Logger
}

func NewHealthzDS(logger *slog.Logger) biz.HealthzRepoInterface {
	return &HealthzDS{
		logger: logger.With("layer", "HealthzDS"),
	}
}

func (r *HealthzDS) Readiness(ctx context.Context) error {
	_, span := otel.Tracer("repo", trace.WithInstrumentationVersion("12"), trace.WithInstrumentationAttributes(attribute.String("a", "d"))).Start(ctx, "Readiness")
	defer span.End()
	return nil
}

func (r *HealthzDS) Liveness(ctx context.Context) error {
	logger := r.logger.With("method", "Liveness", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("repo Liveness")
	return nil
}