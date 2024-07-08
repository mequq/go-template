package gorm

import (
	"application/internal/entity"
	"context"
	"gorm.io/gorm"
)

type SampleEntity struct {
	conn *gorm.DB
}

func NewSampleEntityRepository(db *gorm.DB) *SampleEntity {
	return &SampleEntity{conn: db}
}

func (s SampleEntity) Create(ctx context.Context, sampleEntity *entity.SampleEntity) (id uint64, err error) {
	result := s.conn.WithContext(ctx).Create(sampleEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	return sampleEntity.ID, nil
}

func (s SampleEntity) Update(ctx context.Context, sampleEntity *entity.SampleEntity) error {
	result := s.conn.WithContext(ctx).Save(sampleEntity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s SampleEntity) List(ctx context.Context) ([]*entity.SampleEntity, error) {
	var samples []*entity.SampleEntity
	result := s.conn.WithContext(ctx).Find(&samples)
	if result.Error != nil {
		return nil, result.Error
	}
	return samples, nil
}

func (s SampleEntity) Delete(ctx context.Context, id uint64) error {
	result := s.conn.WithContext(ctx).Delete(&entity.SampleEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
