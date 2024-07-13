package entity

type SampleEntity struct {
	ID   uint64
	Name string
	Text string
}

func NewSampleEntity(name, text string) *SampleEntity {
	return &SampleEntity{Name: name, Text: text}
}
