package handlers

import (
	"net/http"

	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/health", h.HealthCheck).Methods("GET")
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"
	// DB check is omitted for now

	utils.SuccessResponse(w, "Service is healthy", map[string]interface{}{
		"status": "ok",
		"components": map[string]string{
			"database": dbStatus,
		},
	})
}
