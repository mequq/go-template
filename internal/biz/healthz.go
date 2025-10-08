package biz

import (
	"application/app"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type healthz struct {
	logger     *slog.Logger
	tracer     trace.Tracer
	controller app.Controller
}

// NewHealthz creates a new instance of the Healthz use case.
func NewHealthz(logger *slog.Logger, controller app.Controller) *healthz {
	return &healthz{
		logger:     logger.With("layer", "Healthz"),
		tracer:     otel.Tracer("HealthzUseCase"),
		controller: controller,
	}
}

func (uc *healthz) Readiness(ctx context.Context) error {
	return uc.checkers(ctx, uc.controller.GetHealthz())
}

func (uc *healthz) Liveness(ctx context.Context) error {
	return uc.checkers(ctx, uc.controller.GetHealthz())
}

func (uc *healthz) checkers(ctx context.Context, checkFunc map[string]func(ctx context.Context) error) error {
	logger := uc.logger.With("method", "checkers")

	if len(checkFunc) == 0 {
		return nil
	}

	ctx, span := uc.tracer.Start(
		ctx,
		"checkers",
		trace.WithAttributes(attribute.Bool("checkers", true)),
	)
	defer span.End()

	wg := sync.WaitGroup{}

	errCh := make(chan struct {
		name string
		err  error
	}, len(checkFunc))

	for name, check := range checkFunc {
		wg.Add(1)

		go func(name string, check func(ctx context.Context) error) {
			defer wg.Done()

			ctx, span := uc.tracer.Start(
				ctx, fmt.Sprintf("Readiness-%s", name),
				trace.WithAttributes(attribute.String("name", name)),
			)
			defer span.End()

			if err := check(ctx); err != nil {
				errCh <- struct {
					name string
					err  error
				}{name: name, err: err}

				logger.ErrorContext(ctx, "check failed", "name", name, "error", err)
			}
		}(name, check)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		span.AddEvent("check failed", trace.WithStackTrace(true))
		span.SetStatus(codes.Error, "Readiness check failed")

		var err error = nil

		for resp := range errCh {
			span.RecordError(resp.err, trace.WithAttributes(attribute.String("name", resp.name)))
			err = errors.Join(fmt.Errorf("service %s failed.: %w", resp.name, resp.err), err)
		}

		return err
	}
	return nil
}
