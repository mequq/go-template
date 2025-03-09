package utils

type HTTPError struct {
	Message string
	Code    uint
	Error   string
}

func NewHTTPError(code uint, message string, err error) HTTPError {
	if err == nil {
		return HTTPError{
			Message: message,
			Code:    code,
			Error:   "",
		}
	}
	return HTTPError{
		Message: message,
		Code:    code,
		Error:   err.Error(),
	}
}
