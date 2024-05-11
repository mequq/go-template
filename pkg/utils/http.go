package utils

type HttpError struct {
	Message string
	Code    uint
	Error   string
}

func NewHttpError(code uint, message string, err error) HttpError {
	if err == nil {
		return HttpError{
			Message: message,
			Code:    code,
			Error:   "",
		}
	} else {
		return HttpError{
			Message: message,
			Code:    code,
			Error:   err.Error(),
		}
	}
}
