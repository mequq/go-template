package biz

import "context"

type UsecaseHealthzer interface {
	Readiness(ctx context.Context) error
	Liveness(ctx context.Context) error
}
