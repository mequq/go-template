package server

// import (
// 	"application/internal/handler"
// 	"encoding/json"
// 	"fmt"
// 	"log/slog"
// 	"net/http"

// 	"git.abanppc.com/lenz-public/lenz-goapp-sdk/pkg/utils/httpmiddleware"
// 	"github.com/gorilla/mux"
// 	"github.com/spf13/viper"
// 	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
// )

// // New GorilaMuxServer creates a new HTTP server and set up all routes.

// func NewGorillaMuxServer(
// 	cfg *viper.Viper,
// 	logger *slog.Logger,
// 	svcs ...handler.ServiceInterface,

// ) http.Handler {

// 	muxRouter := mux.NewRouter()
// 	middleware := httpmiddleware.NewGorilaMuxMiddleware(
// 		httpmiddleware.WithLevel(slog.LevelDebug),
// 		httpmiddleware.WithLogger(logger),
// 	)
// 	// logger middleware

// 	muxRouter.Use(otelmux.Middleware("my-server"))
// 	muxRouter.Use(middleware.ContextMiddleware)
// 	muxRouter.Use(middleware.RecoverMiddleware)
// 	muxRouter.Use(middleware.LoggerMiddleware)

// 	for _, s := range svcs {
// 		logger.Debug("Registering routes", "handler", s)
// 		s.RegisterRoutes(muxRouter)
// 	}

// 	// walk the mux
// 	_ = muxRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
// 		method, err := route.GetMethods()
// 		if err == nil {
// 			path, err := route.GetPathTemplate()
// 			if err != nil {
// 				logger.Error("Error getting path", "error", err)
// 			}
// 			fmt.Println("Methods:", method, "Route:", path)
// 		}

// 		return nil
// 	})

// 	muxRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		logger.Debug("Not Found", "path", r.URL.Path)
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(
// 			map[string]interface{}{
// 				"message": "Not Found",
// 				"code":    "404",
// 				"path":    r.URL.Path,
// 			},
// 		)
// 	})

// 	// http.Handle("/", mux)
// 	return muxRouter

// }
