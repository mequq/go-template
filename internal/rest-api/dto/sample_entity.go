package dto

import (
	"application/internal/entity"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type SampleEntityRequest struct {
	ID   uint64 `json:"id"`
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

func (req *SampleEntityRequest) ToEntity() *entity.SampleEntity {
	return &entity.SampleEntity{
		ID:   req.ID,
		Name: req.Name,
		Text: req.Text,
	}
}
