package main

import (
	configPKG "application/pkg/initializer/config"
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/swaggest/openapi-go/openapi3"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	sdkotel "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type appConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Description string `mapstructure:"description"`
	Environment string `mapstructure:"environment"`
}

type observabilitConfig struct {
	OTELGrpc struct {
		Enabled bool
		Address string
	} `koanf:"otel-grpc"`

	Tracing struct {
		Enabled bool
	}

	Metrics struct {
		Enabled bool
	}

	Logging struct {
		Level string
	}
}

type HTTPServer struct {
	Port       int    `koanf:"Port"`
	Host       string `mapstructure:"host"`
	Production bool   `mapstructure:"production"`
}

type RuntimeCommand struct {
	// Command to get runtime information
	configAddress string
	saveOpenApi   bool
	openApiPath   string
}

func getRuntimeCommand() (*RuntimeCommand, error) {

	wd, err := os.Getwd()
	if err != nil {
		log.Println("failed to get working directory", "err", err)
		return nil, err
	}

	configAddress := flag.String("config", path.Join(wd, "config.yaml"), "config file address")
	saveOpenApi := flag.Bool("save-openapi", false, "save openapi spec to file")
	openApiPath := flag.String("openapi-path", "./doc/swagger/", "openapi spec file path")
	flag.Parse()

	return &RuntimeCommand{
		configAddress: *configAddress,
		saveOpenApi:   *saveOpenApi,
		openApiPath:   *openApiPath,
	}, nil
}

func logLevel(logLevel string) slog.Level {
	switch logLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func saveOpenApiSpec(oapi3Reflector *openapi3.Reflector, openApiPath string) error {
	yamlPath := openApiPath + "swagger.yaml"
	data, err := oapi3Reflector.Spec.MarshalYAML()
	if err != nil {
		slog.Error("failed to marshal openapi spec", "err", err)
		return err
	}
	err = os.WriteFile(yamlPath, data, 0o644) //nolint
	if err != nil {
		slog.Error("failed to write openapi spec", "err", err)
		return err
	}

	jsonPath := openApiPath + "swagger.json"
	data, err = oapi3Reflector.Spec.MarshalJSON()
	if err != nil {
		slog.Error("failed to marshal openapi spec", "err", err)
		return err
	}
	err = os.WriteFile(jsonPath, data, 0o644) //nolint
	if err != nil {
		slog.Error("failed to write openapi spec", "err", err)
		return err
	}
	return nil

}

func initConfig(confAddress string) configPKG.Config {

	config, err := configPKG.NewKoanfConfig(configPKG.WithYamlConfigPath(confAddress))
	if err != nil {
		panic(err)
	}

	return config
}

func initHTTPServer(ctx context.Context, config configPKG.Config, logger *slog.Logger, openApiReflector *openapi3.Reflector, validate *validator.Validate) *http.Server {
	var httpConfig HTTPServer
	if err := config.Unmarshal("server.http", &httpConfig); err != nil {
		panic(err)
	}

	logger.Debug("starting server", "port", httpConfig.Port, "host", httpConfig.Host)

	engine, err := wireApp(ctx, config, logger, openApiReflector, validate)
	if err != nil {
		logger.Error("failed to init app", "err", err)
		panic(err)
	}

	engine = otelhttp.NewHandler(engine, "")
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

func newResources(ctx context.Context, cfg configPKG.Config) *sdkotel.Resource {
	appInfo := appConfig{}
	err := cfg.Unmarshal("", &appInfo)
	if err != nil {
		log.Println(err, " failed to unmarshal config")
		return nil
	}

	appResource := sdkotel.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(appInfo.Name),
		semconv.ServiceVersionKey.String(appInfo.Version),
	)

	hostResource, err := sdkotel.New(
		ctx,
		sdkotel.WithContainer(),
		sdkotel.WithHost(),
		sdkotel.WithProcess(),
		sdkotel.WithContainerID(),
		sdkotel.WithOS(),
		sdkotel.WithTelemetrySDK(),
	)

	if err != nil {
		log.Println("failed to create host resource", "err", err)
	}

	resource, err := sdkotel.Merge(appResource, hostResource)
	if err != nil {
		log.Println("failed to merge resources", "err", err)
	}

	return resource

}
