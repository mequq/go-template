package service

import (
	"github.com/google/wire"
	"github.com/gorilla/mux"
)

var ServiceProviderSet = wire.NewSet(
	NewGorilaMuxHealthzService,
	NewServiceList,
)

// New ServiceList
func NewServiceList(healthzSvc *GorilaMuxHealthzService) []ServiceInterface {
	return []ServiceInterface{
		healthzSvc,
	}
}

// Service Interface
type ServiceInterface interface {
	RegisterRoutes(mux *mux.Router)
}
