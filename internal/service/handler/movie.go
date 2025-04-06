package handler

import (
	"application/internal/service"
	"log/slog"
	"net/http"

	"github.com/swaggest/openapi-go"
)

type MovieHandler struct {
	logger *slog.Logger
}

var _ service.Handler = (*MovieHandler)(nil)
var _ service.OpenApiHandler = (*MovieHandler)(nil)

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

func (s *MovieHandler) RegisterOpenApi(o service.OAPI) {
	o.Register("GET", "/api/v3/content/v2/movies", s.GetMoviesOAPI)
	o.Register("GET", "/api/v3/content/v2/movies/{movie_id}", s.GetMovieByIDOAPI)
	o.Register("GET", "/api/v3/content/v2/movies/{movie_id}/similar", s.GetSimilarMoviesOAPI)
}

func (s *MovieHandler) GetMoviesOAPI(op openapi.OperationContext) {
	op.SetTags("Movies")
	op.SetSummary("Get Movies")
	op.SetDescription("Get Movies")
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

func (s *MovieHandler) GetMovieByIDOAPI(op openapi.OperationContext) {
	op.SetTags("Movies")
	op.SetSummary("Get Movie By ID")
	op.SetDescription("Get Movie By ID")
	op.AddReqStructure(
		new(
			struct {
				MovieID string `path:"movie_id"`
			},
		),
	)
}

func (s *MovieHandler) GetSimilarMoviesOAPI(op openapi.OperationContext) {
	op.SetTags("Movies")
	op.SetSummary("Get Similar Movies")
	op.SetDescription("Get Similar Movies")
	op.AddReqStructure(
		new(
			struct {
				MovieID string `path:"movie_id"`
			},
		),
	)
}
