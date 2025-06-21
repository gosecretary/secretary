package middleware

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"secretary/alpha/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalOutput)

	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "test response", w.Body.String())
}

func TestCORS(t *testing.T) {
	handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))

	tests := []struct {
		name           string
		method         string
		expectedCode   int
		expectedOrigin string
	}{
		{
			name:           "GET request",
			method:         "GET",
			expectedCode:   200,
			expectedOrigin: "*",
		},
		{
			name:           "OPTIONS preflight request",
			method:         "OPTIONS",
			expectedCode:   200,
			expectedOrigin: "*",
		},
		{
			name:           "POST request",
			method:         "POST",
			expectedCode:   200,
			expectedOrigin: "*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
			assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
			assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
			assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		})
	}
}

func TestRecovery(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalOutput)

	handler := Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "Internal server error")
}

func TestAuth(t *testing.T) {
	// The Auth middleware is basically a pass-through, so test that it works
	handler := Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "success", w.Body.String())
}

func TestRateLimitMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expectedCode int
		hasSession   bool
		isAdmin      bool
	}{
		{
			name:         "health endpoint should pass",
			path:         "/health",
			expectedCode: 200,
			hasSession:   false,
		},
		{
			name:         "login endpoint should pass",
			path:         "/api/login",
			expectedCode: 200,
			hasSession:   false,
		},
		{
			name:         "any endpoint should pass through rate limit middleware",
			path:         "/api/register",
			expectedCode: 200,
			hasSession:   false,
		},
		{
			name:         "other endpoint should pass through rate limit middleware",
			path:         "/api/users",
			expectedCode: 200,
			hasSession:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			}))

			req := httptest.NewRequest("GET", tt.path, nil)

			if tt.hasSession {
				session := &domain.Session{
					ID:       "test-session",
					UserID:   "test-user",
					Username: "testuser",
				}
				if tt.isAdmin {
					session.Username = "admin"
				}
				ctx := context.WithValue(req.Context(), "session", session)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	// Test default status code
	assert.Equal(t, http.StatusOK, rw.statusCode)

	// Test WriteHeader
	rw.WriteHeader(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, rw.statusCode)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHTTPRequest(t *testing.T) {
	// This function just logs, so we can test that it doesn't panic
	HTTPRequest("GET", "/test", "127.0.0.1", 200, time.Millisecond*100)
	// If we get here without panicking, the test passes
}

func TestSessionMiddleware(t *testing.T) {
	// Note: In a real implementation you'd want proper mocking of the session store
	// For now, we'll test the basic flow without mocking

	tests := []struct {
		name         string
		path         string
		hasCookie    bool
		cookieValue  string
		expectedCode int
	}{
		{
			name:         "health endpoint should pass without session",
			path:         "/health",
			hasCookie:    false,
			expectedCode: 200,
		},
		{
			name:         "login endpoint should pass without session",
			path:         "/api/login",
			hasCookie:    false,
			expectedCode: 200,
		},
		{
			name:         "protected endpoint without cookie should fail",
			path:         "/api/register",
			hasCookie:    false,
			expectedCode: 401,
		},
		{
			name:         "protected endpoint with invalid cookie should fail",
			path:         "/api/register",
			hasCookie:    true,
			cookieValue:  "invalid-session",
			expectedCode: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := SessionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			}))

			req := httptest.NewRequest("GET", tt.path, nil)

			if tt.hasCookie {
				req.AddCookie(&http.Cookie{
					Name:  SessionCookieName,
					Value: tt.cookieValue,
				})
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
