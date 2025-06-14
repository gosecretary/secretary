package handlers

import (
	"encoding/json"
	"net/http"

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
	r.HandleFunc("/api/ephemeral-credentials", h.CreateEphemeralCredential).Methods("POST")
	r.HandleFunc("/api/ephemeral-credentials", h.ListEphemeralCredentials).Methods("GET")
	r.HandleFunc("/api/ephemeral-credentials/{id}", h.GetEphemeralCredential).Methods("GET")
	r.HandleFunc("/api/ephemeral-credentials/{id}", h.DeleteEphemeralCredential).Methods("DELETE")
}

type createEphemeralCredentialRequest struct {
	ResourceID string `json:"resource_id"`
	UserID     string `json:"user_id"`
	Duration   string `json:"duration"`
}

func (h *EphemeralCredentialHandler) CreateEphemeralCredential(w http.ResponseWriter, r *http.Request) {
	var req createEphemeralCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	ephemeralCredential := &domain.EphemeralCredential{
		ResourceID: req.ResourceID,
		UserID:     req.UserID,
		Duration:   req.Duration,
	}

	createdCredential, err := h.ephemeralCredentialService.Create(r.Context(), ephemeralCredential)
	if err != nil {
		utils.InternalError(w, "Failed to create ephemeral credential", err.Error())
		return
	}

	utils.SuccessResponse(w, "Ephemeral credential created successfully", createdCredential)
}

func (h *EphemeralCredentialHandler) ListEphemeralCredentials(w http.ResponseWriter, r *http.Request) {
	ephemeralCredentials, err := h.ephemeralCredentialService.List(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list ephemeral credentials", err.Error())
		return
	}

	utils.SuccessResponse(w, "Ephemeral credentials retrieved successfully", ephemeralCredentials)
}

func (h *EphemeralCredentialHandler) GetEphemeralCredential(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ephemeralCredential, err := h.ephemeralCredentialService.GetEphemeralCredential(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Ephemeral credential not found")
		return
	}

	utils.SuccessResponse(w, "Ephemeral credential retrieved successfully", ephemeralCredential)
}

func (h *EphemeralCredentialHandler) DeleteEphemeralCredential(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.ephemeralCredentialService.DeleteEphemeralCredential(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to delete ephemeral credential", err.Error())
		return
	}

	utils.SuccessResponse(w, "Ephemeral credential deleted successfully", nil)
}
