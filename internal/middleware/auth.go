package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/sessions"

	"secretary/alpha/internal"
	"secretary/alpha/internal/audit"
	"secretary/alpha/internal/config"
	"secretary/alpha/utils"
)

var (
	store         *sessions.CookieStore
	rateLimitMap  = make(map[string]*RateLimit)
	rateLimitMux  sync.RWMutex
	csrfTokens    = make(map[string]time.Time)
	csrfTokensMux sync.RWMutex
)

type RateLimit struct {
	requests []time.Time
	mutex    sync.Mutex
}

type contextKey string

const (
	UserContextKey contextKey = "user"
	CSRFContextKey contextKey = "csrf_token"
)

func InitializeMiddleware() {
	if config.GlobalConfig == nil {
		panic("Config not loaded before middleware initialization")
	}

	store = sessions.NewCookieStore(config.GlobalConfig.Security.SessionSecret)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   config.GlobalConfig.Security.SessionMaxAge,
		HttpOnly: true,
		Secure:   config.GlobalConfig.Security.SecureCookies,
		SameSite: http.SameSiteStrictMode,
	}

	// Start cleanup goroutine for expired CSRF tokens
	go cleanupExpiredCSRFTokens()
}

func generateCSRFToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func cleanupExpiredCSRFTokens() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		csrfTokensMux.Lock()
		now := time.Now()
		for token, expiry := range csrfTokens {
			if now.After(expiry) {
				delete(csrfTokens, token)
			}
		}
		csrfTokensMux.Unlock()
	}
}

func SetSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")

	// Only set HSTS if we're using HTTPS
	if config.GlobalConfig.Server.TLSCertPath != "" {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		rateLimitMux.Lock()
		rateLimit, exists := rateLimitMap[clientIP]
		if !exists {
			rateLimit = &RateLimit{
				requests: make([]time.Time, 0),
			}
			rateLimitMap[clientIP] = rateLimit
		}
		rateLimitMux.Unlock()

		rateLimit.mutex.Lock()
		defer rateLimit.mutex.Unlock()

		now := time.Now()
		windowStart := now.Add(-config.GlobalConfig.Security.RateLimitWindow)

		// Remove old requests outside the window
		validRequests := make([]time.Time, 0)
		for _, reqTime := range rateLimit.requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}
		rateLimit.requests = validRequests

		// Check if rate limit exceeded
		if len(rateLimit.requests) >= config.GlobalConfig.Security.RateLimitRequests {
			audit.Audit(fmt.Sprintf("[security] [rate_limit_exceeded] IP: %s", clientIP))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Add current request
		rateLimit.requests = append(rateLimit.requests, now)

		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take the first IP (before any comma)
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetSecurityHeaders(w)

		session, err := store.Get(r, "secretary_session")
		if err != nil {
			utils.Logger("err", fmt.Sprintf("Session error: %v", err))
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		// Check if user is authenticated
		authenticated, ok := session.Values["authenticated"].(bool)
		if !ok || !authenticated {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check session expiration
		sessionTime, ok := session.Values["created_at"].(int64)
		if !ok {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		if time.Now().Unix()-sessionTime > int64(config.GlobalConfig.Security.SessionMaxAge) {
			// Session expired, clear it
			session.Values = make(map[interface{}]interface{})
			session.Save(r, w)
			audit.Audit(fmt.Sprintf("[security] [session_expired] IP: %s", getClientIP(r)))
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		// Get user from session
		username, ok := session.Values["username"].(string)
		if !ok {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// Load user data
		user := &internal.User{}
		user = user.GetUser(username)
		if user == nil || !user.Active {
			audit.Audit(fmt.Sprintf("[security] [invalid_user_access] user: %s, IP: %s", username, getClientIP(r)))
			http.Error(w, "User not found or inactive", http.StatusUnauthorized)
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			// Generate CSRF token for safe methods
			token := generateCSRFToken()
			csrfTokensMux.Lock()
			csrfTokens[token] = time.Now().Add(30 * time.Minute)
			csrfTokensMux.Unlock()

			ctx := context.WithValue(r.Context(), CSRFContextKey, token)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Validate CSRF token for unsafe methods
		token := r.Header.Get("X-CSRF-Token")
		if token == "" {
			audit.Audit(fmt.Sprintf("[security] [csrf_missing] IP: %s, path: %s", getClientIP(r), r.URL.Path))
			http.Error(w, "CSRF token required", http.StatusForbidden)
			return
		}

		csrfTokensMux.RLock()
		expiry, exists := csrfTokens[token]
		csrfTokensMux.RUnlock()

		if !exists || time.Now().After(expiry) {
			audit.Audit(fmt.Sprintf("[security] [csrf_invalid] IP: %s, path: %s", getClientIP(r), r.URL.Path))
			http.Error(w, "Invalid or expired CSRF token", http.StatusForbidden)
			return
		}

		// Remove used token
		csrfTokensMux.Lock()
		delete(csrfTokens, token)
		csrfTokensMux.Unlock()

		next.ServeHTTP(w, r)
	})
}

func PublicEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetSecurityHeaders(w)
		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(ctx context.Context) *internal.User {
	if user, ok := ctx.Value(UserContextKey).(*internal.User); ok {
		return user
	}
	return nil
}

func GetCSRFTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(CSRFContextKey).(string); ok {
		return token
	}
	return ""
}
