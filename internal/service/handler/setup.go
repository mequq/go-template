package handler

import (
	"application/internal/service"

	"github.com/google/wire"
)

var HandlerProviderSet = wire.NewSet(

	NewMuxHealthzHandler,
	NewServiceList,
	NewMuxBuildingHandler,
	NewMuxMovieHandler,
	NewMuxSeriesHandler,
	NewMuxChannelHandler,
)

// New ServiceList
func NewServiceList(
	healthzSvc *HealthzHandler,

	buildingSvc *BuildingHandler,
	movieSvc *MovieHandler,
	seriesSvc *SeriesHandler,
	channelSvc *ChannelHandler,
) []service.Handler {
	return []service.Handler{
		healthzSvc,
		buildingSvc,
		movieSvc,
		seriesSvc,
		channelSvc,
	}
}
