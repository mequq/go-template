package handler

import (
	"application/internal/service"
	"application/internal/service/handler/building"
	"application/internal/service/handler/channel"
	"application/internal/service/handler/healthz"
	"application/internal/service/handler/movie"
	"application/internal/service/handler/sampleentity"
	"application/internal/service/handler/series"

	"github.com/google/wire"
)

var HandlerProviderSet = wire.NewSet(

	healthz.NewMuxHealthzHandler,
	sampleentity.NewSampleEntityHandler,
	NewServiceList,
	building.NewMuxBuildingHandler,
	movie.NewMuxMovieHandler,
	series.NewMuxSeriesHandler,
	channel.NewMuxChannelHandler,
)

// New ServiceList
func NewServiceList(
	healthzSvc *healthz.HealthzHandler,
	sampleEntityHandler *sampleentity.SampleEntityHandler,
	buildingSvc *building.BuildingHandler,
	movieSvc *movie.MovieHandler,
	seriesSvc *series.SeriesHandler,
	channelSvc *channel.ChannelHandler,
) []service.Handler {
	return []service.Handler{
		healthzSvc,
		sampleEntityHandler,
		buildingSvc,
		movieSvc,
		seriesSvc,
		channelSvc,
	}
}
