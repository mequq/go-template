package handler

import (
	"application/internal/service"
	"log/slog"
	"net/http"
)

type ChannelHandler struct {
	logger *slog.Logger
}

var _ service.Handler = (*ChannelHandler)(nil)

func NewMuxChannelHandler(logger *slog.Logger) *ChannelHandler {
	return &ChannelHandler{
		logger: logger.With("layer", "MuxChannelService"),
	}
}

func (s *ChannelHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v3/content/v2/channels", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/channels/{channel_id}", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/channels/{channel_id}/programs", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/channels/{channel_id}/similar", service.NotImplemented)

}
