package gorm

import (
	"application/internal/entity"
	"gorm.io/gorm"
)

type SampleEntity struct {
	conn *gorm.DB
}

func NewSampleEntityRepository(db *gorm.DB) *SampleEntity {
	return &SampleEntity{conn: db}
}

func (s SampleEntity) Create(sample *entity.SampleEntity) (id uint64, err error) {
	result := s.conn.Create(sample)
	if result.Error != nil {
		return 0, result.Error
	}
	return sample.ID, nil
}

func (s SampleEntity) Update(sample *entity.SampleEntity) error {
	result := s.conn.Save(sample)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s SampleEntity) List() ([]*entity.SampleEntity, error) {
	var samples []*entity.SampleEntity
	result := s.conn.Find(&samples)
	if result.Error != nil {
		return nil, result.Error
	}
	return samples, nil
}

func (s SampleEntity) Delete(id uint64) error {
	result := s.conn.Delete(&entity.SampleEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
