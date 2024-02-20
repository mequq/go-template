package main

import (
	"application/config"
	"application/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/signal"
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
	ErrorRequestTimeout = utils.NewAppError(1000, 503, "request timeout")
)

func main() {

	// ...

	ctx := context.Background()
	defer ctx.Done()

	cfg := &config.ViperConfig{}
	conf, err := config.NewConfig("config.yaml")
	if err != nil {
		panic(err)
	}
	// load config
	if err := conf.Load(cfg); err != nil {
		panic(err)
	}

	// init tracer to stdout

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	// exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
	// if err != nil {
	// 	panic(err)
	// }

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

	r2, err := resource.New(context.Background(),
		// resource.WithFromEnv(),   // pull attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables
		// resource.WithProcess(),   // This option configures a set of Detectors that discover process information
		// resource.WithOS(),        // This option configures a set of Detectors that discover OS information
		resource.WithContainer(), // This option configures a set of Detectors that discover container information
		// resource.WithHost(),      // This option configures a set of Detectors that discover host information
		// resource.WithAttributes(attribute.String("foo", "bar")), // Or specify resource attributes directly
		// resource.WithContainerID(),
		// resource.WithSchemaURL(semconv.SchemaURL),
	)
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
	logger := initSlogLogger(cfg)
	logger.Debug("config", "cfg", cfg)

	engine, err := wireApp(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to init app", "err", err)
		panic(err)
	}
	timeoutMSG, err := json.Marshal(ErrorRequestTimeout)
	if err != nil {
		panic(err)
	}

	timeoutHandler := http.TimeoutHandler(engine, 5*time.Second, string(timeoutMSG))

	httpServer := &http.Server{
		Addr:        ":8080",
		Handler:     timeoutHandler,
		ReadTimeout: 3 * time.Second,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("failed to run app", "err", err)
		}

	}()
	logger.Info("microservice started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal := <-quit

	logger.Info("app stopping...")
	// ctx, cancell := context.WithTimeout(ctx, 5*time.Second)
	// defer cancell()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown app", err)
		panic(err)
	}

	logger.Info("app stopped", "signal", signal)

}

func initSlogLogger(conf *config.ViperConfig) *slog.Logger {

	// create list of slog handlers
	slogHandlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	switch conf.Observability.Logging.Level {
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

	// add stdout
	slogHandlers = append(slogHandlers, slog.NewJSONHandler(os.Stdout, slogHandlerOptions))
	// add udp if logstash enabled
	if conf.Observability.Logging.Logstash.Enabled {
		// options := slogsyslog.Option{}
		fmt.Println("logstash enabled")
		// address := conf.Observability.Logging.Logstash.Address
		con, err := net.Dial("udp", conf.Observability.Logging.Logstash.Address)
		if err != nil {
			panic(err)
		}
		slogHandlers = append(slogHandlers, slog.NewJSONHandler(con, slogHandlerOptions))
	}

	logger := slog.New(slogmulti.Fanout(slogHandlers...))
	tages := []string{}
	tages = append(tages, conf.ServerConfig.Name)

	logger = logger.With("service", conf.ServerConfig.Name, "tags", tages)
	slog.SetDefault(logger)

	return logger

}
