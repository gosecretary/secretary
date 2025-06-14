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
				"message": "Session retrieved successfully",
				"data": map[string]interface{}{
					"id":        "test-session-id",
					"username":  "testuser",
					"status":    "active",
					"expiresAt": mock.Anything,
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
				"error": "Session not found",
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
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Session expired",
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

			// Call handler
			handler.GetByID(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Compare response with expected body
			for k, v := range tt.expectedBody {
				if k == "data" && v.(map[string]interface{})["expiresAt"] == mock.Anything {
					// Skip comparing expiresAt as it's time-dependent
					continue
				}
				assert.Equal(t, v, response[k])
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
				"message": "Session deleted successfully",
			},
		},
		{
			name:      "session not found",
			sessionID: "non-existent-session",
			mockSetup: func(m *MockSessionService) {
				m.On("Terminate", mock.Anything, "non-existent-session").Return(errors.New("session not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Session not found",
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

			// Call handler
			handler.Delete(w, req)

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
