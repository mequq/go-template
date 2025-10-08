package app

import (
	"context"
	"errors"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type httpServerConfig struct {
	HTTP struct {
		Addr string `koanf:"addr"`
	} `koanf:"http"`
}

func NewHTTPServerConfig(ctx context.Context, c *kConfig) (*httpServerConfig, error) {
	config := new(httpServerConfig)
	if err := c.Unmarshal("server", config); err != nil {
		return nil, err
	}

	return config, nil
}

type HTTPServer interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

var _ HTTPServer = (*httpServer)(nil)

type httpServer struct {
	config  *httpServerConfig
	handler http.Handler
	se      *http.Server
}

func NewHTTPServer(cfg *httpServerConfig, handler http.Handler) *httpServer {
	s := &httpServer{
		config:  cfg,
		handler: handler,
		se:      nil,
	}

	return s
}

func (s *httpServer) Start(ctx context.Context) error {
	s.se = &http.Server{
		Addr:    s.config.HTTP.Addr,
		Handler: otelhttp.NewHandler(s.handler, "http-server"),
	}

	go func() {
		if err := s.se.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return nil
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	if s.se == nil {
		return errors.New("server not started")
	}

	if err := s.se.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
