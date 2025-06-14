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
	r.HandleFunc("/api/credentials", h.Create).Methods("POST")
	r.HandleFunc("/api/credentials", h.List).Methods("GET")
	r.HandleFunc("/api/credentials/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/api/credentials/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/api/credentials/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/api/credentials/{id}/rotate", h.Rotate).Methods("POST")
	r.HandleFunc("/api/credentials/resource/{resource_id}", h.GetByResourceID).Methods("GET")
}

type createCredentialRequest struct {
	ResourceID string `json:"resource_id"`
	Type       string `json:"type"`
	Secret     string `json:"secret"`
}

func (h *CredentialHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCredentialRequest
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

func (h *CredentialHandler) List(w http.ResponseWriter, r *http.Request) {
	credentials, err := h.credentialService.ListCredentials(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list credentials", err.Error())
		return
	}

	utils.SuccessResponse(w, "Credentials retrieved successfully", credentials)
}

func (h *CredentialHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
	Type   string `json:"type,omitempty"`
	Secret string `json:"secret,omitempty"`
}

func (h *CredentialHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req updateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	credential, err := h.credentialService.GetCredential(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Credential not found")
		return
	}

	if req.Type != "" {
		credential.Type = req.Type
	}
	if req.Secret != "" {
		credential.Secret = req.Secret
	}

	if err := h.credentialService.UpdateCredential(r.Context(), credential); err != nil {
		utils.InternalError(w, "Failed to update credential", err.Error())
		return
	}

	utils.SuccessResponse(w, "Credential updated successfully", credential)
}

func (h *CredentialHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.credentialService.DeleteCredential(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to delete credential", err.Error())
		return
	}

	utils.SuccessResponse(w, "Credential deleted successfully", nil)
}

func (h *CredentialHandler) Rotate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.credentialService.RotateCredential(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to rotate credential", err.Error())
		return
	}

	// Get updated credential
	credential, err := h.credentialService.GetCredential(r.Context(), id)
	if err != nil {
		utils.InternalError(w, "Failed to get updated credential", err.Error())
		return
	}

	utils.SuccessResponse(w, "Credential rotated successfully", credential)
}

func (h *CredentialHandler) GetByResourceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	credentials, err := h.credentialService.GetCredentialByResourceID(r.Context(), resourceID)
	if err != nil {
		utils.InternalError(w, "Failed to get resource credentials", err.Error())
		return
	}

	utils.SuccessResponse(w, "Resource credentials retrieved successfully", credentials)
}
