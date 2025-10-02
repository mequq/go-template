package handler

import (
	"application/internal/service"
	"context"
	"log/slog"
	"net/http"
)

type BuildingHandler struct {
	logger *slog.Logger
	mux    *http.ServeMux
}

var _ service.Handler = (*BuildingHandler)(nil)

func NewMuxBuildingHandler(logger *slog.Logger, mux *http.ServeMux) *BuildingHandler {
	return &BuildingHandler{
		logger: logger.With("layer", "MuxBuildingService"),
		mux:    mux,
	}
}

// register router

func (s *BuildingHandler) RegisterHandler(_ context.Context) error {
	s.mux.HandleFunc("GET /api/building/v1/buildings", service.NotImplemented)
	s.mux.HandleFunc("GET /api/building/v1/buildings/{buildingID}", service.NotImplemented)
	s.mux.HandleFunc("POST /api/building/v1/buildings", service.NotImplemented)
	s.mux.HandleFunc("PUT /api/building/v1/buildings/{buildingID}", service.NotImplemented)
	s.mux.HandleFunc("DELETE /api/building/v1/buildings/{buildingID}", service.NotImplemented)

	return nil
}
