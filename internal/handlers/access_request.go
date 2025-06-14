package handlers

import (
	"encoding/json"
	"net/http"

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
	r.HandleFunc("/api/access-requests", h.CreateAccessRequest).Methods("POST")
	r.HandleFunc("/api/access-requests", h.ListAccessRequests).Methods("GET")
	r.HandleFunc("/api/access-requests/{id}", h.GetAccessRequest).Methods("GET")
	r.HandleFunc("/api/access-requests/{id}/approve", h.ApproveAccessRequest).Methods("POST")
	r.HandleFunc("/api/access-requests/{id}/reject", h.RejectAccessRequest).Methods("POST")
}

type createAccessRequestRequest struct {
	ResourceID string `json:"resource_id"`
	UserID     string `json:"user_id"`
	Reason     string `json:"reason"`
	Duration   string `json:"duration"`
}

func (h *AccessRequestHandler) CreateAccessRequest(w http.ResponseWriter, r *http.Request) {
	var req createAccessRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	accessRequest := &domain.AccessRequest{
		ResourceID: req.ResourceID,
		UserID:     req.UserID,
		Reason:     req.Reason,
		Duration:   utils.ParseDuration(req.Duration),
	}

	if err := h.accessRequestService.CreateAccessRequest(r.Context(), accessRequest); err != nil {
		utils.InternalError(w, "Failed to create access request", err.Error())
		return
	}

	utils.SuccessResponse(w, "Access request created successfully", accessRequest)
}

func (h *AccessRequestHandler) ListAccessRequests(w http.ResponseWriter, r *http.Request) {
	accessRequests, err := h.accessRequestService.ListAccessRequests(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list access requests", err.Error())
		return
	}

	utils.SuccessResponse(w, "Access requests retrieved successfully", accessRequests)
}

func (h *AccessRequestHandler) GetAccessRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	accessRequest, err := h.accessRequestService.GetAccessRequest(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Access request not found")
		return
	}

	utils.SuccessResponse(w, "Access request retrieved successfully", accessRequest)
}

type approveAccessRequestRequest struct {
	ApproverID string `json:"approver_id"`
	Comment    string `json:"comment"`
}

func (h *AccessRequestHandler) ApproveAccessRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req approveAccessRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	accessRequest := &domain.AccessRequest{
		ID:         id,
		ApproverID: req.ApproverID,
		Comment:    req.Comment,
		Status:     "approved",
	}

	if err := h.accessRequestService.UpdateAccessRequest(r.Context(), accessRequest); err != nil {
		utils.InternalError(w, "Failed to approve access request", err.Error())
		return
	}

	utils.SuccessResponse(w, "Access request approved successfully", accessRequest)
}

type rejectAccessRequestRequest struct {
	ApproverID string `json:"approver_id"`
	Comment    string `json:"comment"`
}

func (h *AccessRequestHandler) RejectAccessRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req rejectAccessRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	accessRequest := &domain.AccessRequest{
		ID:         id,
		ApproverID: req.ApproverID,
		Comment:    req.Comment,
		Status:     "rejected",
	}

	if err := h.accessRequestService.UpdateAccessRequest(r.Context(), accessRequest); err != nil {
		utils.InternalError(w, "Failed to reject access request", err.Error())
		return
	}

	utils.SuccessResponse(w, "Access request rejected successfully", accessRequest)
}
