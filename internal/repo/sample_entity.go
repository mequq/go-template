package repo

import (
	"context"
	"sync"

	"application/internal/biz"
	sampleusecasev1 "application/internal/biz"
	entity "application/internal/entity"
)

type sampleEntity struct {
	mu      sync.Mutex
	entries []*entity.Sample
	nextID  uint64
}

func NewSampleEntity() biz.SampleEntityRepoInterface {
	return &sampleEntity{
		entries: make([]*entity.Sample, 0),
		nextID:  1,
	}
}

func (s *sampleEntity) Create(_ context.Context, sampleEntity *entity.Sample) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sampleEntity.ID = s.nextID
	s.nextID++
	s.entries = append(s.entries, sampleEntity)
	return sampleEntity.ID, nil
}

func (s *sampleEntity) Update(_ context.Context, sampleEntity *entity.Sample) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, entry := range s.entries {
		if entry.ID == sampleEntity.ID {
			s.entries[i] = sampleEntity
			return nil
		}
	}
	return sampleusecasev1.ErrNotFound
}

func (s *sampleEntity) List(_ context.Context) ([]*entity.Sample, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]*entity.Sample(nil), s.entries...), nil
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
	return sampleusecasev1.ErrNotFound
}
