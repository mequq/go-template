package service

import (
	"application/internal/biz"
	"application/pkg/middlewares"
	"application/pkg/middlewares/httplogger"
	"application/pkg/middlewares/httprecovery"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
)

type HealthzService struct {
	uc     biz.HealthzUseCaseInterface
	logger *slog.Logger
}

type Response struct {
	Message string `json:"message"`
}

var _ ServiceInterface = (*HealthzService)(nil)

// New GorilaMuxHealthzService
func NewGorilaMuxHealthzService(uc biz.HealthzUseCaseInterface, logger *slog.Logger) *HealthzService {
	return &HealthzService{
		uc:     uc,
		logger: logger.With("layer", "GorilaMuxHealthzService"),
	}
}

// Healthz Liveness
func (s *HealthzService) HealthzLiveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	err := s.uc.Liveness(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	s.logger.Debug("HealthzLiveness", "ctx", ctx)
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

	ctx, span := otel.Tracer("service").Start(ctx, "rediness")
	defer span.End()
	err := s.uc.Readiness(ctx)
	if err != nil {
		logger.Error("HealthzReadiness", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "ok"})
	logger.DebugContext(ctx, "HealthzReadiness", "url", r.Host, "status", http.StatusOK)
}

// panic
func (s *HealthzService) Panic(w http.ResponseWriter, r *http.Request) {
	panic("Panic for test")
}

//  long running request

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

// Healthz Route
// func (s *GorilaMuxHealthzService) RegisterRoutes(r *mux.Router) {
// 	sr := r.PathPrefix("/healthz").Subrouter()
// 	sr.HandleFunc("/liveness", s.HealthzLiveness).Methods(http.MethodGet)
// 	sr.HandleFunc("/readiness", s.HealthzReadiness).Methods(http.MethodGet)
// }

func (s *HealthzService) RegisterMuxRouter(mux *http.ServeMux) {

	recoverMiddleware, err := httprecovery.NewRecoveryMiddleware()
	if err != nil {
		panic(err)
	}

	loggerMiddleware, err := httplogger.NewLoggerMiddleware()
	if err != nil {
		panic(err)
	}
	healthzMiddleware := []middlewares.Middleware{
		recoverMiddleware.RecoverMiddleware,
	}

	otherMiddleware := []middlewares.Middleware{
		loggerMiddleware.LoggerMiddleware,
		recoverMiddleware.RecoverMiddleware,
	}
	mux.HandleFunc("GET /healthz/liveness", middlewares.MultipleMiddleware(s.HealthzLiveness, healthzMiddleware...))
	mux.HandleFunc("GET /healthz/readiness", middlewares.MultipleMiddleware(s.HealthzReadiness, healthzMiddleware...))
	mux.HandleFunc("GET /healthz/panic", middlewares.MultipleMiddleware(s.Panic, otherMiddleware...))
	mux.HandleFunc("GET /healthz/sleep/{time}", middlewares.MultipleMiddleware(s.LongRun, otherMiddleware...))
}
