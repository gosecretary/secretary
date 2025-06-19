package http

import (
	"net/http"

	"secretary/alpha/internal/domain"
	"secretary/alpha/internal/handlers"
	"secretary/alpha/internal/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(
	userService domain.UserService,
	resourceService domain.ResourceService,
	credentialService domain.CredentialService,
	permissionService domain.PermissionService,
	sessionService domain.SessionService,
	accessRequestService domain.AccessRequestService,
	ephemeralCredentialService domain.EphemeralCredentialService,
) *mux.Router {
	router := mux.NewRouter()

	// Create handlers
	userHandler := handlers.NewUserHandler(userService)
	resourceHandler := handlers.NewResourceHandler(resourceService)
	credentialHandler := handlers.NewCredentialHandler(credentialService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)

	// Create new handlers for the additional services
	sessionHandler := handlers.NewSessionHandler(sessionService)
	accessRequestHandler := handlers.NewAccessRequestHandler(accessRequestService)
	ephemeralCredentialHandler := handlers.NewEphemeralCredentialHandler(ephemeralCredentialService)

	// Public routes - ONLY health and login should be public
	router.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	// Protected routes - All other endpoints require authentication
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.SessionMiddleware)
	api.Use(middleware.RateLimitMiddleware)

	// Register endpoint - now protected (requires authentication)
	api.HandleFunc("/register", userHandler.Register).Methods("POST")

	// User routes
	api.HandleFunc("/users/{id}", userHandler.GetByID).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.Update).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.Delete).Methods("DELETE")

	// Resource routes
	api.HandleFunc("/resources", resourceHandler.Create).Methods("POST")
	api.HandleFunc("/resources", resourceHandler.List).Methods("GET")
	api.HandleFunc("/resources/{id}", resourceHandler.GetByID).Methods("GET")
	api.HandleFunc("/resources/{id}", resourceHandler.Update).Methods("PUT")
	api.HandleFunc("/resources/{id}", resourceHandler.Delete).Methods("DELETE")

	// Credential routes
	api.HandleFunc("/credentials", credentialHandler.Create).Methods("POST")
	api.HandleFunc("/credentials/{id}", credentialHandler.GetByID).Methods("GET")
	api.HandleFunc("/resources/{resource_id}/credentials", credentialHandler.GetByResourceID).Methods("GET")
	api.HandleFunc("/credentials/{id}", credentialHandler.Update).Methods("PUT")
	api.HandleFunc("/credentials/{id}", credentialHandler.Delete).Methods("DELETE")

	// Permission routes
	api.HandleFunc("/permissions", permissionHandler.Create).Methods("POST")
	api.HandleFunc("/permissions/{id}", permissionHandler.GetByID).Methods("GET")
	api.HandleFunc("/users/{user_id}/permissions", permissionHandler.GetByUserID).Methods("GET")
	api.HandleFunc("/resources/{resource_id}/permissions", permissionHandler.GetByResourceID).Methods("GET")
	api.HandleFunc("/permissions/{id}", permissionHandler.Update).Methods("PUT")
	api.HandleFunc("/permissions/{id}", permissionHandler.Delete).Methods("DELETE")
	api.HandleFunc("/users/{user_id}/permissions", permissionHandler.DeleteByUserID).Methods("DELETE")
	api.HandleFunc("/resources/{resource_id}/permissions", permissionHandler.DeleteByResourceID).Methods("DELETE")

	// Session routes
	api.HandleFunc("/sessions", sessionHandler.GetActive).Methods("GET")
	api.HandleFunc("/sessions/{id}", sessionHandler.GetByID).Methods("GET")
	api.HandleFunc("/sessions/{id}/terminate", sessionHandler.Terminate).Methods("POST")
	api.HandleFunc("/users/{user_id}/sessions", sessionHandler.GetByUserID).Methods("GET")
	api.HandleFunc("/resources/{resource_id}/sessions", sessionHandler.GetByResourceID).Methods("GET")

	// Access request routes
	api.HandleFunc("/access-requests", accessRequestHandler.Create).Methods("POST")
	api.HandleFunc("/access-requests", accessRequestHandler.GetPending).Methods("GET")
	api.HandleFunc("/access-requests/{id}", accessRequestHandler.GetByID).Methods("GET")
	api.HandleFunc("/access-requests/{id}/approve", accessRequestHandler.Approve).Methods("POST")
	api.HandleFunc("/access-requests/{id}/deny", accessRequestHandler.Deny).Methods("POST")
	api.HandleFunc("/users/{user_id}/access-requests", accessRequestHandler.GetByUserID).Methods("GET")
	api.HandleFunc("/resources/{resource_id}/access-requests", accessRequestHandler.GetByResourceID).Methods("GET")

	// Ephemeral credential routes
	api.HandleFunc("/ephemeral-credentials", ephemeralCredentialHandler.Create).Methods("POST")
	api.HandleFunc("/ephemeral-credentials/{id}", ephemeralCredentialHandler.GetByID).Methods("GET")
	api.HandleFunc("/ephemeral-credentials/{id}/use", ephemeralCredentialHandler.MarkAsUsed).Methods("POST")
	api.HandleFunc("/ephemeral-credentials/token/{token}", ephemeralCredentialHandler.GetByToken).Methods("GET")

	return router
}
