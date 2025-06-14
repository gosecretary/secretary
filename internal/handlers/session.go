package handlers

import (
	"encoding/json"
	"net/http"

	"secretary/alpha/internal/domain"
	"secretary/alpha/internal/middleware"
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
	r.HandleFunc("/api/sessions", h.CreateSession).Methods("POST")
	r.HandleFunc("/api/sessions", h.ListSessions).Methods("GET")
	r.HandleFunc("/api/sessions/{id}", h.GetSession).Methods("GET")
	r.HandleFunc("/api/sessions/{id}", h.DeleteSession).Methods("DELETE")
	r.HandleFunc("/api/sessions/{id}/extend", h.ExtendSession).Methods("POST")
}

type createSessionRequest struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	ExpiresAt string `json:"expires_at"`
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	session := &domain.Session{
		UserID:    req.UserID,
		Username:  req.Username,
		ExpiresAt: utils.ParseTime(req.ExpiresAt),
	}

	if err := h.sessionService.CreateSession(r.Context(), session); err != nil {
		utils.InternalError(w, "Failed to create session", err.Error())
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
	})

	utils.SuccessResponse(w, "Session created successfully", session)
}

func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.sessionService.ListSessions(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to list sessions", err.Error())
		return
	}

	utils.SuccessResponse(w, "Sessions retrieved successfully", sessions)
}

func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	session, err := h.sessionService.GetSession(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Session not found")
		return
	}

	utils.SuccessResponse(w, "Session retrieved successfully", session)
}

func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.sessionService.DeleteSession(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to delete session", err.Error())
		return
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	utils.SuccessResponse(w, "Session deleted successfully", nil)
}

type extendSessionRequest struct {
	ExpiresAt string `json:"expires_at"`
}

func (h *SessionHandler) ExtendSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req extendSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	session := &domain.Session{
		ID:        id,
		ExpiresAt: utils.ParseTime(req.ExpiresAt),
	}

	if err := h.sessionService.UpdateSession(r.Context(), session); err != nil {
		utils.InternalError(w, "Failed to extend session", err.Error())
		return
	}

	// Update session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
	})

	utils.SuccessResponse(w, "Session extended successfully", session)
}
