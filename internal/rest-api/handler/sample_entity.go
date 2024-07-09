package handler

import (
	"application/internal/rest-api/response"
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

type SampleEntityService struct {
	logger *slog.Logger
}

var _ ServiceInterface = (*SampleEntityService)(nil)

func NewMuxSampleEntityService(logger *slog.Logger) *SampleEntityService {
	return &SampleEntityService{
		logger: logger.With("layer", "MuxSampleEntityService"),
	}
}

// SampleEntity Liveness
func (s *SampleEntityService) SampleEntityLiveness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.logger.With("method", "SampleEntityLiveness", "ctx", utils.GetLoggerContext(ctx))
	logger.Debug("Liveness")
	w.Header().Set("Content-Type", "application/json")

	response.ResponseOk(w, nil, "ok")
}

// SampleEntity Readiness
func (s *SampleEntityService) SampleEntityReadiness(w http.ResponseWriter, r *http.Request) {
	// context
	ctx := r.Context()
	// logger
	logger := s.logger.With("method", "SampleEntityReadiness", "ctx", ctx)
	//  application json
	w.Header().Set("Content-Type", "application/json")

	ctx, span := otel.Tracer("handler").Start(ctx, "rediness")
	defer span.End()

	response.ResponseOk(w, nil, "ok")
	logger.DebugContext(ctx, "SampleEntityReadiness", "url", r.Host, "status", http.StatusOK)
}

// panic
func (s *SampleEntityService) Panic(w http.ResponseWriter, r *http.Request) {
	panic("Panic for test")
}

func (s *SampleEntityService) LongRun(w http.ResponseWriter, r *http.Request) {
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

func (s *SampleEntityService) RegisterMuxRouter(mux *http.ServeMux) {

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

	SampleEntityMiddleware := []middlewares.Middleware{
		recoverMiddleware.RecoverMiddleware,
		httplogger.SetRequestContextLogger,
		loggerMiddlewareDebug.LoggerMiddleware,
	}

	otherMiddleware := []middlewares.Middleware{
		loggerMiddleware.LoggerMiddleware,
		recoverMiddleware.RecoverMiddleware,
		httplogger.SetRequestContextLogger,
	}
	mux.HandleFunc("GET /SampleEntity/liveness", middlewares.MultipleMiddleware(s.SampleEntityLiveness, SampleEntityMiddleware...))
	mux.HandleFunc("GET /SampleEntity/readiness", middlewares.MultipleMiddleware(s.SampleEntityReadiness, SampleEntityMiddleware...))
	mux.HandleFunc("GET /SampleEntity/panic", middlewares.MultipleMiddleware(s.Panic, otherMiddleware...))
	mux.HandleFunc("GET /SampleEntity/sleep/{time}", middlewares.MultipleMiddleware(s.LongRun, otherMiddleware...))
}
