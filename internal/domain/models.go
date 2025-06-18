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

// SessionCommand represents a command executed during a session
type SessionCommand struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	UserID      string    `json:"user_id"`
	ResourceID  string    `json:"resource_id"`
	Command     string    `json:"command"`
	CommandType string    `json:"command_type"` // "sql", "ssh", "shell", etc.
	Response    string    `json:"response,omitempty"`
	Status      string    `json:"status"` // "executed", "blocked", "failed"
	Risk        string    `json:"risk"`   // "low", "medium", "high", "critical"
	Timestamp   time.Time `json:"timestamp"`
	Duration    int64     `json:"duration_ms"` // Command execution time in milliseconds
	CreatedAt   time.Time `json:"created_at"`
}

// SessionRecording represents a complete session recording
type SessionRecording struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id"`
	UserID        string    `json:"user_id"`
	ResourceID    string    `json:"resource_id"`
	RecordingPath string    `json:"recording_path"` // Path to the recording file
	Format        string    `json:"format"`         // "asciinema", "text", "binary"
	Size          int64     `json:"size"`           // File size in bytes
	Duration      int64     `json:"duration"`       // Session duration in seconds
	CommandCount  int       `json:"command_count"`  // Total commands executed
	CreatedAt     time.Time `json:"created_at"`
}

// ProxyConnection represents an active proxy connection
type ProxyConnection struct {
	ID           string    `json:"id"`
	SessionID    string    `json:"session_id"`
	UserID       string    `json:"user_id"`
	ResourceID   string    `json:"resource_id"`
	Protocol     string    `json:"protocol"`    // "ssh", "mysql", "postgres", etc.
	LocalPort    int       `json:"local_port"`  // Local proxy port
	RemoteHost   string    `json:"remote_host"` // Target resource host
	RemotePort   int       `json:"remote_port"` // Target resource port
	Status       string    `json:"status"`      // "active", "closed", "error"
	BytesIn      int64     `json:"bytes_in"`    // Bytes received from client
	BytesOut     int64     `json:"bytes_out"`   // Bytes sent to target
	LastActivity time.Time `json:"last_activity"`
	CreatedAt    time.Time `json:"created_at"`
}

// SecurityAlert represents a security alert during a session
type SecurityAlert struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	CommandID   string    `json:"command_id,omitempty"`
	UserID      string    `json:"user_id"`
	ResourceID  string    `json:"resource_id"`
	AlertType   string    `json:"alert_type"` // "suspicious_command", "data_exfiltration", "privilege_escalation"
	Severity    string    `json:"severity"`   // "low", "medium", "high", "critical"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	RawData     string    `json:"raw_data"` // The actual command or data that triggered the alert
	Action      string    `json:"action"`   // "logged", "blocked", "terminated"
	CreatedAt   time.Time `json:"created_at"`
}
