package biz

import (
	"app/config"
	"context"

	"github.com/rs/zerolog"
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
func (u *HealthzUsecase) Readiness(ctx context.Context) error {
	return u.repo.Readiness(ctx)
}

// Liveness checks the liveness of the service.
func (u *HealthzUsecase) Liveness(ctx context.Context) error {
	return u.repo.Liveness(ctx)
}
