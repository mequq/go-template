package biz

import "context"

type RepositoryHealthzer interface {
	Readiness(ctx context.Context) error
	Liveness(ctx context.Context) error
}

type UsecaseHealthzer interface {
	Readiness(ctx context.Context) error
	Liveness(ctx context.Context) error
}
