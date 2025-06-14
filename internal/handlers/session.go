package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"secretary/alpha/internal/domain"
	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

type SessionHandler struct {
	sessionService domain.SessionService
}

func NewSessionHandler(sessionService domain.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

func (h *SessionHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/sessions", h.Create).Methods("POST")
	r.HandleFunc("/api/sessions", h.List).Methods("GET")
	r.HandleFunc("/api/sessions/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/api/sessions/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/api/sessions/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/api/sessions/{id}/terminate", h.Terminate).Methods("POST")
	r.HandleFunc("/api/sessions/user/{user_id}", h.GetByUserID).Methods("GET")
	r.HandleFunc("/api/sessions/resource/{resource_id}", h.GetByResourceID).Methods("GET")
	r.HandleFunc("/api/sessions/active", h.GetActive).Methods("GET")
}

type createSessionRequest struct {
	ResourceID     string `json:"resource_id"`
	ClientIP       string `json:"client_ip"`
	ClientMetadata string `json:"client_metadata,omitempty"`
}

func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createSessionRequest
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

	newSession := &domain.Session{
		UserID:         session.UserID,
		ResourceID:     req.ResourceID,
		StartTime:      time.Now(),
		Status:         "active",
		ClientIP:       req.ClientIP,
		ClientMetadata: req.ClientMetadata,
		ExpiresAt:      time.Now().Add(8 * time.Hour), // Default session duration
	}

	err := h.sessionService.Create(r.Context(), newSession)
	if err != nil {
		utils.InternalError(w, "Failed to create session", err.Error())
		return
	}

	utils.SuccessResponse(w, "Session created successfully", newSession)
}

func (h *SessionHandler) List(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.sessionService.List(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list sessions", err.Error())
		return
	}

	utils.SuccessResponse(w, "Sessions retrieved successfully", sessions)
}

func (h *SessionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	session, err := h.sessionService.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Session not found")
		return
	}

	utils.SuccessResponse(w, "Session retrieved successfully", session)
}

func (h *SessionHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var session domain.Session
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	session.ID = id
	if err := h.sessionService.Update(r.Context(), &session); err != nil {
		utils.InternalError(w, "Failed to update session", err.Error())
		return
	}

	utils.SuccessResponse(w, "Session updated successfully", session)
}

func (h *SessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.sessionService.Terminate(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to delete session", err.Error())
		return
	}

	utils.SuccessResponse(w, "Session deleted successfully", nil)
}

func (h *SessionHandler) Terminate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.sessionService.Terminate(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to terminate session", err.Error())
		return
	}

	utils.SuccessResponse(w, "Session terminated successfully", nil)
}

func (h *SessionHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	sessions, err := h.sessionService.GetByUserID(r.Context(), userID)
	if err != nil {
		utils.InternalError(w, "Failed to get user sessions", err.Error())
		return
	}

	utils.SuccessResponse(w, "User sessions retrieved successfully", sessions)
}

func (h *SessionHandler) GetByResourceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	sessions, err := h.sessionService.GetByResourceID(r.Context(), resourceID)
	if err != nil {
		utils.InternalError(w, "Failed to get resource sessions", err.Error())
		return
	}

	utils.SuccessResponse(w, "Resource sessions retrieved successfully", sessions)
}

func (h *SessionHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.sessionService.GetActive(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to get active sessions", err.Error())
		return
	}

	utils.SuccessResponse(w, "Active sessions retrieved successfully", sessions)
}
