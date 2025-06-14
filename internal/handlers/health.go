package handlers

import (
	"context"
	"net/http"
	"time"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

type HealthHandler struct {
	db *domain.DB
}

func NewHealthHandler(db *domain.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/health", h.HealthCheck).Methods("GET")
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Create a context with timeout for DB check
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check database connection
	dbStatus := "ok"
	if err := h.db.Ping(ctx); err != nil {
		dbStatus = "error"
		utils.InternalError(w, "Service is unhealthy", map[string]interface{}{
			"status": "error",
			"components": map[string]string{
				"database": dbStatus,
			},
			"error": err.Error(),
		})
		return
	}

	utils.SuccessResponse(w, "Service is healthy", map[string]interface{}{
		"status": "ok",
		"components": map[string]string{
			"database": dbStatus,
		},
	})
}
