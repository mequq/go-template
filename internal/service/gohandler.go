package service

import (
	"log/slog"
	"net/http"

	_ "application/docs" // Import generated docs

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggest/swgui/v5emb"
	"github.com/swaggo/swag"
)

func NewHTTPHandler(
	logger *slog.Logger,
	svcs ...Handler,
) (http.Handler, error) {
	mux := http.NewServeMux()

	for _, svc := range svcs {

		svc.RegisterMuxRouter(mux)

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
		w.Write([]byte(doc)) //nolint // ignore error
	})

	mux.Handle("/swagger/", v5emb.New(
		"swagger",
		"/docs/swagger/swagger.json",
		"/swagger/",
	))

	return mux, nil
}
