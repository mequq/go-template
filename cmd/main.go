package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	configPKG "application/config"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var ErrorRequestTimeout = errors.New("request take too long to respond")

func main() {
	ctx := context.Background()
	defer ctx.Done()

	// Initialize configuration
	config := initConfig()

	// Initialize tracing
	initTracing(ctx)

	// Initialize logger
	logger := initLogger(config)

	// Initialize and start HTTP server
	httpServer := initHTTPServer(ctx, config, logger)

	// Handle graceful shutdown
	handleGracefulShutdown(ctx, httpServer, logger)
}

func initConfig() configPKG.Config {
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

	config, err := configPKG.NewKoanfConfig(configPKG.WithYamlConfigPath(confAddress))
	if err != nil {
		panic(err)
	}

	return config
}

func initTracing(ctx context.Context) {
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

	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("myService"),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String("abcdef12345"),
		semconv.ContainerName("myContainer"),
	)

	r2, err := resource.New(context.Background())
	if err != nil {
		panic(err)
	}

	resource, err := resource.Merge(r, r2)
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)

	p := b3.New()
	otel.SetTextMapPropagator(p)
	otel.SetTracerProvider(tp)
}

func initLogger(config configPKG.Config) *slog.Logger {
	var cfg configPKG.LogingConfig
	if err := config.Unmarshal("", &cfg); err != nil {
		log.Fatal(err)
	}
	logger := initSlogLogger(cfg)
	logger.Info("logger started", "config", cfg)
	return logger
}

func initHTTPServer(ctx context.Context, config configPKG.Config, logger *slog.Logger) *http.Server {
	var httpConfig configPKG.HTTPServer
	if err := config.Unmarshal("http_server", &httpConfig); err != nil {
		log.Fatal(err)
	}

	engine, err := wireApp(ctx, config, logger)
	if err != nil {
		logger.Error("failed to init app", "err", err)
		panic(err)
	}

	serviceAddr := fmt.Sprintf("%s:%d", httpConfig.Host, httpConfig.Port)
	httpServer := &http.Server{
		Addr:        serviceAddr,
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
	logger.Info(fmt.Sprintf("microservice started at %s", serviceAddr))

	return httpServer
}

func handleGracefulShutdown(ctx context.Context, httpServer *http.Server, logger *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("app stopping...")

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown app", "err", err)
		panic(err)
	}

	logger.Info("app stopped", "signal", sig)
}

func initSlogLogger(cfg configPKG.LogingConfig) *slog.Logger {
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
