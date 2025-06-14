package http

import (
	"encoding/json"
	"net/http"

	"secretary/alpha/internal/domain"

	"github.com/gorilla/mux"
)

type PermissionHandler struct {
	permissionService domain.PermissionService
}

func NewPermissionHandler(permissionService domain.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService}
}

type createPermissionRequest struct {
	UserID     string `json:"user_id"`
	ResourceID string `json:"resource_id"`
	Role       string `json:"role"`
}

func (h *PermissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	permission := &domain.Permission{
		UserID:     req.UserID,
		ResourceID: req.ResourceID,
		Role:       req.Role,
	}

	if err := h.permissionService.CreatePermission(r.Context(), permission); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permission)
}

func (h *PermissionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	permission, err := h.permissionService.GetPermission(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permission)
}

func (h *PermissionHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	permissions, err := h.permissionService.GetPermissionByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
}

func (h *PermissionHandler) GetByResourceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	permissions, err := h.permissionService.GetPermissionByResourceID(r.Context(), resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
}

func (h *PermissionHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var permission domain.Permission
	if err := json.NewDecoder(r.Body).Decode(&permission); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	permission.ID = id
	if err := h.permissionService.UpdatePermission(r.Context(), &permission); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permission)
}

func (h *PermissionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.permissionService.DeletePermission(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PermissionHandler) DeleteByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := h.permissionService.DeleteByUserID(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PermissionHandler) DeleteByResourceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	if err := h.permissionService.DeleteByResourceID(r.Context(), resourceID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
