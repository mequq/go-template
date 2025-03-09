package metrics

import (
	"context"
	"log/slog"
	"net"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	meterapiotel "go.opentelemetry.io/otel/metric"
	metrnoopotel "go.opentelemetry.io/otel/metric/noop"
	metersdkotel "go.opentelemetry.io/otel/sdk/metric"
	resoucesdkotel "go.opentelemetry.io/otel/sdk/resource"
)

// NewProvider returns a new metrics provider.

func NewProvider(
	ctx context.Context,
	otelCollectorGrpcAddress *net.TCPAddr,
	resources *resoucesdkotel.Resource,
) (meterapiotel.MeterProvider, func(), error) {

	if otelCollectorGrpcAddress != nil {
		metricExporter, err := otlpmetricgrpc.New(
			ctx, otlpmetricgrpc.WithEndpoint(otelCollectorGrpcAddress.String()),
			otlpmetricgrpc.WithInsecure(),
		)
		if err != nil {
			slog.Info("Failed to create OTLP metric exporter", "error", err)
			return nil, nil, err
		}
		mp := metersdkotel.NewMeterProvider(
			metersdkotel.WithReader(metersdkotel.NewPeriodicReader(
				metricExporter,
				metersdkotel.WithInterval(3*time.Second),
			)),
			metersdkotel.WithResource(resources),
		)

		return mp, func() {
			if err := mp.Shutdown(ctx); err != nil {
				slog.Info("Failed to stop OTLP metric exporter", "error", err)
			}
		}, nil

	} else {
		mp := metrnoopotel.NewMeterProvider()
		return mp, func() {}, nil
	}
}
