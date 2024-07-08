package sample_entity

import (
	"application/internal/datasource/sample_entitiy"
	"application/internal/entity"
	"application/pkg/utils"
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"log/slog"
)

type SampleEntity interface {
	Create(ctx context.Context, entity *entity.SampleEntity) (*entity.SampleEntity, error)
	Update(ctx context.Context, entity *entity.SampleEntity) error
	List(ctx context.Context, entity *entity.SampleEntity) ([]*entity.SampleEntity, error)
	Delete(ctx context.Context, entity *entity.SampleEntity) error
}

type sampleEntity struct {
	logger       *slog.Logger
	seDataSource sample_entitiy.DataSource
}

func NewSampleEntity(seDataSource sample_entitiy.DataSource, logger *slog.Logger) SampleEntity {
	return &sampleEntity{
		seDataSource: seDataSource,
		logger:       logger.With("layer", "SampleEntityBiz"),
	}
}

func (s *sampleEntity) Create(ctx context.Context, entity *entity.SampleEntity) (*entity.SampleEntity, error) {
	s.logger.With("method", "Create", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.Create")
	defer span.End()
	id, err := s.seDataSource.Create(ctx, entity)
	if err != nil {
		if errors.Is(err, sample_entitiy.ErrAlreadyExist) {
			return nil, err
		}
		s.logger.Error("error creating sample entity", "error", err.Error())
		return nil, err
	}

	entity.ID = id
	return entity, nil
}

func (s *sampleEntity) Update(ctx context.Context, entity *entity.SampleEntity) error {
	s.logger.With("method", "Update", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.Update")
	defer span.End()
	if err := s.seDataSource.Update(ctx, entity); err != nil {
		if errors.Is(err, sample_entitiy.ErrNotFound) || errors.Is(err, sample_entitiy.ErrAlreadyExist) {
			return err
		}
		s.logger.Error("error creating sample entity", "error", err.Error())
		return err
	}
	return nil
}

func (s *sampleEntity) List(ctx context.Context, entity *entity.SampleEntity) ([]*entity.SampleEntity, error) {
	s.logger.With("method", "List", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.List")
	defer span.End()
	id, err := s.seDataSource.Create(ctx, entity)
	if err != nil {
		s.logger.Error("error creating sample entity", "error", err.Error())
		return nil, err
	}
	entity.ID = id

	return nil, err
}

func (s *sampleEntity) Delete(ctx context.Context, entity *entity.SampleEntity) error {
	s.logger.With("method", "Delete", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.Delete")
	defer span.End()
	if err := s.seDataSource.Update(ctx, entity); err != nil {
		if errors.Is(err, sample_entitiy.ErrNotFound) {
			return err
		}
		s.logger.Error("error creating sample entity", "error", err.Error())
		return err
	}
	return nil
}
