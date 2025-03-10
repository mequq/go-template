package channel

import (
	"application/internal/service"
	"log/slog"
	"net/http"

	"github.com/swaggest/openapi-go"
)

type ChannelHandler struct {
	logger *slog.Logger
}

var _ service.OpenApiHandler = (*ChannelHandler)(nil)
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

func (s *ChannelHandler) RegisterOpenApi(o service.OAPI) {
	o.Register("GET", "/api/v3/content/v2/channels", s.GetChannelsOAPI)
	o.Register("GET", "/api/v3/content/v2/channels/{channel_id}", s.GetChannelByIDOAPI)
	o.Register("GET", "/api/v3/content/v2/channels/{channel_id}/programs", s.GetProgramsOAPI)
	o.Register("GET", "/api/v3/content/v2/channels/{channel_id}/similar", s.GetSimilarChannelsOAPI)
}

func (s *ChannelHandler) GetChannelsOAPI(op openapi.OperationContext) {
	op.SetTags("Channels")
	op.SetSummary("Get Channels")
	op.SetDescription("Get Channels")
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

func (s *ChannelHandler) GetChannelByIDOAPI(op openapi.OperationContext) {
	op.SetTags("Channels")
	op.SetSummary("Get Channel By ID")
	op.SetDescription("Get Channel By ID")
	op.AddReqStructure(
		new(
			struct {
				ChannelID string `path:"channel_id"`
			},
		),
	)
}

func (s *ChannelHandler) GetProgramsOAPI(op openapi.OperationContext) {
	op.SetTags("Channels")
	op.SetSummary("Get Programs")
	op.SetDescription("Get Programs")
	op.AddReqStructure(
		new(
			struct {
				ChannelID string `path:"channel_id"`
				Limit     int    `query:"limit"`
				Offset    int    `query:"offset"`
			},
		),
	)
}

func (s *ChannelHandler) GetSimilarChannelsOAPI(op openapi.OperationContext) {
	op.SetTags("Channels")
	op.SetSummary("Get Similar Channels")
	op.SetDescription("Get Similar Channels")
	op.AddReqStructure(
		new(
			struct {
				ChannelID string `path:"channel_id"`
			},
		),
	)
}
