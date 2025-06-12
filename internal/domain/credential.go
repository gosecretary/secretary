package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Credential struct {
	ID         uuid.UUID `json:"id"`
	ResourceID uuid.UUID `json:"resource_id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CredentialRepository interface {
	Create(credential *Credential) error
	FindByID(id uuid.UUID) (*Credential, error)
	FindByResourceID(resourceID uuid.UUID) ([]*Credential, error)
	Update(credential *Credential) error
	Delete(id uuid.UUID) error
}

type CredentialService interface {
	Create(resourceID uuid.UUID, username, password string) (*Credential, error)
	GetByID(id uuid.UUID) (*Credential, error)
	GetByResourceID(resourceID uuid.UUID) ([]*Credential, error)
	Update(id uuid.UUID, username, password string) (*Credential, error)
	Delete(id uuid.UUID) error
}

// NewCredential creates a new credential with encrypted password
func NewCredential(resourceID uuid.UUID, username, password string) (*Credential, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Credential{
		ID:         uuid.New(),
		ResourceID: resourceID,
		Username:   username,
		Password:   string(hashedPassword),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

// ValidatePassword checks if the provided password matches the credential's password
func (c *Credential) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
	return err == nil
}

// UpdatePassword updates the credential's password with a new encrypted password
func (c *Credential) UpdatePassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(hashedPassword)
	c.UpdatedAt = time.Now()
	return nil
} 