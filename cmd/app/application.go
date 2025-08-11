package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type application struct {
	logger *slog.Logger
	ctx    context.Context
	server *http.Server
}

func NewApplication(ctx context.Context, mux *http.ServeMux, logger *slog.Logger) *application {
	return &application{
		logger: logger,
		ctx:    ctx,
		server: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}
}

func (app *application) Start(ctx context.Context) error {
	// Register the CloudEvents client with the HTTP server

	app.logger.Info("Starting application...")
	go func() {
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("Could not start server", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	app.logger.Info("Shutting down application...", slog.String("signal", sig.String()))
	if err := app.ctx.Err(); err != nil {
		app.logger.Error("Context error during shutdown", "err", err)
		return err
	}

	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.Error("Server forced to shutdown", "err", err)
		return err
	}
	return nil
}
