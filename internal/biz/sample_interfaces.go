package biz

import (
	"context"
	"errors"

	sampleentity "application/internal/entity"
)

var (
	ErrNotFound     = errors.New("simple entity not found")
	ErrAlreadyExist = errors.New("simple entity already exist")
)

type RepositorySampleer interface {
	Create(ctx context.Context, sampleEntity *sampleentity.Sample) (id uint64, err error)
	Update(ctx context.Context, sampleEntity *sampleentity.Sample) error
	List(ctx context.Context) ([]*sampleentity.Sample, error)
	Delete(ctx context.Context, id uint64) error
}

type UsecaseSampleer interface {
	Create(ctx context.Context, sampEnt *sampleentity.Sample) (*sampleentity.Sample, error)
	Update(ctx context.Context, sampEnt *sampleentity.Sample) error
	List(ctx context.Context) ([]*sampleentity.Sample, error)
	Delete(ctx context.Context, id uint64) error
}
