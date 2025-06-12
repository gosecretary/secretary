# Secretary Project

A secure credential management system with role-based access control.

## Features

- User authentication and authorization
- Role-based access control (RBAC)
- Secure credential storage
- Resource management
- API-first design

## Project Structure

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

2. Set up environment variables:
   ```bash
   export SERVER_PORT=8080
   export DB_DRIVER=sqlite3
   export DB_NAME=secretary.db
   export JWT_SECRET=your-secret-key
   export SESSION_SECRET=your-session-secret
   ```

3. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

## API Endpoints

### Public Endpoints
- `POST /api/register` - Register a new user
- `POST /api/login` - Login user
- `GET /health` - Health check

### Protected Endpoints
- `GET /api/users/{id}` - Get user details
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user

## Security Features

- Password hashing using bcrypt
- JWT-based authentication
- SQL injection protection
- Input validation
- Role-based access control
- Secure session management

## Development

1. Run tests:
   ```bash
   go test ./...
   ```

2. Run linter:
   ```bash
   go vet ./...
   ```

## License

MIT License - see LICENSE file for details
