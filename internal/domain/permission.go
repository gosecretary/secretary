package domain

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ResourceID uuid.UUID `json:"resource_id"`
	Action     string    `json:"action"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PermissionRepository interface {
	Create(permission *Permission) error
	FindByID(id uuid.UUID) (*Permission, error)
	FindByUserID(userID uuid.UUID) ([]*Permission, error)
	FindByResourceID(resourceID uuid.UUID) ([]*Permission, error)
	Update(permission *Permission) error
	Delete(id uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error
	DeleteByResourceID(resourceID uuid.UUID) error
}

type PermissionService interface {
	Create(userID, resourceID uuid.UUID, action string) (*Permission, error)
	GetByID(id uuid.UUID) (*Permission, error)
	GetByUserID(userID uuid.UUID) ([]*Permission, error)
	GetByResourceID(resourceID uuid.UUID) ([]*Permission, error)
	Update(id uuid.UUID, action string) (*Permission, error)
	Delete(id uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error
	DeleteByResourceID(resourceID uuid.UUID) error
}

// NewPermission creates a new permission
func NewPermission(userID, resourceID uuid.UUID, action string) *Permission {
	return &Permission{
		ID:         uuid.New(),
		UserID:     userID,
		ResourceID: resourceID,
		Action:     action,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
} 