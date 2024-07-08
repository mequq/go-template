package sample_entity

import "context"

type SampleEntityBiz interface {
	Create(ctx context.Context, entity SampleEntity) error
}

type SampleEntity struct {
}

func NewSampleEntity() *SampleEntity {
	return &SampleEntity{}
}
