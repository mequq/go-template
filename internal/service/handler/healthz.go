package handler

import (
	"application/internal/biz"
	"application/internal/service"
	"application/internal/service/dto"
	"application/pkg/middlewares"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	otlpsemconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

type HealthzHandler struct {
	logger *slog.Logger
	uc     biz.UsecaseHealthzer
	tracer trace.Tracer
	mux    *http.ServeMux
}

var _ service.Handler = (*HealthzHandler)(nil)

func NewMuxHealthzHandler(
	uc biz.UsecaseHealthzer,
	logger *slog.Logger,
	mux *http.ServeMux,
) *HealthzHandler {
	return &HealthzHandler{
		logger: logger.With("layer", "MuxHealthzService"),
		uc:     uc,
		tracer: otel.Tracer(
			reflect.TypeOf(HealthzHandler{}).String(),
		),
		mux: mux,
	}
}

func (s *HealthzHandler) RegisterHandler(_ context.Context) error {
	recoverMiddleware := middlewares.NewRecoveryMiddleware()

	loggerMiddleware := middlewares.NewHTTPLoggerMiddleware()

	loggerMiddlewareDebug := middlewares.NewHTTPLoggerMiddleware(
		middlewares.WithLevel[*middlewares.HTTPLoggerMiddleware](slog.LevelDebug),
	)

	healthzMiddleware := []middlewares.Middleware{
		recoverMiddleware.RecoverMiddleware,
		loggerMiddlewareDebug.LoggerMiddleware,
	}

	otherMiddleware := []middlewares.Middleware{
		loggerMiddleware.LoggerMiddleware,
	}

	s.mux.HandleFunc(
		"GET /healthz/liveness",
		middlewares.MultipleMiddleware(s.healthzLiveness, healthzMiddleware...),
	)
	s.mux.HandleFunc(
		"GET /healthz/readiness",
		middlewares.MultipleMiddleware(s.healthzReadiness, healthzMiddleware...),
	)
	s.mux.HandleFunc(
		"GET /healthz/panic",
		middlewares.MultipleMiddleware(s.panic, otherMiddleware...),
	)
	s.mux.HandleFunc(
		"GET /healthz/sleep/{time}",
		middlewares.MultipleMiddleware(s.longRun, otherMiddleware...),
	)

	return nil
}

// HealthzLiveness
//
//	@Summary		Healthz Liveness
//	@Description	Check the liveness of the service
//	@ID				healthz-liveness
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response[bool]	"ok"
//	@Router			/healthz/liveness [get]
//	@Tags			healthz
func (s *HealthzHandler) healthzLiveness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.logger.With("method", "liveness")

	ctx, span := s.tracer.Start(
		ctx,
		"liveness",

		trace.WithAttributes(otlpsemconv.AppWidgetName("liveness")),
	)
	defer span.End()

	span.SetName("liveness")

	w.Header().Set("Content-Type", "application/json")

	logger.DebugContext(ctx, "Liveness")

	err := s.uc.Liveness(ctx)
	if err != nil {
		dto.HandleError(errors.New("service not available"), w)

		return
	}

	span.SetStatus(otelCodes.Ok, "ok")
	dto.HandleError(nil, w)
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
func (s *HealthzHandler) healthzReadiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.logger.With("method", "HealthzReadiness")
	ctx, span := s.tracer.Start(
		ctx,
		"readiness",
		trace.WithAttributes(attribute.String("name", "readiness2")),
	)

	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	err := s.uc.Readiness(ctx)
	if err != nil {
		dto.HandleError(errors.New("service not available"), w)

		return
	}

	span.SetStatus(otelCodes.Ok, "ok")
	logger.InfoContext(ctx, "Readiness ok")
	dto.HandleError(nil, w)
}

// panic
//
//	@Router		/healthz/panic [get]
//	@Summary	Panic for test
//	@Success	500	{object}	response.Response[string]	"panic"
//	@Tags		healthz
func (s *HealthzHandler) panic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.logger.With("method", "Panic")

	ctx, span := s.tracer.Start(ctx, "Panic", trace.WithAttributes(attribute.Bool("panic", true)))
	defer span.End()

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
func (s *HealthzHandler) longRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// add baggage
	// sleep 30 second
	timeString := r.PathValue("time")
	logger := s.logger.With("method", "LongRun")
	logger.DebugContext(ctx, "LongRun", "time", timeString)

	ctx, span := otel.Tracer("handler").Start(ctx, "LongRun")
	defer span.End()

	// sleep to int
	duration, err := time.ParseDuration(timeString)
	if err != nil {
		logger.Error("LongRun", "err", err)

		span.SetStatus(otelCodes.Error, "error")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", err.Error()))

		dto.HandleError(biz.ErrInvalidResource, w)

		return
	}

	time.Sleep(duration)

	span.SetStatus(otelCodes.Ok, "ok")
	logger.InfoContext(ctx, "LongRun Test")

	dto.HandleError(nil, w)
}
