package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

type AccessRequestHandler struct {
	accessRequestService domain.AccessRequestService
}

func NewAccessRequestHandler(accessRequestService domain.AccessRequestService) *AccessRequestHandler {
	return &AccessRequestHandler{
		accessRequestService: accessRequestService,
	}
}

func (h *AccessRequestHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/access-requests", h.Create).Methods("POST")
	r.HandleFunc("/api/access-requests", h.List).Methods("GET")
	r.HandleFunc("/api/access-requests/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/api/access-requests/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/api/access-requests/{id}/approve", h.Approve).Methods("POST")
	r.HandleFunc("/api/access-requests/{id}/deny", h.Deny).Methods("POST")
	r.HandleFunc("/api/access-requests/user/{user_id}", h.GetByUserID).Methods("GET")
	r.HandleFunc("/api/access-requests/resource/{resource_id}", h.GetByResourceID).Methods("GET")
	r.HandleFunc("/api/access-requests/pending", h.GetPending).Methods("GET")
}

type createAccessRequestRequest struct {
	ResourceID string        `json:"resource_id"`
	Reason     string        `json:"reason"`
	Duration   time.Duration `json:"duration"`
}

type approveAccessRequestRequest struct {
	Notes     string    `json:"notes"`
	ExpiresAt time.Time `json:"expires_at"`
}

type denyAccessRequestRequest struct {
	Notes string `json:"notes"`
}

func (h *AccessRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAccessRequestRequest
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

	request := &domain.AccessRequest{
		UserID:      session.UserID,
		ResourceID:  req.ResourceID,
		Reason:      req.Reason,
		Duration:    req.Duration,
		Status:      "pending",
		RequestedAt: time.Now(),
		ExpiresAt:   time.Now().Add(req.Duration),
	}

	if err := h.accessRequestService.CreateAccessRequest(r.Context(), request); err != nil {
		utils.InternalError(w, "Failed to create access request", err.Error())
		return
	}

	utils.SuccessResponse(w, "Access request created successfully", request)
}

func (h *AccessRequestHandler) List(w http.ResponseWriter, r *http.Request) {
	requests, err := h.accessRequestService.ListAccessRequests(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list access requests", err.Error())
		return
	}

	utils.SuccessResponse(w, "Access requests retrieved successfully", requests)
}

func (h *AccessRequestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	request, err := h.accessRequestService.GetAccessRequest(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Access request not found")
		return
	}

	utils.SuccessResponse(w, "Access request retrieved successfully", request)
}

func (h *AccessRequestHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var request domain.AccessRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	request.ID = id
	if err := h.accessRequestService.UpdateAccessRequest(r.Context(), &request); err != nil {
		utils.InternalError(w, "Failed to update access request", err.Error())
		return
	}

	utils.SuccessResponse(w, "Access request updated successfully", request)
}

func (h *AccessRequestHandler) Approve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req approveAccessRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	// Get reviewer ID from session
	session := r.Context().Value("session").(*domain.Session)
	if session == nil {
		utils.Unauthorized(w, "No active session")
		return
	}

	if err := h.accessRequestService.Approve(r.Context(), id, session.UserID, req.Notes, req.ExpiresAt); err != nil {
		utils.InternalError(w, "Failed to approve access request", err.Error())
		return
	}

	utils.SuccessResponse(w, "Access request approved successfully", nil)
}

func (h *AccessRequestHandler) Deny(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req denyAccessRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	// Get reviewer ID from session
	session := r.Context().Value("session").(*domain.Session)
	if session == nil {
		utils.Unauthorized(w, "No active session")
		return
	}

	if err := h.accessRequestService.Deny(r.Context(), id, session.UserID, req.Notes); err != nil {
		utils.InternalError(w, "Failed to deny access request", err.Error())
		return
	}

	utils.SuccessResponse(w, "Access request denied successfully", nil)
}

func (h *AccessRequestHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	requests, err := h.accessRequestService.GetAccessRequestByUserID(r.Context(), userID)
	if err != nil {
		utils.InternalError(w, "Failed to get user access requests", err.Error())
		return
	}

	utils.SuccessResponse(w, "User access requests retrieved successfully", requests)
}

func (h *AccessRequestHandler) GetByResourceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	requests, err := h.accessRequestService.GetAccessRequestByResourceID(r.Context(), resourceID)
	if err != nil {
		utils.InternalError(w, "Failed to get resource access requests", err.Error())
		return
	}

	utils.SuccessResponse(w, "Resource access requests retrieved successfully", requests)
}

func (h *AccessRequestHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	requests, err := h.accessRequestService.GetPendingAccessRequests(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to get pending access requests", err.Error())
		return
	}

	utils.SuccessResponse(w, "Pending access requests retrieved successfully", requests)
}
