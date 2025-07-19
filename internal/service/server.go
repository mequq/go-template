package service

import (
	"net/http"

	"github.com/google/wire"
)

// @BasePath /api/v1
var ServerProviderSet = wire.NewSet(NewHTTPHandler)

// Service Interface
type Handler interface {
	RegisterMuxRouter(mux *http.ServeMux)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
}
