package biz

import "errors"

var (
	ErrResourceAccessDenied = errors.New("not authorized")
	ErrResourceNotFound     = errors.New("placeholder resource not found")
	ErrResourceExists       = errors.New("placeholder resource already exists")
	ErrResourceInvalid      = errors.New("invalid placeholder resource")
)
