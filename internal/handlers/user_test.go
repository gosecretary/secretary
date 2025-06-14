package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"secretary/alpha/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		request        domain.User
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful user creation",
			request: domain.User{
				Username: "testuser",
				Password: "testpass",
				Email:    "test@example.com",
				Name:     "Test User",
			},
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
					return user.Username == "testuser" && user.Email == "test@example.com" && user.Name == "Test User"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "User created successfully",
				"data": map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
					"name":     "Test User",
				},
			},
		},
		{
			name: "invalid request body",
			request: domain.User{
				Username: "testuser",
				// Missing required fields
			},
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
		{
			name: "username already exists",
			request: domain.User{
				Username: "existinguser",
				Password: "testpass",
				Email:    "test@example.com",
				Name:     "Test User",
			},
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("username already exists"))
			},
			expectedStatus: http.StatusConflict,
			expectedBody: map[string]interface{}{
				"error": "Username already exists",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			// Create handler
			handler := NewUserHandler(mockService)

			// Create request
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			// Call handler
			handler.Register(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Compare response with expected body
			for k, v := range tt.expectedBody {
				assert.Equal(t, v, response[k])
			}

			// Verify all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful get user",
			userID: "test-user-id",
			mockSetup: func(m *MockUserService) {
				user := &domain.User{
					ID:        "test-user-id",
					Username:  "testuser",
					Email:     "test@example.com",
					Name:      "Test User",
					Role:      "user",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				m.On("GetByID", mock.Anything, "test-user-id").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "User retrieved successfully",
				"data": map[string]interface{}{
					"id":        "test-user-id",
					"username":  "testuser",
					"email":     "test@example.com",
					"name":      "Test User",
					"role":      "user",
					"createdAt": mock.Anything,
					"updatedAt": mock.Anything,
				},
			},
		},
		{
			name:   "user not found",
			userID: "non-existent-user",
			mockSetup: func(m *MockUserService) {
				m.On("GetByID", mock.Anything, "non-existent-user").Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "User not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			// Create handler
			handler := NewUserHandler(mockService)

			// Create request
			req := httptest.NewRequest("GET", "/api/users/"+tt.userID, nil)
			w := httptest.NewRecorder()

			// Call handler
			handler.GetByID(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Compare response with expected body
			for k, v := range tt.expectedBody {
				if k == "data" {
					data := v.(map[string]interface{})
					responseData := response[k].(map[string]interface{})
					for dk, dv := range data {
						if dk == "createdAt" || dk == "updatedAt" {
							// Skip comparing timestamps
							continue
						}
						assert.Equal(t, dv, responseData[dk])
					}
				} else {
					assert.Equal(t, v, response[k])
				}
			}

			// Verify all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
