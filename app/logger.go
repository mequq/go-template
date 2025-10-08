package app

import (
	"application/pkg/utils"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

func NewSlogLogger(appLogger AppLogger) *slog.Logger {
	return appLogger.GetLogger()
}

type appLoggerConfig struct {
	Slog struct {
		Level    string `koanf:"level"`
		Encoding string `koanf:"encoding"`
	} `koanf:"slog"`
}

func NewAppLoggerConfig(c *KConfig) (*appLoggerConfig, error) {
	config := new(appLoggerConfig)

	if err := c.Unmarshal("logger", config); err != nil {
		return nil, err
	}

	return config, nil
}

type AppLogger interface {
	GetLogger() *slog.Logger
	AppendHandler(handler slog.Handler)
}

type appLogger struct {
	mu sync.Mutex

	logger       *slog.Logger
	loggerConfig *appLoggerConfig
	handlers     []slog.Handler
}

func NewAppLogger(ctx context.Context, appConfig *appConfig, config *appLoggerConfig, otlp OTLP) (*appLogger, error) {
	l := &appLogger{
		loggerConfig: config,
		logger:       nil,
		handlers: []slog.Handler{
			defaultHandler(config.Slog.Level, config.Slog.Encoding),
			otelslog.NewHandler(
				fmt.Sprintf("otlp/%s", appConfig.Title),
				otelslog.WithLoggerProvider(otlp.GetLoggerProvider()),
			),
		},
	}

	logger := slog.New(utils.NewContextLoggerHandler(slogmulti.Fanout(l.handlers...)))
	l.logger = logger

	logger.Info("Logger initialized", "level", config.Slog.Level, "encoding", config.Slog.Encoding)

	return l, nil
}

func (l *appLogger) GetLogger() *slog.Logger {
	return l.logger
}

func (l *appLogger) AppendHandler(handler slog.Handler) {
	l.handlers = append(l.handlers, handler)
}

func defaultHandler(level, encoder string) slog.Handler {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(strings.ToLower(level))); err != nil {
		lvl = slog.LevelInfo
	}

	if strings.ToLower(encoder) == "json" {
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     lvl,
		})
	}

	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
	})
}
