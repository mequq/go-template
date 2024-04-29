package httprecovery

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type RecoverMiddleware struct {
	logger       *slog.Logger
	consolePanic bool
	logLevel     slog.Level
}

type RecoverMiddlewareOptions func(*RecoverMiddleware) error

func NewRecoveryMiddleware(opts ...RecoverMiddlewareOptions) (*RecoverMiddleware, error) {
	r := &RecoverMiddleware{
		logger:       slog.Default(),
		consolePanic: true,
		logLevel:     slog.LevelError,
	}
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func WithLogger(logger *slog.Logger) RecoverMiddlewareOptions {
	return func(rm *RecoverMiddleware) error {
		rm.logger = logger
		return nil
	}
}

func WithConsolePanic(printPanic bool) RecoverMiddlewareOptions {
	return func(r *RecoverMiddleware) error {
		r.consolePanic = printPanic
		return nil
	}
}

func WithLogLevel(level slog.Level) RecoverMiddlewareOptions {
	return func(r *RecoverMiddleware) error {
		r.logLevel = level
		return nil
	}
}

func (rm *RecoverMiddleware) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rm.logger.Log(req.Context(), rm.logLevel, "Panic Recovered", "panic", err, "stack", debug.Stack())
				if rm.consolePanic {
					fmt.Println(err)
					debug.PrintStack()
				}
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, req)
	})
}
