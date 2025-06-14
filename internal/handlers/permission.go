package handlers

import (
	"encoding/json"
	"net/http"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

type PermissionHandler struct {
	permissionService domain.PermissionService
}

func NewPermissionHandler(permissionService domain.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

func (h *PermissionHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/permissions", h.CreatePermission).Methods("POST")
	r.HandleFunc("/api/permissions", h.ListPermissions).Methods("GET")
	r.HandleFunc("/api/permissions/{id}", h.GetPermission).Methods("GET")
	r.HandleFunc("/api/permissions/{id}", h.UpdatePermission).Methods("PUT")
	r.HandleFunc("/api/permissions/{id}", h.DeletePermission).Methods("DELETE")
}

type createPermissionRequest struct {
	ResourceID string `json:"resource_id"`
	UserID     string `json:"user_id"`
	Role       string `json:"role"`
}

func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req createPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	permission := &domain.Permission{
		ResourceID: req.ResourceID,
		UserID:     req.UserID,
		Role:       req.Role,
	}

	if err := h.permissionService.CreatePermission(r.Context(), permission); err != nil {
		utils.InternalError(w, "Failed to create permission", err.Error())
		return
	}

	utils.SuccessResponse(w, "Permission created successfully", permission)
}

func (h *PermissionHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, err := h.permissionService.ListPermissions(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list permissions", err.Error())
		return
	}

	utils.SuccessResponse(w, "Permissions retrieved successfully", permissions)
}

func (h *PermissionHandler) GetPermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	permission, err := h.permissionService.GetPermission(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Permission not found")
		return
	}

	utils.SuccessResponse(w, "Permission retrieved successfully", permission)
}

type updatePermissionRequest struct {
	ResourceID string `json:"resource_id"`
	UserID     string `json:"user_id"`
	Role       string `json:"role"`
}

func (h *PermissionHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req updatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	permission := &domain.Permission{
		ID:         id,
		ResourceID: req.ResourceID,
		UserID:     req.UserID,
		Role:       req.Role,
	}

	if err := h.permissionService.UpdatePermission(r.Context(), permission); err != nil {
		utils.InternalError(w, "Failed to update permission", err.Error())
		return
	}

	utils.SuccessResponse(w, "Permission updated successfully", permission)
}

func (h *PermissionHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.permissionService.DeletePermission(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to delete permission", err.Error())
		return
	}

	utils.SuccessResponse(w, "Permission deleted successfully", nil)
}
