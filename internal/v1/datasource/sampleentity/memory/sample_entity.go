package memory

import (
	"context"
	"sync"

	se "application/internal/v1/datasource/sampleentity"
	"application/internal/v1/entity"
)

type sampleEntity struct {
	mu      sync.Mutex
	entries []*entity.SampleEntity
	nextID  uint64
}

func NewSampleEntity() se.DataSource {
	return &sampleEntity{
		entries: make([]*entity.SampleEntity, 0),
		nextID:  1,
	}
}

func (s *sampleEntity) Create(_ context.Context, sampleEntity *entity.SampleEntity) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sampleEntity.ID = s.nextID
	s.nextID++
	s.entries = append(s.entries, sampleEntity)
	return sampleEntity.ID, nil
}

func (s *sampleEntity) Update(_ context.Context, sampleEntity *entity.SampleEntity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, entry := range s.entries {
		if entry.ID == sampleEntity.ID {
			s.entries[i] = sampleEntity
			return nil
		}
	}
	return se.ErrNotFound
}

func (s *sampleEntity) List(_ context.Context) ([]*entity.SampleEntity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]*entity.SampleEntity(nil), s.entries...), nil
}

func (s *sampleEntity) Delete(_ context.Context, id uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, entry := range s.entries {
		if entry.ID == id {
			s.entries = append(s.entries[:i], s.entries[i+1:]...)
			return nil
		}
	}
	return se.ErrNotFound
}
