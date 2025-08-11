package biz

import (
	"context"
	"errors"
	"log/slog"

	sampleentity "application/internal/entity"
	"application/pkg/utils"

	"go.opentelemetry.io/otel"
)

type sampleEntity struct {
	logger *slog.Logger
	repo   RepositorySampleer
}

func NewSampleEntity(repo RepositorySampleer, logger *slog.Logger) *sampleEntity {
	return &sampleEntity{
		repo:   repo,
		logger: logger.With("layer", "SampleEntityBiz"),
	}
}

func (uc *sampleEntity) Create(ctx context.Context, sampEnt *sampleentity.Sample) (*sampleentity.Sample, error) {
	logger := uc.logger.With("method", "Create", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.Create")
	defer span.End()
	id, err := uc.repo.Create(ctx, sampEnt)
	if err != nil {
		if errors.Is(err, ErrAlreadyExist) {
			return nil, err
		}
		logger.Error("error creating sample entity", "error", err.Error())
		return nil, err
	}

	sampEnt.ID = id
	return sampEnt, nil
}

func (uc *sampleEntity) Update(ctx context.Context, sampEnt *sampleentity.Sample) error {
	logger := uc.logger.With("method", "Update", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.Update")
	defer span.End()
	if err := uc.repo.Update(ctx, sampEnt); err != nil {
		if errors.Is(err, ErrNotFound) || errors.Is(err, ErrAlreadyExist) {
			return err
		}
		logger.Error("error creating sample entity", "error", err.Error())
		return err
	}
	return nil
}

func (uc *sampleEntity) List(ctx context.Context) ([]*sampleentity.Sample, error) {
	logger := uc.logger.With("method", "List", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.List")
	defer span.End()
	es, err := uc.repo.List(ctx)
	if err != nil {
		logger.Error("error creating sample entity", "error", err.Error())
		return nil, err
	}

	return es, err
}

func (s *sampleEntity) Delete(ctx context.Context, id uint64) error {
	s.logger.With("method", "Delete", "ctx", utils.GetLoggerContext(ctx))

	ctx, span := otel.Tracer("biz").Start(ctx, "SampleEntity.Delete")
	defer span.End()
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return err
		}
		s.logger.Error("error creating sample entity", "error", err.Error())
		return err
	}
	return nil
}
