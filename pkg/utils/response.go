package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JSONResponse sends a JSON response with the given status code and data
func JSONResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}) {
	response := Response{
		Success: success,
		Code:    statusCode,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// SuccessResponse sends a success response with data
func SuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	JSONResponse(w, http.StatusOK, true, message, data)
}

// ErrorResponse sends an error response
func ErrorResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	JSONResponse(w, statusCode, false, message, data)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(w http.ResponseWriter, message string, err interface{}) {
	ErrorResponse(w, http.StatusBadRequest, message, map[string]interface{}{
		"error": err,
	})
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusUnauthorized, message, map[string]interface{}{
		"error": "Authentication required",
	})
}

// Forbidden sends a 403 Forbidden response
func Forbidden(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusForbidden, message, map[string]interface{}{
		"error": "Insufficient permissions",
	})
}

// NotFound sends a 404 Not Found response
func NotFound(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusNotFound, message, map[string]interface{}{
		"error": "Resource not found",
	})
}

// InternalError sends a 500 Internal Server Error response
func InternalError(w http.ResponseWriter, message string, err interface{}) {
	ErrorResponse(w, http.StatusInternalServerError, message, map[string]interface{}{
		"error": err,
	})
}
