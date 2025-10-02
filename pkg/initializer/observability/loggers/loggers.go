package loggers

import (
	"context"
	"log/slog"
	"net"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	logapiotel "go.opentelemetry.io/otel/log"
	lognoopotel "go.opentelemetry.io/otel/log/noop"
	logsdkotel "go.opentelemetry.io/otel/sdk/log"
	resoucesdkotel "go.opentelemetry.io/otel/sdk/resource"
)

func NewProvider(
	ctx context.Context,
	otelCollectorGrpcAddress *net.TCPAddr,
	resources *resoucesdkotel.Resource,
) (logapiotel.LoggerProvider, error) {
	if otelCollectorGrpcAddress != nil {
		bgctx := context.Background()

		batchExporter, err := otlploggrpc.New(
			bgctx,
			otlploggrpc.WithEndpoint(otelCollectorGrpcAddress.String()),
			otlploggrpc.WithInsecure(),
		)
		if err != nil {
			return nil, err
		}

		lp := logsdkotel.NewLoggerProvider(
			logsdkotel.WithProcessor(
				logsdkotel.NewBatchProcessor(batchExporter),
			),
			logsdkotel.WithResource(resources),
		)

		go func() {
			<-ctx.Done()

			if err := lp.Shutdown(bgctx); err != nil {
				slog.Error("failed to shutdown logger provider", "error", err)
			}

			if err := batchExporter.Shutdown(bgctx); err != nil {
				slog.Error("failed to shutdown batch exporter", "error", err)
			}
		}()

		return lp, nil
	} else {
		lp := lognoopotel.NewLoggerProvider()

		return lp, nil
	}
}
