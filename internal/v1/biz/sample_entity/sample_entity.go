package sample_entity

import (
	"context"
	"errors"
	"log/slog"

	"application/internal/v1/datasource/sample_entitiy"
	ent "application/internal/v1/entity"
	"application/pkg/utils"

	"go.opentelemetry.io/otel"
)

type SampleEntity interface {
	Create(ctx context.Context, sampEnt *ent.SampleEntity) (*ent.SampleEntity, error)
	Update(ctx context.Context, sampEnt *ent.SampleEntity) error
	List(ctx context.Context) ([]*ent.SampleEntity, error)
	Delete(ctx context.Context, id uint64) error
}

var _ SampleEntity = (*sampleEntity)(nil)

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

func (s *sampleEntity) Create(ctx context.Context, sampEnt *ent.SampleEntity) (*ent.SampleEntity, error) {
	s.logger.With("method", "Create", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.Create")
	defer span.End()
	id, err := s.seDataSource.Create(ctx, sampEnt)
	if err != nil {
		if errors.Is(err, sample_entitiy.ErrAlreadyExist) {
			return nil, err
		}
		s.logger.Error("error creating sample entity", "error", err.Error())
		return nil, err
	}

	sampEnt.ID = id
	return sampEnt, nil
}

func (s *sampleEntity) Update(ctx context.Context, sampEnt *ent.SampleEntity) error {
	s.logger.With("method", "Update", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.Update")
	defer span.End()
	if err := s.seDataSource.Update(ctx, sampEnt); err != nil {
		if errors.Is(err, sample_entitiy.ErrNotFound) || errors.Is(err, sample_entitiy.ErrAlreadyExist) {
			return err
		}
		s.logger.Error("error creating sample entity", "error", err.Error())
		return err
	}
	return nil
}

func (s *sampleEntity) List(ctx context.Context) ([]*ent.SampleEntity, error) {
	s.logger.With("method", "List", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.List")
	defer span.End()
	es, err := s.seDataSource.List(ctx)
	if err != nil {
		s.logger.Error("error creating sample entity", "error", err.Error())
		return nil, err
	}

	return es, err
}

func (s *sampleEntity) Delete(ctx context.Context, id uint64) error {
	s.logger.With("method", "Delete", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.Delete")
	defer span.End()
	if err := s.seDataSource.Delete(ctx, id); err != nil {
		if errors.Is(err, sample_entitiy.ErrNotFound) {
			return err
		}
		s.logger.Error("error creating sample entity", "error", err.Error())
		return err
	}
	return nil
}
