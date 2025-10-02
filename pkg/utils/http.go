package utils

import (
	"net/http"
	"strings"
)

func GetUserIPAddress(r *http.Request) string {
	if r.Header.Get("X-Forwarded-For") != "" {
		return r.Header.Get("X-Forwarded-For")
	}

	return strings.Split(r.RemoteAddr, ":")[0] // Split to get the IP address without port
}
