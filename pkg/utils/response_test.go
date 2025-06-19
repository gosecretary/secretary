package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		success        bool
		message        string
		data           interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "successful response with data",
			statusCode:     http.StatusOK,
			success:        true,
			message:        "Operation successful",
			data:           map[string]string{"key": "value"},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"code":    float64(200),
				"message": "Operation successful",
				"data":    map[string]interface{}{"key": "value"},
			},
		},
		{
			name:           "error response without data",
			statusCode:     http.StatusBadRequest,
			success:        false,
			message:        "Bad request",
			data:           nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(400),
				"message": "Bad request",
				"data":    nil,
			},
		},
		{
			name:           "successful response with array data",
			statusCode:     http.StatusOK,
			success:        true,
			message:        "List retrieved",
			data:           []string{"item1", "item2", "item3"},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"code":    float64(200),
				"message": "List retrieved",
				"data":    []interface{}{"item1", "item2", "item3"},
			},
		},
		{
			name:           "internal server error",
			statusCode:     http.StatusInternalServerError,
			success:        false,
			message:        "Internal server error",
			data:           nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(500),
				"message": "Internal server error",
				"data":    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a ResponseRecorder to record the response
			recorder := httptest.NewRecorder()

			// Call the JSONResponse function
			JSONResponse(recorder, tt.statusCode, tt.success, tt.message, tt.data)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, recorder.Code)

			// Check the Content-Type header
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

			// Parse the response body
			var responseBody map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			// Check the response structure
			assert.Equal(t, tt.expectedBody["success"], responseBody["success"])
			assert.Equal(t, tt.expectedBody["code"], responseBody["code"])
			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])
			assert.Equal(t, tt.expectedBody["data"], responseBody["data"])
		})
	}
}

func TestSuccessResponse(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		data           interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "success with data",
			message:        "User created successfully",
			data:           map[string]string{"id": "123", "username": "testuser"},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"code":    float64(200),
				"message": "User created successfully",
				"data":    map[string]interface{}{"id": "123", "username": "testuser"},
			},
		},
		{
			name:           "success without data",
			message:        "Operation completed",
			data:           nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"code":    float64(200),
				"message": "Operation completed",
				"data":    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			SuccessResponse(recorder, tt.message, tt.data)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

			var responseBody map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedBody["success"], responseBody["success"])
			assert.Equal(t, tt.expectedBody["code"], responseBody["code"])
			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])
			assert.Equal(t, tt.expectedBody["data"], responseBody["data"])
		})
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		message        string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "bad request error",
			statusCode:     http.StatusBadRequest,
			message:        "Invalid input data",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(400),
				"message": "Invalid input data",
				"data":    nil,
			},
		},
		{
			name:           "unauthorized error",
			statusCode:     http.StatusUnauthorized,
			message:        "Authentication required",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(401),
				"message": "Authentication required",
				"data":    nil,
			},
		},
		{
			name:           "not found error",
			statusCode:     http.StatusNotFound,
			message:        "Resource not found",
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(404),
				"message": "Resource not found",
				"data":    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			ErrorResponse(recorder, tt.statusCode, tt.message, nil)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

			var responseBody map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedBody["success"], responseBody["success"])
			assert.Equal(t, tt.expectedBody["code"], responseBody["code"])
			assert.Equal(t, tt.expectedBody["message"], responseBody["message"])
			assert.Equal(t, tt.expectedBody["data"], responseBody["data"])
		})
	}
}

func TestInternalError(t *testing.T) {
	recorder := httptest.NewRecorder()

	InternalError(recorder, "Database connection failed", "connection timeout")

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	assert.Equal(t, false, responseBody["success"])
	assert.Equal(t, float64(500), responseBody["code"])
	assert.Equal(t, "Database connection failed", responseBody["message"])
	assert.NotNil(t, responseBody["data"])
}

func TestBadRequest(t *testing.T) {
	recorder := httptest.NewRecorder()

	BadRequest(recorder, "Missing required field", "validation error")

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	assert.Equal(t, false, responseBody["success"])
	assert.Equal(t, float64(400), responseBody["code"])
	assert.Equal(t, "Missing required field", responseBody["message"])
	assert.NotNil(t, responseBody["data"])
}

func TestUnauthorized(t *testing.T) {
	recorder := httptest.NewRecorder()

	Unauthorized(recorder, "Invalid session")

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	assert.Equal(t, false, responseBody["success"])
	assert.Equal(t, float64(401), responseBody["code"])
	assert.Equal(t, "Invalid session", responseBody["message"])
	// The Unauthorized function creates a data object with error field
	dataMap, ok := responseBody["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Authentication required", dataMap["error"])
}

func TestNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()

	NotFound(recorder, "User not found")

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	assert.Equal(t, false, responseBody["success"])
	assert.Equal(t, float64(404), responseBody["code"])
	assert.Equal(t, "User not found", responseBody["message"])
	assert.NotNil(t, responseBody["data"])
}

func TestForbidden(t *testing.T) {
	recorder := httptest.NewRecorder()

	Forbidden(recorder, "Access denied")

	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	assert.Equal(t, false, responseBody["success"])
	assert.Equal(t, float64(403), responseBody["code"])
	assert.Equal(t, "Access denied", responseBody["message"])
	assert.NotNil(t, responseBody["data"])
}
