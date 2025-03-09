package trace

import (
	"context"
	"log/slog"
	"net"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	resoucesdkotel "go.opentelemetry.io/otel/sdk/resource"
	tracesdkotel "go.opentelemetry.io/otel/sdk/trace"
	traceapiotel "go.opentelemetry.io/otel/trace"
	tracenoopotel "go.opentelemetry.io/otel/trace/noop"
)

func NewProvider(
	ctx context.Context,
	otelCollectorGrpcAddress *net.TCPAddr,
	resources *resoucesdkotel.Resource,
) (traceapiotel.TracerProvider, func(), error) {

	if otelCollectorGrpcAddress != nil {
		traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(otelCollectorGrpcAddress.String()), otlptracegrpc.WithInsecure())
		if err != nil {
			slog.Info("Failed to create OTLP trace exporter", "error", err)
			return nil, nil, err
		}
		tp := tracesdkotel.NewTracerProvider(
			tracesdkotel.WithBatcher(traceExporter),
			tracesdkotel.WithResource(resources),
		)
		return tp, func() {
			if err := traceExporter.Shutdown(ctx); err != nil {
				slog.Info("Failed to stop OTLP trace exporter", "error", err)
			}
		}, nil

	} else {
		tp := tracenoopotel.NewTracerProvider()
		return tp, func() {}, nil
	}

}
