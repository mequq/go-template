package service

import (
	"net/http"

	"github.com/google/wire"
)

var ServiceProviderSet = wire.NewSet(
	NewMuxHealthzService,
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
	// RegisterRoutes(mux *mux.Router)
	RegisterMuxRouter(mux *http.ServeMux)
}
