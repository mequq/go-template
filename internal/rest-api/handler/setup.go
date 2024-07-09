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
func NewServiceList(healthzSvc *HealthzHandler.sampleEntityHandler) []HandlerInterface {
	return []HandlerInterface{
		healthzSvc,
	}
}

// Service Interface
type HandlerInterface interface {
	RegisterMuxRouter(mux *http.ServeMux)
}
