package dto

import (
	"encoding/json"
	"net/http"

	sampleentity "application/internal/entity"

	"github.com/go-playground/validator/v10"
)

type SampleEntityRequest struct {
	Name string `json:"name" validate:"required"`
	Text string `json:"text" validate:"required"`
}

func (r *SampleEntityRequest) FromRequest(req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	return nil
}

func (r *SampleEntityRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(r)
}

func (r *SampleEntityRequest) ToEntity() *sampleentity.Sample {
	return &sampleentity.Sample{
		Name: r.Name,
		Text: r.Text,
	}
}

type SampleEntityResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Text string `json:"text"`
}

func FromEntity(e *sampleentity.Sample) SampleEntityResponse {
	return SampleEntityResponse{
		ID:   e.ID,
		Name: e.Name,
		Text: e.Text,
	}
}

func SampleEntityListResponses(entities []*sampleentity.Sample) []SampleEntityResponse {
	responses := make([]SampleEntityResponse, len(entities))
	for i, e := range entities {
		responses[i] = FromEntity(e)
	}

	return responses
}
