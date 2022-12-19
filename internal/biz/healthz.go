package biz

import (
	"app/config"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
)

// HealthzRepo is the repository interface for healthz.
type HealthzRepo interface {
	// Readyness checks the readyness of the service.
	Readiness(ctx context.Context) error
	// Liveness checks the liveness of the service.
	Liveness(ctx context.Context) error
}

// HealthzUsecase is usecase
type HealthzUsecase struct {
	repo HealthzRepo
	conf *config.Config
}

// NewHealthzUsecase creates a new healthz usecase.
func NewHealthzUsecase(repo HealthzRepo, conf *config.Config, logger zerolog.Logger) *HealthzUsecase {
	return &HealthzUsecase{
		repo: repo,
		conf: conf,
	}
}

// Readyness checks the readyness of the service.
func (u *HealthzUsecase) Readiness(c *gin.Context) error {
	ctx, span := otel.Tracer("healthz").Start(c.Request.Context(), "readiness-check")
	defer span.End()
	return u.repo.Readiness(ctx)
}

// Liveness checks the liveness of the service.
func (u *HealthzUsecase) Liveness(ctx context.Context) error {
	return u.repo.Liveness(ctx)
}
