package http

import (
	"encoding/json"
	"net/http"

	"secretary/alpha/internal/domain"

	"github.com/gorilla/mux"
)

type ResourceHandler struct {
	resourceService domain.ResourceService
}

func NewResourceHandler(resourceService domain.ResourceService) *ResourceHandler {
	return &ResourceHandler{resourceService: resourceService}
}

type createResourceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type updateResourceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *ResourceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resource := &domain.Resource{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
	}

	if err := h.resourceService.Create(r.Context(), resource); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

func (h *ResourceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resource, err := h.resourceService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

func (h *ResourceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resources, err := h.resourceService.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}

func (h *ResourceHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var resource domain.Resource
	if err := json.NewDecoder(r.Body).Decode(&resource); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resource.ID = id
	if err := h.resourceService.Update(r.Context(), &resource); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

func (h *ResourceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.resourceService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
