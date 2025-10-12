package repo

import (
	"application/internal/biz"
	"application/internal/datasource"
	"application/internal/entity"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type placeholder struct {
	logger *slog.Logger
	tracer trace.Tracer
	db     *datasource.PostgresDB
}

var _ biz.RepositoryPlaceholder = (*placeholder)(nil)

func NewPlaceholder(logger *slog.Logger, db *datasource.PostgresDB) *placeholder {
	return &placeholder{
		logger: logger.With("layer", "Placeholder"),
		tracer: otel.Tracer("HealthzUseCase"),
		db:     db,
	}
}

// Create implements biz.RepositoryPlaceholder.
func (r *placeholder) Create(ctx context.Context, name string) (uuid.UUID, error) {
	logger := r.logger.With("method", "Create")
	query := `INSERT INTO placeholder (name) VALUES ($1) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, name)

	if row.Err() != nil {
		logger.WarnContext(ctx, "failed to execute query", "error", row.Err())

		return uuid.Nil, row.Err()
	}

	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// List implements biz.RepositoryPlaceholder.
func (r *placeholder) List(ctx context.Context) ([]entity.Placeholder, error) {
	logger := r.logger.With("method", "List")
	query := `SELECT id, name FROM placeholder`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.WarnContext(ctx, "failed to execute query", "error", err)

		return nil, err
	}
	defer rows.Close()

	var placeholders []entity.Placeholder

	for rows.Next() {
		var p entity.Placeholder
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			logger.WarnContext(ctx, "failed to scan row", "error", err)

			continue
		}

		placeholders = append(placeholders, p)
	}

	if err := rows.Err(); err != nil {
		logger.WarnContext(ctx, "rows iteration error", "error", err)
	}

	return placeholders, nil
}

// Get implements biz.RepositoryPlaceholder.
func (r *placeholder) Get(ctx context.Context, id uuid.UUID) (entity.Placeholder, error) {
	logger := r.logger.With("method", "Get")
	query := `SELECT id, name FROM placeholder WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	if row.Err() != nil {
		logger.WarnContext(ctx, "failed to execute query", "error", row.Err())

		return entity.Placeholder{}, row.Err()
	}

	var p entity.Placeholder
	if err := row.Scan(&p.ID, &p.Name); err != nil {
		logger.WarnContext(ctx, "failed to scan row", "error", err)

		return entity.Placeholder{}, err
	}

	return p, nil
}

// Update implements biz.RepositoryPlaceholder.
func (r *placeholder) Update(ctx context.Context, id uuid.UUID, name string) error {
	logger := r.logger.With("method", "Update")
	query := `UPDATE placeholder SET name = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, name, id)
	if err != nil {
		logger.WarnContext(ctx, "failed to execute query", "error", err)

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WarnContext(ctx, "failed to get rows affected", "error", err)

		return err
	}

	if rowsAffected == 0 {
		logger.WarnContext(ctx, "no rows updated, placeholder not found", "id", id)

		return biz.ErrResourceNotFound
	}

	return nil
}

// Delete implements biz.RepositoryPlaceholder.
func (r *placeholder) Delete(ctx context.Context, id uuid.UUID) error {
	logger := r.logger.With("method", "Delete")
	query := `DELETE FROM placeholder WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.WarnContext(ctx, "failed to execute query", "error", err)

		return err
	}

	return nil
}
