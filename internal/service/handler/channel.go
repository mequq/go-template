package handler

import (
	"application/internal/service"
	"context"
	"log/slog"
	"net/http"
)

type ChannelHandler struct {
	logger *slog.Logger
	mux    *http.ServeMux
}

var _ service.Handler = (*ChannelHandler)(nil)

func NewMuxChannelHandler(logger *slog.Logger, mux *http.ServeMux) *ChannelHandler {
	return &ChannelHandler{
		logger: logger.With("layer", "MuxChannelService"),
		mux:    mux,
	}
}

func (s *ChannelHandler) RegisterHandler(_ context.Context) error {
	s.mux.HandleFunc("GET /api/v3/content/v2/channels", service.NotImplemented)
	s.mux.HandleFunc("GET /api/v3/content/v2/channels/{channel_id}", service.NotImplemented)
	s.mux.HandleFunc("GET /api/v3/content/v2/channels/{channel_id}/programs", service.NotImplemented)
	s.mux.HandleFunc("GET /api/v3/content/v2/channels/{channel_id}/similar", service.NotImplemented)

	return nil
}
