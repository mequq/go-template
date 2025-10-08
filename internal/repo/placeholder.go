package repo

import (
	"application/internal/biz"
	"application/internal/datasource"
	"application/internal/entity"
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type placeholder struct {
	logger *slog.Logger
	tracer trace.Tracer
	sqlDb  *datasource.InmemoryDB
}

var _ biz.RepositoryPlaceholder = (*placeholder)(nil)

func NewPlaceholder(logger *slog.Logger, sqlDb *datasource.InmemoryDB) *placeholder {
	return &placeholder{
		logger: logger.With("layer", "Placeholder"),
		tracer: otel.Tracer("HealthzUseCase"),
		sqlDb:  sqlDb,
	}
}

// Create implements biz.RepositoryPlaceholder.
func (r *placeholder) Create(ctx context.Context, name string) error {
	return nil
}

// List implements biz.RepositoryPlaceholder.
func (r *placeholder) List(ctx context.Context) ([]entity.Placeholder, error) {
	return []entity.Placeholder{{ID: 1, Name: "placeholder"}}, nil
}

// Get implements biz.RepositoryPlaceholder.
func (r *placeholder) Get(ctx context.Context, id int64) (entity.Placeholder, error) {
	return entity.Placeholder{ID: id, Name: "placeholder"}, nil
}

// Update implements biz.RepositoryPlaceholder.
func (r *placeholder) Update(ctx context.Context, id int64, name string) error {
	return nil
}

// Delete implements biz.RepositoryPlaceholder.
func (r *placeholder) Delete(ctx context.Context, id int64) error {
	return nil
}
