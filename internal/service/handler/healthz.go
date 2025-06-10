package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"time"

	healthzusecase "application/internal/biz"
	"application/internal/service"
	"application/internal/service/response"
	"application/pkg/middlewares"
	"application/pkg/middlewares/httplogger"
	"application/pkg/middlewares/httprecovery"
	"application/pkg/utils"

	"github.com/swaggest/openapi-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type HealthzHandler struct {
	logger *slog.Logger
	uc     healthzusecase.HealthzUseCaseInterface
	tracer trace.Tracer
}

var _ service.Handler = (*HealthzHandler)(nil)
var _ service.OpenApiHandler = (*HealthzHandler)(nil)

func NewMuxHealthzHandler(uc healthzusecase.HealthzUseCaseInterface, logger *slog.Logger) *HealthzHandler {
	return &HealthzHandler{
		logger: logger.With("layer", "MuxHealthzService"),
		uc:     uc,
		tracer: otel.Tracer("handler"),
	}
}

// Healthz Liveness
func (s *HealthzHandler) HealthzLiveness(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	ctx, span := s.tracer.Start(ctx, "HealthzLiveness", trace.WithAttributes(attribute.Bool("liveness", true)))
	defer span.End()
	w.Header().Set("Content-Type", "application/json")
	logger := s.logger.With("method", "HealthzLiveness", "ctx", utils.GetLoggerContext(r.Context()))
	logger.DebugContext(ctx, "Liveness")
	err := s.uc.Liveness(ctx)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Ok(w, nil, "ok")
}

// Healthz Readiness
func (s *HealthzHandler) HealthzReadiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "rediness")
	defer span.End()
	logger := s.logger.With("method", "HealthzReadiness", "ctx", ctx)
	w.Header().Set("Content-Type", "application/json")

	err := s.uc.Readiness(ctx)
	if err != nil {
		response.InternalError(w)
		return
	}
	span.SetStatus(otelCodes.Ok, "ok")

	response.Ok(w, nil, "ok")
	logger.DebugContext(ctx, "HealthzReadiness", "url", r.Host, "status", http.StatusOK)
}

// panic
func (s *HealthzHandler) Panic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.tracer.Start(ctx, "Panic", trace.WithAttributes(attribute.Bool("panic", true)))
	defer span.End()
	logger := s.logger.With("method", "Panic", "ctx", ctx)
	logger.ErrorContext(ctx, "Panic")
	span.SetStatus(otelCodes.Error, "Panic")
	span.RecordError(errors.New("Panic"))
	span.AddEvent("Panic", trace.WithStackTrace(true))

	panic("Panic for test")
}

func (s *HealthzHandler) LongRun(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "LongRun")
	defer span.End()

	// add baggage
	// sleep 30 second
	timeString := r.PathValue("time")
	logger := s.logger.With("method", "LongRun", "ctx", ctx)
	logger.Debug("LongRun", "time", timeString)

	// sleep to int
	duration, err := time.ParseDuration(timeString)
	if err != nil {
		logger.Error("LongRun", "err", err)
		response.InternalError(w)
		return
	}
	logger.InfoContext(ctx, "LongRun Test")
	time.Sleep(duration)
	response.Ok(w, nil, "ok")
}

func (s *HealthzHandler) RegisterMuxRouter(mux *http.ServeMux) {

	recoverMiddleware, err := httprecovery.NewRecoveryMiddleware()
	if err != nil {
		panic(err)
	}

	loggerMiddleware, err := httplogger.NewLoggerMiddleware()
	if err != nil {
		panic(err)
	}
	loggerMiddlewareDebug, err := httplogger.NewLoggerMiddleware(httplogger.WithLevel(slog.LevelDebug))
	if err != nil {
		panic(err)
	}

	healthzMiddleware := []middlewares.Middleware{
		recoverMiddleware.RecoverMiddleware,
		httplogger.SetRequestContextLogger,
		loggerMiddlewareDebug.LoggerMiddleware,
	}

	otherMiddleware := []middlewares.Middleware{
		loggerMiddleware.LoggerMiddleware,
		recoverMiddleware.RecoverMiddleware,
		httplogger.SetRequestContextLogger,
	}
	mux.HandleFunc("GET /healthz/liveness", middlewares.MultipleMiddleware(s.HealthzLiveness, healthzMiddleware...))
	mux.HandleFunc("GET /healthz/readiness", middlewares.MultipleMiddleware(s.HealthzReadiness, healthzMiddleware...))
	mux.HandleFunc("GET /healthz/panic", middlewares.MultipleMiddleware(s.Panic, otherMiddleware...))
	mux.HandleFunc("GET /healthz/sleep/{time}", middlewares.MultipleMiddleware(s.LongRun, otherMiddleware...))
}

func (s *HealthzHandler) OHealthzLiveness(op openapi.OperationContext) {
	op.SetSummary("Healthz Liveness")
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// rediness
func (s *HealthzHandler) OHealthzReadiness(op openapi.OperationContext) {
	op.SetSummary("Healthz Readiness")
	op.AddRespStructure(new(response.Response[string]))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

func (s *HealthzHandler) RegisterOpenApi(o service.OAPI) {
	healthzTag := service.WithTags("Healthz(Internal)")

	o.Register(http.MethodGet, "/healthz/liveness", s.OHealthzLiveness, healthzTag)
	o.Register(http.MethodGet, "/healthz/readiness", s.OHealthzReadiness, healthzTag)
}
