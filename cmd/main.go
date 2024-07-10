package main

import (
	"application/config"
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"path"
	"syscall"
	"time"

	"os"

	"log/slog"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	ErrorRequestTimeout = errors.New("request take to longs to response")
)

func main() {
	ctx := context.Background()
	defer ctx.Done()

	configAddress := flag.String("config", "", "config file address")
	flag.Parse()
	confAddress := *configAddress

	if confAddress == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		confAddress = path.Join(wd, "config.yaml")
	}

	config, err := config.NewKoanfConfig(config.WithYamlConfigPath(confAddress))
	if err != nil {
		panic(err)
	}

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}
	defer func() {
		err := exporter.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}()

	// init tracer
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("myService"),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String("abcdef12345"),
		semconv.ContainerName("myContainer"),
	)

	r2, err := resource.New(context.Background()) // resource.WithFromEnv(),   // pull attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables
	if err != nil {
		panic(err)
	}

	resource, err := resource.Merge(r, r2)
	if err != nil {
		panic(err)
	}

	// init tracer

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(nil),
		sdktrace.WithResource(resource),
	)

	p := b3.New()
	// Register the B3 propagator globally.
	otel.SetTextMapPropagator(p)

	// init tracer global
	otel.SetTracerProvider(tp)

	// init logger

	var cfg LogingConfig
	config.Unmarshal(&cfg)
	logger := initSlogLogger(cfg)
	logger.Info("logger started", "config", cfg)

	engine, err := wireApp(ctx, config, logger)
	if err != nil {
		logger.Error("failed to init app", "err", err)
		panic(err)
	}

	httpServer := &http.Server{
		Addr:        ":8080",
		Handler:     engine,
		ReadTimeout: 3 * time.Second,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("failed to run app", "err", err)
			panic(err)
		}

	}()
	logger.Info("microservice started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal := <-quit

	logger.Info("app stopping...")

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown app", err)
		panic(err)
	}

	logger.Info("app stopped", "signal", signal)

}

type (
	LogingConfig struct {
		Observability Observability `mapstructure:"observability"`
	}
	Observability struct {
		Logging Logging `mapstructure:"logging"`
	}
	Logging struct {
		Level    string   `mapstructure:"level" `
		Logstash Logstash `mapstructure:"logstash"`
	}
	Logstash struct {
		Enabled bool   `mapstructure:"enabled"`
		Address string `mapstructure:"address"`
	}
)

func initSlogLogger(cfg LogingConfig) *slog.Logger {
	slogHandlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	level := cfg.Observability.Logging.Level

	switch level {
	case "debug":
		slogHandlerOptions.Level = slog.LevelDebug
	case "info":
		slogHandlerOptions.Level = slog.LevelInfo
	case "warn":
		slogHandlerOptions.Level = slog.LevelWarn
	case "error":
		slogHandlerOptions.Level = slog.LevelError
	default:
		slogHandlerOptions.Level = slog.LevelInfo
	}

	slogHandlers := []slog.Handler{}

	slogHandlers = append(slogHandlers, slog.NewJSONHandler(os.Stdout, slogHandlerOptions))
	if cfg.Observability.Logging.Logstash.Enabled {
		fmt.Println("logstash enabled")
		con, err := net.Dial("udp", cfg.Observability.Logging.Logstash.Address)
		if err != nil {
			panic(err)
		}
		slogHandlers = append(slogHandlers, slog.NewJSONHandler(con, slogHandlerOptions))
	}

	logger := slog.New(slogmulti.Fanout(slogHandlers...))
	slog.SetDefault(logger)

	return logger

}
