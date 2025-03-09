package service

import (
	"net/http"
	"reflect"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/swgui/v5emb"
)

func NewHTTPHandler(
	r *openapi3.Reflector,
	o OAPI,
	svcs ...Handler,
) (http.Handler, error) {
	mux := http.NewServeMux()

	for _, svc := range svcs {

		if _, err := r.Reflect(svc); err == nil {
			svc.RegisterMuxRouter(mux)
		} else {
			return nil, err
		}

		if reflect.TypeOf(svc).Implements(reflect.TypeFor[OpenApiHandler]()) {
			svc.(OpenApiHandler).RegisterOpenApi(o)
		}

	}

	mux.Handle("/metrics", promhttp.Handler())

	sw := NewSwagger(o)
	mux.HandleFunc("/docs/swagger/swagger.json", sw.swagerjson)

	mux.Handle("/swagger/", v5emb.New(
		"swagger",
		"/docs/swagger/swagger.json",
		"/swagger/",
	))

	return mux, nil
}

type Swagger struct {
	OpenAPI OAPI
}

func NewSwagger(o OAPI) *Swagger {
	return &Swagger{
		OpenAPI: o,
	}
}

func (s *Swagger) swagerjson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json, err := s.OpenAPI.GetJsonData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json) //nolint // ignore error
}
