package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	message := "test error message"
	err := NewError(message)

	assert.Error(t, err)
	assert.Equal(t, message, err.Error())
}

func TestUnauthorized(t *testing.T) {
	message := "unauthorized access"
	err := Unauthorized(message)

	assert.Error(t, err)
	assert.Equal(t, message, err.Error())
}

func TestBadRequest(t *testing.T) {
	message := "bad request error"
	err := BadRequest(message)

	assert.Error(t, err)
	assert.Equal(t, message, err.Error())
}

func TestNotFound(t *testing.T) {
	message := "resource not found"
	err := NotFound(message)

	assert.Error(t, err)
	assert.Equal(t, message, err.Error())
}

func TestInternalServerError(t *testing.T) {
	message := "internal server error"
	err := InternalServerError(message)

	assert.Error(t, err)
	assert.Equal(t, message, err.Error())
}
