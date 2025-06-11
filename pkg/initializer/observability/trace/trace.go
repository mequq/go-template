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
) (traceapiotel.TracerProvider, error) {

	bgctx := context.Background()
	if otelCollectorGrpcAddress != nil {
		traceExporter, err := otlptracegrpc.New(bgctx, otlptracegrpc.WithEndpoint(otelCollectorGrpcAddress.String()), otlptracegrpc.WithInsecure())
		if err != nil {
			slog.Info("Failed to create OTLP trace exporter", "error", err)
			return nil, err
		}
		tp := tracesdkotel.NewTracerProvider(
			tracesdkotel.WithBatcher(traceExporter),
			tracesdkotel.WithResource(resources),
		)
		go func() {
			<-ctx.Done()
			if err := tp.Shutdown(bgctx); err != nil {
				slog.Error("failed to shutdown tracer provider", "error", err)
			}
			if err := traceExporter.Shutdown(bgctx); err != nil {
				slog.Error("failed to shutdown trace exporter", "error", err)
			}
		}()
		return tp, nil

	} else {
		tp := tracenoopotel.NewTracerProvider()
		return tp, nil
	}

}
