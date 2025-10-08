package app

import "github.com/google/wire"

var AppProviderSet = wire.NewSet(
	NewRunTimeFlags,
	NewKoanfConfig,

	NewController,
	wire.Bind(new(Controller), new(*controller)),

	NewAppLoggerConfig,
	NewAppLogger,
	wire.Bind(new(AppLogger), new(*appLogger)),

	NewCollectorConfig,
	NewOTLP,
	wire.Bind(new(OTLP), new(*otlp)),

	NewHTTPServerConfig,
	NewHTTPServer,
	wire.Bind(new(HTTPServer), new(*httpServer)),

	NewAppConfig,
	NewApp,
	wire.Bind(new(Application), new(*app)),
	NewSlogLogger,
)
