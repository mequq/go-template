package sample_entitiy

import (
	"application/internal/entity"
	"context"
)

type DataSource interface {
	Create(ctx context.Context, sampleEntity *entity.SampleEntity) (id uint64, err error)
	Update(ctx context.Context, sampleEntity *entity.SampleEntity) error
	List(ctx context.Context) []*entity.SampleEntity
	Delete(ctx context.Context, id uint64) error
}
