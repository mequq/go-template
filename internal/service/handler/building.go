package handler

import (
	"application/internal/service"
	"application/internal/service/response"
	"log/slog"
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

type BuildingHandler struct {
	logger *slog.Logger
}

var _ service.Handler = (*BuildingHandler)(nil)

func NewMuxBuildingHandler(logger *slog.Logger) *BuildingHandler {
	return &BuildingHandler{
		logger: logger.With("layer", "MuxBuildingService"),
	}
}

// register router

func (s *BuildingHandler) RegisterMuxRouter(r *http.ServeMux) {
	r.HandleFunc("GET /api/building/v1/buildings", response.NotIplemented)
	r.HandleFunc("GET /api/building/v1/buildings/{buildingID}", response.NotIplemented)
	r.HandleFunc("POST /api/building/v1/buildings", response.NotIplemented)
	r.HandleFunc("PUT /api/building/v1/buildings/{buildingID}", response.NotIplemented)
	r.HandleFunc("DELETE /api/building/v1/buildings/{buildingID}", response.NotIplemented)
}

func (s *BuildingHandler) OpenApiSpec(r *openapi3.Reflector) error {
	return nil
}
