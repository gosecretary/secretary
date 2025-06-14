package handlers

import (
	"encoding/json"
	"net/http"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

type ResourceHandler struct {
	resourceService domain.ResourceService
}

func NewResourceHandler(resourceService domain.ResourceService) *ResourceHandler {
	return &ResourceHandler{
		resourceService: resourceService,
	}
}

func (h *ResourceHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/resources", h.CreateResource).Methods("POST")
	r.HandleFunc("/api/resources", h.ListResources).Methods("GET")
	r.HandleFunc("/api/resources/{id}", h.GetResource).Methods("GET")
	r.HandleFunc("/api/resources/{id}", h.UpdateResource).Methods("PUT")
	r.HandleFunc("/api/resources/{id}", h.DeleteResource).Methods("DELETE")
}

type createResourceRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (h *ResourceHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	var req createResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	resource := &domain.Resource{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
	}

	if err := h.resourceService.CreateResource(r.Context(), resource); err != nil {
		utils.InternalError(w, "Failed to create resource", err.Error())
		return
	}

	utils.SuccessResponse(w, "Resource created successfully", resource)
}

func (h *ResourceHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	resources, err := h.resourceService.ListResources(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list resources", err.Error())
		return
	}

	utils.SuccessResponse(w, "Resources retrieved successfully", resources)
}

func (h *ResourceHandler) GetResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resource, err := h.resourceService.GetResource(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Resource not found")
		return
	}

	utils.SuccessResponse(w, "Resource retrieved successfully", resource)
}

type updateResourceRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (h *ResourceHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req updateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	resource := &domain.Resource{
		ID:          id,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
	}

	if err := h.resourceService.UpdateResource(r.Context(), resource); err != nil {
		utils.InternalError(w, "Failed to update resource", err.Error())
		return
	}

	utils.SuccessResponse(w, "Resource updated successfully", resource)
}

func (h *ResourceHandler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.resourceService.DeleteResource(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to delete resource", err.Error())
		return
	}

	utils.SuccessResponse(w, "Resource deleted successfully", nil)
}
