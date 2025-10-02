//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"
	"log/slog"
	"net/http"

	"application/internal/biz"
	"application/internal/datasource"
	"application/internal/repo"
	rest_api "application/internal/service"
	"application/internal/service/handler"
	"application/pkg/initializer/config"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func wireApp(
	ctx context.Context,
	cfg config.Config,
	logger *slog.Logger,
	validate *validator.Validate,
) (http.Handler, error) {
	panic(wire.Build(
		datasource.DataProviderSet,
		biz.BizProviderSet,
		rest_api.ServerProviderSet,
		handler.HandlerProviderSet,
		repo.RepoProvider,
	))
}
