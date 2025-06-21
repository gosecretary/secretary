-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS resources (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS credentials (
    id UUID PRIMARY KEY,
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    secret TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(user_id, resource_id, action)
);

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL,
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    client_ip VARCHAR(45) NOT NULL,
    client_metadata TEXT,
    audit_path TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS access_requests (
    id UUID PRIMARY KEY,
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    duration BIGINT NOT NULL, -- Duration in nanoseconds
    approver_id UUID REFERENCES users(id) ON DELETE SET NULL,
    comment TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    requested_at TIMESTAMP NOT NULL,
    reviewed_at TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    reviewer_id UUID REFERENCES users(id) ON DELETE SET NULL,
    review_notes TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS ephemeral_credentials (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    token VARCHAR(255) UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    duration VARCHAR(50) NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    action VARCHAR(255) NOT NULL,
    details TEXT,
    ip VARCHAR(45) NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS ephemeral_credentials;
DROP TABLE IF EXISTS access_requests;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS users; 