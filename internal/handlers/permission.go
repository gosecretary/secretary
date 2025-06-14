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
	r.HandleFunc("/api/permissions", h.Create).Methods("POST")
	r.HandleFunc("/api/permissions", h.List).Methods("GET")
	r.HandleFunc("/api/permissions/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/api/permissions/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/api/permissions/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/api/permissions/user/{user_id}", h.GetByUserID).Methods("GET")
	r.HandleFunc("/api/permissions/resource/{resource_id}", h.GetByResourceID).Methods("GET")
	r.HandleFunc("/api/permissions/user/{user_id}", h.DeleteByUserID).Methods("DELETE")
	r.HandleFunc("/api/permissions/resource/{resource_id}", h.DeleteByResourceID).Methods("DELETE")
}

type createPermissionRequest struct {
	UserID     string `json:"user_id"`
	ResourceID string `json:"resource_id"`
	Role       string `json:"role"`
	Action     string `json:"action"`
}

func (h *PermissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	permission := &domain.Permission{
		UserID:     req.UserID,
		ResourceID: req.ResourceID,
		Role:       req.Role,
		Action:     req.Action,
	}

	if err := h.permissionService.CreatePermission(r.Context(), permission); err != nil {
		utils.InternalError(w, "Failed to create permission", err.Error())
		return
	}

	utils.SuccessResponse(w, "Permission created successfully", permission)
}

func (h *PermissionHandler) List(w http.ResponseWriter, r *http.Request) {
	permissions, err := h.permissionService.ListPermissions(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list permissions", err.Error())
		return
	}

	utils.SuccessResponse(w, "Permissions retrieved successfully", permissions)
}

func (h *PermissionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
	Role   string `json:"role,omitempty"`
	Action string `json:"action,omitempty"`
}

func (h *PermissionHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req updatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	permission, err := h.permissionService.GetPermission(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Permission not found")
		return
	}

	if req.Role != "" {
		permission.Role = req.Role
	}
	if req.Action != "" {
		permission.Action = req.Action
	}

	if err := h.permissionService.UpdatePermission(r.Context(), permission); err != nil {
		utils.InternalError(w, "Failed to update permission", err.Error())
		return
	}

	utils.SuccessResponse(w, "Permission updated successfully", permission)
}

func (h *PermissionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.permissionService.DeletePermission(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to delete permission", err.Error())
		return
	}

	utils.SuccessResponse(w, "Permission deleted successfully", nil)
}

func (h *PermissionHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	permissions, err := h.permissionService.GetPermissionByUserID(r.Context(), userID)
	if err != nil {
		utils.InternalError(w, "Failed to get user permissions", err.Error())
		return
	}

	utils.SuccessResponse(w, "User permissions retrieved successfully", permissions)
}

func (h *PermissionHandler) GetByResourceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	permissions, err := h.permissionService.GetPermissionByResourceID(r.Context(), resourceID)
	if err != nil {
		utils.InternalError(w, "Failed to get resource permissions", err.Error())
		return
	}

	utils.SuccessResponse(w, "Resource permissions retrieved successfully", permissions)
}

func (h *PermissionHandler) DeleteByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := h.permissionService.DeleteByUserID(r.Context(), userID); err != nil {
		utils.InternalError(w, "Failed to delete permissions by user ID", err.Error())
		return
	}

	utils.SuccessResponse(w, "Permissions deleted successfully", nil)
}

func (h *PermissionHandler) DeleteByResourceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	if err := h.permissionService.DeleteByResourceID(r.Context(), resourceID); err != nil {
		utils.InternalError(w, "Failed to delete permissions by resource ID", err.Error())
		return
	}

	utils.SuccessResponse(w, "Permissions deleted successfully", nil)
}
