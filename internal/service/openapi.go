package service

import (
	"log/slog"

	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
)

type OAPI interface {
	GenerateOperationContext(method string, path string) openapi.OperationContext
	GetYamlData() ([]byte, error)
	GetJsonData() ([]byte, error)
	Register(method string, path string, handler OPHandlerFunc, opt ...RegisterOPT)
}

type RegisterOPT func(o openapi.OperationContext)

// with tags for operation
func WithTags(tags ...string) RegisterOPT {
	return func(o openapi.OperationContext) {
		o.SetTags(tags...)
	}
}

type OAPIMP struct {
	reflector *openapi3.Reflector
	logger    *slog.Logger
}

func NewOAPI(reflector *openapi3.Reflector, logger *slog.Logger) OAPI {
	return &OAPIMP{
		reflector: reflector,
		logger:    logger,
	}
}

func (s *OAPIMP) GenerateOperationContext(method, path string) openapi.OperationContext {
	op, err := s.reflector.NewOperationContext(method, path)
	if err != nil {
		s.logger.Error("GenerateOperationContext", "error", err)
		panic(err)
	}
	return op
}

// GetYamlData return yaml data
func (s *OAPIMP) GetYamlData() ([]byte, error) {
	return s.reflector.Spec.MarshalYAML()
}

// GetJsonData return json data
func (s *OAPIMP) GetJsonData() ([]byte, error) {
	return s.reflector.Spec.MarshalJSON()
}

func (s *OAPIMP) Register(method, path string, handler OPHandlerFunc, opt ...RegisterOPT) {
	op := s.GenerateOperationContext(method, path)

	for _, o := range opt {
		if o == nil {
			continue
		}
		o(op)
	}

	handler(op)
	err := s.reflector.AddOperation(op)
	if err != nil {
		s.logger.Error("Register", "error", err)
	}
}

type OPHandlerFunc func(openapi.OperationContext)
