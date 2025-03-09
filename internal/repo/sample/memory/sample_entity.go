package samplememrepo

import (
	"context"
	"sync"

	sampleusecasev1 "application/internal/biz/sample"
	sampleentity "application/internal/entity/sample"
)

type sampleEntity struct {
	mu      sync.Mutex
	entries []*sampleentity.Sample
	nextID  uint64
}

func NewSampleEntity() sampleusecasev1.SampleEntityRepoInterface {
	return &sampleEntity{
		entries: make([]*sampleentity.Sample, 0),
		nextID:  1,
	}
}

func (s *sampleEntity) Create(_ context.Context, sampleEntity *sampleentity.Sample) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sampleEntity.ID = s.nextID
	s.nextID++
	s.entries = append(s.entries, sampleEntity)
	return sampleEntity.ID, nil
}

func (s *sampleEntity) Update(_ context.Context, sampleEntity *sampleentity.Sample) error {
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

func (s *sampleEntity) List(_ context.Context) ([]*sampleentity.Sample, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]*sampleentity.Sample(nil), s.entries...), nil
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
