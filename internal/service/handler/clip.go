package handler

import (
	"application/internal/service"
	"context"
	"log/slog"
	"net/http"
)

type ClipHandler struct {
	logger *slog.Logger
	mux    *http.ServeMux
}

var _ service.Handler = (*ClipHandler)(nil)

func NewMuxClipHandler(logger *slog.Logger, mux *http.ServeMux) *ClipHandler {
	return &ClipHandler{
		logger: logger.With("layer", "MuxClipService"),
		mux:    mux,
	}
}

func (s *ClipHandler) RegisterHandler(_ context.Context) error {
	s.mux.HandleFunc("GET /api/v3/content/v2/clips", service.NotImplemented)
	s.mux.HandleFunc("GET /api/v3/content/v2/clips/{clip_id}", service.NotImplemented)

	return nil
}
