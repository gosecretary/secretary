package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"secretary/alpha/internal/domain"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSessionHandler_GetSession(t *testing.T) {
	tests := []struct {
		name           string
		sessionID      string
		mockSetup      func(*MockSessionService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "successful get session",
			sessionID: "test-session-id",
			mockSetup: func(m *MockSessionService) {
				session := &domain.Session{
					ID:        "test-session-id",
					UserID:    uuid.New().String(),
					Username:  "testuser",
					ExpiresAt: time.Now().Add(time.Hour),
					Status:    "active",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				m.On("GetByID", mock.Anything, "test-session-id").Return(session, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"code":    float64(200),
				"message": "Session retrieved successfully",
				"data": map[string]interface{}{
					"id":         "test-session-id",
					"username":   "testuser",
					"status":     "active",
					"expires_at": mock.Anything,
					"created_at": mock.Anything,
					"updated_at": mock.Anything,
				},
			},
		},
		{
			name:      "session not found",
			sessionID: "non-existent-session",
			mockSetup: func(m *MockSessionService) {
				m.On("GetByID", mock.Anything, "non-existent-session").Return(nil, errors.New("session not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(404),
				"message": "Session not found",
				"data": map[string]interface{}{
					"error": "Resource not found",
				},
			},
		},
		{
			name:      "expired session",
			sessionID: "expired-session",
			mockSetup: func(m *MockSessionService) {
				session := &domain.Session{
					ID:        "expired-session",
					UserID:    uuid.New().String(),
					Username:  "testuser",
					ExpiresAt: time.Now().Add(-time.Hour), // Expired
					Status:    "active",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				m.On("GetByID", mock.Anything, "expired-session").Return(session, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"code":    float64(200),
				"message": "Session retrieved successfully",
				"data": map[string]interface{}{
					"id":         "expired-session",
					"username":   "testuser",
					"status":     "active",
					"expires_at": mock.Anything,
					"created_at": mock.Anything,
					"updated_at": mock.Anything,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockSessionService)
			tt.mockSetup(mockService)

			// Create handler
			handler := NewSessionHandler(mockService)

			// Create request
			req := httptest.NewRequest("GET", "/api/sessions/"+tt.sessionID, nil)
			w := httptest.NewRecorder()

			// Set up mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.sessionID})

			// Call handler
			handler.GetByID(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Compare response with expected body
			if tt.name == "successful get session" || tt.name == "expired session" {
				assert.Equal(t, tt.expectedBody["success"], response["success"])
				assert.Equal(t, tt.expectedBody["code"], response["code"])
				assert.Equal(t, tt.expectedBody["message"], response["message"])

				data := response["data"].(map[string]interface{})
				expectedData := tt.expectedBody["data"].(map[string]interface{})

				assert.Equal(t, expectedData["id"], data["id"])
				assert.Equal(t, expectedData["username"], data["username"])
				assert.Equal(t, expectedData["status"], data["status"])
				assert.NotNil(t, data["expires_at"])
				assert.NotNil(t, data["created_at"])
				assert.NotNil(t, data["updated_at"])
			} else if tt.name == "session not found" {
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

func TestSessionHandler_DeleteSession(t *testing.T) {
	tests := []struct {
		name           string
		sessionID      string
		mockSetup      func(*MockSessionService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "successful delete session",
			sessionID: "test-session-id",
			mockSetup: func(m *MockSessionService) {
				m.On("Terminate", mock.Anything, "test-session-id").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"code":    float64(200),
				"message": "Session deleted successfully",
				"data":    nil,
			},
		},
		{
			name:      "session not found",
			sessionID: "non-existent-session",
			mockSetup: func(m *MockSessionService) {
				m.On("Terminate", mock.Anything, "non-existent-session").Return(errors.New("session not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"success": false,
				"code":    float64(500),
				"message": "Failed to delete session",
				"data": map[string]interface{}{
					"error": mock.Anything,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockSessionService)
			tt.mockSetup(mockService)

			// Create handler
			handler := NewSessionHandler(mockService)

			// Create request
			req := httptest.NewRequest("DELETE", "/api/sessions/"+tt.sessionID, nil)
			w := httptest.NewRecorder()

			// Set up mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.sessionID})

			// Call handler
			handler.Delete(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Compare response with expected body
			if tt.name == "session not found" {
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
