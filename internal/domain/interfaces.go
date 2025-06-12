package domain

import (
	"context"
)

// UserService defines the interface for user-related operations
type UserService interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, email, password string) (string, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// ResourceService defines the interface for resource-related operations
type ResourceService interface {
	Create(ctx context.Context, resource *Resource) error
	GetByID(ctx context.Context, id string) (*Resource, error)
	GetAll(ctx context.Context) ([]*Resource, error)
	Update(ctx context.Context, resource *Resource) error
	Delete(ctx context.Context, id string) error
}

// CredentialService defines the interface for credential-related operations
type CredentialService interface {
	Create(ctx context.Context, credential *Credential) error
	GetByID(ctx context.Context, id string) (*Credential, error)
	GetByResourceID(ctx context.Context, resourceID string) ([]*Credential, error)
	Update(ctx context.Context, credential *Credential) error
	Delete(ctx context.Context, id string) error
}

// PermissionService defines the interface for permission-related operations
type PermissionService interface {
	Create(ctx context.Context, permission *Permission) error
	GetByID(ctx context.Context, id string) (*Permission, error)
	GetByUserID(ctx context.Context, userID string) ([]*Permission, error)
	GetByResourceID(ctx context.Context, resourceID string) ([]*Permission, error)
	Update(ctx context.Context, permission *Permission) error
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteByResourceID(ctx context.Context, resourceID string) error
} 