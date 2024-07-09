package http

import (
	"application/internal/http/handler"
	"net/http"
)

func NewHttpHandler(
	svcs ...handler.HandlerInterface,

) http.Handler {
	mux := http.NewServeMux()

	for _, svc := range svcs {
		svc.RegisterMuxRouter(mux)
	}

	return mux

}
