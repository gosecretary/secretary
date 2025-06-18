# Secretary API Documentation

This directory contains the API documentation for the Secretary privileged access management system.

## Files

- `swagger.yaml` - OpenAPI 3.0 specification for the Secretary API
- `swagger-ui.html` - HTML file that displays the Swagger documentation using Swagger UI
- `README.md` - This file

## Viewing the Documentation

### Option 1: Through the Secretary Server (Recommended)

1. Start the Secretary server:
   ```bash
   ./scripts/run.sh server --dev
   ```

2. Open your web browser and navigate to:
   ```
   http://localhost:6080/docs
   ```

This will display the interactive Swagger UI where you can:
- Browse all API endpoints
- View request/response schemas
- Test API endpoints directly from the browser
- Download the OpenAPI specification

### Option 2: Local File Viewing

1. Open the `swagger-ui.html` file directly in your web browser
2. Make sure the `swagger.yaml` file is in the same directory

### Option 3: Using External Tools

You can use the `swagger.yaml` file with various tools:

#### Swagger Editor
1. Go to [editor.swagger.io](https://editor.swagger.io/)
2. Copy and paste the contents of `swagger.yaml`

#### Swagger Codegen
Generate client libraries in various languages:
```bash
# Install swagger-codegen
npm install -g @apidevtools/swagger-cli

# Generate client for different languages
swagger-codegen generate -i docs/swagger.yaml -l go -o ./client/go
swagger-codegen generate -i docs/swagger.yaml -l python -o ./client/python
swagger-codegen generate -i docs/swagger.yaml -l javascript -o ./client/javascript
```

#### Postman
1. Open Postman
2. Click "Import"
3. Select the `swagger.yaml` file to import all endpoints as a Postman collection

## API Overview

The Secretary API provides the following functionality:

### Authentication
- **POST /api/register** - Register a new user account
- **POST /api/login** - Authenticate and create a session

### User Management
- **GET /api/users/{id}** - Get user information
- **PUT /api/users/{id}** - Update user details
- **DELETE /api/users/{id}** - Delete user account

### Resource Management
- **GET /api/resources** - List all available resources
- **POST /api/resources** - Create a new resource
- **GET /api/resources/{id}** - Get resource details
- **PUT /api/resources/{id}** - Update resource information
- **DELETE /api/resources/{id}** - Delete a resource

### Session Management
- **GET /api/sessions** - List active sessions
- **GET /api/sessions/{id}** - Get session details
- **POST /api/sessions/{id}/terminate** - Terminate a session

### Access Request Workflow
- **GET /api/access-requests** - List pending access requests
- **POST /api/access-requests** - Submit a new access request
- **POST /api/access-requests/{id}/approve** - Approve an access request
- **POST /api/access-requests/{id}/deny** - Deny an access request

### Ephemeral Credentials
- **POST /api/ephemeral-credentials** - Generate temporary credentials
- **GET /api/ephemeral-credentials/{id}** - Get credential details

### Health Check
- **GET /health** - Check system health

## Authentication

Most endpoints require authentication via session cookies. The typical workflow is:

1. **Register** a user account (or use the dev admin account)
2. **Login** to obtain a session cookie
3. Use the session cookie for subsequent API calls

## Response Format

All API responses follow a consistent format:

```json
{
  "success": true,
  "code": 200,
  "message": "Operation completed successfully",
  "data": {
    // Response payload varies by endpoint
  }
}
```

Error responses follow the same format but with `success: false`:

```json
{
  "success": false,
  "code": 400,
  "message": "Bad request",
  "data": {
    "error": "Detailed error information"
  }
}
```

## Example Usage

### 1. Login to Get Session Cookie

```bash
curl -v -X POST "localhost:6080/api/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "your_password"}'
```

### 2. List Resources (using session cookie)

```bash
curl -v "localhost:6080/api/resources" \
  -H "Content-Type: application/json" \
  --cookie "session_id=your_session_id"
```

### 3. Create an Access Request

```bash
curl -v -X POST "localhost:6080/api/access-requests" \
  -H "Content-Type: application/json" \
  --cookie "session_id=your_session_id" \
  -d '{
    "resource_id": "resource_uuid",
    "reason": "Need access for debugging",
    "duration": 3600000000000
  }'
```

## Development

To update the API documentation:

1. Modify the `swagger.yaml` file
2. Restart the Secretary server to serve the updated documentation
3. Test the documentation by visiting `http://localhost:6080/docs`

## Security Considerations

- The documentation endpoint is publicly accessible for development
- In production, consider restricting access to the `/docs` endpoint
- Never expose sensitive information in the API documentation
- Use HTTPS in production environments 