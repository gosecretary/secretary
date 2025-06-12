package domain

import (
	"context"
	"time"
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

// ResourceRepository defines the interface for resource-related data operations
type ResourceRepository interface {
	Create(resource *Resource) error
	FindByID(id string) (*Resource, error)
	FindAll() ([]*Resource, error)
	Update(resource *Resource) error
	Delete(id string) error
}

// CredentialRepository defines the interface for credential-related data operations
type CredentialRepository interface {
	Create(credential *Credential) error
	FindByID(id string) (*Credential, error)
	FindByResourceID(resourceID string) ([]*Credential, error)
	Update(credential *Credential) error
	Delete(id string) error
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

// PermissionRepository defines the interface for permission-related data operations
type PermissionRepository interface {
	Create(permission *Permission) error
	FindByID(id string) (*Permission, error)
	FindByUserID(userID string) ([]*Permission, error)
	FindByResourceID(resourceID string) ([]*Permission, error)
	Update(permission *Permission) error
	Delete(id string) error
	DeleteByUserID(userID string) error
	DeleteByResourceID(resourceID string) error
}

// UserRepository defines the interface for user-related data operations
type UserRepository interface {
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
}

// SessionService defines the interface for session-related operations
type SessionService interface {
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id string) (*Session, error)
	GetByUserID(ctx context.Context, userID string) ([]*Session, error)
	GetByResourceID(ctx context.Context, resourceID string) ([]*Session, error)
	GetActive(ctx context.Context) ([]*Session, error)
	Update(ctx context.Context, session *Session) error
	Terminate(ctx context.Context, id string) error
}

// SessionRepository defines the interface for session-related data operations
type SessionRepository interface {
	Create(session *Session) error
	FindByID(id string) (*Session, error)
	FindByUserID(userID string) ([]*Session, error)
	FindByResourceID(resourceID string) ([]*Session, error)
	FindActive() ([]*Session, error)
	Update(session *Session) error
	Delete(id string) error
}

// AccessRequestService defines the interface for access request operations
type AccessRequestService interface {
	Create(ctx context.Context, request *AccessRequest) error
	GetByID(ctx context.Context, id string) (*AccessRequest, error)
	GetByUserID(ctx context.Context, userID string) ([]*AccessRequest, error)
	GetByResourceID(ctx context.Context, resourceID string) ([]*AccessRequest, error)
	GetPending(ctx context.Context) ([]*AccessRequest, error)
	Approve(ctx context.Context, id string, reviewerID string, notes string, expiresAt time.Time) error
	Deny(ctx context.Context, id string, reviewerID string, notes string) error
}

// AccessRequestRepository defines the interface for access request data operations
type AccessRequestRepository interface {
	Create(request *AccessRequest) error
	FindByID(id string) (*AccessRequest, error)
	FindByUserID(userID string) ([]*AccessRequest, error)
	FindByResourceID(resourceID string) ([]*AccessRequest, error)
	FindByStatus(status string) ([]*AccessRequest, error)
	Update(request *AccessRequest) error
}

// EphemeralCredentialService defines the interface for ephemeral credential operations
type EphemeralCredentialService interface {
	Generate(ctx context.Context, userID string, resourceID string, duration time.Duration) (*EphemeralCredential, error)
	GetByID(ctx context.Context, id string) (*EphemeralCredential, error)
	GetByToken(ctx context.Context, token string) (*EphemeralCredential, error)
	MarkAsUsed(ctx context.Context, id string) error
	RevokeByUserID(ctx context.Context, userID string) error
	RevokeByResourceID(ctx context.Context, resourceID string) error
}

// EphemeralCredentialRepository defines the interface for ephemeral credential data operations
type EphemeralCredentialRepository interface {
	Create(credential *EphemeralCredential) error
	FindByID(id string) (*EphemeralCredential, error)
	FindByToken(token string) (*EphemeralCredential, error)
	FindByUserID(userID string) ([]*EphemeralCredential, error)
	Update(credential *EphemeralCredential) error
	DeleteExpired() error
	DeleteByUserID(userID string) error
	DeleteByResourceID(resourceID string) error
}
