package http

import (
	"encoding/json"
	"net/http"

	"secretary/alpha/internal/domain"

	"github.com/gorilla/mux"
)

type CredentialHandler struct {
	credentialService domain.CredentialService
}

func NewCredentialHandler(credentialService domain.CredentialService) *CredentialHandler {
	return &CredentialHandler{credentialService: credentialService}
}

type createCredentialRequest struct {
	ResourceID string `json:"resource_id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type updateCredentialRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *CredentialHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	credential := &domain.Credential{
		ResourceID: req.ResourceID,
		Username:   req.Username,
		Password:   req.Password,
	}

	if err := h.credentialService.Create(r.Context(), credential); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credential)
}

func (h *CredentialHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	credential, err := h.credentialService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credential)
}

func (h *CredentialHandler) GetByResourceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	credentials, err := h.credentialService.GetByResourceID(r.Context(), resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credentials)
}

func (h *CredentialHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var credential domain.Credential
	if err := json.NewDecoder(r.Body).Decode(&credential); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	credential.ID = id
	if err := h.credentialService.Update(r.Context(), &credential); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credential)
}

func (h *CredentialHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.credentialService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
