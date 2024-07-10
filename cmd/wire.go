//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"application/internal/biz"
	"application/internal/datasource"
	rest_api "application/internal/http"
	"application/internal/http/handler"
	"net/http"

	"github.com/google/wire"

	"application/config"

	"context"
	"log/slog"
)

func wireApp(ctx context.Context, cfg config.ConfigInterface, logger *slog.Logger) (http.Handler, error) {
	panic(wire.Build(
		datasource.DataProviderSet,
		biz.BizProviderSet,
		rest_api.ServerProviderSet,
		handler.HandlerProviderSet,
	))
}
