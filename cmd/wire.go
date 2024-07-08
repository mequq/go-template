//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"net/http"

	"github.com/google/wire"

	"application/config"
	"application/internal/biz"
	"application/internal/data"
	"application/internal/handler"
	"application/internal/server"

	"context"
	"log/slog"
)

func wireApp(ctx context.Context, cfg config.ConfigInterface, logger *slog.Logger) (http.Handler, error) {
	panic(wire.Build(
		server.ServerProviderSet,
		handler.ServiceProviderSet,
		biz.BizProviderSet,
		data.DataProviderSet,
	))
}
