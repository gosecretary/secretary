#!/bin/bash -e

BASE_URL="http://localhost:6080"

echo "Testing Secretary API..."

# Login
echo "1. Logging in..."
LOGIN_RESPONSE=$(curl -s --fail -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "securepassword123"}')
echo $LOGIN_RESPONSE | jq .

# Extract user ID for further requests
USER_ID=$(echo $LOGIN_RESPONSE | jq -r '.id')
echo "User ID: $USER_ID"

echo ""

# Create a resource
echo "2. Creating a resource..."
RESOURCE_RESPONSE=$(curl -s --fail -X POST "$BASE_URL/api/resources" \
  -H "Content-Type: application/json" \
  -d '{"name": "Production Database", "description": "Main production PostgreSQL database"}')
echo $RESOURCE_RESPONSE | jq .

RESOURCE_ID=$(echo $RESOURCE_RESPONSE | jq -r '.id')
echo "Resource ID: $RESOURCE_ID"

echo ""

# Create credentials for the resource
echo "3. Creating credentials..."
curl -s --fail -X POST "$BASE_URL/api/credentials" \
  -H "Content-Type: application/json" \
  -d "{\"resource_id\": \"$RESOURCE_ID\", \"username\": \"dbuser\", \"password\": \"dbpassword123\"}" | jq .

echo ""

# Create a permission
echo "4. Creating a permission..."
curl -s --fail -X POST "$BASE_URL/api/permissions" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": \"$USER_ID\", \"resource_id\": \"$RESOURCE_ID\", \"action\": \"read\"}" | jq .

echo ""

# List all resources
echo "5. Listing all resources..."
curl -s --fail -X GET "$BASE_URL/api/resources" | jq .

echo ""

# Session API example
echo "6. Creating a session..."
SESSION_RESPONSE=$(curl -s --fail -X POST "$BASE_URL/api/sessions" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": \"$USER_ID\", \"resource_id\": \"$RESOURCE_ID\", \"client_ip\": \"127.0.0.1\"}")
echo $SESSION_RESPONSE | jq .
SESSION_ID=$(echo $SESSION_RESPONSE | jq -r '.id')
echo "Session ID: $SESSION_ID"

echo ""

# Access Request API example
echo "7. Creating an access request..."
ACCESS_REQUEST_RESPONSE=$(curl -s --fail -X POST "$BASE_URL/api/access-requests" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": \"$USER_ID\", \"resource_id\": \"$RESOURCE_ID\", \"reason\": \"Need access for maintenance\"}")
echo $ACCESS_REQUEST_RESPONSE | jq .
ACCESS_REQUEST_ID=$(echo $ACCESS_REQUEST_RESPONSE | jq -r '.id')
echo "Access Request ID: $ACCESS_REQUEST_ID"

echo ""

# Ephemeral Credential API example
echo "8. Generating ephemeral credentials..."
EPH_CRED_RESPONSE=$(curl -s --fail -X POST "$BASE_URL/api/ephemeral-credentials" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": \"$USER_ID\", \"resource_id\": \"$RESOURCE_ID\"}")
echo $EPH_CRED_RESPONSE | jq .
EPH_CRED_ID=$(echo $EPH_CRED_RESPONSE | jq -r '.id')
echo "Ephemeral Credential ID: $EPH_CRED_ID"

echo ""

# Health check
echo "9. Health check..."
curl -s --fail -X GET "$BASE_URL/health"

echo ""
echo "API testing complete!"
