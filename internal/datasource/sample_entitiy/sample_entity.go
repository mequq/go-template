package sample_entitiy

import (
	"application/internal/entity"
	"context"
	"errors"
)

var (
	ErrNotFound     = errors.New("simple entity not found")
	ErrAlreadyExist = errors.New("simple entity already exist")
)

type DataSource interface {
	Create(ctx context.Context, sampleEntity *entity.SampleEntity) (id uint64, err error)
	Update(ctx context.Context, sampleEntity *entity.SampleEntity) error
	List(ctx context.Context) ([]*entity.SampleEntity, error)
	Delete(ctx context.Context, id uint64) error
}
