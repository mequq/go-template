package handler

import (
	"application/pkg/middlewares"
	"application/pkg/middlewares/httplogger"
	"application/pkg/middlewares/httprecovery"
	"application/pkg/utils"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
)

type HealthzService struct {
	logger *slog.Logger
}

type Response struct {
	Message string `json:"message"`
}

var _ ServiceInterface = (*HealthzService)(nil)

func NewMuxHealthzService(logger *slog.Logger) *HealthzService {
	return &HealthzService{
		logger: logger.With("layer", "MuxHealthzService"),
	}
}

// Healthz Liveness
func (s *HealthzService) HealthzLiveness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.logger.With("method", "HealthzLiveness", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("Liveness")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "ok"})
}

// Healthz Readiness
func (s *HealthzService) HealthzReadiness(w http.ResponseWriter, r *http.Request) {
	// context
	ctx := r.Context()
	// logger
	logger := s.logger.With("method", "HealthzReadiness", "ctx", ctx)
	//  application json
	w.Header().Set("Content-Type", "application/json")

	ctx, span := otel.Tracer("handler").Start(ctx, "rediness")
	defer span.End()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "ok"})
	logger.DebugContext(ctx, "HealthzReadiness", "url", r.Host, "status", http.StatusOK)
}

// panic
func (s *HealthzService) Panic(w http.ResponseWriter, r *http.Request) {
	panic("Panic for test")
}

func (s *HealthzService) LongRun(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(Response{Message: "ok"})
	w.WriteHeader(http.StatusOK)
}

func (s *HealthzService) RegisterMuxRouter(mux *http.ServeMux) {

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
