// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"application/config"
	seGorm "application/internal/datasource/sample_entitiy/memory"
	"application/internal/rest-api"
	restHandlers "application/internal/rest-api/handler"
	mock_sample_entitiy "application/mocks/datasource"
	"context"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

// Injectors from wire.go:

func wireApp(ctx context.Context, cfg config.ConfigInterface, logger *slog.Logger) http.Handler {
	dataSource := seGorm.NewSampleEntity()

	healthzHandler := restHandlers.NewMuxHealthzHandler(logger)
	v := restHandlers.NewServiceList(healthzHandler)
	handler := rest_api.NewHttpHandler(cfg, logger, v...)
	return handler, nil
}
