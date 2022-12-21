package app

import (
	"app/config"
	"app/internal/service"
	"context"
	"net/http"

	ginzerolog "github.com/easonlin404/gin-zerolog"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// App is the application
type App struct {
	config  *config.Config
	Service *service.Service
	Engine  *gin.Engine
	Router  *gin.RouterGroup
	srv     *http.Server
}

// NewApp creates a new application.
func NewApp(config *config.Config, svc *service.Service) *App {
	// defer tracingShoutdon(context.Background())
	// gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(ginzerolog.Logger())
	engine.Use(otelgin.Middleware("app"))
	p := ginprometheus.NewPrometheus("gin")
	p.Use(engine)

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
		srv: &http.Server{
			Addr: config.Server.Host + ":" + config.Server.Port,
		},
	}
	return app

}

// register routes
func (a *App) RegisterRoutes() {
	a.Service.RegisterRoutes(a.Router)
}

// shutdown server
func (a *App) Shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}

// Run runs the application.
func (a *App) Run() error {
	a.srv.Handler = a.Engine
	return a.srv.ListenAndServe()
}
