package healthzusecase

import "context"

type HealthzRepoInterface interface {
	Readiness(ctx context.Context) error
	Liveness(ctx context.Context) error
}

type HealthzUseCaseInterface interface {
	Readiness(ctx context.Context) error
	Liveness(ctx context.Context) error
}
