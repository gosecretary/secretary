package utils

import "errors"

// NewError creates a new error with the given message
func NewError(message string) error {
	return errors.New(message)
}

// Unauthorized returns an unauthorized error
func Unauthorized(message string) error {
	return errors.New(message)
}

// BadRequest returns a bad request error
func BadRequest(message string) error {
	return errors.New(message)
}

// NotFound returns a not found error
func NotFound(message string) error {
	return errors.New(message)
}

// InternalServerError returns an internal server error
func InternalServerError(message string) error {
	return errors.New(message)
}
