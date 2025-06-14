package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"secretary/alpha/internal/domain"
	"secretary/alpha/internal/middleware"
	"secretary/alpha/pkg/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	userService domain.UserService
}

func NewAuthHandler(userService domain.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/register", h.Register).Methods("POST")
	r.HandleFunc("/api/login", h.Login).Methods("POST")
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	user := &domain.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	if err := h.userService.CreateUser(r.Context(), user); err != nil {
		utils.InternalError(w, "Failed to create user", err.Error())
		return
	}

	utils.SuccessResponse(w, "User registered successfully", map[string]string{
		"username": user.Username,
		"email":    user.Email,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.Authenticate(r.Context(), req.Username, req.Password)
	if err != nil {
		utils.Unauthorized(w, "Invalid credentials")
		return
	}

	// Create session
	session := &domain.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Username:  user.Username,
		Status:    "active",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store session
	if err := domain.GetSessionStore().Set(session); err != nil {
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

	utils.SuccessResponse(w, "Login successful", map[string]interface{}{
		"user":    user,
		"session": session,
	})
}
