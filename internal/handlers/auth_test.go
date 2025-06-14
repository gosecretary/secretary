package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"secretary/alpha/internal/domain"
	"secretary/alpha/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful registration",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "testpass",
			},
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
					return user.Username == "testuser"
				})).Return(nil)
			},
			expectedStatus: 200,
			expectedBody: map[string]interface{}{
				"success": true,
				"code":    float64(200),
				"message": "User registered successfully",
				"data": map[string]interface{}{
					"username": "testuser",
					"email":    "",
				},
			},
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: 400,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(400),
				"message": "Invalid request body",
				"error":   "invalid request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			// Create handler with mock service
			handler := NewAuthHandler(mockService)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			// Call handler
			handler.Register(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockUserService, *MockSessionService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful login",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "testpass",
			},
			mockSetup: func(us *MockUserService, ss *MockSessionService) {
				user := &domain.User{
					ID:       "user123",
					Username: "testuser",
				}
				us.On("Authenticate", mock.Anything, "testuser", "testpass").Return(user, nil)
				ss.On("Create", mock.Anything, mock.MatchedBy(func(s *domain.Session) bool {
					return s.UserID == "user123"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"user": map[string]interface{}{
					"id":       "user123",
					"username": "testuser",
				},
			},
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"username": "testuser",
			},
			mockSetup:      func(us *MockUserService, ss *MockSessionService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
		{
			name: "authentication failed",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "wrongpass",
			},
			mockSetup: func(us *MockUserService, ss *MockSessionService) {
				us.On("Authenticate", mock.Anything, "testuser", "wrongpass").Return(nil, utils.NewError("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "invalid credentials",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(MockUserService)
			mockSessionService := new(MockSessionService)
			tt.mockSetup(mockUserService, mockSessionService)

			handler := NewAuthHandler(mockUserService)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedBody, response)
			mockUserService.AssertExpectations(t)
			mockSessionService.AssertExpectations(t)
		})
	}
}
