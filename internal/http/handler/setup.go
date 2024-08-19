package handler

import (
	"net/http"

	"github.com/google/wire"
)

var HandlerProviderSet = wire.NewSet(
	NewMuxHealthzHandler,
	NewSampleEntityHandler,
	NewServiceList,
)

// New ServiceList
func NewServiceList(healthzSvc *HealthzHandler, sampleEntityHandler *SampleEntityHandler) []Handler {
	return []Handler{
		healthzSvc,
		sampleEntityHandler,
	}
}

// Service Interface
type Handler interface {
	RegisterMuxRouter(mux *http.ServeMux)
}
