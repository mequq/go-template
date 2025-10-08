package dto

import (
	"application/internal/biz"
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type ErrorMap map[error]struct {
	Message string
	Code    int
}

var ErrorsMap = ErrorMap{
	biz.ErrResouceExists: {
		Message: "یافت نشد",
		Code:    http.StatusNotFound,
	},
	biz.ErrResouceExists: {
		Message: "از قبل وجود دارد",
		Code:    http.StatusConflict,
	},
	biz.ErrInvalidResource: {
		Message: "منبع نامعتبر",
		Code:    http.StatusBadRequest,
	},

	biz.ErrNotAuthorized: {
		Message: "دسترسی غیرمجاز",
		Code:    http.StatusUnauthorized,
	},
}

func HandleError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	if err == nil {
		w.WriteHeader(http.StatusOK)

		_ = encoder.Encode(ErrorResponse{
			Message: "ok",
			Details: "no error",
		})

		return
	}

	for e, v := range ErrorsMap {
		if errors.Is(err, e) {
			w.WriteHeader(v.Code)
			_ = encoder.Encode(ErrorResponse{
				Message: v.Message,
				Details: err.Error(),
			})

			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)

	_ = encoder.Encode(ErrorResponse{
		Message: "خطای ناشناخته",
		Details: err.Error(),
	})
}
