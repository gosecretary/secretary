package domain

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ResourceRepository interface {
	Create(resource *Resource) error
	FindByID(id uuid.UUID) (*Resource, error)
	FindAll() ([]*Resource, error)
	Update(resource *Resource) error
	Delete(id uuid.UUID) error
}

type ResourceService interface {
	Create(name, description string) (*Resource, error)
	GetByID(id uuid.UUID) (*Resource, error)
	GetAll() ([]*Resource, error)
	Update(id uuid.UUID, name, description string) (*Resource, error)
	Delete(id uuid.UUID) error
}

// NewResource creates a new resource
func NewResource(name, description string) *Resource {
	return &Resource{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
} 