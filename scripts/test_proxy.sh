#!/bin/bash

# Secretary TCP Proxy Test Script
# This script demonstrates the proxy functionality by creating a session,
# setting up a proxy, and testing SSH connections through Secretary.

set -e

# Configuration
SECRETARY_URL="http://localhost:8080"
ADMIN_USERNAME="admin"
ADMIN_PASSWORD=""
SESSION_TOKEN=""
PROXY_ID=""
LOCAL_PORT=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Secretary is running
check_secretary_running() {
    log_info "Checking if Secretary is running..."
    if curl -s "$SECRETARY_URL/health" > /dev/null; then
        log_success "Secretary is running"
    else
        log_error "Secretary is not running. Please start it first."
        exit 1
    fi
}

# Get admin password from user
get_admin_password() {
    if [ -z "$ADMIN_PASSWORD" ]; then
        echo -n "Enter admin password: "
        read -s ADMIN_PASSWORD
        echo
    fi
}

# Login and get session token
login() {
    log_info "Logging in as admin..."
    get_admin_password
    
    local login_response=$(curl -s -X POST "$SECRETARY_URL/api/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$ADMIN_USERNAME\",
            \"password\": \"$ADMIN_PASSWORD\"
        }")
    
    if echo "$login_response" | grep -q '"success":true'; then
        SESSION_TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
        log_success "Login successful"
        log_info "Session token: ${SESSION_TOKEN:0:20}..."
    else
        log_error "Login failed: $login_response"
        exit 1
    fi
}

# Create a test user
create_test_user() {
    log_info "Creating test user..."
    local user_response=$(curl -s -X POST "$SECRETARY_URL/api/users" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $SESSION_TOKEN" \
        -d '{
            "username": "testuser",
            "email": "testuser@example.com",
            "password": "TestPass123!",
            "name": "Test User",
            "role": "user"
        }')
    
    if echo "$user_response" | grep -q '"success":true'; then
        log_success "Test user created"
    else
        log_warning "Test user creation failed (may already exist): $user_response"
    fi
}

# Create a test resource
create_test_resource() {
    log_info "Creating test resource..."
    local resource_response=$(curl -s -X POST "$SECRETARY_URL/api/resources" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $SESSION_TOKEN" \
        -d '{
            "name": "test-server",
            "description": "Test SSH server for proxy testing",
            "type": "ssh",
            "host": "localhost",
            "port": 2222
        }')
    
    if echo "$resource_response" | grep -q '"success":true'; then
        log_success "Test resource created"
    else
        log_warning "Test resource creation failed (may already exist): $resource_response"
    fi
}

# Get user and resource IDs
get_ids() {
    log_info "Getting user and resource IDs..."
    
    # Get user ID
    local users_response=$(curl -s -X GET "$SECRETARY_URL/api/users" \
        -H "Authorization: Bearer $SESSION_TOKEN")
    USER_ID=$(echo "$users_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    
    # Get resource ID
    local resources_response=$(curl -s -X GET "$SECRETARY_URL/api/resources" \
        -H "Authorization: Bearer $SESSION_TOKEN")
    RESOURCE_ID=$(echo "$resources_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    
    log_info "User ID: $USER_ID"
    log_info "Resource ID: $RESOURCE_ID"
}

# Create a session
create_session() {
    log_info "Creating session..."
    local session_response=$(curl -s -X POST "$SECRETARY_URL/api/sessions" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $SESSION_TOKEN" \
        -d "{
            \"user_id\": \"$USER_ID\",
            \"resource_id\": \"$RESOURCE_ID\"
        }")
    
    if echo "$session_response" | grep -q '"success":true'; then
        SESSION_ID=$(echo "$session_response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        log_success "Session created: $SESSION_ID"
    else
        log_error "Session creation failed: $session_response"
        exit 1
    fi
}

# Create proxy
create_proxy() {
    log_info "Creating proxy..."
    local proxy_response=$(curl -s -X POST "$SECRETARY_URL/api/sessions/$SESSION_ID/proxy" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $SESSION_TOKEN" \
        -d '{
            "protocol": "ssh",
            "remote_host": "localhost",
            "remote_port": 2222
        }')
    
    if echo "$proxy_response" | grep -q '"success":true'; then
        PROXY_ID=$(echo "$proxy_response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        LOCAL_PORT=$(echo "$proxy_response" | grep -o '"local_port":[0-9]*' | cut -d':' -f2)
        log_success "Proxy created: $PROXY_ID"
        log_info "Local port: $LOCAL_PORT"
    else
        log_error "Proxy creation failed: $proxy_response"
        exit 1
    fi
}

# Start proxy
start_proxy() {
    log_info "Starting proxy..."
    local start_response=$(curl -s -X POST "$SECRETARY_URL/api/proxies/$PROXY_ID/start" \
        -H "Authorization: Bearer $SESSION_TOKEN")
    
    if echo "$start_response" | grep -q '"success":true'; then
        log_success "Proxy started successfully"
    else
        log_error "Proxy start failed: $start_response"
        exit 1
    fi
}

# Test SSH connection through proxy
test_ssh_connection() {
    log_info "Testing SSH connection through proxy..."
    
    # Check if SSH server is running on port 2222
    if ! nc -z localhost 2222 2>/dev/null; then
        log_warning "SSH server not running on port 2222. Starting test SSH server..."
        start_test_ssh_server
    fi
    
    # Test connection through proxy
    if nc -z localhost "$LOCAL_PORT" 2>/dev/null; then
        log_success "Proxy is listening on port $LOCAL_PORT"
        log_info "You can now connect via: ssh -p $LOCAL_PORT testuser@localhost"
    else
        log_error "Proxy is not listening on port $LOCAL_PORT"
    fi
}

# Start test SSH server
start_test_ssh_server() {
    log_info "Starting test SSH server..."
    
    # Check if Docker is available
    if command -v docker >/dev/null 2>&1; then
        log_info "Using Docker to start SSH server..."
        docker run -d --name secretary-test-ssh \
            -p 2222:2222 \
            -e USER_NAME=testuser \
            -e USER_PASSWORD=testpass123 \
            -e PASSWORD_ACCESS=true \
            -e SUDO_ACCESS=true \
            linuxserver/openssh-server:latest
        
        # Wait for SSH server to start
        log_info "Waiting for SSH server to start..."
        sleep 10
        
        if nc -z localhost 2222 2>/dev/null; then
            log_success "SSH server started on port 2222"
        else
            log_error "Failed to start SSH server"
        fi
    else
        log_warning "Docker not available. Please start an SSH server on port 2222 manually."
    fi
}

# Monitor commands
monitor_commands() {
    log_info "Monitoring commands for session..."
    local commands_response=$(curl -s -X GET "$SECRETARY_URL/api/sessions/$SESSION_ID/commands" \
        -H "Authorization: Bearer $SESSION_TOKEN")
    
    echo "$commands_response" | jq '.' 2>/dev/null || echo "$commands_response"
}

# Monitor alerts
monitor_alerts() {
    log_info "Monitoring alerts for session..."
    local alerts_response=$(curl -s -X GET "$SECRETARY_URL/api/sessions/$SESSION_ID/alerts" \
        -H "Authorization: Bearer $SESSION_TOKEN")
    
    echo "$alerts_response" | jq '.' 2>/dev/null || echo "$alerts_response"
}

# Get active proxies
get_active_proxies() {
    log_info "Getting active proxies..."
    local proxies_response=$(curl -s -X GET "$SECRETARY_URL/api/proxies/active" \
        -H "Authorization: Bearer $SESSION_TOKEN")
    
    echo "$proxies_response" | jq '.' 2>/dev/null || echo "$proxies_response"
}

# Cleanup
cleanup() {
    log_info "Cleaning up..."
    
    # Stop proxy
    if [ -n "$PROXY_ID" ]; then
        curl -s -X POST "$SECRETARY_URL/api/proxies/$PROXY_ID/stop" \
            -H "Authorization: Bearer $SESSION_TOKEN" > /dev/null
        log_info "Proxy stopped"
    fi
    
    # Stop SSH server
    if docker ps | grep -q secretary-test-ssh; then
        docker stop secretary-test-ssh > /dev/null
        docker rm secretary-test-ssh > /dev/null
        log_info "Test SSH server stopped"
    fi
}

# Main test function
run_proxy_test() {
    log_info "Starting Secretary TCP Proxy Test"
    log_info "=================================="
    
    # Check prerequisites
    check_secretary_running
    
    # Run test steps
    login
    create_test_user
    create_test_resource
    get_ids
    create_session
    create_proxy
    start_proxy
    test_ssh_connection
    
    log_success "Proxy test completed successfully!"
    log_info "Proxy is running on port $LOCAL_PORT"
    log_info "Connect via: ssh -p $LOCAL_PORT testuser@localhost"
    log_info "Password: testpass123"
    
    # Show monitoring options
    echo
    log_info "Monitoring Options:"
    echo "1. Monitor commands: $0 --monitor-commands"
    echo "2. Monitor alerts: $0 --monitor-alerts"
    echo "3. Get active proxies: $0 --active-proxies"
    echo "4. Cleanup: $0 --cleanup"
}

# Handle command line arguments
case "${1:-}" in
    --monitor-commands)
        login
        get_ids
        create_session
        monitor_commands
        ;;
    --monitor-alerts)
        login
        get_ids
        create_session
        monitor_alerts
        ;;
    --active-proxies)
        login
        get_active_proxies
        ;;
    --cleanup)
        login
        cleanup
        ;;
    --help|-h)
        echo "Usage: $0 [OPTION]"
        echo
        echo "Options:"
        echo "  --monitor-commands  Monitor commands for current session"
        echo "  --monitor-alerts    Monitor alerts for current session"
        echo "  --active-proxies    Get list of active proxies"
        echo "  --cleanup          Clean up test resources"
        echo "  --help, -h         Show this help message"
        echo
        echo "Default: Run full proxy test"
        ;;
    *)
        run_proxy_test
        ;;
esac 