package handler

import (
	"application/internal/v1/biz/sample_entity"
	"application/internal/v1/datasource/sample_entitiy"
	"application/internal/v1/http/dto"
	"application/internal/v1/http/response"
	_ "application/internal/v1/http/swagger"
	"application/pkg/middlewares"
	"application/pkg/middlewares/httplogger"
	"application/pkg/middlewares/httprecovery"
	"application/pkg/utils"
	"errors"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	"go.opentelemetry.io/otel"
	"log/slog"
	"net/http"
	"strconv"
)

var _ HandlerInterface = (*SampleEntityHandler)(nil)

type SampleEntityHandler struct {
	sampleEntityBiz sample_entity.SampleEntity
	logger          *slog.Logger
}

func NewSampleEntityHandler(logger *slog.Logger, sampleEntityBiz sample_entity.SampleEntity) *SampleEntityHandler {
	return &SampleEntityHandler{
		logger: logger.With("layer", "MuxSampleEntityService"), sampleEntityBiz: sampleEntityBiz,
	}
}

// @Summary Create Sample Entity
// @Schemes
// @Description Create Sample Entity
// @Tags SampleEntity
// @Accept json
// @Produce json
// @Param body body dto.SampleEntityRequest true "request body"
// @Success 201 {object} response.Response[dto.SampleEntityResponse]
// @failure 400 {object} response.Response[swagger.EmptyObject]
// @failure 500 {object} response.Response[swagger.EmptyObject]
// @Router /sample-entities [post]
func (s *SampleEntityHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "SampleEntityHandler.Create")
	defer span.End()

	logger := s.logger.With("method", "Create", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("Create called with ctx")

	var request dto.SampleEntityRequest
	if err := request.FromRequest(r); err != nil {
		response.ResponseBadRequest(w, "invalid-request")
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusBadRequest)
		return
	}

	if err := request.Validate(); err != nil {
		response.ResponseBadRequest(w, "invalid-request")
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusBadRequest)
		return
	}

	_, err := s.sampleEntityBiz.Create(ctx, request.ToEntity())
	if err != nil {
		if errors.Is(err, sample_entitiy.ErrAlreadyExist) {
			response.ResponseBadRequest(w, "already-exist")
			return
		}
		response.ResponseInternalError(w)
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusInternalServerError)
		return
	}

	response.ResponseCreated(w)
	logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusOK)
}

// @Summary Update Sample Entity
// @Schemes
// @Description Update Sample Entity
// @Tags SampleEntity
// @Accept json
// @Produce json
// @Param body body dto.SampleEntityRequest true "request body"
// @Param id path int true "id"
// @Success 200 {object} response.Response[swagger.EmptyObject]
// @failure 400 {object} response.Response[swagger.EmptyObject]
// @failure 404 {object} response.Response[swagger.EmptyObject]
// @failure 500 {object} response.Response[swagger.EmptyObject]
// @Router /sample-entities/{id} [PUT]
func (s *SampleEntityHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "SampleEntityHandler.Update")
	defer span.End()

	logger := s.logger.With("method", "Create", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("Update called with ctx")

	var request dto.SampleEntityRequest
	if err := request.FromRequest(r); err != nil {
		response.ResponseBadRequest(w, "invalid-request")
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusBadRequest)
		return
	}

	if err := request.Validate(); err != nil {
		response.ResponseBadRequest(w, "invalid-request")
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseBadRequest(w, "invalid-request")
		return
	}
	ent := request.ToEntity()
	ent.ID = uint64(id)

	if err := s.sampleEntityBiz.Update(ctx, ent); err != nil {
		if errors.Is(err, sample_entitiy.ErrAlreadyExist) {
			response.ResponseBadRequest(w, "already-exist")
			return
		}
		if errors.Is(err, sample_entitiy.ErrNotFound) {
			response.ResponseNotFound(w)
			return
		}
		response.ResponseInternalError(w)
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusInternalServerError)
		return
	}

	response.ResponseOk(w, nil, "Updated successfully")
	logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusOK)
}

// @Summary Delete Sample Entity
// @Schemes
// @Description Delete Sample Entity
// @Tags SampleEntity
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} response.Response[swagger.EmptyObject]
// @failure 400 {object} response.Response[swagger.EmptyObject]
// @failure 404 {object} response.Response[swagger.EmptyObject]
// @failure 500 {object} response.Response[swagger.EmptyObject]
// @Router /sample-entities/{id} [DELETE]
func (s *SampleEntityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "SampleEntityHandler.Delete")
	defer span.End()

	logger := s.logger.With("method", "Delete", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("Delete called with ctx")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseBadRequest(w, "invalid-request")
		return
	}

	if err := s.sampleEntityBiz.Delete(ctx, uint64(id)); err != nil {
		if errors.Is(err, sample_entitiy.ErrNotFound) {
			response.ResponseNotFound(w)
			return
		}
		response.ResponseInternalError(w)
		logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusInternalServerError)
		return
	}

	response.ResponseOk(w, nil, "sample entity deleted")
	logger.DebugContext(ctx, "SampleEntityHandler.", "url", r.Host, "status", http.StatusOK)
}

// @Summary Get all sample entities
// @Schemes
// @Description Get all sample entities
// @Tags SampleEntity
// @Accept json
// @Produce json
// @Success 200 {object} response.Response[[]dto.SampleEntityResponse]
// @failure 500 {object} response.Response[swagger.EmptyObject]
// @Router /sample-entities [GET]
func (s *SampleEntityHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "SampleEntityHandler.List")
	defer span.End()

	logger := s.logger.With("method", "List", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("List called with ctx")

	es, err := s.sampleEntityBiz.List(ctx)
	if err != nil {
		response.ResponseInternalError(w)
		return
	}

	response.ResponseOk(w, dto.SampleEntityListResponses(es), "")
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

	mux.HandleFunc("POST /api/sample-entities", middlewares.MultipleMiddleware(s.Create, middles...))
	mux.HandleFunc("GET /api/sample-entities", middlewares.MultipleMiddleware(s.List, middles...))
	mux.HandleFunc("PUT /api/sample-entities/{id}", middlewares.MultipleMiddleware(s.Update, middles...))
	mux.HandleFunc("DELETE /api/sample-entities/{id}", middlewares.MultipleMiddleware(s.Delete, middles...))
}
