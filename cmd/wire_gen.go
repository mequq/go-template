// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"application/config"
	"application/internal/biz"
	"application/internal/data"
	"application/internal/rest-api"
	handler2 "application/internal/rest-api/handler"
	"context"
	"log/slog"
	"net/http"
)

// Injectors from wire.go:

func wireApp(ctx context.Context, cfg config.ConfigInterface, logger *slog.Logger) (http.Handler, error) {
	dataSource, err := data.NewDataSource(ctx, logger, cfg)
	if err != nil {
		return nil, err
	}
	healthzRepoInterface := data.NewHealthzRepo(logger, dataSource)
	healthzUseCaseInterface := biz.NewHealthzUseCase(healthzRepoInterface, logger)
	healthzService := handler2.NewMuxHealthzService(healthzUseCaseInterface, logger)
	v := handler2.NewServiceList(healthzService)
	handler := rest_api.NewHttpHandler(cfg, logger, v...)
	return handler, nil
}
