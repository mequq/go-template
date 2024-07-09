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
func NewServiceList(healthzSvc *HealthzService) []ServiceInterface {
	return []ServiceInterface{
		healthzSvc,
	}
}

// Service Interface
type ServiceInterface interface {
	RegisterMuxRouter(mux *http.ServeMux)
}
