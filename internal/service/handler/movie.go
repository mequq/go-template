package handler

import (
	"application/internal/service"
	"log/slog"
	"net/http"
)

type MovieHandler struct {
	logger *slog.Logger
}

var _ service.Handler = (*MovieHandler)(nil)

func NewMuxMovieHandler(logger *slog.Logger) *MovieHandler {
	return &MovieHandler{
		logger: logger.With("layer", "MuxMovieService"),
	}
}

func (s *MovieHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v3/content/v2/movies", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/movies/{movie_id}", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/movies/{movie_id}/similar", service.NotImplemented)

}
