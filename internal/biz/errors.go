package biz

import "errors"

var (
	ErrNotAuthorized   = errors.New("not authorized")
	ErrResouceNotFound = errors.New("placeholder resource not found")
	ErrResouceExists   = errors.New("placeholder resource already exists")
	ErrInvalidResource = errors.New("invalid placeholder resource")
)
