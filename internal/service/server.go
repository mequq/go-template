package service

import (
	"context"
	"net/http"

	"github.com/google/wire"
)

// @BasePath /api/v1
var ServerProviderSet = wire.NewSet(NewHTTPHandler, http.NewServeMux)

// Service Interface
type Handler interface {
	RegisterHandler(ctx context.Context) error
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
}
