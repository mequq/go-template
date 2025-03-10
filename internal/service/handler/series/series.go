package series

import (
	"application/internal/service"
	"log/slog"
	"net/http"

	"github.com/swaggest/openapi-go"
)

type SeriesHandler struct {
	logger *slog.Logger
}

var _ service.Handler = (*SeriesHandler)(nil)
var _ service.OpenApiHandler = (*SeriesHandler)(nil)

func NewMuxSeriesHandler(logger *slog.Logger) *SeriesHandler {
	return &SeriesHandler{
		logger: logger.With("layer", "MuxSeriesService"),
	}
}

func (s *SeriesHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v3/content/v2/series", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/series/{series_id}", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/series/{series_id}/seasons", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/series/{series_id}/seasons/{season_id}", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/series/{series_id}/seasons/{season_id}/episodes", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/series/{series_id}/seasons/{season_id}/episodes/{episode_id}", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/series/{series_id}/similar", service.NotImplemented)

}

func (s *SeriesHandler) RegisterOpenApi(o service.OAPI) {
	o.Register("GET", "/api/v3/content/v2/series", s.GetSeriesOAPI)
	o.Register("GET", "/api/v3/content/v2/series/{series_id}", s.GetSeriesByIDOAPI)
	o.Register("GET", "/api/v3/content/v2/series/{series_id}/seasons", s.GetSeasonsOAPI)
	o.Register("GET", "/api/v3/content/v2/series/{series_id}/seasons/{season_id}", s.GetSeasonByIDOAPI)
	o.Register("GET", "/api/v3/content/v2/series/{series_id}/seasons/{season_id}/episodes", s.GetEpisodesOAPI)
	o.Register("GET", "/api/v3/content/v2/series/{series_id}/seasons/{season_id}/episodes/{episode_id}", s.GetEpisodeByIDOAPI)
	o.Register("GET", "/api/v3/content/v2/series/{series_id}/similar", s.GetSeriesOAPI)
}

func (s *SeriesHandler) GetSeriesOAPI(op openapi.OperationContext) {
	op.SetTags("Series")
	op.SetSummary("Get Series")
	op.SetDescription("Get Series")
	op.AddReqStructure(
		new(
			struct {
				Limit    int    `query:"limit"`
				Offset   int    `query:"offset"`
				Category string `query:"category"`
			},
		),
	)
}

func (s *SeriesHandler) GetSeriesByIDOAPI(op openapi.OperationContext) {
	op.SetTags("Series")
	op.SetSummary("Get Series By ID")
	op.SetDescription("Get Series By ID")
	op.AddReqStructure(
		new(
			struct {
				SeriesID string `path:"series_id"`
			},
		),
	)
}

func (s *SeriesHandler) GetSeasonsOAPI(op openapi.OperationContext) {
	op.SetTags("Series")
	op.SetSummary("Get Seasons")
	op.SetDescription("Get Seasons")
	op.AddReqStructure(
		new(
			struct {
				SeriesID string `path:"series_id"`
			},
		),
	)
}

func (s *SeriesHandler) GetSeasonByIDOAPI(op openapi.OperationContext) {
	op.SetTags("Series")
	op.SetSummary("Get Season By ID")
	op.SetDescription("Get Season By ID")
	op.AddReqStructure(
		new(
			struct {
				SeriesID string `path:"series_id"`
				SeasonID string `path:"season_id"`
			},
		),
	)
}

func (s *SeriesHandler) GetEpisodesOAPI(op openapi.OperationContext) {
	op.SetTags("Series")
	op.SetSummary("Get Episodes")
	op.SetDescription("Get Episodes")
	op.AddReqStructure(
		new(
			struct {
				SeriesID string `path:"series_id"`
				SeasonID string `path:"season_id"`
			},
		),
	)
}

func (s *SeriesHandler) GetEpisodeByIDOAPI(op openapi.OperationContext) {
	op.SetTags("Series")
	op.SetSummary("Get Episode By ID")
	op.SetDescription("Get Episode By ID")
	op.AddReqStructure(
		new(
			struct {
				SeriesID  string `path:"series_id"`
				SeasonID  string `path:"season_id"`
				EpisodeID string `path:"episode_id"`
			},
		),
	)
}
