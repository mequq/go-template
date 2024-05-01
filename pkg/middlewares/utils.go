package middlewares

import "net/http"

type Middleware func(http.Handler) http.Handler

func MultipleMiddleware(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {

	if len(m) < 1 {
		return h
	}

	wrapped := h

	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped).ServeHTTP
	}

	return wrapped

}
