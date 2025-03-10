package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"

	"application/pkg/initializer/observability/loggers"
	"application/pkg/initializer/observability/metrics"
	"application/pkg/initializer/observability/trace"

	slogmulti "github.com/samber/slog-multi"
	"github.com/swaggest/openapi-go/openapi3"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {

	ctx := context.Background()
	defer ctx.Done()

	runtimeCommand, err := getRuntimeCommand()
	if err != nil {
		panic(err)
	}

	// Initialize configuration
	config := initConfig(runtimeCommand.configAddress)

	resources := newResources(ctx, config)

	observabilitConfig := observabilitConfig{}
	err = config.Unmarshal("observability", &observabilitConfig)
	if err != nil {
		log.Println(err, " failed to unmarshal observability config")
	}

	var otelAddress *net.TCPAddr
	if observabilitConfig.OTELGrpc.Enabled {
		otelAddress, err = net.ResolveTCPAddr("tcp", observabilitConfig.OTELGrpc.Address)
		if err != nil {
			log.Println(err, " failed to resolve otel grpc address")
		}
	}

	if observabilitConfig.Tracing.Enabled {
		traceProvider, shutdown, err := trace.NewProvider(ctx, otelAddress, resources)
		if err != nil {
			log.Println(err, " failed to create trace provider")
			return
		}
		defer shutdown()
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{}),
		)
		otel.SetTracerProvider(traceProvider)
	}

	if observabilitConfig.Metrics.Enabled {
		metricProvider, shutdown, err := metrics.NewProvider(ctx, otelAddress, resources)
		if err != nil {
			panic(err)
		}
		defer shutdown()
		otel.SetMeterProvider(metricProvider)
		if err := host.Start(); err != nil {
			slog.Info("Failed to start host observer", "error", err)
		}
		if err := runtime.Start(); err != nil {
			slog.Info("Failed to start runtime observer", "error", err)
		}
	}

	logProvider, shutdown, err := loggers.NewProvider(ctx, otelAddress, resources)
	if err != nil {
		panic(err)
	}
	defer shutdown()
	handlers := []slog.Handler{}
	handlers = append(
		handlers,
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel(observabilitConfig.Logging.Level),
		}),
		otelslog.NewHandler("otel", otelslog.WithLoggerProvider(logProvider), otelslog.WithSource(true)),
	)

	logger := slog.New(slogmulti.Fanout(handlers...))
	slog.SetDefault(logger)

	// init open api
	oapi3Reflector := openapi3.NewReflector()

	// Initialize and start HTTP server
	httpServer := initHTTPServer(ctx, config, logger, oapi3Reflector)

	if runtimeCommand.saveOpenApi {
		err := saveOpenApiSpec(oapi3Reflector, runtimeCommand.openApiPath)
		if err != nil {
			logger.Error("failed to save openapi spec", "err", err)
			return
		}
	}
	// Handle graceful shutdown
	handleGracefulShutdown(ctx, httpServer, logger)
}
