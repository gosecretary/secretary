package domain

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Password is never exposed in JSON
	Name      string    `json:"name"`
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
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
} 