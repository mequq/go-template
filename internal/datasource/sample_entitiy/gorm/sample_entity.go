package gorm

import (
	"application/internal/datasource/sample_entitiy"
	"application/internal/entity"
	"context"
	"gorm.io/gorm"
)

type sampleEntity struct {
	conn *gorm.DB
}

func NewSampleEntity(db *gorm.DB) sample_entitiy.DataSource {
	return &sampleEntity{conn: db}
}

func (s sampleEntity) Create(ctx context.Context, sampleEntity *entity.SampleEntity) (id uint64, err error) {
	result := s.conn.WithContext(ctx).Create(sampleEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	return sampleEntity.ID, nil
}

func (s sampleEntity) Update(ctx context.Context, sampleEntity *entity.SampleEntity) error {
	result := s.conn.WithContext(ctx).Save(sampleEntity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s sampleEntity) List(ctx context.Context) ([]*entity.SampleEntity, error) {
	var samples []*entity.SampleEntity
	result := s.conn.WithContext(ctx).Find(&samples)
	if result.Error != nil {
		return nil, result.Error
	}
	return samples, nil
}

func (s sampleEntity) Delete(ctx context.Context, id uint64) error {
	result := s.conn.WithContext(ctx).Delete(&entity.SampleEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
