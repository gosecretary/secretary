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

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/register", h.Register).Methods("POST")
	r.HandleFunc("/api/login", h.Login).Methods("POST")
	r.HandleFunc("/api/users/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/api/users/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/api/users/{id}", h.Delete).Methods("DELETE")
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.userService.CreateUser(r.Context(), user); err != nil {
		utils.InternalError(w, "Failed to create user", err.Error())
		return
	}

	utils.SuccessResponse(w, "User created successfully", user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
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
		ExpiresAt: time.Now().Add(1 * time.Hour),
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

	utils.SuccessResponse(w, "Login successful", user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "User not found")
		return
	}

	utils.SuccessResponse(w, "User retrieved successfully", user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "User not found")
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	user.UpdatedAt = time.Now()

	if err := h.userService.Update(r.Context(), user); err != nil {
		utils.InternalError(w, "Failed to update user", err.Error())
		return
	}

	utils.SuccessResponse(w, "User updated successfully", user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.userService.Delete(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to delete user", err.Error())
		return
	}

	utils.SuccessResponse(w, "User deleted successfully", nil)
}
