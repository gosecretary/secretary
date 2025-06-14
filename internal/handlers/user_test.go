package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"secretary/alpha/internal/domain"

	"github.com/gorilla/mux"
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
				"success": true,
				"code":    float64(200),
				"message": "User created successfully",
				"data": map[string]interface{}{
					"id":         "",
					"username":   "testuser",
					"email":      "test@example.com",
					"name":       "Test User",
					"role":       "",
					"created_at": mock.Anything,
					"updated_at": mock.Anything,
				},
			},
		},
		{
			name: "invalid request body",
			request: domain.User{
				Username: "testuser",
				// Missing required fields
			},
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
					return user.Username == "testuser" && user.Password == ""
				})).Return(fmt.Errorf("password is required"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(500),
				"message": "Failed to create user",
				"data": map[string]interface{}{
					"error": mock.Anything,
				},
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
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(500),
				"message": "Failed to create user",
				"data": map[string]interface{}{
					"error": mock.Anything,
				},
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
			if tt.name == "successful user creation" {
				assert.Equal(t, tt.expectedBody["success"], response["success"])
				assert.Equal(t, tt.expectedBody["code"], response["code"])
				assert.Equal(t, tt.expectedBody["message"], response["message"])

				data := response["data"].(map[string]interface{})
				expectedData := tt.expectedBody["data"].(map[string]interface{})

				assert.Equal(t, expectedData["id"], data["id"])
				assert.Equal(t, expectedData["username"], data["username"])
				assert.Equal(t, expectedData["email"], data["email"])
				assert.Equal(t, expectedData["name"], data["name"])
				assert.Equal(t, expectedData["role"], data["role"])
				assert.NotNil(t, data["created_at"])
				assert.NotNil(t, data["updated_at"])
			} else if tt.name == "invalid request body" || tt.name == "username already exists" {
				assert.Equal(t, tt.expectedBody["success"], response["success"])
				assert.Equal(t, tt.expectedBody["code"], response["code"])
				assert.Equal(t, tt.expectedBody["message"], response["message"])

				data := response["data"].(map[string]interface{})
				assert.NotNil(t, data["error"])
			} else {
				for k, v := range tt.expectedBody {
					assert.Equal(t, v, response[k])
				}
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
				"success": true,
				"code":    float64(200),
				"message": "User retrieved successfully",
				"data": map[string]interface{}{
					"id":         "test-user-id",
					"username":   "testuser",
					"email":      "test@example.com",
					"name":       "Test User",
					"role":       "user",
					"created_at": mock.Anything,
					"updated_at": mock.Anything,
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
				"success": false,
				"code":    float64(404),
				"message": "User not found",
				"data": map[string]interface{}{
					"error": "Resource not found",
				},
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

			// Set up mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.userID})

			// Call handler
			handler.GetByID(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Compare response with expected body
			if tt.name == "successful get user" {
				assert.Equal(t, tt.expectedBody["success"], response["success"])
				assert.Equal(t, tt.expectedBody["code"], response["code"])
				assert.Equal(t, tt.expectedBody["message"], response["message"])

				data := response["data"].(map[string]interface{})
				expectedData := tt.expectedBody["data"].(map[string]interface{})

				assert.Equal(t, expectedData["id"], data["id"])
				assert.Equal(t, expectedData["username"], data["username"])
				assert.Equal(t, expectedData["email"], data["email"])
				assert.Equal(t, expectedData["name"], data["name"])
				assert.Equal(t, expectedData["role"], data["role"])
				assert.NotNil(t, data["created_at"])
				assert.NotNil(t, data["updated_at"])
			} else if tt.name == "user not found" {
				assert.Equal(t, tt.expectedBody["success"], response["success"])
				assert.Equal(t, tt.expectedBody["code"], response["code"])
				assert.Equal(t, tt.expectedBody["message"], response["message"])

				data := response["data"].(map[string]interface{})
				assert.NotNil(t, data["error"])
			} else {
				for k, v := range tt.expectedBody {
					assert.Equal(t, v, response[k])
				}
			}

			// Verify all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
