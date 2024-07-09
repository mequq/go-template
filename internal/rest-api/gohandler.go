package rest_api

import (
	"application/internal/rest-api/handler"
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
