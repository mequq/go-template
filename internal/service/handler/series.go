package handler

import (
	"application/internal/service"
	"log/slog"
	"net/http"
)

type SeriesHandler struct {
	logger *slog.Logger
}

var _ service.Handler = (*SeriesHandler)(nil)

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
