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

// MockResourceRepository implements domain.ResourceRepository for testing
type MockResourceRepository struct {
	mock.Mock
}

func (m *MockResourceRepository) Create(resource *domain.Resource) error {
	args := m.Called(resource)
	return args.Error(0)
}

func (m *MockResourceRepository) FindByID(id string) (*domain.Resource, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Resource), args.Error(1)
}

func (m *MockResourceRepository) FindAll() ([]*domain.Resource, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Resource), args.Error(1)
}

func (m *MockResourceRepository) Update(resource *domain.Resource) error {
	args := m.Called(resource)
	return args.Error(0)
}

func (m *MockResourceRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestResourceService_CreateResource(t *testing.T) {
	tests := []struct {
		name      string
		resource  *domain.Resource
		mockSetup func(*MockResourceRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful resource creation",
			resource: &domain.Resource{
				Name:        "Test Database",
				Description: "A test database resource",
				Type:        "mysql",
			},
			mockSetup: func(m *MockResourceRepository) {
				m.On("Create", mock.AnythingOfType("*domain.Resource")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "empty name",
			resource: &domain.Resource{
				Name:        "",
				Description: "A test database resource",
				Type:        "mysql",
			},
			mockSetup: func(m *MockResourceRepository) {
				m.On("Create", mock.AnythingOfType("*domain.Resource")).Return(nil)
			},
			wantErr: false, // The service doesn't validate, it just calls the repo
		},
		{
			name: "repository error",
			resource: &domain.Resource{
				Name:        "Test Database",
				Description: "A test database resource",
				Type:        "mysql",
			},
			mockSetup: func(m *MockResourceRepository) {
				m.On("Create", mock.Anything).Return(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockResourceRepository)
			tt.mockSetup(mockRepo)

			service := NewResourceService(mockRepo)
			err := service.CreateResource(context.Background(), tt.resource)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				// Note: The service doesn't set ID or timestamps, the repository would do that
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestResourceService_ListResources(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockResourceRepository)
		wantResources []*domain.Resource
		wantErr       bool
	}{
		{
			name: "successful list resources",
			mockSetup: func(m *MockResourceRepository) {
				resources := []*domain.Resource{
					{
						ID:          "resource1",
						Name:        "Database 1",
						Description: "First database",
						Type:        "mysql",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          "resource2",
						Name:        "Database 2",
						Description: "Second database",
						Type:        "postgresql",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				m.On("FindAll").Return(resources, nil)
			},
			wantResources: []*domain.Resource{
				{
					ID:          "resource1",
					Name:        "Database 1",
					Description: "First database",
					Type:        "mysql",
				},
				{
					ID:          "resource2",
					Name:        "Database 2",
					Description: "Second database",
					Type:        "postgresql",
				},
			},
			wantErr: false,
		},
		{
			name: "empty list",
			mockSetup: func(m *MockResourceRepository) {
				m.On("FindAll").Return([]*domain.Resource{}, nil)
			},
			wantResources: []*domain.Resource{},
			wantErr:       false,
		},
		{
			name: "repository error",
			mockSetup: func(m *MockResourceRepository) {
				m.On("FindAll").Return(nil, errors.New("database error"))
			},
			wantResources: nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockResourceRepository)
			tt.mockSetup(mockRepo)

			service := NewResourceService(mockRepo)
			resources, err := service.ListResources(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resources)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.wantResources), len(resources))
				for i, expectedResource := range tt.wantResources {
					assert.Equal(t, expectedResource.ID, resources[i].ID)
					assert.Equal(t, expectedResource.Name, resources[i].Name)
					assert.Equal(t, expectedResource.Description, resources[i].Description)
					assert.Equal(t, expectedResource.Type, resources[i].Type)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestResourceService_GetResource(t *testing.T) {
	tests := []struct {
		name         string
		resourceID   string
		mockSetup    func(*MockResourceRepository)
		wantResource *domain.Resource
		wantErr      bool
	}{
		{
			name:       "successful get resource",
			resourceID: "resource123",
			mockSetup: func(m *MockResourceRepository) {
				resource := &domain.Resource{
					ID:          "resource123",
					Name:        "Test Database",
					Description: "A test database",
					Type:        "mysql",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("FindByID", "resource123").Return(resource, nil)
			},
			wantResource: &domain.Resource{
				ID:          "resource123",
				Name:        "Test Database",
				Description: "A test database",
				Type:        "mysql",
			},
			wantErr: false,
		},
		{
			name:       "resource not found",
			resourceID: "nonexistent",
			mockSetup: func(m *MockResourceRepository) {
				m.On("FindByID", "nonexistent").Return(nil, errors.New("resource not found"))
			},
			wantResource: nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockResourceRepository)
			tt.mockSetup(mockRepo)

			service := NewResourceService(mockRepo)
			resource, err := service.GetResource(context.Background(), tt.resourceID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resource)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resource)
				assert.Equal(t, tt.wantResource.ID, resource.ID)
				assert.Equal(t, tt.wantResource.Name, resource.Name)
				assert.Equal(t, tt.wantResource.Description, resource.Description)
				assert.Equal(t, tt.wantResource.Type, resource.Type)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestResourceService_UpdateResource(t *testing.T) {
	tests := []struct {
		name      string
		resource  *domain.Resource
		mockSetup func(*MockResourceRepository)
		wantErr   bool
	}{
		{
			name: "successful update",
			resource: &domain.Resource{
				ID:          "resource123",
				Name:        "Updated Database",
				Description: "Updated description",
				Type:        "postgresql",
			},
			mockSetup: func(m *MockResourceRepository) {
				m.On("Update", mock.AnythingOfType("*domain.Resource")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			resource: &domain.Resource{
				ID:          "resource123",
				Name:        "Updated Database",
				Description: "Updated description",
			},
			mockSetup: func(m *MockResourceRepository) {
				m.On("Update", mock.Anything).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockResourceRepository)
			tt.mockSetup(mockRepo)

			service := NewResourceService(mockRepo)
			err := service.UpdateResource(context.Background(), tt.resource)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Note: The service doesn't set timestamps, the repository would do that
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestResourceService_DeleteResource(t *testing.T) {
	tests := []struct {
		name       string
		resourceID string
		mockSetup  func(*MockResourceRepository)
		wantErr    bool
	}{
		{
			name:       "successful delete",
			resourceID: "resource123",
			mockSetup: func(m *MockResourceRepository) {
				m.On("Delete", "resource123").Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "repository error",
			resourceID: "resource123",
			mockSetup: func(m *MockResourceRepository) {
				m.On("Delete", "resource123").Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockResourceRepository)
			tt.mockSetup(mockRepo)

			service := NewResourceService(mockRepo)
			err := service.DeleteResource(context.Background(), tt.resourceID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
