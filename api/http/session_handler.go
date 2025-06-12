package http

import (
	"encoding/json"
	"net/http"
	"time"

	"secretary/alpha/internal/domain"

	"github.com/gorilla/mux"
)

type SessionHandler struct {
	sessionService domain.SessionService
}

func NewSessionHandler(sessionService domain.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

type createSessionRequest struct {
	UserID         string `json:"user_id"`
	ResourceID     string `json:"resource_id"`
	ClientIP       string `json:"client_ip"`
	ClientMetadata string `json:"client_metadata"`
}

type sessionResponse struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	ResourceID     string    `json:"resource_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time,omitempty"`
	Status         string    `json:"status"`
	ClientIP       string    `json:"client_ip"`
	ClientMetadata string    `json:"client_metadata,omitempty"`
	AuditPath      string    `json:"audit_path,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session := &domain.Session{
		UserID:         req.UserID,
		ResourceID:     req.ResourceID,
		ClientIP:       req.ClientIP,
		ClientMetadata: req.ClientMetadata,
	}

	if err := h.sessionService.Create(r.Context(), session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toSessionResponse(session))
}

func (h *SessionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	session, err := h.sessionService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toSessionResponse(session))
}

func (h *SessionHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	sessions, err := h.sessionService.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]sessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = toSessionResponse(session)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *SessionHandler) GetByResourceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	sessions, err := h.sessionService.GetByResourceID(r.Context(), resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]sessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = toSessionResponse(session)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *SessionHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.sessionService.GetActive(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]sessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = toSessionResponse(session)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *SessionHandler) Terminate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.sessionService.Terminate(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function to convert domain model to response
func toSessionResponse(s *domain.Session) sessionResponse {
	return sessionResponse{
		ID:             s.ID,
		UserID:         s.UserID,
		ResourceID:     s.ResourceID,
		StartTime:      s.StartTime,
		EndTime:        s.EndTime,
		Status:         s.Status,
		ClientIP:       s.ClientIP,
		ClientMetadata: s.ClientMetadata,
		AuditPath:      s.AuditPath,
		CreatedAt:      s.CreatedAt,
	}
}
