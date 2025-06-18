package domain

import (
	"context"
	"time"
)

// UserService defines the interface for user-related operations
type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	Authenticate(ctx context.Context, username, password string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// ResourceService defines the interface for resource-related operations
type ResourceService interface {
	CreateResource(ctx context.Context, resource *Resource) error
	ListResources(ctx context.Context) ([]*Resource, error)
	GetResource(ctx context.Context, id string) (*Resource, error)
	UpdateResource(ctx context.Context, resource *Resource) error
	DeleteResource(ctx context.Context, id string) error
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
	CreateCredential(ctx context.Context, credential *Credential) error
	ListCredentials(ctx context.Context) ([]*Credential, error)
	GetCredential(ctx context.Context, id string) (*Credential, error)
	GetCredentialByResourceID(ctx context.Context, resourceID string) ([]*Credential, error)
	UpdateCredential(ctx context.Context, credential *Credential) error
	DeleteCredential(ctx context.Context, id string) error
	RotateCredential(ctx context.Context, id string) error
}

// PermissionService defines the interface for permission-related operations
type PermissionService interface {
	CreatePermission(ctx context.Context, permission *Permission) error
	ListPermissions(ctx context.Context) ([]*Permission, error)
	GetPermission(ctx context.Context, id string) (*Permission, error)
	GetPermissionByUserID(ctx context.Context, userID string) ([]*Permission, error)
	GetPermissionByResourceID(ctx context.Context, resourceID string) ([]*Permission, error)
	UpdatePermission(ctx context.Context, permission *Permission) error
	DeletePermission(ctx context.Context, id string) error
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
	List(ctx context.Context) ([]*Session, error)
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
	CreateAccessRequest(ctx context.Context, request *AccessRequest) error
	ListAccessRequests(ctx context.Context) ([]*AccessRequest, error)
	GetAccessRequest(ctx context.Context, id string) (*AccessRequest, error)
	GetAccessRequestByUserID(ctx context.Context, userID string) ([]*AccessRequest, error)
	GetAccessRequestByResourceID(ctx context.Context, resourceID string) ([]*AccessRequest, error)
	GetPendingAccessRequests(ctx context.Context) ([]*AccessRequest, error)
	UpdateAccessRequest(ctx context.Context, request *AccessRequest) error
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
	Create(ctx context.Context, credential *EphemeralCredential) (*EphemeralCredential, error)
	List(ctx context.Context) ([]*EphemeralCredential, error)
	GetEphemeralCredential(ctx context.Context, id string) (*EphemeralCredential, error)
	DeleteEphemeralCredential(ctx context.Context, id string) error
	MarkAsUsedEphemeralCredential(ctx context.Context, id string) error
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

// SessionStore defines the interface for session storage operations
type SessionStore interface {
	Get(id string) (*Session, error)
	Set(session *Session) error
	Delete(id string) error
}

// AuditLogService defines the interface for audit log operations
type AuditLogService interface {
	List(ctx context.Context) ([]*AuditLog, error)
	GetByID(ctx context.Context, id string) (*AuditLog, error)
	GetByUserID(ctx context.Context, userID string) ([]*AuditLog, error)
	GetByResourceID(ctx context.Context, resourceID string) ([]*AuditLog, error)
	GetByAction(ctx context.Context, action string) ([]*AuditLog, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*AuditLog, error)
}

// SessionCommandService defines the interface for session command operations
type SessionCommandService interface {
	RecordCommand(ctx context.Context, command *SessionCommand) error
	GetSessionCommands(ctx context.Context, sessionID string) ([]*SessionCommand, error)
	GetCommandsByUser(ctx context.Context, userID string) ([]*SessionCommand, error)
	GetCommandsByResource(ctx context.Context, resourceID string) ([]*SessionCommand, error)
	GetHighRiskCommands(ctx context.Context) ([]*SessionCommand, error)
	AnalyzeCommand(ctx context.Context, command string, commandType string) (risk string, shouldBlock bool, err error)
}

// SessionRecordingService defines the interface for session recording operations
type SessionRecordingService interface {
	StartRecording(ctx context.Context, sessionID string) (*SessionRecording, error)
	StopRecording(ctx context.Context, sessionID string) error
	GetRecording(ctx context.Context, sessionID string) (*SessionRecording, error)
	GetRecordingFile(ctx context.Context, recordingID string) ([]byte, error)
	DeleteRecording(ctx context.Context, recordingID string) error
	ListRecordings(ctx context.Context, userID string) ([]*SessionRecording, error)
}

// ProxyService defines the interface for proxy operations
type ProxyService interface {
	CreateProxy(ctx context.Context, sessionID string, protocol string, remoteHost string, remotePort int) (*ProxyConnection, error)
	StartProxy(ctx context.Context, proxyID string) (localPort int, err error)
	StopProxy(ctx context.Context, proxyID string) error
	GetActiveProxies(ctx context.Context) ([]*ProxyConnection, error)
	GetProxyBySession(ctx context.Context, sessionID string) (*ProxyConnection, error)
	UpdateProxyStats(ctx context.Context, proxyID string, bytesIn, bytesOut int64) error
}

// SecurityAlertService defines the interface for security alert operations
type SecurityAlertService interface {
	CreateAlert(ctx context.Context, alert *SecurityAlert) error
	GetAlerts(ctx context.Context, sessionID string) ([]*SecurityAlert, error)
	GetAlertsByUser(ctx context.Context, userID string) ([]*SecurityAlert, error)
	GetAlertsBySeverity(ctx context.Context, severity string) ([]*SecurityAlert, error)
	MarkAlertAsReviewed(ctx context.Context, alertID string) error
}

// SessionMonitorService defines the interface for real-time session monitoring
type SessionMonitorService interface {
	StartMonitoring(ctx context.Context, sessionID string) error
	StopMonitoring(ctx context.Context, sessionID string) error
	GetLiveSession(ctx context.Context, sessionID string) (*Session, error)
	InterruptSession(ctx context.Context, sessionID string, reason string) error
	GetSessionMetrics(ctx context.Context, sessionID string) (map[string]interface{}, error)
}
