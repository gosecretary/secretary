package handlers

import (
	"encoding/json"
	"net/http"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

type CredentialHandler struct {
	credentialService domain.CredentialService
}

func NewCredentialHandler(credentialService domain.CredentialService) *CredentialHandler {
	return &CredentialHandler{
		credentialService: credentialService,
	}
}

func (h *CredentialHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/credentials", h.CreateCredential).Methods("POST")
	r.HandleFunc("/api/credentials", h.ListCredentials).Methods("GET")
	r.HandleFunc("/api/credentials/{id}", h.GetCredential).Methods("GET")
	r.HandleFunc("/api/credentials/{id}", h.UpdateCredential).Methods("PUT")
	r.HandleFunc("/api/credentials/{id}", h.DeleteCredential).Methods("DELETE")
	r.HandleFunc("/api/credentials/{id}/rotate", h.RotateCredential).Methods("POST")
}

func (h *CredentialHandler) CreateCredential(w http.ResponseWriter, r *http.Request) {
	var req domain.Credential
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	credential := &domain.Credential{
		ResourceID: req.ResourceID,
		Type:       req.Type,
		Secret:     req.Secret,
	}

	if err := h.credentialService.CreateCredential(r.Context(), credential); err != nil {
		utils.InternalError(w, "Failed to create credential", err.Error())
		return
	}

	utils.SuccessResponse(w, "Credential created successfully", credential)
}

func (h *CredentialHandler) ListCredentials(w http.ResponseWriter, r *http.Request) {
	credentials, err := h.credentialService.ListCredentials(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list credentials", err.Error())
		return
	}

	utils.SuccessResponse(w, "Credentials retrieved successfully", credentials)
}

func (h *CredentialHandler) GetCredential(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	credential, err := h.credentialService.GetCredential(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Credential not found")
		return
	}

	utils.SuccessResponse(w, "Credential retrieved successfully", credential)
}

type updateCredentialRequest struct {
	ResourceID string `json:"resource_id"`
	Type       string `json:"type"`
	Secret     string `json:"secret"`
}

func (h *CredentialHandler) UpdateCredential(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req updateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	credential := &domain.Credential{
		ID:         id,
		ResourceID: req.ResourceID,
		Type:       req.Type,
		Secret:     req.Secret,
	}

	if err := h.credentialService.UpdateCredential(r.Context(), credential); err != nil {
		utils.InternalError(w, "Failed to update credential", err.Error())
		return
	}

	utils.SuccessResponse(w, "Credential updated successfully", credential)
}

func (h *CredentialHandler) DeleteCredential(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.credentialService.DeleteCredential(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to delete credential", err.Error())
		return
	}

	utils.SuccessResponse(w, "Credential deleted successfully", nil)
}

func (h *CredentialHandler) RotateCredential(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.credentialService.RotateCredential(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to rotate credential", err.Error())
		return
	}

	utils.SuccessResponse(w, "Credential rotated successfully", nil)
}
