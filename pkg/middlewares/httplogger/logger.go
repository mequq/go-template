package httplogger

import (
	"log/slog"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

type LoggerMiddleware struct {
	logger *slog.Logger
	level  slog.Level
}

type LoggerMiddlewareOptions func(*LoggerMiddleware) error

func NewLoggerMiddleware(opts ...LoggerMiddlewareOptions) (*LoggerMiddleware, error) {
	r := &LoggerMiddleware{
		logger: slog.Default(),
		level:  slog.LevelInfo,
	}
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func WithLogger(logger *slog.Logger) LoggerMiddlewareOptions {
	return func(rm *LoggerMiddleware) error {
		rm.logger = logger
		return nil
	}
}

func WithLevel(level slog.Level) LoggerMiddlewareOptions {
	return func(r *LoggerMiddleware) error {
		r.level = level
		return nil
	}
}

func (lm *LoggerMiddleware) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		defer func() {
			ctx := req.Context()

			remoteAddr := req.Header.Get("X-Forwarded-For")
			if remoteAddr == "" {
				remoteAddr = req.RemoteAddr
			}

			lm.logger.Log(ctx, lm.level, "request fulfilled",
				slog.Group(
					"request-info",
					slog.String("method", req.Method),
					slog.String("url", req.URL.String()),
					slog.Int("status", recorder.Status),
					slog.String("duration", time.Since(startTime).String()),
					slog.String("remote", remoteAddr),
					slog.String("request-id", req.Header.Get("x-request-id")),
				),
			)

		}()
		next.ServeHTTP(recorder, req)

	})
}
