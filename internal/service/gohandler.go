package service

import (
	"context"
	"log/slog"
	"net/http"

	_ "application/docs" // Import generated docs

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggest/swgui/v5emb"
	"github.com/swaggo/swag"
)

func NewHTTPHandler(
	ctx context.Context,
	logger *slog.Logger,
	mux *http.ServeMux,
	svcs ...Handler,
) (http.Handler, error) {
	for _, svc := range svcs {
		if err := svc.RegisterHandler(ctx); err != nil {
			logger.Error("failed to register handler", "err", err)

			return nil, err
		}
	}

	mux.Handle("/metrics", promhttp.Handler())

	doc, err := swag.ReadDoc("")
	if err != nil {
		logger.Error("failed to read swagger doc", "err", err)

		return nil, err
	}

	mux.HandleFunc("/docs/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(doc))
	})

	mux.Handle("/swagger/", v5emb.New(
		"swagger",
		"/docs/swagger/swagger.json",
		"/swagger/",
	))

	return mux, nil
}
