package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// main function.
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := wireApp(ctx)
	if err != nil {
		panic(err)
	}

	logger := app.GetLogger().With("component", "main")
	slog.SetDefault(logger)

	logger.Info("app starting...")

	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	defer app.Shutdown(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("app stopping...", "signal", sig)
}
