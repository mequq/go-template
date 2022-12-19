package app

import (
	"app/config"
	"app/internal/service"

	ginzerolog "github.com/easonlin404/gin-zerolog"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// App is the application
type App struct {
	config  *config.Config
	Service *service.Service
	Engine  *gin.Engine
	Router  *gin.RouterGroup
}

// NewApp creates a new application.
func NewApp(config *config.Config, svc *service.Service) *App {
	// defer tracingShoutdon(context.Background())
	// gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(ginzerolog.Logger())
	engine.Use(otelgin.Middleware("app"))

	// tmplName := "user"
	// tmplStr := "user {{ .name }} (id {{ .id }})\n"
	// tmpl := template.Must(template.New(tmplName).Parse(tmplStr))
	// engine.SetHTMLTemplate(tmpl)

	// create a new router
	router := engine.Group("/")
	// create a new app
	app := &App{
		config:  config,
		Service: svc,
		Engine:  engine,
		Router:  router,
	}
	return app

}

// register routes
func (a *App) RegisterRoutes() {
	a.Service.RegisterRoutes(a.Router)
}

// Run runs the application.
func (a *App) Run() error {
	return a.Engine.Run(a.config.Server.Host + ":" + a.config.Server.Port)
}
