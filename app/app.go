package app

import (
	"context"
	"log/slog"
	"net/http"
)

type appConfig struct {
	Title       string `koanf:"title"`
	Version     string `koanf:"version"`
	Description string `koanf:"description"`
	Environment string `koanf:"environment"`
}

func NewAppConfig(ctx context.Context, c *KConfig) (*appConfig, error) {
	config := new(appConfig)
	if err := c.Unmarshal("app", config); err != nil {
		return nil, err
	}

	return config, nil
}

type Application interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	GetLogger() *slog.Logger
}

var _ Application = (*app)(nil)

type app struct {
	logger      *slog.Logger
	httpHandler *http.ServeMux

	appConfig *appConfig

	httpServer HTTPServer
	appLogger  AppLogger
	controller Controller
}

func NewApp(
	run *runTimeFlags,
	appConfig *appConfig,
	httpServer HTTPServer,
	appLogger AppLogger,
	controller Controller,
) *app {
	a := &app{
		appConfig:  appConfig,
		logger:     appLogger.GetLogger().With("component", "app"),
		httpServer: httpServer,
		appLogger:  appLogger,
		controller: controller,
	}

	return a
}

func (a *app) Start(ctx context.Context) error {
	logger := a.logger.With("component", "app", "version", a.appConfig.Version)

	if err := a.httpServer.Start(ctx); err != nil {
		logger.Error("Failed to start HTTP server", "error", err)

		return err
	}

	for name, startup := range a.controller.GetStarters() {
		logger.Info("Starting component...", "component", name)

		if err := startup(ctx); err != nil {
			logger.Error("Failed to start component", "component", name, "error", err)

			return err
		}
	}

	return nil
}

func (a *app) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down application...")

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Error("Failed to shutdown HTTP server", "error", err)
	}

	for name, shutdown := range a.controller.GetSutdowners() {
		a.logger.Info("Shutting down component...", "component", name)

		if err := shutdown(ctx); err != nil {
			a.logger.Error("Failed to shutdown component", "component", name, "error", err)
		}
	}

	return nil
}

func (a *app) GetLogger() *slog.Logger {
	return a.appLogger.GetLogger()
}
