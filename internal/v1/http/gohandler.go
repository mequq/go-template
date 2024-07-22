package http

import (
	"net/http"

	"application/internal/v1/http/handler"
)

func NewHTTPHandler(
	svcs ...handler.Handler,
) http.Handler {
	mux := http.NewServeMux()

	for _, svc := range svcs {
		svc.RegisterMuxRouter(mux)
	}

	return mux
}
