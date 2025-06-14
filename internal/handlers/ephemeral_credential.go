package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

type EphemeralCredentialHandler struct {
	ephemeralCredentialService domain.EphemeralCredentialService
}

func NewEphemeralCredentialHandler(ephemeralCredentialService domain.EphemeralCredentialService) *EphemeralCredentialHandler {
	return &EphemeralCredentialHandler{
		ephemeralCredentialService: ephemeralCredentialService,
	}
}

func (h *EphemeralCredentialHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/ephemeral-credentials", h.Create).Methods("POST")
	r.HandleFunc("/api/ephemeral-credentials", h.List).Methods("GET")
	r.HandleFunc("/api/ephemeral-credentials/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/api/ephemeral-credentials/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/api/ephemeral-credentials/token/{token}", h.GetByToken).Methods("GET")
	r.HandleFunc("/api/ephemeral-credentials/{id}/use", h.MarkAsUsed).Methods("POST")
}

type createEphemeralCredentialRequest struct {
	ResourceID string        `json:"resource_id"`
	Duration   time.Duration `json:"duration"`
}

func (h *EphemeralCredentialHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createEphemeralCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	// Get user ID from session
	session := r.Context().Value("session").(*domain.Session)
	if session == nil {
		utils.Unauthorized(w, "No active session")
		return
	}

	credential := &domain.EphemeralCredential{
		UserID:     session.UserID,
		ResourceID: req.ResourceID,
		ExpiresAt:  time.Now().Add(req.Duration),
	}

	createdCredential, err := h.ephemeralCredentialService.Create(r.Context(), credential)
	if err != nil {
		utils.InternalError(w, "Failed to create ephemeral credential", err.Error())
		return
	}

	utils.SuccessResponse(w, "Ephemeral credential created successfully", createdCredential)
}

func (h *EphemeralCredentialHandler) List(w http.ResponseWriter, r *http.Request) {
	credentials, err := h.ephemeralCredentialService.List(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list ephemeral credentials", err.Error())
		return
	}

	utils.SuccessResponse(w, "Ephemeral credentials retrieved successfully", credentials)
}

func (h *EphemeralCredentialHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	credential, err := h.ephemeralCredentialService.GetEphemeralCredential(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Ephemeral credential not found")
		return
	}

	utils.SuccessResponse(w, "Ephemeral credential retrieved successfully", credential)
}

func (h *EphemeralCredentialHandler) GetByToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	credential, err := h.ephemeralCredentialService.GetEphemeralCredential(r.Context(), token)
	if err != nil {
		utils.NotFound(w, "Ephemeral credential not found")
		return
	}

	utils.SuccessResponse(w, "Ephemeral credential retrieved successfully", credential)
}

func (h *EphemeralCredentialHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.ephemeralCredentialService.DeleteEphemeralCredential(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to delete ephemeral credential", err.Error())
		return
	}

	utils.SuccessResponse(w, "Ephemeral credential deleted successfully", nil)
}

func (h *EphemeralCredentialHandler) MarkAsUsed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.ephemeralCredentialService.MarkAsUsedEphemeralCredential(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to mark ephemeral credential as used", err.Error())
		return
	}

	utils.SuccessResponse(w, "Ephemeral credential marked as used successfully", nil)
}
