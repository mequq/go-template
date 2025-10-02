package middlewares

import (
	"application/pkg/utils"
	"log/slog"
	"net/http"
	"strings"
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
		ctx := req.Context()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		ctx = utils.SetLoggerContext(ctx, slog.String("request-id", req.Header.Get("x-request-id")))
		ctx = utils.SetLoggerContext(ctx, slog.String("request-ip", utils.GetUserIPAddress(req)))
		ctx = utils.SetLoggerContext(ctx, slog.String("method", req.Method))
		ctx = utils.SetLoggerContext(ctx, slog.String("url", req.URL.String()))

		defer func() {
			attrs := utils.GetLoggerContextAsAttrs(ctx)

			attrs = append(attrs,
				slog.Int("status", recorder.Status),
				slog.String("duration", time.Since(startTime).String()),
			)
			lm.logger.LogAttrs(ctx, lm.level, "request fulfilled", attrs...)
		}()

		next.ServeHTTP(recorder, req.WithContext(ctx))
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
			ctx = utils.SetLoggerContext(
				ctx,
				slog.String("request-id", req.Header.Get("x-request-id")),
			)
		}

		reqIP := req.Header.Get("x-forwarded-for")
		if reqIP == "" {
			reqIP = strings.Split(req.RemoteAddr, ":")[0]
		}

		ctx = utils.SetLoggerContext(ctx, slog.String("request-ip", reqIP))
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
