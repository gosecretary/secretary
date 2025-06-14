package middleware

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

const (
	SessionCookieName = "session_id"
)

// RateLimitMiddleware implements rate limiting for API endpoints
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow public access to login and health endpoints
		if r.URL.Path == "/api/login" || r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		// For register endpoint, check if user is admin
		if r.URL.Path == "/api/register" {
			session := r.Context().Value("session")
			if session == nil {
				utils.Unauthorized(w, "Authentication required")
				return
			}

			s, ok := session.(*domain.Session)
			if !ok || s.Username != "admin" {
				utils.Forbidden(w, "Admin access required")
				return
			}
		}

		// For all other endpoints, require authentication
		session := r.Context().Value("session")
		if session == nil {
			utils.Unauthorized(w, "Authentication required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SessionMiddleware handles session management
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip session check for login and health endpoints
		if r.URL.Path == "/api/login" || r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		// Get session from cookie
		cookie, err := r.Cookie(SessionCookieName)
		if err != nil {
			utils.Unauthorized(w, "No session cookie found")
			return
		}

		// Get session from store
		session, err := domain.GetSessionStore().Get(cookie.Value)
		if err != nil {
			utils.Unauthorized(w, "Invalid session")
			return
		}

		// Check if session is expired
		if session.ExpiresAt.Before(time.Now()) {
			utils.Unauthorized(w, "Session expired")
			return
		}

		// Add session to context
		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger is a middleware that logs all requests
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		// Log the request
		utils.GetStandardLogger().Printf(
			"%s %s %s %d %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			rw.statusCode,
			time.Since(start),
		)
	})
}

// Recovery is a middleware that recovers from panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				utils.GetStandardLogger().Printf("panic: %v\n%s", err, debug.Stack())
				utils.InternalError(w, "Internal server error", "A panic occurred")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Auth middleware for authentication
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is now handled by SessionMiddleware
		next.ServeHTTP(w, r)
	})
}

// RBAC middleware for role-based access control
func RBAC(userService domain.UserService, roles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := r.Context().Value("session")
			if session == nil {
				utils.Unauthorized(w, "Authentication required")
				return
			}

			s, ok := session.(*domain.Session)
			if !ok {
				utils.Unauthorized(w, "Invalid session")
				return
			}

			// Get user from service
			user, err := userService.GetByID(r.Context(), s.UserID)
			if err != nil {
				utils.Unauthorized(w, "User not found")
				return
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				utils.Forbidden(w, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
