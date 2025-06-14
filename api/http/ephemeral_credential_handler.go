package http

import (
	"encoding/json"
	"net/http"
	"time"

	"secretary/alpha/internal/domain"

	"github.com/gorilla/mux"
)

type EphemeralCredentialHandler struct {
	ephemeralCredentialService domain.EphemeralCredentialService
}

func NewEphemeralCredentialHandler(ephemeralCredentialService domain.EphemeralCredentialService) *EphemeralCredentialHandler {
	return &EphemeralCredentialHandler{ephemeralCredentialService: ephemeralCredentialService}
}

type generateCredentialRequest struct {
	UserID     string        `json:"user_id"`
	ResourceID string        `json:"resource_id"`
	Duration   time.Duration `json:"duration,omitempty"` // in seconds
}

type ephemeralCredentialResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ResourceID string    `json:"resource_id"`
	Username   string    `json:"username"`
	Password   string    `json:"password,omitempty"` // Only included in Generate response
	Token      string    `json:"token,omitempty"`    // Only included in Generate response
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	UsedAt     time.Time `json:"used_at,omitempty"`
}

func (h *EphemeralCredentialHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req generateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert duration from seconds to time.Duration
	var duration time.Duration
	if req.Duration > 0 {
		duration = req.Duration * time.Second
	}

	credential := &domain.EphemeralCredential{
		UserID:     req.UserID,
		ResourceID: req.ResourceID,
		ExpiresAt:  time.Now().Add(duration),
	}

	credential, err := h.ephemeralCredentialService.Create(r.Context(), credential)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For credential generation, we include the sensitive fields
	response := toEphemeralCredentialResponse(credential, true)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *EphemeralCredentialHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	credential, err := h.ephemeralCredentialService.GetEphemeralCredential(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// For retrieval, we do not include sensitive fields
	response := toEphemeralCredentialResponse(credential, false)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *EphemeralCredentialHandler) GetByToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	credential, err := h.ephemeralCredentialService.GetEphemeralCredential(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// For retrieval, we do not include sensitive fields
	response := toEphemeralCredentialResponse(credential, false)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *EphemeralCredentialHandler) MarkAsUsed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.ephemeralCredentialService.MarkAsUsedEphemeralCredential(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated credential
	credential, err := h.ephemeralCredentialService.GetEphemeralCredential(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For marking as used, we do not include sensitive fields
	response := toEphemeralCredentialResponse(credential, false)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to convert domain model to response
func toEphemeralCredentialResponse(ec *domain.EphemeralCredential, includeSensitive bool) ephemeralCredentialResponse {
	response := ephemeralCredentialResponse{
		ID:         ec.ID,
		UserID:     ec.UserID,
		ResourceID: ec.ResourceID,
		Username:   ec.Username,
		ExpiresAt:  ec.ExpiresAt,
		CreatedAt:  ec.CreatedAt,
		UsedAt:     ec.UsedAt,
	}

	// Only include sensitive fields when generating a new credential
	if includeSensitive {
		response.Password = ec.Password
		response.Token = ec.Token
	}

	return response
}
