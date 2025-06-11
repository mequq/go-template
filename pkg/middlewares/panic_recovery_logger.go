package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"application/internal/service/response"
)

type RecoverMiddleware struct {
	MiddlewareGeneral
	consolePanic bool
}

func NewRecoveryMiddleware(opts ...Options[*RecoverMiddleware]) *RecoverMiddleware {
	r := &RecoverMiddleware{
		MiddlewareGeneral: MiddlewareGeneral{
			logger: slog.Default(),
			level:  slog.LevelError,
		},
		consolePanic: true,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func WithConsolePanic(printPanic bool) Options[*RecoverMiddleware] {
	return func(r *RecoverMiddleware) {
		r.consolePanic = printPanic
	}
}

func (rm *RecoverMiddleware) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rm.MiddlewareGeneral.logger.Log(req.Context(), rm.MiddlewareGeneral.level, "Panic Recovered", "panic", err, "stack", string(debug.Stack()))
				if rm.consolePanic {
					fmt.Println(err) //nolint:all
					debug.PrintStack()
				}
				response.InternalError(w)
			}
		}()
		next.ServeHTTP(w, req)
	})
}
