package main

import (
	"app/config"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/zipkin"

	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
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

	// init metrics
	mp, mpshut, err := initMetrics(conf)
	if err != nil {
		panic(err)
	}
	defer mpshut(context.Background())

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Starting the app")

	// print the config
	logger.Debug().Interface("config", conf).Msg("config")
	app, clenaup, err := wireApp(conf, logger, mp)
	if err != nil {
		panic(err)
	}
	defer clenaup()
	app.RegisterRoutes()
	go func() {
		if err := app.Run(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()
	// create channel to block the main thread
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info().Msgf("Shutdown the app with signal %s", sig.String())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Error shutting down the app")
	}
	logger.Info().Msg("App shutdown")

	<-ctx.Done()
	log.Println("timeout of 5 seconds.")

}

// init metrics for prometheus
func initMetrics(config *config.Config) (metric.MeterProvider, func(context.Context) error, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, nil, err
	}
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.Tracing.ServiceName),
		),
		),
	)
	return provider, provider.Shutdown, nil

}

// func intiMetrics(config *config.Config) error {
// 	exporter, err := prometheus.New()
// 	if err != nil {
// 		return err
// 	}
// 	provider := metric.NewMeterProvider(
// 		metric.WithReader(exporter),
// 		metric.WithResource(resource.NewWithAttributes(
// 			semconv.SchemaURL,
// 			semconv.ServiceNameKey.String(config.Tracing.ServiceName),
// 		),
// 		),
// 	)

// 	meter := provider.Meter("")

// }

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
