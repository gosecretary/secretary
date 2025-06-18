package handlers

import (
	"net/http"

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
) {
	// Register routes for each handler
	authHandler.RegisterRoutes(r.router)
	resourceHandler.RegisterRoutes(r.router)
	credentialHandler.RegisterRoutes(r.router)
	permissionHandler.RegisterRoutes(r.router)
	accessRequestHandler.RegisterRoutes(r.router)

	// Add health check handler
	healthHandler := NewHealthHandler()
	healthHandler.RegisterRoutes(r.router)

	// Add documentation handler
	docsHandler := NewDocsHandler()
	docsHandler.RegisterRoutes(r.router)
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
