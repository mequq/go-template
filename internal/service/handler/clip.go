package handler

import (
	"application/internal/service"
	"log/slog"
	"net/http"

	"github.com/swaggest/openapi-go"
)

type ClipHandler struct {
	logger *slog.Logger
}

var _ service.Handler = (*ClipHandler)(nil)
var _ service.OpenApiHandler = (*ClipHandler)(nil)

func NewMuxClipHandler(logger *slog.Logger) *ClipHandler {
	return &ClipHandler{
		logger: logger.With("layer", "MuxClipService"),
	}
}

func (s *ClipHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v3/content/v2/clips", service.NotImplemented)
	mux.HandleFunc("GET /api/v3/content/v2/clips/{clip_id}", service.NotImplemented)

}

func (s *ClipHandler) RegisterOpenApi(o service.OAPI) {
	o.Register("GET", "/api/v3/content/v2/clips", s.GetClipsOAPI)
	o.Register("GET", "/api/v3/content/v2/clips/{clip_id}", s.GetClipByIDOAPI)
}

func (s *ClipHandler) GetClipsOAPI(op openapi.OperationContext) {
	op.SetTags("Clips")
	op.SetSummary("Get Clips")
	op.SetDescription("Get Clips")
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

func (s *ClipHandler) GetClipByIDOAPI(op openapi.OperationContext) {
	op.SetTags("Clips")
	op.SetSummary("Get Clip By ID")
	op.SetDescription("Get Clip By ID")
	op.AddReqStructure(
		new(
			struct {
				ClipID string `path:"clip_id"`
			},
		),
	)
}
