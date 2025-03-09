package handler

import (
	"application/internal/service"

	"github.com/google/wire"
)

var HandlerProviderSet = wire.NewSet(

	NewMuxHealthzHandler,
	NewSampleEntityHandler,
	NewServiceList,
	NewMuxBuildingHandler,
)

// New ServiceList
func NewServiceList(healthzSvc *HealthzHandler, sampleEntityHandler *SampleEntityHandler, buildingSvc *BuildingHandler) []service.Handler {
	return []service.Handler{
		healthzSvc,
		sampleEntityHandler,
		buildingSvc,
	}
}
