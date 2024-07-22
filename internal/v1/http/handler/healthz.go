package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	biz "application/internal/v1/biz/healthz"
	"application/internal/v1/http/response"
	"application/pkg/middlewares"
	"application/pkg/middlewares/httplogger"
	"application/pkg/middlewares/httprecovery"
	"application/pkg/utils"

	"go.opentelemetry.io/otel"
)

type HealthzHandler struct {
	logger *slog.Logger
	uc     biz.HealthzUseCaseInterface
}

var _ HandlerInterface = (*HealthzHandler)(nil)

func NewMuxHealthzHandler(uc biz.HealthzUseCaseInterface, logger *slog.Logger) *HealthzHandler {
	return &HealthzHandler{
		logger: logger.With("layer", "MuxHealthzService"),
		uc:     uc,
	}
}

// Healthz Liveness
func (s *HealthzHandler) HealthzLiveness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "rediness")
	defer span.End()
	logger := s.logger.With("method", "HealthzLiveness", "ctx", utils.GetLoggerContext(r.Context()))
	logger.Debug("Liveness")
	err := s.uc.Liveness(ctx)
	if err != nil {
		response.ResponseInternalError(w)
		return
	}

	response.ResponseOk(w, nil, "ok")
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
		response.ResponseInternalError(w)
		return
	}

	response.ResponseOk(w, nil, "ok")
	logger.DebugContext(ctx, "HealthzReadiness", "url", r.Host, "status", http.StatusOK)
}

// panic
func (s *HealthzHandler) Panic(w http.ResponseWriter, r *http.Request) {
	panic("Panic for test")
}

func (s *HealthzHandler) LongRun(w http.ResponseWriter, r *http.Request) {
	// sleep 30 second
	timeString := r.PathValue("time")
	ctx := r.Context()
	logger := s.logger.With("method", "LongRun", "ctx", ctx)
	logger.Debug("LongRun", "time", timeString)

	// sleep to int
	duration, err := time.ParseDuration(timeString)
	if err != nil {
		logger.Error("LongRun", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	time.Sleep(duration)
	response.ResponseOk(w, nil, "ok")
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
