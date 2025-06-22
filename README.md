# Secretary Project: Open Infrastructure Access Management [![CodeQL](https://github.com/gosecretary/secretary/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/gosecretary/secretary/actions/workflows/github-code-scanning/codeql) [![Secretary Test](https://github.com/gosecretary/secretary/actions/workflows/test-pipeline.yml/badge.svg)](https://github.com/gosecretary/secretary/actions/workflows/test-pipeline.yml)

Secretary is an open-source solution for infrastructure access management, combining core ideas from tools like HashiCorp Boundary and Teleport—but with all features available for free and fully open.

## Overview

Secretary acts as a secure gateway and policy enforcement point for accessing infrastructure resources such as databases and SSH servers. It introduces a centralized model for managing, auditing, and controlling access across distributed environments.

## Key Features

* **Session Auditing**: All connections to target resources (e.g., databases, SSH servers) are routed through the system and fully audited. Session metadata is persisted in the application database for traceability.
* **Ephemeral Credentials**: When access is granted, Secretary dynamically generates temporary user credentials (e.g., DB username/password) that expire after use, eliminating the need for static secrets.
* **Session Moderation**: Supports real-time session monitoring and moderation, allowing administrators or reviewers to observe or terminate active sessions if needed.
* **Access Request Workflow**: Users can request access to specific resources, and designated reviewers can approve or deny these requests, adding a layer of human-in-the-loop approval before access is granted.
* **Role-Based Access Control (RBAC)**: Fine-grained permissioning model based on roles, ensuring users only see and access resources relevant to their privileges.

## Security & Observability

* All connections are secured, either via mutual TLS or encrypted tunnels
* Session activity is monitored, logged, and tamper-resistant
* Designed with zero trust principles, assuming no inherent trust in the network or client

## Architecture

```
secretary/
├── cmd/                    # Application entry points
│   └── server/            # Main server application
├── internal/              # Private application code
│   ├── domain/           # Domain models and interfaces
│   ├── service/          # Business logic
│   ├── repository/       # Data access layer
│   ├── middleware/       # HTTP middleware
│   └── config/           # Configuration
├── pkg/                   # Public libraries
│   ├── auth/             # Authentication utilities
│   └── validator/        # Input validation
├── api/                  # API handlers
│   ├── http/            # HTTP handlers
│   └── rest/            # REST API definitions
└── scripts/             # Build and deployment scripts
```

## Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Create a config file or set environment variables:
   ```bash
   cp config.example.env .env
   # Edit .env with your settings
   ```

3. Build the application:
   ```bash
   # Using the build script
   ./scripts/build.sh
   
   # Or manually
   go build -o bin/secretary cmd/server/main.go
   ```

4. Run the application:
   ```bash
   # Using the run script
   ./scripts/run.sh
   
   # Or manually
   ./bin/secretary
   ```

## Usage

### Starting the Application

1. Make sure you have the required environment variables set in your `.env` file:
   ```bash
   # Database configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=secretary
   
   # Server configuration
   SERVER_PORT=8080
   SERVER_HOST=localhost
   
   # JWT configuration
   JWT_SECRET=your_secret_key
   JWT_EXPIRATION=24h
   ```

2. Start the application using the run script:
   ```bash
   ./scripts/run.sh
   ```

### Making Access Requests

1. Use the request script to create an access request:
   ```bash
   ./scripts/request.sh create \
     --user-id "user123" \
     --resource-id "db123" \
     --reason "Need access for maintenance"
   ```

2. Check the status of your request:
   ```bash
   ./scripts/request.sh status --request-id "req123"
   ```

3. Approve or deny requests (admin only):
   ```bash
   ./scripts/request.sh approve --request-id "req123"
   ./scripts/request.sh deny --request-id "req123"
   ```

### Session Management

1. List active sessions:
   ```bash
   curl -H "Authorization: Bearer $JWT_TOKEN" http://localhost:8080/api/sessions
   ```

2. Monitor a specific session:
   ```bash
   curl -H "Authorization: Bearer $JWT_TOKEN" http://localhost:8080/api/sessions/{session_id}
   ```

3. Terminate a session if needed:
   ```bash
   curl -X POST -H "Authorization: Bearer $JWT_TOKEN" http://localhost:8080/api/sessions/{session_id}/terminate
   ```

### Ephemeral Credentials

1. Generate ephemeral credentials for a resource:
   ```bash
   curl -X POST -H "Authorization: Bearer $JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"resource_id": "db123", "duration": "1h"}' \
     http://localhost:8080/api/ephemeral-credentials
   ```

2. Use the credentials:
   ```bash
   # The credentials will be automatically marked as used when accessed
   # through the Secretary proxy
   ```

### TCP Proxy with Monitoring

Secretary provides a comprehensive TCP proxy solution that allows clients to connect to target servers through Secretary while providing real-time monitoring, command interception, and security auditing.

#### Quick Start

1. Create a session for accessing a resource:
   ```bash
   curl -X POST -H "Authorization: Bearer $JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"user_id": "user123", "resource_id": "resource456"}' \
     http://localhost:8080/api/sessions
   ```

2. Create a proxy for the session:
   ```bash
   curl -X POST -H "Authorization: Bearer $JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"protocol": "ssh", "remote_host": "target-server", "remote_port": 22}' \
     http://localhost:8080/api/sessions/{session_id}/proxy
   ```

3. Start the proxy:
   ```bash
   curl -X POST -H "Authorization: Bearer $JWT_TOKEN" \
     http://localhost:8080/api/proxies/{proxy_id}/start
   ```

4. Connect through the proxy:
   ```bash
   # SSH through proxy
   ssh -p {local_port} user@localhost
   
   # MySQL through proxy
   mysql -h localhost -P {local_port} -u username -p
   
   # PostgreSQL through proxy
   psql -h localhost -p {local_port} -U username -d database
   ```

#### Supported Protocols

- **SSH**: Full command interception and analysis
- **MySQL**: SQL query monitoring and risk assessment
- **PostgreSQL**: SQL query monitoring and risk assessment
- **Generic TCP**: Basic traffic monitoring for other protocols

#### Security Features

- **Real-time Command Analysis**: All commands are analyzed for risk levels (low, medium, high, critical)
- **Automatic Blocking**: Critical commands are automatically blocked
- **Session Recording**: All sessions are automatically recorded
- **Security Alerts**: High-risk activities trigger security alerts
- **Audit Logging**: Complete audit trail of all activities

#### Testing the Proxy

Use the provided test script to quickly test the proxy functionality:

```bash
# Run the full proxy test
./scripts/test_proxy.sh

# Monitor commands in real-time
./scripts/test_proxy.sh --monitor-commands

# Monitor security alerts
./scripts/test_proxy.sh --monitor-alerts

# Get active proxies
./scripts/test_proxy.sh --active-proxies

# Clean up test resources
./scripts/test_proxy.sh --cleanup
```

For detailed proxy documentation, see [docs/PROXY_GUIDE.md](docs/PROXY_GUIDE.md).

## API Documentation

Interactive API documentation is available when running the Secretary server:

- **Swagger UI**: http://localhost:6080/docs
- **OpenAPI Spec**: http://localhost:6080/docs/swagger.yaml

For detailed API documentation, examples, and usage instructions, see [docs/README.md](docs/README.md).

## API Endpoints

### Public Endpoints
- `POST /api/register` - Register a new user
- `POST /api/login` - Login user
- `GET /health` - Health check

### Protected Endpoints (User Management)
- `GET /api/users/{id}` - Get user details
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user

### Protected Endpoints (Resource Management)
- `GET /api/resources` - List all resources
- `POST /api/resources` - Create resource
- `GET /api/resources/{id}` - Get resource details
- `PUT /api/resources/{id}` - Update resource
- `DELETE /api/resources/{id}` - Delete resource

### Protected Endpoints (Session Management)
- `GET /api/sessions` - List active sessions
- `GET /api/sessions/{id}` - Get session details
- `POST /api/sessions/{id}/terminate` - Terminate a session
- `GET /api/users/{user_id}/sessions` - Get user's sessions
- `GET /api/resources/{resource_id}/sessions` - Get resource's sessions

### Protected Endpoints (Access Request Flow)
- `POST /api/access-requests` - Create access request
- `GET /api/access-requests` - List pending access requests
- `GET /api/access-requests/{id}` - Get access request details
- `POST /api/access-requests/{id}/approve` - Approve access request
- `POST /api/access-requests/{id}/deny` - Deny access request

### Protected Endpoints (Ephemeral Credentials)
- `POST /api/ephemeral-credentials` - Generate ephemeral credentials
- `GET /api/ephemeral-credentials/{id}` - Get credential details
- `POST /api/ephemeral-credentials/{id}/use` - Mark credentials as used

### Protected Endpoints (TCP Proxy)
- `POST /api/sessions/{session_id}/proxy` - Create proxy for session
- `POST /api/proxies/{proxy_id}/start` - Start proxy
- `POST /api/proxies/{proxy_id}/stop` - Stop proxy
- `GET /api/proxies/active` - Get active proxies
- `GET /api/sessions/{session_id}/proxy` - Get session proxy

### Protected Endpoints (Session Monitoring)
- `GET /api/sessions/{session_id}/commands` - Get session commands
- `GET /api/commands/high-risk` - Get high-risk commands
- `POST /api/sessions/{session_id}/recording/start` - Start session recording
- `POST /api/sessions/{session_id}/recording/stop` - Stop session recording
- `GET /api/sessions/{session_id}/recording` - Get session recording
- `GET /api/sessions/{session_id}/alerts` - Get session alerts

## Security Features

- Password hashing using bcrypt
- JWT-based authentication
- SQL injection protection
- Input validation
- Role-based access control
- Secure session management
- Dynamic credential generation
- Audit logging

## Development

1. Run tests:
   ```bash
   go test ./...
   ```

2. Run linter:
   ```bash
   go vet ./...
   ```

3. Build the application:
   ```bash
   ./scripts/build.sh
   ```

4. Run the application in development mode:
   ```bash
   ./scripts/run.sh --dev
   ```

5. Generate API documentation:
   ```bash
   # The API documentation is automatically generated in the docs directory
   # when you run the build script
   ```

## Troubleshooting

1. If the application fails to start, check:
   - Database connection settings in `.env`
   - Required environment variables
   - Port availability
   - File permissions on scripts

2. If access requests fail:
   - Verify JWT token is valid
   - Check user permissions
   - Ensure resource exists
   - Check request status

3. If sessions are not working:
   - Verify network connectivity
   - Check resource availability
   - Ensure proper authentication
   - Check session logs

## License

MIT License - see LICENSE file for details
