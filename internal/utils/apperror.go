package utils

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap/zapcore"
)

type AppError interface {
	Error() string
	AddError(error) AppError
	Errors() []error
	ErrorCode() int
	StatusCode() int
	Message() string
	MarshalLogObject(s zapcore.ObjectEncoder) error
	MarshalJSON() ([]byte, error)
	CleanDetail() AppError
}

type appError struct {
	errorCode  int
	statusCode int
	message    string
	errors     []error
}

func NewAppError(errorCode int, statusCode int, message string) AppError {
	return &appError{
		errorCode:  errorCode,
		statusCode: statusCode,
		message:    message,
		errors:     []error{},
	}
}

func (e *appError) Error() string {
	return fmt.Sprintf("errorCode: %d, statusCode: %d, message: %s", e.errorCode, e.statusCode, e.message)
}

func (e appError) AddError(err error) AppError {
	e.errors = append(e.errors, err)
	return &e
}

func (e appError) Errors() []error {
	return e.errors
}

func (e appError) ErrorCode() int {
	return e.errorCode
}

func (e appError) StatusCode() int {
	return e.statusCode
}

func (e appError) Message() string {
	return e.message
}

func (e appError) Is(err error) bool {
	if err == nil {
		return false
	}
	if e.errorCode == err.(AppError).ErrorCode() {
		return true
	}
	return false
}

func (e appError) As(target interface{}) bool {
	if target == nil {
		return false
	}
	if e.errorCode == target.(AppError).ErrorCode() {
		return true
	}
	return false
}

func (e *appError) MarshalLogObject(s zapcore.ObjectEncoder) error {
	s.AddInt("errorCode", e.errorCode)
	s.AddInt("statusCode", e.statusCode)
	s.AddString("message", e.message)
	s.AddArray("errors", zapcore.ArrayMarshalerFunc(func(a zapcore.ArrayEncoder) error {
		for _, err := range e.errors {
			a.AppendString(fmt.Sprintf("%v", err))
		}
		return nil
	}))
	return nil
}

func (e *appError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		StatusCode int      `json:"status_code"`
		Code       int      `json:"code"`
		Message    string   `json:"message"`
		Detail     []string `json:"errors"`
	}{
		StatusCode: e.StatusCode(),
		Code:       e.errorCode,
		Message:    e.message,
		Detail:     e.listErrors(),
	})

}

func (e appError) listErrors() []string {
	errs := make([]string, 0, len(e.errors))
	for _, err := range e.errors {
		errs = append(errs, fmt.Sprintf("%v", err))
	}
	return errs
}

func (e appError) CleanDetail() AppError {
	e.errors = []error{}
	return &e
}

func ConvertError(err error) AppError {
	if err == nil {
		return nil
	}
	if apperr, ok := err.(AppError); ok {
		return apperr
	}
	return NewAppError(0, 500, "internal server error").AddError(err)
}
