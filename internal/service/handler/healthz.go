package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"time"

	healthzusecase "application/internal/biz"
	"application/internal/service"
	"application/internal/service/response"
	"application/pkg/middlewares"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type HealthzHandler struct {
	logger *slog.Logger
	uc     healthzusecase.HealthzUseCaseInterface
	tracer trace.Tracer
	mux    *http.ServeMux
}

var _ service.Handler = (*HealthzHandler)(nil)

func NewMuxHealthzHandler(uc healthzusecase.HealthzUseCaseInterface, logger *slog.Logger, mux *http.ServeMux) *HealthzHandler {
	return &HealthzHandler{
		logger: logger.With("layer", "MuxHealthzService"),
		uc:     uc,
		tracer: otel.Tracer("handler"),
		mux:    mux,
	}
}

// Healthz Liveness
//
//	@Summary		Healthz Liveness
//	@Description	Check the liveness of the service
//	@ID				healthz-liveness
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response[bool]	"ok"
//	@Router			/healthz/liveness [get]
//	@Tags			healthz
func (s *HealthzHandler) HealthzLiveness(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	ctx, span := s.tracer.Start(ctx, "HealthzLiveness", trace.WithAttributes(attribute.Bool("liveness", true)))
	defer span.End()
	w.Header().Set("Content-Type", "application/json")
	logger := s.logger.With("method", "HealthzLiveness")
	logger.DebugContext(ctx, "Liveness")
	err := s.uc.Liveness(ctx)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Ok(w, nil, "ok")
}

// Healthz Readiness
//
//	@Summary		Healthz Readiness
//	@Description	Check the readiness of the service
//	@ID				healthz-rediness
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response[string]	"ok"
//	@Router			/healthz/rediness [get]
//	@Tags			healthz
func (s *HealthzHandler) HealthzReadiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "rediness")
	defer span.End()
	logger := s.logger.With("method", "HealthzReadiness")
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
//
//	@Router		/healthz/panic [get]
//	@Summary	Panic for test
//	@Success	500	{object}	response.Response[string]	"panic"
//	@Tags		healthz
func (s *HealthzHandler) Panic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.tracer.Start(ctx, "Panic", trace.WithAttributes(attribute.Bool("panic", true)))
	defer span.End()
	logger := s.logger.With("method", "Panic")
	logger.ErrorContext(ctx, "Panic")
	span.SetStatus(otelCodes.Error, "Panic")
	span.RecordError(errors.New("Panic"))
	span.AddEvent("Panic", trace.WithStackTrace(true))

	panic("Panic for test")
}

// longRun for test
//
//	@Router		/healthz/sleep/{time} [get]
//	@Summary	Long Run for test
//	@Success	200		{object}	response.Response[string]	"ok"
//	@Param		time	path		string						true	"Time to sleep, e.g. 30s"
//	@Tags		healthz
func (s *HealthzHandler) LongRun(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "LongRun")
	defer span.End()

	// add baggage
	// sleep 30 second
	timeString := r.PathValue("time")
	logger := s.logger.With("method", "LongRun")
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

func (s *HealthzHandler) RegisterHandler(_ context.Context) error {

	recoverMiddleware := middlewares.NewRecoveryMiddleware()

	loggerMiddleware := middlewares.NewHTTPLoggerMiddleware()

	loggerMiddlewareDebug := middlewares.NewHTTPLoggerMiddleware(middlewares.WithLevel[*middlewares.HTTPLoggerMiddleware](slog.LevelDebug))

	healthzMiddleware := []middlewares.Middleware{
		recoverMiddleware.RecoverMiddleware,
		loggerMiddlewareDebug.LoggerMiddleware,
	}

	otherMiddleware := []middlewares.Middleware{
		loggerMiddleware.LoggerMiddleware,
		recoverMiddleware.RecoverMiddleware,
	}
	s.mux.HandleFunc("GET /healthz/liveness", middlewares.MultipleMiddleware(s.HealthzLiveness, healthzMiddleware...))
	s.mux.HandleFunc("GET /healthz/readiness", middlewares.MultipleMiddleware(s.HealthzReadiness, healthzMiddleware...))
	s.mux.HandleFunc("GET /healthz/panic", middlewares.MultipleMiddleware(s.Panic, otherMiddleware...))
	s.mux.HandleFunc("GET /healthz/sleep/{time}", middlewares.MultipleMiddleware(s.LongRun, otherMiddleware...))
	return nil
}
