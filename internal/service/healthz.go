package service

import (
	"application/internal/biz"
	"application/internal/utils"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type GorilaMuxHealthzService struct {
	uc     biz.HealthzUseCaseInterface
	logger *slog.Logger
}

type Response struct {
	Message string `json:"message"`
}

// New GorilaMuxHealthzService
func NewGorilaMuxHealthzService(uc biz.HealthzUseCaseInterface, logger *slog.Logger) *GorilaMuxHealthzService {
	return &GorilaMuxHealthzService{
		uc:     uc,
		logger: logger.With("layer", "GorilaMuxHealthzService"),
	}
}

// Healthz Liveness
func (s *GorilaMuxHealthzService) HealthzLiveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	err := s.uc.Liveness(ctx)
	if err != nil {
		apperr := utils.ConvertError(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apperr.CleanDetail())
		return
	}

	w.WriteHeader(http.StatusOK)

	s.logger.Debug("HealthzLiveness", "ctx", r)
	json.NewEncoder(w).Encode(Response{Message: "ok"})
}

// Healthz Readiness
func (s *GorilaMuxHealthzService) HealthzReadiness(w http.ResponseWriter, r *http.Request) {
	// context
	ctx := r.Context()
	// logger
	logger := s.logger.With("method", "HealthzReadiness", "ctx", utils.LogContext(ctx))
	//  application json
	w.Header().Set("Content-Type", "application/json")

	ctx, span := otel.Tracer("service").Start(ctx, "rediness")
	defer span.End()
	err := s.uc.Readiness(ctx)
	if err != nil {
		logger.Error("HealthzReadiness", "err", err)
		apperr := utils.ConvertError(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apperr.CleanDetail())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "ok"})
	logger.DebugContext(ctx, "HealthzReadiness", "url", r.Host, "status", http.StatusOK)
}

// Healthz Route

func (s *GorilaMuxHealthzService) RegisterRoutes(r *mux.Router) {
	sr := r.PathPrefix("/healthz").Subrouter()
	sr.HandleFunc("/liveness", s.HealthzLiveness).Methods(http.MethodGet)
	sr.HandleFunc("/readiness", s.HealthzReadiness).Methods(http.MethodGet)
}
