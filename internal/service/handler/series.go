package handler

import (
	"application/internal/service"
	"context"
	"log/slog"
	"net/http"
)

type SeriesHandler struct {
	logger *slog.Logger
	mux    *http.ServeMux
}

var _ service.Handler = (*SeriesHandler)(nil)

func NewMuxSeriesHandler(logger *slog.Logger) *SeriesHandler {
	return &SeriesHandler{
		logger: logger.With("layer", "MuxSeriesService"),
	}
}

func (s *SeriesHandler) RegisterHandler(_ context.Context) error {
	s.mux.HandleFunc("GET /api/v3/content/v2/series", service.NotImplemented)
	s.mux.HandleFunc("GET /api/v3/content/v2/series/{series_id}", service.NotImplemented)
	s.mux.HandleFunc("GET /api/v3/content/v2/series/{series_id}/seasons", service.NotImplemented)
	s.mux.HandleFunc(
		"GET /api/v3/content/v2/series/{series_id}/seasons/{season_id}",
		service.NotImplemented,
	)
	s.mux.HandleFunc(
		"GET /api/v3/content/v2/series/{series_id}/seasons/{season_id}/episodes",
		service.NotImplemented,
	)
	s.mux.HandleFunc(
		"GET /api/v3/content/v2/series/{series_id}/seasons/{season_id}/episodes/{episode_id}",
		service.NotImplemented,
	)
	s.mux.HandleFunc("GET /api/v3/content/v2/series/{series_id}/similar", service.NotImplemented)

	return nil
}
