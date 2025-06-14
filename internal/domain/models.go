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

// Credential represents a credential in the system
type Credential struct {
	ID         string    `json:"id"`
	ResourceID string    `json:"resource_id"`
	Type       string    `json:"type"`
	Secret     string    `json:"secret"`
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
	Username       string    `json:"username"`
	ResourceID     string    `json:"resource_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time,omitempty"`
	ExpiresAt      time.Time `json:"expires_at"`
	Status         string    `json:"status"` // "active", "completed", "terminated"
	ClientIP       string    `json:"client_ip"`
	ClientMetadata string    `json:"client_metadata,omitempty"`
	AuditPath      string    `json:"audit_path,omitempty"` // Path to session recording
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AccessRequest represents a request for access to a resource
type AccessRequest struct {
	ID          string        `json:"id"`
	ResourceID  string        `json:"resource_id"`
	UserID      string        `json:"user_id"`
	Reason      string        `json:"reason"`
	Duration    time.Duration `json:"duration"`
	ApproverID  string        `json:"approver_id"`
	Comment     string        `json:"comment"`
	Status      string        `json:"status"`
	RequestedAt time.Time     `json:"requested_at"`
	ReviewedAt  time.Time     `json:"reviewed_at"`
	ExpiresAt   time.Time     `json:"expires_at"`
	ReviewerID  string        `json:"reviewer_id"`
	ReviewNotes string        `json:"review_notes"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
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
	Duration   string    `json:"duration"`
	Used       bool      `json:"used"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ResourceID string    `json:"resource_id"`
	Action     string    `json:"action"`
	Details    string    `json:"details"`
	IP         string    `json:"ip"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}
