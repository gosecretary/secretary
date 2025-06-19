package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"secretary/alpha/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAccessRequestRepository implements domain.AccessRequestRepository for testing
type MockAccessRequestRepository struct {
	mock.Mock
}

func (m *MockAccessRequestRepository) Create(request *domain.AccessRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockAccessRequestRepository) FindByID(id string) (*domain.AccessRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AccessRequest), args.Error(1)
}

func (m *MockAccessRequestRepository) FindByUserID(userID string) ([]*domain.AccessRequest, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.AccessRequest), args.Error(1)
}

func (m *MockAccessRequestRepository) FindByResourceID(resourceID string) ([]*domain.AccessRequest, error) {
	args := m.Called(resourceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.AccessRequest), args.Error(1)
}

func (m *MockAccessRequestRepository) FindByStatus(status string) ([]*domain.AccessRequest, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.AccessRequest), args.Error(1)
}

func (m *MockAccessRequestRepository) Update(request *domain.AccessRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func TestAccessRequestService_CreateAccessRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *domain.AccessRequest
		mockSetup func(*MockAccessRequestRepository)
		wantErr   bool
	}{
		{
			name: "successful creation",
			request: &domain.AccessRequest{
				UserID:     "user123",
				ResourceID: "resource123",
				Reason:     "Need access for debugging",
			},
			mockSetup: func(m *MockAccessRequestRepository) {
				m.On("Create", mock.AnythingOfType("*domain.AccessRequest")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			request: &domain.AccessRequest{
				UserID:     "user123",
				ResourceID: "resource123",
				Reason:     "Need access",
			},
			mockSetup: func(m *MockAccessRequestRepository) {
				m.On("Create", mock.Anything).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAccessRequestRepository)
			tt.mockSetup(mockRepo)

			service := NewAccessRequestService(mockRepo)
			err := service.CreateAccessRequest(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.request.ID)
				assert.Equal(t, "pending", tt.request.Status)
				assert.False(t, tt.request.RequestedAt.IsZero())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAccessRequestService_GetAccessRequest(t *testing.T) {
	tests := []struct {
		name        string
		requestID   string
		mockSetup   func(*MockAccessRequestRepository)
		wantRequest *domain.AccessRequest
		wantErr     bool
	}{
		{
			name:      "successful get",
			requestID: "request123",
			mockSetup: func(m *MockAccessRequestRepository) {
				request := &domain.AccessRequest{
					ID:         "request123",
					UserID:     "user123",
					ResourceID: "resource123",
					Reason:     "Test reason",
					Status:     "pending",
				}
				m.On("FindByID", "request123").Return(request, nil)
			},
			wantRequest: &domain.AccessRequest{
				ID:         "request123",
				UserID:     "user123",
				ResourceID: "resource123",
				Reason:     "Test reason",
				Status:     "pending",
			},
			wantErr: false,
		},
		{
			name:      "not found",
			requestID: "nonexistent",
			mockSetup: func(m *MockAccessRequestRepository) {
				m.On("FindByID", "nonexistent").Return(nil, errors.New("not found"))
			},
			wantRequest: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAccessRequestRepository)
			tt.mockSetup(mockRepo)

			service := NewAccessRequestService(mockRepo)
			request, err := service.GetAccessRequest(context.Background(), tt.requestID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, request)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, request)
				assert.Equal(t, tt.wantRequest.ID, request.ID)
				assert.Equal(t, tt.wantRequest.UserID, request.UserID)
				assert.Equal(t, tt.wantRequest.ResourceID, request.ResourceID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAccessRequestService_Approve(t *testing.T) {
	tests := []struct {
		name       string
		requestID  string
		reviewerID string
		notes      string
		expiresAt  time.Time
		mockSetup  func(*MockAccessRequestRepository)
		wantErr    bool
	}{
		{
			name:       "successful approval",
			requestID:  "request123",
			reviewerID: "reviewer123",
			notes:      "Approved for maintenance",
			expiresAt:  time.Now().Add(time.Hour),
			mockSetup: func(m *MockAccessRequestRepository) {
				request := &domain.AccessRequest{
					ID:     "request123",
					Status: "pending",
				}
				m.On("FindByID", "request123").Return(request, nil)
				m.On("Update", mock.AnythingOfType("*domain.AccessRequest")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "request not found",
			requestID:  "nonexistent",
			reviewerID: "reviewer123",
			notes:      "Notes",
			expiresAt:  time.Now().Add(time.Hour),
			mockSetup: func(m *MockAccessRequestRepository) {
				m.On("FindByID", "nonexistent").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAccessRequestRepository)
			tt.mockSetup(mockRepo)

			service := NewAccessRequestService(mockRepo)
			err := service.Approve(context.Background(), tt.requestID, tt.reviewerID, tt.notes, tt.expiresAt)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAccessRequestService_Deny(t *testing.T) {
	tests := []struct {
		name       string
		requestID  string
		reviewerID string
		notes      string
		mockSetup  func(*MockAccessRequestRepository)
		wantErr    bool
	}{
		{
			name:       "successful denial",
			requestID:  "request123",
			reviewerID: "reviewer123",
			notes:      "Insufficient justification",
			mockSetup: func(m *MockAccessRequestRepository) {
				request := &domain.AccessRequest{
					ID:     "request123",
					Status: "pending",
				}
				m.On("FindByID", "request123").Return(request, nil)
				m.On("Update", mock.AnythingOfType("*domain.AccessRequest")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "request not found",
			requestID:  "nonexistent",
			reviewerID: "reviewer123",
			notes:      "Notes",
			mockSetup: func(m *MockAccessRequestRepository) {
				m.On("FindByID", "nonexistent").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAccessRequestRepository)
			tt.mockSetup(mockRepo)

			service := NewAccessRequestService(mockRepo)
			err := service.Deny(context.Background(), tt.requestID, tt.reviewerID, tt.notes)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAccessRequestService_GetPendingAccessRequests(t *testing.T) {
	tests := []struct {
		name         string
		mockSetup    func(*MockAccessRequestRepository)
		wantRequests []*domain.AccessRequest
		wantErr      bool
	}{
		{
			name: "successful get pending",
			mockSetup: func(m *MockAccessRequestRepository) {
				requests := []*domain.AccessRequest{
					{ID: "req1", Status: "pending", UserID: "user1"},
					{ID: "req2", Status: "pending", UserID: "user2"},
				}
				m.On("FindByStatus", "pending").Return(requests, nil)
			},
			wantRequests: []*domain.AccessRequest{
				{ID: "req1", Status: "pending", UserID: "user1"},
				{ID: "req2", Status: "pending", UserID: "user2"},
			},
			wantErr: false,
		},
		{
			name: "repository error",
			mockSetup: func(m *MockAccessRequestRepository) {
				m.On("FindByStatus", "pending").Return(nil, errors.New("database error"))
			},
			wantRequests: nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAccessRequestRepository)
			tt.mockSetup(mockRepo)

			service := NewAccessRequestService(mockRepo)
			requests, err := service.GetPendingAccessRequests(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, requests)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.wantRequests), len(requests))
				for i, expectedReq := range tt.wantRequests {
					assert.Equal(t, expectedReq.ID, requests[i].ID)
					assert.Equal(t, expectedReq.Status, requests[i].Status)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
