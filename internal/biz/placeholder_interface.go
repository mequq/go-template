package biz

import (
	"application/internal/entity"
	"context"
)

type UsecasePlaceholder interface {
	Get(ctx context.Context, id int64) (entity.Placeholder, error)
	List(ctx context.Context) ([]entity.Placeholder, error)
	Create(ctx context.Context, name string) error
	Update(ctx context.Context, id int64, name string) error
	Delete(ctx context.Context, id int64) error
}

type RepositoryPlaceholder interface {
	Get(ctx context.Context, id int64) (entity.Placeholder, error)
	List(ctx context.Context) ([]entity.Placeholder, error)
	Create(ctx context.Context, name string) error
	Update(ctx context.Context, id int64, name string) error
	Delete(ctx context.Context, id int64) error
}
