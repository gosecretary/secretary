package domain

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Password is never exposed in JSON
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Resource represents a resource that can be accessed
type Resource struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Credential represents credentials for accessing a resource
type Credential struct {
	ID         string    `json:"id"`
	ResourceID string    `json:"resource_id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"` // Password is never exposed in JSON
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Permission represents a user's permission to access a resource
type Permission struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ResourceID string    `json:"resource_id"`
	Role       string    `json:"role"`
	Action     string    `json:"action"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Session represents an active or completed connection to a resource
type Session struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	ResourceID     string    `json:"resource_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time,omitempty"`
	Status         string    `json:"status"` // "active", "completed", "terminated"
	ClientIP       string    `json:"client_ip"`
	ClientMetadata string    `json:"client_metadata,omitempty"`
	AuditPath      string    `json:"audit_path,omitempty"` // Path to session recording
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AccessRequest represents a request for access to a resource
type AccessRequest struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ResourceID  string    `json:"resource_id"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"` // "pending", "approved", "denied"
	ReviewerID  string    `json:"reviewer_id,omitempty"`
	ReviewNotes string    `json:"review_notes,omitempty"`
	RequestedAt time.Time `json:"requested_at"`
	ReviewedAt  time.Time `json:"reviewed_at,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"` // When approved access expires
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EphemeralCredential represents a temporary credential for accessing a resource
type EphemeralCredential struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ResourceID string    `json:"resource_id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"` // Password is never exposed in JSON
	Token      string    `json:"token,omitempty"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	UsedAt     time.Time `json:"used_at,omitempty"`
}

// ValidatePassword checks if the provided password matches the user's password
func (u *User) ValidatePassword(password string) bool {
	// Implement password validation logic here
	return u.Password == password
}
