package handlers

import (
	"net/http"

	"secretary/alpha/internal/middleware"

	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
}

func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
	}
}

func (r *Router) RegisterHandlers(
	authHandler *UserHandler,
	resourceHandler *ResourceHandler,
	credentialHandler *CredentialHandler,
	permissionHandler *PermissionHandler,
	accessRequestHandler *AccessRequestHandler,
	sessionHandler *SessionHandler,
	ephemeralCredentialHandler *EphemeralCredentialHandler,
	sessionMonitorHandler *SessionMonitorHandler,
) {
	// PUBLIC ROUTES - No authentication required (ONLY health and login)
	// Health check
	healthHandler := NewHealthHandler()
	r.router.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	// Login endpoint - public
	r.router.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	// PROTECTED ROUTES - Require authentication
	api := r.router.PathPrefix("/api").Subrouter()
	api.Use(middleware.SessionMiddleware)
	api.Use(middleware.RateLimitMiddleware)

	// Register endpoint - protected (requires authentication)
	// In production, this should require admin role
	api.HandleFunc("/register", authHandler.Register).Methods("POST")

	// User routes (authenticated users can manage their own profile)
	api.HandleFunc("/users/{id}", authHandler.GetByID).Methods("GET")
	api.HandleFunc("/users/{id}", authHandler.Update).Methods("PUT")
	api.HandleFunc("/users/{id}", authHandler.Delete).Methods("DELETE")

	// Resource routes
	api.HandleFunc("/resources", resourceHandler.Create).Methods("POST")
	api.HandleFunc("/resources", resourceHandler.List).Methods("GET")
	api.HandleFunc("/resources/{id}", resourceHandler.GetByID).Methods("GET")
	api.HandleFunc("/resources/{id}", resourceHandler.Update).Methods("PUT")
	api.HandleFunc("/resources/{id}", resourceHandler.Delete).Methods("DELETE")

	// Credential routes
	api.HandleFunc("/credentials", credentialHandler.Create).Methods("POST")
	api.HandleFunc("/credentials", credentialHandler.List).Methods("GET")
	api.HandleFunc("/credentials/{id}", credentialHandler.GetByID).Methods("GET")
	api.HandleFunc("/credentials/{id}", credentialHandler.Update).Methods("PUT")
	api.HandleFunc("/credentials/{id}", credentialHandler.Delete).Methods("DELETE")
	api.HandleFunc("/credentials/{id}/rotate", credentialHandler.Rotate).Methods("POST")
	api.HandleFunc("/credentials/resource/{resource_id}", credentialHandler.GetByResourceID).Methods("GET")

	// Permission routes
	api.HandleFunc("/permissions", permissionHandler.Create).Methods("POST")
	api.HandleFunc("/permissions", permissionHandler.List).Methods("GET")
	api.HandleFunc("/permissions/{id}", permissionHandler.GetByID).Methods("GET")
	api.HandleFunc("/permissions/{id}", permissionHandler.Update).Methods("PUT")
	api.HandleFunc("/permissions/{id}", permissionHandler.Delete).Methods("DELETE")
	api.HandleFunc("/permissions/user/{user_id}", permissionHandler.GetByUserID).Methods("GET")
	api.HandleFunc("/permissions/resource/{resource_id}", permissionHandler.GetByResourceID).Methods("GET")
	api.HandleFunc("/permissions/user/{user_id}", permissionHandler.DeleteByUserID).Methods("DELETE")
	api.HandleFunc("/permissions/resource/{resource_id}", permissionHandler.DeleteByResourceID).Methods("DELETE")

	// Access request routes
	api.HandleFunc("/access-requests", accessRequestHandler.Create).Methods("POST")
	api.HandleFunc("/access-requests", accessRequestHandler.List).Methods("GET")
	api.HandleFunc("/access-requests/{id}", accessRequestHandler.GetByID).Methods("GET")
	api.HandleFunc("/access-requests/{id}", accessRequestHandler.Update).Methods("PUT")
	api.HandleFunc("/access-requests/{id}/approve", accessRequestHandler.Approve).Methods("POST")
	api.HandleFunc("/access-requests/{id}/deny", accessRequestHandler.Deny).Methods("POST")
	api.HandleFunc("/access-requests/user/{user_id}", accessRequestHandler.GetByUserID).Methods("GET")
	api.HandleFunc("/access-requests/resource/{resource_id}", accessRequestHandler.GetByResourceID).Methods("GET")
	api.HandleFunc("/access-requests/pending", accessRequestHandler.GetPending).Methods("GET")

	// Session routes
	api.HandleFunc("/sessions", sessionHandler.Create).Methods("POST")
	api.HandleFunc("/sessions", sessionHandler.List).Methods("GET")
	api.HandleFunc("/sessions/{id}", sessionHandler.GetByID).Methods("GET")
	api.HandleFunc("/sessions/{id}", sessionHandler.Update).Methods("PUT")
	api.HandleFunc("/sessions/{id}", sessionHandler.Delete).Methods("DELETE")
	api.HandleFunc("/sessions/{id}/terminate", sessionHandler.Terminate).Methods("POST")
	api.HandleFunc("/sessions/user/{user_id}", sessionHandler.GetByUserID).Methods("GET")
	api.HandleFunc("/sessions/resource/{resource_id}", sessionHandler.GetByResourceID).Methods("GET")
	api.HandleFunc("/sessions/active", sessionHandler.GetActive).Methods("GET")

	// Ephemeral credential routes
	api.HandleFunc("/ephemeral-credentials", ephemeralCredentialHandler.Create).Methods("POST")
	api.HandleFunc("/ephemeral-credentials", ephemeralCredentialHandler.List).Methods("GET")
	api.HandleFunc("/ephemeral-credentials/{id}", ephemeralCredentialHandler.GetByID).Methods("GET")
	api.HandleFunc("/ephemeral-credentials/{id}", ephemeralCredentialHandler.Delete).Methods("DELETE")
	api.HandleFunc("/ephemeral-credentials/{id}/use", ephemeralCredentialHandler.MarkAsUsed).Methods("POST")
	api.HandleFunc("/ephemeral-credentials/token/{token}", ephemeralCredentialHandler.GetByToken).Methods("GET")

	// Session monitoring routes
	sessionMonitorHandler.RegisterRoutes(api)

	// Add documentation handler - protected for security
	docsHandler := NewDocsHandler()
	api.HandleFunc("/docs", docsHandler.SwaggerUI).Methods("GET")
	api.HandleFunc("/docs/", docsHandler.SwaggerUI).Methods("GET")
	api.HandleFunc("/docs/swagger.yaml", docsHandler.SwaggerYAML).Methods("GET")
}

func (r *Router) Use(middleware func(http.Handler) http.Handler) {
	r.router.Use(middleware)
}

func (r *Router) Start(addr string) error {
	return http.ListenAndServe(addr, r.router)
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
