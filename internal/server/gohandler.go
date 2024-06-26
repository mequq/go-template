package server

import (
	"application/config"
	"application/internal/service"
	"log/slog"
	"net/http"
)

func NewHttpHandler(
	cfg config.ConfigInterface,
	logger *slog.Logger,
	svcs ...service.ServiceInterface,

) http.Handler {
	// recoverMiddleware, err := httprecovery.NewRecoveryMiddleware()
	// if err != nil {
	// 	panic(err)
	// }

	// loggerMiddleware, err := httplogger.NewLoggerMiddleware()
	// if err != nil {
	// 	panic(err)
	// }

	// context middleware

	// logger middleware

	mux := http.NewServeMux()

	for _, svc := range svcs {
		svc.RegisterMuxRouter(mux)
	}

	// m := loggerMiddleware.LoggerMiddleware(mux)
	// m = recoverMiddleware.RecoverMiddleware(m)

	return mux

}
