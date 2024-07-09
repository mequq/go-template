package memory

import (
	"application/internal/datasource/sample_entitiy"
	"application/internal/entity"
	"context"
	"errors"
	"sync"
)

type sampleEntity struct {
	mu      sync.Mutex
	entries []*entity.SampleEntity
	nextID  uint64
}

func NewSampleEntity() sample_entitiy.DataSource {
	return &sampleEntity{
		entries: make([]*entity.SampleEntity, 0),
		nextID:  1,
	}
}

func (s *sampleEntity) Create(ctx context.Context, sampleEntity *entity.SampleEntity) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sampleEntity.ID = s.nextID
	s.nextID++
	s.entries = append(s.entries, sampleEntity)
	return sampleEntity.ID, nil
}

func (s *sampleEntity) Update(ctx context.Context, sampleEntity *entity.SampleEntity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, entry := range s.entries {
		if entry.ID == sampleEntity.ID {
			s.entries[i] = sampleEntity
			return nil
		}
	}
	return errors.New("entity not found")
}

func (s *sampleEntity) List(ctx context.Context) ([]*entity.SampleEntity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]*entity.SampleEntity(nil), s.entries...), nil
}

func (s *sampleEntity) Delete(ctx context.Context, id uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, entry := range s.entries {
		if entry.ID == id {
			s.entries = append(s.entries[:i], s.entries[i+1:]...)
			return nil
		}
	}
	return errors.New("entity not found")
}