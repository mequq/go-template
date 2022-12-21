//go:build wireinject
// +build wireinject

package main

import (
	"app/app"
	"app/config"
	"app/internal/biz"
	"app/internal/data"
	"app/internal/service"
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
)

// go wireapp wire

func wireApp(config *config.Config, logger zerolog.Logger, metricProvider metric.MeterProvider) (*app.App, func(), error) {
	panic(wire.Build(app.NewApp, service.ServiceProviderSet, biz.BizProviderSet, data.DataProviderSet))
}
