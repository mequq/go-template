package main

import (
	"app/config"
	"context"
	"os"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	// load the config
	conf, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	tpshut, err := initTracer(conf)
	if err != nil {
		panic(err)
	}
	defer tpshut(context.Background())
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Starting the app")

	// print the config
	logger.Debug().Interface("config", conf).Msg("config")
	app, clenaup, err := wireApp(conf, logger)
	if err != nil {
		panic(err)
	}
	defer clenaup()
	app.RegisterRoutes()
	if err := app.Run(); err != nil {
		panic(err)
	}

}

func initTracer(config *config.Config) (func(context.Context) error, error) {
	// create a new zipkin exporter
	exporter, err := zipkin.New(
		config.Tracing.Endpoint,
		// zipkin.WithLogger(log.New(os.Stdout, "zipkin: ", log.LstdFlags)),
	)
	if err != nil {
		return nil, err
	}
	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	// create a new trace provider
	p := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader | b3.B3SingleHeader))
	otel.SetTextMapPropagator(p)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.Tracing.ServiceName),
			semconv.DeploymentEnvironmentKey.String(config.Tracing.Environment),
		),
		),
	)
	// register the trace provider
	otel.SetTracerProvider(tp)

	// otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// register the global propagator
	return tp.Shutdown, nil

}
