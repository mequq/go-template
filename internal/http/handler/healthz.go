package handler

import (
	"log/slog"
	"net/http"
	"time"

	healthzusecase "application/internal/biz/healthz"
	"application/internal/http/response"
	"application/pkg/middlewares"
	"application/pkg/middlewares/httplogger"
	"application/pkg/middlewares/httprecovery"
	"application/pkg/utils"
	"go.opentelemetry.io/otel"
)

type HealthzHandler struct {
	logger *slog.Logger
	uc     healthzusecase.HealthzUseCaseInterface
}

var _ Handler = (*HealthzHandler)(nil)

func NewMuxHealthzHandler(uc healthzusecase.HealthzUseCaseInterface, logger *slog.Logger) *HealthzHandler {
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

	response.Ok(w, nil, "ok")
	logger.DebugContext(ctx, "HealthzReadiness", "url", r.Host, "status", http.StatusOK)
}

// panic
func (s *HealthzHandler) Panic(_ http.ResponseWriter, _ *http.Request) {
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
		response.InternalError(w)
		return
	}
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
