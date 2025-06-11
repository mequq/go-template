package middlewares

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"application/pkg/utils"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

type HTTPLoggerMiddleware struct {
	MiddlewareGeneral
}

func NewHTTPLoggerMiddleware(opts ...Options[*HTTPLoggerMiddleware]) *HTTPLoggerMiddleware {
	r := &HTTPLoggerMiddleware{
		MiddlewareGeneral: MiddlewareGeneral{
			logger: slog.Default(),
			level:  slog.LevelDebug,
		},
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (lm *HTTPLoggerMiddleware) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		defer func() {
			ctx := req.Context()

			lm.logger.Log(ctx, lm.level, "request fulfilled",
				slog.Group(
					"request-info",
					slog.String("method", req.Method),
					slog.String("url", req.URL.String()),
					slog.Int("status", recorder.Status),
					slog.String("duration", time.Since(startTime).String()),
				),
				"ctx", utils.GetLoggerContext(ctx),
			)
		}()
		next.ServeHTTP(recorder, req)
	})
}

// func WithLogger[T any](logger *slog.Logger) Options[*T] {
// 	return func(lm *T) {
// 		lm.MiddlewareGeneral.logger = logger
// 	}
// }

func SetRequestContextLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		reqID := req.Header.Get("x-request-id")
		if reqID != "" {
			ctx = utils.SetLoggerContext(ctx, slog.String("request-id", req.Header.Get("x-request-id")))
		}
		reqIP := req.Header.Get("x-forwarded-for")
		if reqIP == "" {
			reqIP = strings.Split(req.RemoteAddr, ":")[0]
		}
		ctx = utils.SetLoggerContext(ctx, slog.String("request-ip", reqIP))
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
