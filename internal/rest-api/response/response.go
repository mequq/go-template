package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}

func ResponseOk(w http.ResponseWriter, data any, message string) {
	ResponseCustom(w, http.StatusOK, data, message)
}

func ResponseCreated(w http.ResponseWriter) {
	ResponseCustom(w, http.StatusCreated, nil, "Created Successfully")
}

func ResponseNotFound(w http.ResponseWriter) {
	ResponseCustom(w, http.StatusNotFound, nil, "not-found")
}

func ResponseInternalError(w http.ResponseWriter) {
	ResponseCustom(w, http.StatusInternalServerError, nil, "internal-error")
}

func ResponseBadRequest(w http.ResponseWriter, message string) {
	ResponseCustom(w, http.StatusBadRequest, nil, message)
}

func ResponseCustom(w http.ResponseWriter, statusCode int, data any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Response{
		Message: message,
		Status:  statusCode,
		Data:    data,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal-server-error"))
	}
}

func PureResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal-server-error"))
	}
}
