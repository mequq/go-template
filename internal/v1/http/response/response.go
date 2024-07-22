package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response[T any] struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    T      `json:"data"`
}

func Ok(w http.ResponseWriter, data any, message string) {
	Custom(w, http.StatusOK, data, message)
}

func Created(w http.ResponseWriter, data any) {
	Custom(w, http.StatusCreated, data, "Created Successfully")
}

func NotFound(w http.ResponseWriter) {
	Custom(w, http.StatusNotFound, nil, "not-found")
}

func InternalError(w http.ResponseWriter) {
	Custom(w, http.StatusInternalServerError, nil, "internal-error")
}

func BadRequest(w http.ResponseWriter, message string) {
	Custom(w, http.StatusBadRequest, nil, message)
}

func Custom(w http.ResponseWriter, statusCode int, data any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Response[any]{
		Message: message,
		Status:  statusCode,
		Data:    data,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("internal-server-error")); err != nil {
			log.Fatal(err)
		}
	}
}

func Pure(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("internal-server-error")); err != nil {
			log.Fatal(err)
		}
	}
}
