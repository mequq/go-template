package handler

import (
	"application/internal/service"
	"log/slog"
	"net/http"
)

type ClipHandler struct {
	logger *slog.Logger
}

var _ service.Handler = (*ClipHandler)(nil)

func NewMuxClipHandler(logger *slog.Logger) *ClipHandler {
	return &ClipHandler{
		logger: logger.With("layer", "MuxClipService"),
	}
}

func (s *ClipHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v3/content/v2/clips", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/clips/{clip_id}", service.NotImplemented)

}
