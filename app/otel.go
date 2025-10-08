package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otlplogapi "go.opentelemetry.io/otel/log"
	otlplognoop "go.opentelemetry.io/otel/log/noop"
	otlpmetricapi "go.opentelemetry.io/otel/metric"
	otlpmericnoop "go.opentelemetry.io/otel/metric/noop"
	otelPropagator "go.opentelemetry.io/otel/propagation"
	otlplogsdk "go.opentelemetry.io/otel/sdk/log"
	otlpmetricsdk "go.opentelemetry.io/otel/sdk/metric"
	otlpresouesdk "go.opentelemetry.io/otel/sdk/resource"
	otlptracesdk "go.opentelemetry.io/otel/sdk/trace"
	otlpsemconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	otlptraceapi "go.opentelemetry.io/otel/trace"
	otlptracenoop "go.opentelemetry.io/otel/trace/noop"
)

type collectorConfig struct {
	Exporters struct {
		OTLP struct {
			Endpoint string `koanf:"endpoint"`
			Insecure bool   `koanf:"insecure"`
		} `koanf:"otlp"`
	} `koanf:"exporters"`

	Traces struct {
		Enabled bool `koanf:"enabled"`
	} `koanf:"traces"`

	Metrics struct {
		Enabled bool `koanf:"enabled"`
	} `koanf:"metrics"`

	Logs struct {
		Enabled bool `koanf:"enabled"`
	} `koanf:"logs"`
}

func NewCollectorConfig(ctx context.Context, c *KConfig) (*collectorConfig, error) {
	config := new(collectorConfig)
	if err := c.Unmarshal("collector", config); err != nil {
		return nil, err
	}

	return config, nil
}

type OTLP interface {
	GetTracerProvider() otlptraceapi.TracerProvider
	GetMeterProvider() otlpmetricapi.MeterProvider
	GetLoggerProvider() otlplogapi.LoggerProvider
}

var _ OTLP = (*otlp)(nil)

type otlp struct {
	logger       *slog.Logger
	ctx          context.Context
	obsConfig    *collectorConfig
	appConfig    *appConfig
	httpConfig   *httpServerConfig
	otlpResource *otlpresouesdk.Resource
	otlpTracer   *otlptracesdk.TracerProvider
	otlpMeter    *otlpmetricsdk.MeterProvider
	otlpLogger   *otlplogsdk.LoggerProvider
}

func NewOTLP(
	ctx context.Context,
	cf *collectorConfig,
	appConfig *appConfig,
	controller Controller,
	httpConfig *httpServerConfig,
) (*otlp, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	o := &otlp{
		logger:     logger,
		obsConfig:  cf,
		appConfig:  appConfig,
		httpConfig: httpConfig,
	}

	logger.Warn("Creating OTLP instance", "version", appConfig.Version)

	o.logger.Info("Initializing OTLP", "version", appConfig.Version)

	if err := o.initOTLPResource(ctx); err != nil {
		return nil, err
	}

	if err := o.startLogs(ctx); err != nil {
		return nil, err
	}

	if err := o.startMetrics(ctx); err != nil {
		return nil, err
	}

	if err := o.startTraces(ctx); err != nil {
		return nil, err
	}

	controller.RegisterShutdown("otlp", o.shutdown)

	return o, nil
}

func (o *otlp) GetTracerProvider() otlptraceapi.TracerProvider {
	if o.otlpTracer != nil {
		return o.otlpTracer
	}

	return otlptracenoop.NewTracerProvider()
}

func (o *otlp) GetMeterProvider() otlpmetricapi.MeterProvider {
	if o.otlpMeter != nil {
		return o.otlpMeter
	}

	return otlpmericnoop.NewMeterProvider()
}

func (o *otlp) GetLoggerProvider() otlplogapi.LoggerProvider {
	if o.otlpLogger != nil {
		return o.otlpLogger
	}

	return otlplognoop.NewLoggerProvider()
}

// initOTLPResource initializes the OTLP resource with default attributes.
func (o *otlp) initOTLPResource(ctx context.Context) error {
	host, port, err := net.SplitHostPort(o.httpConfig.HTTP.Addr)
	if err != nil {
		return err
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	appResource, err := otlpresouesdk.New(
		ctx,
		otlpresouesdk.WithSchemaURL(otlpsemconv.SchemaURL),

		otlpresouesdk.WithAttributes(
			otlpsemconv.ServiceName(o.appConfig.Title),
			otlpsemconv.ServiceVersion(o.appConfig.Version),
			otlpsemconv.DeploymentEnvironmentName(o.appConfig.Environment),
			otlpsemconv.SourceAddress(host),
			otlpsemconv.SourcePort(p),
		),
		otlpresouesdk.WithFromEnv(),      // pull attributes from OTEL_RESOURCE_ATTRIBUTES env var
		otlpresouesdk.WithProcess(),      // pull attributes from the current process
		otlpresouesdk.WithHost(),         // pull attributes from the host
		otlpresouesdk.WithContainer(),    // pull attributes from the container
		otlpresouesdk.WithTelemetrySDK(), // pull attributes from the telemetry SDK
	)
	if err != nil {
		return err
	}

	o.otlpResource = appResource

	return nil
}

// initMetricsProvider initializes the metrics provider if metrics are enabled in the config.
// If not enabled, it returns a Noop meter provider.
func (o *otlp) initMetricProvider(
	ctx context.Context,
	otlpResource *otlpresouesdk.Resource,
) (*otlpmetricsdk.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(o.obsConfig.Exporters.OTLP.Endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		o.logger.Info("Failed to create OTLP metric exporter", "error", err)

		return nil, err
	}

	meterProvider := otlpmetricsdk.NewMeterProvider(
		otlpmetricsdk.WithReader(
			otlpmetricsdk.NewPeriodicReader(
				metricExporter,
			),
		),
		otlpmetricsdk.WithResource(otlpResource),
	)

	return meterProvider, nil
}

// initTracerProvider initializes the tracer provider if tracing is enabled in the config.
// If not enabled, it returns nil.
func (o *otlp) initTraceProvider(
	ctx context.Context,
	otlpResource *otlpresouesdk.Resource,
) (*otlptracesdk.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(o.obsConfig.Exporters.OTLP.Endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		o.logger.Info("Failed to create OTLP trace exporter", "error", err)

		return nil, err
	}

	trp := otlptracesdk.NewTracerProvider(
		otlptracesdk.WithBatcher(traceExporter),
		otlptracesdk.WithResource(otlpResource),
	)

	return trp, nil
}

// initLoggerProvider initializes the logger provider if logging is enabled in the config.
func (o *otlp) initLogProvider(
	ctx context.Context,
	otlpResource *otlpresouesdk.Resource,
) (*otlplogsdk.LoggerProvider, error) {
	logExporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(o.obsConfig.Exporters.OTLP.Endpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		o.logger.Info("Failed to create OTLP log exporter", "error", err)

		return nil, err
	}

	loggerProvider := otlplogsdk.NewLoggerProvider(
		otlplogsdk.WithProcessor(
			otlplogsdk.NewBatchProcessor(logExporter),
		),
		otlplogsdk.WithResource(otlpResource),
	)

	return loggerProvider, nil
}

// StartLogs initializes and starts the OTLP logging components.
func (o *otlp) startLogs(ctx context.Context) error {
	logger := o.logger.With("method", "startLogs", "version", o.appConfig.Version)

	if o.obsConfig.Logs.Enabled && o.otlpLogger == nil {
		logger.InfoContext(ctx, "Logging is enabled, initializing OTLP logger provider")

		var err error
		if o.otlpLogger, err = o.initLogProvider(ctx, o.otlpResource); err != nil {
			logger.Error("failed to create OTLP logger provider", "error", err)

			return err
		}
	}

	return nil
}

func (o *otlp) startMetrics(ctx context.Context) error {
	logger := o.logger.With("method", "startMetrics", "version", o.appConfig.Version)

	if o.obsConfig.Metrics.Enabled && o.otlpMeter == nil {
		logger.InfoContext(ctx, "Metrics is enabled, initializing OTLP meter provider")

		var err error
		if o.otlpMeter, err = o.initMetricProvider(ctx, o.otlpResource); err != nil {
			logger.Error("failed to create OTLP metric provider", "error", err)

			return err
		}

		otel.SetMeterProvider(o.otlpMeter)
	}

	return nil
}

func (o *otlp) startTraces(ctx context.Context) error {
	logger := o.logger.With("method", "startTraces", "version", o.appConfig.Version)

	if o.obsConfig.Traces.Enabled && o.otlpTracer == nil {
		logger.InfoContext(ctx, "Tracing is enabled, initializing OTLP tracer provider")

		var err error
		if o.otlpTracer, err = o.initTraceProvider(ctx, o.otlpResource); err != nil {
			logger.Error("failed to create OTLP tracer provider", "error", err)

			return err
		}

		otel.SetTextMapPropagator(otelPropagator.NewCompositeTextMapPropagator(
			otelPropagator.TraceContext{},
			otelPropagator.Baggage{},
		))

		otel.SetTracerProvider(o.otlpTracer)
	}

	return nil
}

// Shutdown gracefully shuts down the OTLP components.
func (o *otlp) shutdown(ctx context.Context) error {
	logger := o.logger.With("method", "shutdown", "version", o.appConfig.Version)

	var err error
	if o.otlpTracer != nil {
		if err := o.otlpTracer.Shutdown(ctx); err != nil {
			err = errors.Join(err, err)
			logger.Error("failed to shutdown OTLP tracer provider", "error", err)
		}
	}

	if o.otlpMeter != nil {
		if err := o.otlpMeter.Shutdown(ctx); err != nil {
			err = errors.Join(err, err)
			logger.Error("failed to shutdown OTLP meter provider", "error", err)
		}
	}

	if o.otlpLogger != nil {
		if err := o.otlpLogger.Shutdown(ctx); err != nil {
			err = errors.Join(err, err)
			logger.Error("failed to shutdown OTLP logger provider", "error", err)
		}
	}

	return err
}
