package http

import (
	"encoding/json"
	"net/http"
	"time"

	"secretary/alpha/internal/domain"

	"github.com/gorilla/mux"
)

type AccessRequestHandler struct {
	accessRequestService domain.AccessRequestService
}

func NewAccessRequestHandler(accessRequestService domain.AccessRequestService) *AccessRequestHandler {
	return &AccessRequestHandler{accessRequestService: accessRequestService}
}

type createAccessRequestRequest struct {
	UserID     string `json:"user_id"`
	ResourceID string `json:"resource_id"`
	Reason     string `json:"reason"`
}

type approveAccessRequestRequest struct {
	ReviewerID string    `json:"reviewer_id"`
	Notes      string    `json:"notes"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type denyAccessRequestRequest struct {
	ReviewerID string `json:"reviewer_id"`
	Notes      string `json:"notes"`
}

type accessRequestResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ResourceID  string    `json:"resource_id"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"`
	ReviewerID  string    `json:"reviewer_id,omitempty"`
	ReviewNotes string    `json:"review_notes,omitempty"`
	RequestedAt time.Time `json:"requested_at"`
	ReviewedAt  time.Time `json:"reviewed_at,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (h *AccessRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAccessRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	accessRequest := &domain.AccessRequest{
		UserID:     req.UserID,
		ResourceID: req.ResourceID,
		Reason:     req.Reason,
	}

	if err := h.accessRequestService.CreateAccessRequest(r.Context(), accessRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toAccessRequestResponse(accessRequest))
}

func (h *AccessRequestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	accessRequest, err := h.accessRequestService.GetAccessRequest(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toAccessRequestResponse(accessRequest))
}

func (h *AccessRequestHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	accessRequests, err := h.accessRequestService.GetAccessRequestByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]accessRequestResponse, len(accessRequests))
	for i, req := range accessRequests {
		responses[i] = toAccessRequestResponse(req)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *AccessRequestHandler) GetByResourceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	accessRequests, err := h.accessRequestService.GetAccessRequestByResourceID(r.Context(), resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]accessRequestResponse, len(accessRequests))
	for i, req := range accessRequests {
		responses[i] = toAccessRequestResponse(req)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *AccessRequestHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	accessRequests, err := h.accessRequestService.GetPendingAccessRequests(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]accessRequestResponse, len(accessRequests))
	for i, req := range accessRequests {
		responses[i] = toAccessRequestResponse(req)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *AccessRequestHandler) Approve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req approveAccessRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.accessRequestService.Approve(r.Context(), id, req.ReviewerID, req.Notes, req.ExpiresAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated request
	accessRequest, err := h.accessRequestService.GetAccessRequest(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toAccessRequestResponse(accessRequest))
}

func (h *AccessRequestHandler) Deny(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req denyAccessRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.accessRequestService.Deny(r.Context(), id, req.ReviewerID, req.Notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated request
	accessRequest, err := h.accessRequestService.GetAccessRequest(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toAccessRequestResponse(accessRequest))
}

// Helper function to convert domain model to response
func toAccessRequestResponse(ar *domain.AccessRequest) accessRequestResponse {
	return accessRequestResponse{
		ID:          ar.ID,
		UserID:      ar.UserID,
		ResourceID:  ar.ResourceID,
		Reason:      ar.Reason,
		Status:      ar.Status,
		ReviewerID:  ar.ReviewerID,
		ReviewNotes: ar.ReviewNotes,
		RequestedAt: ar.RequestedAt,
		ReviewedAt:  ar.ReviewedAt,
		ExpiresAt:   ar.ExpiresAt,
		CreatedAt:   ar.CreatedAt,
	}
}
