//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"
	"log/slog"
	"net/http"

	"application/config"
	"application/internal/biz"
	"application/internal/datasource"
	rest_api "application/internal/http"
	"application/internal/http/handler"
	"github.com/google/wire"
)

func wireApp(ctx context.Context, cfg config.Config, logger *slog.Logger) (http.Handler, error) {
	panic(wire.Build(
		datasource.DataProviderSet,
		biz.BizProviderSet,
		rest_api.ServerProviderSet,
		handler.HandlerProviderSet,
	))
}
