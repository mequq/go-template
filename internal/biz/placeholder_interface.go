package biz

import (
	"application/internal/entity"
	"context"

	"github.com/google/uuid"
)

type UsecasePlaceholder interface {
	Get(ctx context.Context, id uuid.UUID) (entity.Placeholder, error)
	List(ctx context.Context) ([]entity.Placeholder, error)
	Create(ctx context.Context, name string) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, name string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RepositoryPlaceholder interface {
	Get(ctx context.Context, id uuid.UUID) (entity.Placeholder, error)
	List(ctx context.Context) ([]entity.Placeholder, error)
	Create(ctx context.Context, name string) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, name string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
