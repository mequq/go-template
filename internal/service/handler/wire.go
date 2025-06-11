package handler

import (
	"application/internal/service"

	"github.com/google/wire"
)

var HandlerProviderSet = wire.NewSet(

	NewServiceList,
	NewMuxHealthzHandler,
)

// New ServiceList
func NewServiceList(
	healthzSvc *HealthzHandler,

) []service.Handler {
	return []service.Handler{
		healthzSvc,
	}
}
