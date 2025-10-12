package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type httpServerConfig struct {
	HTTP struct {
		Addr string `koanf:"addr"`
	} `koanf:"http"`
}

func NewHTTPServerConfig(ctx context.Context, c *KConfig) (*httpServerConfig, error) {
	config := new(httpServerConfig)
	if err := c.Unmarshal("server", config); err != nil {
		return nil, err
	}

	return config, nil
}

type HTTPServer interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

var _ HTTPServer = (*httpServer)(nil)

type httpServer struct {
	config  *httpServerConfig
	handler http.Handler
	se      *http.Server
	logger  *slog.Logger
}

var ErrorServerNotStarted = errors.New("server not started")

func NewHTTPServer(cfg *httpServerConfig, handler http.Handler, appLogger AppLogger) *httpServer {
	s := &httpServer{
		config:  cfg,
		handler: handler,
		se:      nil,
		logger:  appLogger.GetLogger(),
	}

	return s
}

func (s *httpServer) Start(ctx context.Context) error {
	s.se = &http.Server{
		Addr:    s.config.HTTP.Addr,
		Handler: otelhttp.NewHandler(NewRecoveryMiddleware(s.handler), "http-server"),
	}

	go func() {
		if err := s.se.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return nil
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	if s.se == nil {
		return ErrorServerNotStarted
	}

	if err := s.se.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

type recoverMiddleware struct {
	logger  slog.Logger
	level   slog.Level
	console bool
	tracer  trace.Tracer
}

type recoveryOption func(*recoverMiddleware)

func WithRecoveryLogger(logger slog.Logger) recoveryOption {
	return func(r *recoverMiddleware) {
		r.logger = logger
	}
}

func WithRecoveryLevel(level slog.Level) recoveryOption {
	return func(r *recoverMiddleware) {
		r.level = level
	}
}

func WithRecoveryConsole(console bool) recoveryOption {
	return func(r *recoverMiddleware) {
		r.console = console
	}
}

func NewRecoveryMiddleware(next http.Handler, opts ...recoveryOption) http.Handler {
	r := &recoverMiddleware{
		logger:  *slog.Default(),
		level:   slog.LevelError,
		console: true,
		tracer:  otel.Tracer("http-server"),
	}
	for _, opt := range opts {
		opt(r)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		defer func() {
			if err := recover(); err != nil {
				ctx, span := r.tracer.Start(ctx, "RecoverMiddleware")
				defer span.End()
				span.RecordError(fmt.Errorf("%v", err)) //nolint:err113
				span.SetStatus(codes.Error, fmt.Sprintf("%v", err))
				span.SetAttributes(attribute.String("panic", fmt.Sprintf("%v", err)))
				span.SetAttributes(attribute.String("stack", string(debug.Stack())))
				r.logger.Log(
					ctx,
					r.level,
					"Panic Recovered",
					"panic",
					err,
					"stack",
					string(debug.Stack()),
				)

				if r.console {
					fmt.Println(err)
					debug.PrintStack()
				}

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, req)
	})
}
