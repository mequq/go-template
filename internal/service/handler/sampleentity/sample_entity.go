package sampleentity

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	sampleusecasev1 "application/internal/biz/sample"
	"application/internal/service"
	"application/internal/service/dto"
	"application/internal/service/response"
	"application/pkg/middlewares"
	"application/pkg/middlewares/httplogger"
	"application/pkg/middlewares/httprecovery"
	"application/pkg/utils"

	"github.com/swaggest/openapi-go"
	"go.opentelemetry.io/otel"
)

var _ service.Handler = (*SampleEntityHandler)(nil)
var _ service.OpenApiHandler = (*SampleEntityHandler)(nil)

type SampleEntityHandler struct {
	sampleEntityBiz sampleusecasev1.SampleEntityUsecaseInterface
	logger          *slog.Logger
}

func NewSampleEntityHandler(logger *slog.Logger, sampleEntityBiz sampleusecasev1.SampleEntityUsecaseInterface) *SampleEntityHandler {
	return &SampleEntityHandler{
		logger: logger.With("layer", "MuxSampleEntityService"), sampleEntityBiz: sampleEntityBiz,
	}
}

func (s *SampleEntityHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "SampleEntityHandler.Create")
	defer span.End()

	logger := s.logger.With("method", "Create", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("Create called with ctx")

	var request dto.SampleEntityRequest
	if err := request.FromRequest(r); err != nil {
		response.BadRequest(w, "invalid-request")
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusBadRequest)
		return
	}

	if err := request.Validate(); err != nil {
		response.BadRequest(w, "invalid-request")
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusBadRequest)
		return
	}

	sampleEnt, err := s.sampleEntityBiz.Create(ctx, request.ToEntity())
	if err != nil {
		if errors.Is(err, sampleusecasev1.ErrAlreadyExist) {
			response.BadRequest(w, "already-exist")
			return
		}
		response.InternalError(w)
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusInternalServerError)
		return
	}

	response.Created(w, sampleEnt)
	logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusOK)
}

func (s *SampleEntityHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "SampleEntityHandler.Update")
	defer span.End()

	logger := s.logger.With("method", "Create", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("Update called with ctx")

	var request dto.SampleEntityRequest
	if err := request.FromRequest(r); err != nil {
		response.BadRequest(w, "invalid-request")
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusBadRequest)
		return
	}

	if err := request.Validate(); err != nil {
		response.BadRequest(w, "invalid-request")
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid-request")
		return
	}
	ent := request.ToEntity()
	ent.ID = uint64(id) //nolint exclusive

	if err := s.sampleEntityBiz.Update(ctx, ent); err != nil {
		if errors.Is(err, sampleusecasev1.ErrAlreadyExist) {
			response.BadRequest(w, "already-exist")
			return
		}
		if errors.Is(err, sampleusecasev1.ErrNotFound) {
			response.NotFound(w)
			return
		}
		response.InternalError(w)
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusInternalServerError)
		return
	}

	response.Ok(w, nil, "Updated successfully")
	logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusOK)
}

func (s *SampleEntityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "SampleEntityHandler.Delete")
	defer span.End()

	logger := s.logger.With("method", "Delete", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("Delete called with ctx")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid-request")
		return
	}

	i := uint64(id) //nolint exclusive

	if err := s.sampleEntityBiz.Delete(ctx, i); err != nil {

		if errors.Is(err, sampleusecasev1.ErrNotFound) {
			response.NotFound(w)
			return
		}
		response.InternalError(w)
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusInternalServerError)
		return
	}

	response.Ok(w, nil, "sample entity deleted")
	logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusOK)
}

func (s *SampleEntityHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "SampleEntityHandler.List")
	defer span.End()

	logger := s.logger.With("method", "List", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("List called with ctx")

	es, err := s.sampleEntityBiz.List(ctx)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Ok(w, dto.SampleEntityListResponses(es), "")
	logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusOK)
}

func (s *SampleEntityHandler) RegisterMuxRouter(mux *http.ServeMux) {
	recoverMiddleware, err := httprecovery.NewRecoveryMiddleware()
	if err != nil {
		panic(err)
	}

	loggerMiddlewareDebug, err := httplogger.NewLoggerMiddleware(httplogger.WithLevel(slog.LevelDebug))
	if err != nil {
		panic(err)
	}

	middles := []middlewares.Middleware{
		recoverMiddleware.RecoverMiddleware,
		httplogger.SetRequestContextLogger,
		loggerMiddlewareDebug.LoggerMiddleware,
	}

	mux.HandleFunc("POST /api/v1/sample-entities", middlewares.MultipleMiddleware(s.Create, middles...))
	mux.HandleFunc("GET /api/v1/sample-entities", middlewares.MultipleMiddleware(s.List, middles...))
	mux.HandleFunc("PUT /api/v1/sample-entities/{id}", middlewares.MultipleMiddleware(s.Update, middles...))
	mux.HandleFunc("DELETE /api/v1/sample-entities/{id}", middlewares.MultipleMiddleware(s.Delete, middles...))
}

func (s *SampleEntityHandler) RegisterOpenApi(o service.OAPI) {
	o.Register("POST", "/api/sample/v1/sample_entities", s.OpenApiCreateSampleEntity, service.WithTags("SampleEntity"))
	o.Register("GET", "/api/sample/v1/sample_entities", s.OpenApiListSampleEntity, service.WithTags("SampleEntity"))

}

// create sample entity for open api
func (s *SampleEntityHandler) OpenApiCreateSampleEntity(o openapi.OperationContext) {
	o.SetSummary("Create SampleEntity")
	o.SetDescription("Create SampleEntity")
	o.AddReqStructure(dto.SampleEntityRequest{})
	o.AddRespStructure(dto.SampleEntityResponse{})
}

// list sample entity for open api
func (s *SampleEntityHandler) OpenApiListSampleEntity(o openapi.OperationContext) {
	o.SetSummary("List SampleEntity")
	o.SetDescription("List SampleEntity")
	o.AddRespStructure([]dto.SampleEntityResponse{})
}
