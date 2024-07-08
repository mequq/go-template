package server

// import (
// 	"application/config"
// 	"application/internal/handler"
// 	"net/http"

// 	"log/slog"

// 	"github.com/gin-gonic/gin"
// )

// // NewServer creates a new HTTP server and set up all routes.
// func NewGinServer(
// 	cfg *config.ViperConfig,
// 	logger *slog.Logger,
// 	healthzSvc *handler.HealthzService,

// ) http.Handler {

// 	// gin.SetMode(gin.ReleaseMode)
// 	engine := gin.Default()
// 	healthzSvc.RegisterRoutes(engine)
// 	return engine

// }
