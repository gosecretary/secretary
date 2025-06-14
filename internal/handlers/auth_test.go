package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"secretary/alpha/internal/domain"

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
				"message": "User created successfully",
				"data": map[string]interface{}{
					"id":         "",
					"username":   "testuser",
					"email":      "",
					"name":       "",
					"role":       "",
					"created_at": mock.Anything,
					"updated_at": mock.Anything,
				},
			},
		},
		{
			name: "missing required fields",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			mockSetup: func(m *MockUserService) {
				// Handler will try to create user with empty values
				m.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
					return user.Username == "" && user.Password == ""
				})).Return(fmt.Errorf("username and password are required"))
			},
			expectedStatus: 500,
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

			// Create handler with mock service
			handler := NewUserHandler(mockService)

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

			// For successful registration, check fields individually due to timestamps
			if tt.name == "successful registration" {
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
			} else if tt.name == "missing required fields" {
				// Check individual fields for missing required fields test
				assert.Equal(t, tt.expectedBody["success"], response["success"])
				assert.Equal(t, tt.expectedBody["code"], response["code"])
				assert.Equal(t, tt.expectedBody["message"], response["message"])

				// Error is nested in data field
				data := response["data"].(map[string]interface{})
				assert.NotNil(t, data["error"]) // Error message can vary
			} else {
				assert.Equal(t, tt.expectedBody, response)
			}

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
				// Session service is not used in the handler anymore
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"code":    float64(200),
				"message": "Login successful",
				"data": map[string]interface{}{
					"id":         "user123",
					"username":   "testuser",
					"email":      "",
					"name":       "",
					"role":       "",
					"created_at": mock.Anything,
					"updated_at": mock.Anything,
				},
			},
		},
		{
			name: "missing password field",
			requestBody: map[string]interface{}{
				"username": "testuser",
			},
			mockSetup: func(us *MockUserService, ss *MockSessionService) {
				// Handler will try to authenticate with empty password
				us.On("Authenticate", mock.Anything, "testuser", "").Return(nil, fmt.Errorf("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(401),
				"message": "Invalid credentials",
				"data": map[string]interface{}{
					"error": "Authentication required",
				},
			},
		},
		{
			name: "authentication failed",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "wrongpass",
			},
			mockSetup: func(us *MockUserService, ss *MockSessionService) {
				us.On("Authenticate", mock.Anything, "testuser", "wrongpass").Return(nil, fmt.Errorf("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(401),
				"message": "Invalid credentials",
				"data": map[string]interface{}{
					"error": "Authentication required",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(MockUserService)
			mockSessionService := new(MockSessionService)
			tt.mockSetup(mockUserService, mockSessionService)

			handler := NewUserHandler(mockUserService)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Handle different test cases
			if tt.name == "successful login" {
				expectedBody := tt.expectedBody.(map[string]interface{})
				assert.Equal(t, expectedBody["success"], response["success"])
				assert.Equal(t, expectedBody["code"], response["code"])
				assert.Equal(t, expectedBody["message"], response["message"])

				data := response["data"].(map[string]interface{})
				expectedData := expectedBody["data"].(map[string]interface{})

				assert.Equal(t, expectedData["id"], data["id"])
				assert.Equal(t, expectedData["username"], data["username"])
				assert.Equal(t, expectedData["email"], data["email"])
				assert.Equal(t, expectedData["name"], data["name"])
				assert.Equal(t, expectedData["role"], data["role"])
				assert.NotNil(t, data["created_at"])
				assert.NotNil(t, data["updated_at"])
			} else if tt.name == "missing password field" || tt.name == "authentication failed" {
				expectedBody := tt.expectedBody.(map[string]interface{})
				assert.Equal(t, expectedBody["success"], response["success"])
				assert.Equal(t, expectedBody["code"], response["code"])
				assert.Equal(t, expectedBody["message"], response["message"])

				data := response["data"].(map[string]interface{})
				assert.NotNil(t, data["error"])
			} else {
				assert.Equal(t, tt.expectedBody, response)
			}
			mockUserService.AssertExpectations(t)
			mockSessionService.AssertExpectations(t)
		})
	}
}
