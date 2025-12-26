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

	otellogapi "go.opentelemetry.io/otel/log"
	otellognoop "go.opentelemetry.io/otel/log/noop"
	otelmetricapi "go.opentelemetry.io/otel/metric"
	otelmericnoop "go.opentelemetry.io/otel/metric/noop"
	otelpropagator "go.opentelemetry.io/otel/propagation"
	otellogsdk "go.opentelemetry.io/otel/sdk/log"
	otelmetricsdk "go.opentelemetry.io/otel/sdk/metric"
	otelresouesdk "go.opentelemetry.io/otel/sdk/resource"
	oteltracesdk "go.opentelemetry.io/otel/sdk/trace"
	otelsemconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	oteltraceapi "go.opentelemetry.io/otel/trace"
	oteltracenoop "go.opentelemetry.io/otel/trace/noop"
)

type ExporterType string

const (
	ExporterUnknown ExporterType = "unknown"
	ExporterOTLP    ExporterType = "otlp"
	ExporterDebug
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
	GetTracerProvider() oteltraceapi.TracerProvider
	GetMeterProvider() otelmetricapi.MeterProvider
	GetLoggerProvider() otellogapi.LoggerProvider
}

var _ OTLP = (*otlp)(nil)

type otlp struct {
	logger       *slog.Logger
	ctx          context.Context
	obsConfig    *collectorConfig
	appConfig    *appConfig
	httpConfig   *httpServerConfig
	otlpResource *otelresouesdk.Resource
	otlpTracer   *oteltracesdk.TracerProvider
	otlpMeter    *otelmetricsdk.MeterProvider
	otlpLogger   *otellogsdk.LoggerProvider
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

func (o *otlp) GetTracerProvider() oteltraceapi.TracerProvider {
	if o.otlpTracer != nil {
		return o.otlpTracer
	}

	return oteltracenoop.NewTracerProvider()
}

func (o *otlp) GetMeterProvider() otelmetricapi.MeterProvider {
	if o.otlpMeter != nil {
		return o.otlpMeter
	}

	return otelmericnoop.NewMeterProvider()
}

func (o *otlp) GetLoggerProvider() otellogapi.LoggerProvider {
	if o.otlpLogger != nil {
		return o.otlpLogger
	}

	return otellognoop.NewLoggerProvider()
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

	appResource, err := otelresouesdk.New(
		ctx,
		otelresouesdk.WithSchemaURL(otelsemconv.SchemaURL),

		otelresouesdk.WithAttributes(
			otelsemconv.ServiceName(o.appConfig.Title),
			otelsemconv.ServiceVersion(o.appConfig.Version),
			otelsemconv.DeploymentEnvironmentName(o.appConfig.Environment),
			otelsemconv.SourceAddress(host),
			otelsemconv.SourcePort(p),
		),
		otelresouesdk.WithFromEnv(),      // pull attributes from OTEL_RESOURCE_ATTRIBUTES env var
		otelresouesdk.WithProcess(),      // pull attributes from the current process
		otelresouesdk.WithHost(),         // pull attributes from the host
		otelresouesdk.WithContainer(),    // pull attributes from the container
		otelresouesdk.WithTelemetrySDK(), // pull attributes from the telemetry SDK
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
	otlpResource *otelresouesdk.Resource,
) (*otelmetricsdk.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(o.obsConfig.Exporters.OTLP.Endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		o.logger.Info("Failed to create OTLP metric exporter", "error", err)

		return nil, err
	}

	meterProvider := otelmetricsdk.NewMeterProvider(
		otelmetricsdk.WithReader(
			otelmetricsdk.NewPeriodicReader(
				metricExporter,
			),
		),
		otelmetricsdk.WithResource(otlpResource),
	)

	return meterProvider, nil
}

// initTracerProvider initializes the tracer provider if tracing is enabled in the config.
// If not enabled, it returns nil.
func (o *otlp) initTraceProvider(
	ctx context.Context,
	otlpResource *otelresouesdk.Resource,
) (*oteltracesdk.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(o.obsConfig.Exporters.OTLP.Endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		o.logger.Info("Failed to create OTLP trace exporter", "error", err)

		return nil, err
	}

	trp := oteltracesdk.NewTracerProvider(
		oteltracesdk.WithBatcher(traceExporter),
		oteltracesdk.WithResource(otlpResource),
	)

	return trp, nil
}

// initLoggerProvider initializes the logger provider if logging is enabled in the config.
func (o *otlp) initLogProvider(
	ctx context.Context,
	otlpResource *otelresouesdk.Resource,
) (*otellogsdk.LoggerProvider, error) {
	logExporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(o.obsConfig.Exporters.OTLP.Endpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		o.logger.Info("Failed to create OTLP log exporter", "error", err)

		return nil, err
	}

	loggerProvider := otellogsdk.NewLoggerProvider(
		otellogsdk.WithProcessor(
			otellogsdk.NewBatchProcessor(logExporter),
		),
		otellogsdk.WithResource(otlpResource),
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

		otel.SetTextMapPropagator(otelpropagator.NewCompositeTextMapPropagator(
			otelpropagator.TraceContext{},
			otelpropagator.Baggage{},
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
