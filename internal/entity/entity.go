package sampleentity

type Sample struct {
	ID   uint64
	Name string
	Text string
}

func NewSampleEntity(name, text string) *Sample {
	return &Sample{Name: name, Text: text}
}
