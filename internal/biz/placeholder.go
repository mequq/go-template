package biz

import (
	"application/internal/entity"
	"context"
	"log/slog"
)

type placeholder struct {
	logger          *slog.Logger
	placeholderRepo RepositoryPlaceholder
}

var _ UsecasePlaceholder = (*placeholder)(nil)

func NewPlaceholder(logger *slog.Logger, placeholderRepo RepositoryPlaceholder) *placeholder {
	return &placeholder{
		logger:          logger.With("layer", "Placeholder"),
		placeholderRepo: placeholderRepo,
	}
}

func (uc *placeholder) Get(ctx context.Context, id int64) (entity.Placeholder, error) {
	return uc.placeholderRepo.Get(ctx, id)
}

func (uc *placeholder) List(ctx context.Context) ([]entity.Placeholder, error) {
	return uc.placeholderRepo.List(ctx)
}

func (uc *placeholder) Create(ctx context.Context, name string) error {
	return uc.placeholderRepo.Create(ctx, name)
}

func (uc *placeholder) Update(ctx context.Context, id int64, name string) error {
	return uc.placeholderRepo.Update(ctx, id, name)
}

func (uc *placeholder) Delete(ctx context.Context, id int64) error {
	return uc.placeholderRepo.Delete(ctx, id)
}
