package middlewares

import (
	"log/slog"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func MultipleMiddleware(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h

	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped).ServeHTTP
	}

	return wrapped
}

type Options[T GeneralConfigInterface] func(T)

type GeneralConfigInterface interface {
	setLogger(*slog.Logger)
	setLevel(slog.Level)
}

type MiddlewareGeneral struct {
	logger *slog.Logger
	level  slog.Level
}

func (mg *MiddlewareGeneral) setLogger(logger *slog.Logger) {
	if logger == nil {
		return
	}
	mg.logger = logger
}

func (mg *MiddlewareGeneral) setLevel(level slog.Level) {
	if level < slog.LevelDebug || level > slog.LevelError {
		mg.logger.Error("invalid log level", "level", level)
		return
	}
	mg.level = level
}

func WithLogger[T GeneralConfigInterface](logger *slog.Logger) Options[T] {
	return func(m T) {
		if mg, ok := any(m).(GeneralConfigInterface); ok {
			mg.setLogger(logger)
		} else {
			slog.Error("middleware does not support logger setting", "middleware", m)
		}
	}
}

func WithLevel[T GeneralConfigInterface](level slog.Level) Options[T] {
	return func(m T) {

		if mg, ok := any(m).(GeneralConfigInterface); ok {
			mg.setLevel(level)
		} else {
			slog.Error("middleware does not support level setting", "middleware", m)
		}
	}
}
