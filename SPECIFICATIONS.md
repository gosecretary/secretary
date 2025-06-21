# Secretary Project Specifications

## Table of Contents
1. [System Overview](#system-overview)
2. [Architecture Specifications](#architecture-specifications)
3. [Security Rules and Policies](#security-rules-and-policies)
4. [API Specifications](#api-specifications)
5. [Data Models and Validation](#data-models-and-validation)
6. [Authentication and Authorization](#authentication-and-authorization)
7. [Session Management](#session-management)
8. [Access Control Rules](#access-control-rules)
9. [Audit and Compliance](#audit-and-compliance)
10. [Operational Rules](#operational-rules)
11. [Deployment Specifications](#deployment-specifications)
12. [Testing Requirements](#testing-requirements)

## System Overview

### Purpose
Secretary is an open-source privileged access management (PAM) system designed to provide secure, auditable access to infrastructure resources. It acts as a centralized gateway for managing access to databases, SSH servers, and other critical infrastructure components.

### Core Principles
- **Zero Trust**: No inherent trust in network or client
- **Just-In-Time Access**: Temporary, time-limited access credentials
- **Complete Audit Trail**: All access attempts and activities logged
- **Human-in-the-Loop**: Access requests require approval workflow
- **Session Monitoring**: Real-time monitoring and recording of sessions

### Key Features
- User authentication and session management
- Resource management with role-based access control
- Access request workflow with approval process
- Ephemeral credential generation
- Session monitoring and recording
- Comprehensive audit logging
- Real-time session moderation capabilities

## Architecture Specifications

### System Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client Apps   │    │   Admin Portal  │    │   API Gateway   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                └───────────────────────┼───────────────────────┐
                                                        │                       │
                                ┌───────────────────────┼───────────────────────┘
                                │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Audit System   │    │  Session Store  │    │  Credential     │
│                 │    │                 │    │  Generator      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                └───────────────────────┼───────────────────────┐
                                                        │                       │
                                ┌───────────────────────┼───────────────────────┘
                                │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Database      │    │   Rate Limiter  │    │   Proxy         │
│   (SQLite)      │    │                 │    │   Service       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Component Specifications

#### 1. API Gateway
- **Protocol**: HTTP/HTTPS (TLS 1.2+)
- **Port**: 6080 (configurable)
- **Framework**: Gorilla Mux
- **Middleware Stack**: CORS, Rate Limiting, Session Management, Logging, Recovery

#### 2. Database Layer
- **Primary Database**: SQLite3 (production: PostgreSQL recommended)
- **Schema Versioning**: Migration-based
- **Connection Pooling**: Built-in Go database/sql
- **Encryption**: At-rest encryption recommended

#### 3. Session Store
- **Type**: In-memory with persistence
- **Expiration**: Configurable (default: 24 hours)
- **Cleanup**: Automatic expired session removal
- **Security**: HttpOnly, Secure, SameSite cookies

#### 4. Credential Generator
- **Algorithm**: Cryptographically secure random generation
- **Length**: Username (8-16 chars), Password (24 chars), Token (32 chars)
- **Character Set**: Alphanumeric + special characters
- **Expiration**: Configurable (default: 8 hours)

## Security Rules and Policies

### 1. Authentication Rules

#### Password Policy
- **Minimum Length**: 8 characters
- **Maximum Length**: 128 characters
- **Complexity Requirements**: At least 3 of:
  - Uppercase letters (A-Z)
  - Lowercase letters (a-z)
  - Digits (0-9)
  - Special characters (!@#$%^&*()_+-=[]{}|;':",./<>?~`)
- **Hashing**: bcrypt with configurable cost (default: 12)
- **History**: No password reuse within 5 generations
- **Expiration**: 90 days (configurable)

#### Username Policy
- **Length**: 3-32 characters
- **Characters**: Letters, numbers, dots, hyphens, underscores only
- **Validation**: UTF-8 encoding required
- **Reserved**: No system-reserved usernames (admin, root, system, etc.)

#### Session Policy
- **Duration**: Maximum 24 hours
- **Inactivity Timeout**: 30 minutes (configurable)
- **Concurrent Sessions**: Maximum 3 per user
- **IP Binding**: Session bound to originating IP
- **Rotation**: Automatic session rotation on privilege escalation

### 2. Authorization Rules

#### Role-Based Access Control (RBAC)
- **Roles**: user, admin, reviewer
- **Hierarchy**: admin > reviewer > user
- **Permission Inheritance**: Higher roles inherit lower role permissions
- **Resource Access**: Users can only access resources they have explicit permissions for

#### Permission Matrix
```
Role          | Resources | Users | Sessions | Credentials | Audit
--------------|-----------|-------|----------|-------------|-------
user          | Read Own  | Read  | Read Own | Generate    | Read Own
reviewer      | Read All  | Read  | Read All | Generate    | Read All
admin         | Full      | Full  | Full     | Full        | Full
```

#### Access Request Rules
- **Approval Required**: All access requests require human approval
- **Escalation**: Automatic escalation after 4 hours if no response
- **Expiration**: Maximum 24 hours for approved requests
- **Justification**: Required reason field (max 1000 characters)
- **Audit**: All request/approval actions logged

### 3. Input Validation Rules

#### General Validation
- **Encoding**: UTF-8 required for all text inputs
- **Length Limits**: Enforced on all string fields
- **Pattern Matching**: Regex validation for structured fields
- **Sanitization**: HTML/script tag removal
- **SQL Injection**: Pattern-based detection and blocking

#### Field-Specific Rules
- **Hostnames**: Valid DNS names or IP addresses only
- **Ports**: 1-65535 range
- **UUIDs**: Valid UUID v4 format
- **Emails**: RFC 5322 compliant email format
- **Timestamps**: ISO 8601 format

### 4. Rate Limiting Rules
- **Default Limit**: 100 requests per hour per IP
- **Burst Limit**: 10 requests per minute
- **Authentication Endpoints**: 5 attempts per 15 minutes
- **Admin Endpoints**: 50 requests per hour
- **Exceeded Action**: HTTP 429 with retry-after header

## API Specifications

### Base URL
- **Development**: `http://localhost:6080`
- **Production**: `https://your-domain.com`

### Authentication
- **Method**: Session-based (HTTP cookies)
- **Cookie Name**: `session_id`
- **Attributes**: HttpOnly, Secure, SameSite=Strict
- **Expiration**: 24 hours from login

### Response Format
```json
{
  "success": true,
  "code": 200,
  "message": "Operation completed successfully",
  "data": {
    // Response payload
  }
}
```

### Error Response Format
```json
{
  "success": false,
  "code": 400,
  "message": "Bad request",
  "data": {
    "error": "Detailed error information",
    "field": "field_name"
  }
}
```

### HTTP Status Codes
- **200**: Success
- **201**: Created
- **400**: Bad Request
- **401**: Unauthorized
- **403**: Forbidden
- **404**: Not Found
- **429**: Too Many Requests
- **500**: Internal Server Error

### API Endpoints

#### Public Endpoints
- `POST /api/login` - User authentication
- `GET /health` - Health check

#### Protected Endpoints
- `POST /api/register` - User registration (requires authentication)
- `GET /api/users/{id}` - Get user details
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user
- `GET /api/resources` - List resources
- `POST /api/resources` - Create resource
- `GET /api/resources/{id}` - Get resource details
- `PUT /api/resources/{id}` - Update resource
- `DELETE /api/resources/{id}` - Delete resource
- `GET /api/sessions` - List active sessions
- `POST /api/sessions/{id}/terminate` - Terminate session
- `POST /api/access-requests` - Create access request
- `GET /api/access-requests` - List pending requests
- `POST /api/access-requests/{id}/approve` - Approve request
- `POST /api/access-requests/{id}/deny` - Deny request
- `POST /api/ephemeral-credentials` - Generate credentials
- `GET /api/ephemeral-credentials/{id}` - Get credential details

## Data Models and Validation

### User Model
```go
type User struct {
    ID        string    `json:"id" validate:"required,uuid"`
    Username  string    `json:"username" validate:"required,min=3,max=32,alphanum"`
    Email     string    `json:"email" validate:"required,email"`
    Password  string    `json:"-" validate:"required,min=8,max=128"`
    Name      string    `json:"name" validate:"required,max=100"`
    Role      string    `json:"role" validate:"required,oneof=user admin reviewer"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Resource Model
```go
type Resource struct {
    ID          string    `json:"id" validate:"required,uuid"`
    Name        string    `json:"name" validate:"required,max=64,alphanum"`
    Description string    `json:"description" validate:"max=500"`
    Type        string    `json:"type" validate:"required,oneof=mysql postgresql ssh redis mongodb"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Session Model
```go
type Session struct {
    ID             string    `json:"id" validate:"required,uuid"`
    UserID         string    `json:"user_id" validate:"required,uuid"`
    ResourceID     string    `json:"resource_id" validate:"required,uuid"`
    StartTime      time.Time `json:"start_time"`
    EndTime        time.Time `json:"end_time,omitempty"`
    Status         string    `json:"status" validate:"required,oneof=active completed terminated"`
    ClientIP       string    `json:"client_ip" validate:"required,ip"`
    ClientMetadata string    `json:"client_metadata,omitempty"`
    AuditPath      string    `json:"audit_path,omitempty"`
    CreatedAt      time.Time `json:"created_at"`
}
```

### Access Request Model
```go
type AccessRequest struct {
    ID          string        `json:"id" validate:"required,uuid"`
    UserID      string        `json:"user_id" validate:"required,uuid"`
    ResourceID  string        `json:"resource_id" validate:"required,uuid"`
    Reason      string        `json:"reason" validate:"required,max=1000"`
    Status      string        `json:"status" validate:"required,oneof=pending approved denied"`
    ReviewerID  string        `json:"reviewer_id,omitempty" validate:"omitempty,uuid"`
    ReviewNotes string        `json:"review_notes,omitempty" validate:"max=1000"`
    RequestedAt time.Time     `json:"requested_at"`
    ReviewedAt  time.Time     `json:"reviewed_at,omitempty"`
    ExpiresAt   time.Time     `json:"expires_at,omitempty"`
    CreatedAt   time.Time     `json:"created_at"`
}
```

## Authentication and Authorization

### Authentication Flow
1. **Login Request**: Client sends username/password
2. **Validation**: Server validates credentials against database
3. **Session Creation**: Server creates session with unique ID
4. **Cookie Setting**: Server sets HttpOnly session cookie
5. **Response**: Server returns user information and success status

### Authorization Flow
1. **Request**: Client makes authenticated request with session cookie
2. **Session Validation**: Server validates session and retrieves user
3. **Permission Check**: Server checks user permissions for requested resource
4. **Access Decision**: Server grants or denies access based on permissions
5. **Audit Logging**: Server logs the access attempt

### Session Management
- **Creation**: On successful login
- **Validation**: On every protected request
- **Expiration**: Automatic cleanup of expired sessions
- **Termination**: Manual termination or automatic on logout
- **Security**: Session ID rotation on privilege changes

## Access Control Rules

### Resource Access Rules
1. **Default Deny**: All access denied unless explicitly granted
2. **Least Privilege**: Users receive minimum necessary permissions
3. **Temporal Limits**: Access expires after specified duration
4. **Approval Required**: All access requires approval workflow
5. **Audit Trail**: All access attempts logged with full context

### Permission Inheritance
- **User Level**: Individual user permissions
- **Role Level**: Role-based permissions
- **Resource Level**: Resource-specific permissions
- **Temporal Level**: Time-based permissions

### Access Request Workflow
1. **Request Submission**: User submits access request with justification
2. **Notification**: System notifies designated reviewers
3. **Review Process**: Reviewers evaluate request based on business rules
4. **Decision**: Approve, deny, or request additional information
5. **Implementation**: System grants or denies access based on decision
6. **Audit**: Complete audit trail maintained

## Audit and Compliance

### Audit Logging Requirements
- **Authentication Events**: All login/logout attempts
- **Authorization Events**: All access decisions
- **Data Access**: All resource access attempts
- **Configuration Changes**: All system configuration modifications
- **Session Events**: Session creation, modification, termination
- **Error Events**: All security-related errors

### Log Format
```json
{
  "timestamp": "2023-01-01T12:00:00Z",
  "event_type": "authentication",
  "user_id": "uuid",
  "resource_id": "uuid",
  "action": "login_success",
  "ip_address": "192.168.1.1",
  "user_agent": "Mozilla/5.0...",
  "details": {
    "session_id": "uuid",
    "duration_ms": 150
  }
}
```

### Compliance Standards
- **SOC 2**: Access controls, change management, monitoring
- **ISO 27001**: Information security management
- **NIST Cybersecurity Framework**: Authentication, access control, logging
- **OWASP Top 10**: SQL injection, XSS, CSRF prevention

### Data Retention
- **Audit Logs**: 7 years minimum
- **Session Data**: 90 days
- **User Data**: Until account deletion
- **Resource Data**: Until resource deletion
- **Access Requests**: 3 years

## Operational Rules

### Monitoring Requirements
- **System Health**: CPU, memory, disk usage
- **Security Events**: Failed logins, rate limit violations
- **Performance Metrics**: Response times, throughput
- **Error Rates**: Application and system errors
- **Session Metrics**: Active sessions, session duration

### Alerting Rules
- **Critical**: System unavailable, security breach detected
- **High**: High error rate, authentication failures
- **Medium**: Performance degradation, resource usage high
- **Low**: Informational events, successful operations

### Backup Requirements
- **Database**: Daily automated backups
- **Configuration**: Version-controlled configuration files
- **Audit Logs**: Real-time replication to secure storage
- **Recovery**: 4-hour RTO, 1-hour RPO

### Maintenance Windows
- **Security Updates**: Monthly, 2-hour window
- **Feature Updates**: Quarterly, 4-hour window
- **Database Maintenance**: Weekly, 1-hour window
- **Emergency**: As needed with 1-hour notice

## Deployment Specifications

### Environment Requirements
- **Operating System**: Linux (Ubuntu 20.04+ or Alpine 3.19+)
- **Go Version**: 1.23.0+
- **Memory**: Minimum 2GB RAM
- **Storage**: Minimum 10GB available space
- **Network**: HTTPS access required

### Security Configuration
- **TLS**: TLS 1.2+ with secure cipher suites
- **Certificates**: Valid SSL certificates required
- **Firewall**: Restrict access to necessary ports only
- **Network**: Isolated network segment recommended
- **Updates**: Regular security updates required

### Docker Deployment
```dockerfile
# Security-focused Dockerfile
FROM golang:1.23-alpine AS build
RUN apk update && apk upgrade
RUN apk add --no-cache gcc musl-dev ca-certificates
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-w -s" -o secretary cmd/secretary/main.go

FROM alpine:3.19
RUN apk update && apk upgrade
RUN adduser -D -s /bin/sh -u 1000 secretary
USER secretary
COPY --from=build --chown=secretary:secretary /workspace/secretary /app/secretary
WORKDIR /app
EXPOSE 6080
ENTRYPOINT ["/app/secretary"]
```

### Environment Variables
```bash
# Required
SECRETARY_HOST=0.0.0.0
SECRETARY_PORT=6080
SECRETARY_SESSION_SECRET=<random-32-byte-hex>
SECRETARY_CSRF_SECRET=<random-32-byte-hex>

# Security
SECRETARY_TLS_CERT_PATH=/path/to/cert.pem
SECRETARY_TLS_KEY_PATH=/path/to/key.pem
SECRETARY_SECURE_COOKIES=true
SECRETARY_BCRYPT_COST=12

# Rate Limiting
SECRETARY_RATE_LIMIT_REQUESTS=100
SECRETARY_RATE_LIMIT_WINDOW=1h

# Database
SECRETARY_DB_DRIVER=sqlite3
SECRETARY_DB_PATH=./data/secretary.db

# Audit
SECRETARY_AUDIT_ENABLED=true
SECRETARY_AUDIT_DIR=./data/audit/
```

## Testing Requirements

### Unit Testing
- **Coverage**: Minimum 80% code coverage
- **Authentication**: All authentication flows tested
- **Authorization**: All permission checks tested
- **Validation**: All input validation rules tested
- **Error Handling**: All error scenarios tested

### Integration Testing
- **API Endpoints**: All API endpoints tested
- **Database Operations**: All database operations tested
- **Session Management**: Session lifecycle tested
- **Access Control**: Complete access control flow tested

### Security Testing
- **Penetration Testing**: Annual security assessment
- **Vulnerability Scanning**: Monthly automated scans
- **Code Review**: Security-focused code reviews
- **Dependency Scanning**: Regular dependency updates

### Performance Testing
- **Load Testing**: 1000 concurrent users
- **Stress Testing**: System limits under load
- **Endurance Testing**: 24-hour continuous operation
- **Scalability Testing**: Horizontal scaling validation

### Compliance Testing
- **Audit Logging**: Complete audit trail validation
- **Data Protection**: Data encryption and protection
- **Access Controls**: Permission enforcement validation
- **Session Security**: Session management security

---

## Version History
- **v1.0.0**: Initial specification
- **Date**: 2024-01-01
- **Author**: Secretary Development Team
- **Status**: Draft for Review

## Approval Process
1. **Technical Review**: Development team review
2. **Security Review**: Security team validation
3. **Compliance Review**: Legal/compliance team approval
4. **Stakeholder Approval**: Business stakeholder sign-off
5. **Final Approval**: CTO/Architecture team approval

---

*This specification document is a living document and should be updated as the system evolves. All changes must go through the approval process outlined above.* 