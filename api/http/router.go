package http

import (
	"net/http"

	"secretary/alpha/internal/domain"
	"secretary/alpha/internal/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(
	userService domain.UserService,
	resourceService domain.ResourceService,
	credentialService domain.CredentialService,
	permissionService domain.PermissionService,
) *mux.Router {
	router := mux.NewRouter()

	// Create handlers
	userHandler := NewUserHandler(userService)
	resourceHandler := NewResourceHandler(resourceService)
	credentialHandler := NewCredentialHandler(credentialService)
	permissionHandler := NewPermissionHandler(permissionService)

	// Public routes
	router.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	// Protected routes
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.Auth)

	// User routes
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Resource routes
	api.HandleFunc("/resources", resourceHandler.Create).Methods("POST")
	api.HandleFunc("/resources", resourceHandler.GetAll).Methods("GET")
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

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	return router
} 