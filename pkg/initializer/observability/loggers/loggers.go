package loggers

import (
	"context"
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
) (logapiotel.LoggerProvider, func(), error) {
	if otelCollectorGrpcAddress != nil {
		batchExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithEndpoint(otelCollectorGrpcAddress.String()), otlploggrpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}

		lp := logsdkotel.NewLoggerProvider(
			logsdkotel.WithProcessor(
				logsdkotel.NewBatchProcessor(batchExporter),
			),
			logsdkotel.WithResource(resources),
		)
		return lp, func() {
			if err := lp.Shutdown(ctx); err != nil {
				return
			}
		}, nil
	} else {
		lp := lognoopotel.NewLoggerProvider()
		return lp, func() {}, nil

	}

}
