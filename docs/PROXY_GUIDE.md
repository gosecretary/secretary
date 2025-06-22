# Secretary TCP Proxy Guide

## Overview

Secretary provides a comprehensive TCP proxy solution that allows clients to connect to target servers through Secretary while providing real-time monitoring, command interception, and security auditing. This proxy acts as a secure gateway that monitors all traffic and can block malicious commands.

## Architecture

```
Client → Secretary Proxy → Target Server
         ↓
    Monitoring & Auditing
         ↓
    Security Analysis
         ↓
    Command Recording
```

## Supported Protocols

Secretary's proxy supports multiple protocols with specialized monitoring:

- **SSH**: Full command interception and analysis
- **MySQL**: SQL query monitoring and risk assessment
- **PostgreSQL**: SQL query monitoring and risk assessment
- **Generic TCP**: Basic traffic monitoring for other protocols

## Proxy Workflow

### 1. Session Creation
First, create a session for the user accessing a resource:

```bash
curl -X POST http://localhost:8080/api/sessions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "user_id": "user123",
    "resource_id": "resource456"
  }'
```

### 2. Create Proxy
Create a proxy connection for the session:

```bash
curl -X POST http://localhost:8080/api/sessions/{session_id}/proxy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "protocol": "ssh",
    "remote_host": "target-server.example.com",
    "remote_port": 22
  }'
```

Response:
```json
{
  "success": true,
  "message": "Proxy created successfully",
  "data": {
    "id": "proxy-uuid",
    "session_id": "session-uuid",
    "protocol": "ssh",
    "local_port": 10001,
    "remote_host": "target-server.example.com",
    "remote_port": 22,
    "status": "created"
  }
}
```

### 3. Start Proxy
Start the proxy to begin accepting connections:

```bash
curl -X POST http://localhost:8080/api/proxies/{proxy_id}/start \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Response:
```json
{
  "success": true,
  "message": "Proxy started successfully",
  "data": {
    "proxy_id": "proxy-uuid",
    "local_port": 10001,
    "status": "started"
  }
}
```

### 4. Connect Through Proxy
Now clients can connect to the target server through Secretary:

```bash
# SSH through proxy
ssh -p 10001 user@localhost

# MySQL through proxy
mysql -h localhost -P 10001 -u username -p

# PostgreSQL through proxy
psql -h localhost -p 10001 -U username -d database
```

## Security Features

### Command Analysis
Secretary analyzes all commands in real-time and categorizes them by risk level:

- **Low Risk**: Basic commands like `ls`, `pwd`
- **Medium Risk**: Commands that access system information
- **High Risk**: Administrative commands like `sudo`, `rm -rf`
- **Critical Risk**: Dangerous commands that are automatically blocked

### Automatic Blocking
Critical commands are automatically blocked:

**SSH Commands:**
- `rm -rf /` (filesystem destruction)
- `mkfs.*` (filesystem formatting)
- `dd if=* of=/dev/` (direct disk access)
- Fork bombs and other dangerous patterns

**SQL Commands:**
- `DROP DATABASE`
- `DROP SCHEMA`
- `TRUNCATE`
- Mass deletion patterns

### Security Alerts
When high-risk commands are detected, Secretary creates security alerts:

```json
{
  "id": "alert-uuid",
  "session_id": "session-uuid",
  "command_id": "command-uuid",
  "alert_type": "blocked_command",
  "severity": "critical",
  "title": "Blocked High-Risk Command",
  "description": "Command blocked due to critical risk level",
  "raw_data": "rm -rf /",
  "action": "blocked"
}
```

## Monitoring and Auditing

### Real-time Command Monitoring
Monitor commands in real-time:

```bash
# Get all commands for a session
curl -X GET http://localhost:8080/api/sessions/{session_id}/commands \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get high-risk commands
curl -X GET http://localhost:8080/api/commands/high-risk \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Session Recording
All sessions are automatically recorded:

```bash
# Start recording
curl -X POST http://localhost:8080/api/sessions/{session_id}/recording/start \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get recording info
curl -X GET http://localhost:8080/api/sessions/{session_id}/recording \
  -H "Authorization: Bearer YOUR_TOKEN"

# Download recording
curl -X GET http://localhost:8080/api/recordings/{recording_id}/download \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Security Alerts
Monitor security alerts:

```bash
# Get alerts for a session
curl -X GET http://localhost:8080/api/sessions/{session_id}/alerts \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get alerts by severity
curl -X GET http://localhost:8080/api/alerts/severity/critical \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Configuration

### Proxy Settings
Configure proxy behavior in your environment:

```bash
# Proxy port range (default: 10000-20000)
export SECRETARY_PROXY_PORT_MIN=10000
export SECRETARY_PROXY_PORT_MAX=20000

# Command analysis settings
export SECRETARY_BLOCK_CRITICAL_COMMANDS=true
export SECRETARY_LOG_ALL_COMMANDS=true
```

### Security Policies
Customize security policies:

```bash
# Enable/disable automatic blocking
export SECRETARY_AUTO_BLOCK=true

# Set risk thresholds
export SECRETARY_BLOCK_HIGH_RISK=false
export SECRETARY_BLOCK_CRITICAL_RISK=true
```

## Use Cases

### 1. SSH Access Control
```bash
# Create SSH proxy
curl -X POST http://localhost:8080/api/sessions/{session_id}/proxy \
  -d '{"protocol": "ssh", "remote_host": "prod-server", "remote_port": 22}'

# Connect through proxy
ssh -p 10001 user@localhost
```

### 2. Database Access Monitoring
```bash
# Create MySQL proxy
curl -X POST http://localhost:8080/api/sessions/{session_id}/proxy \
  -d '{"protocol": "mysql", "remote_host": "db-server", "remote_port": 3306}'

# Connect through proxy
mysql -h localhost -P 10001 -u username -p
```

### 3. PostgreSQL Access Control
```bash
# Create PostgreSQL proxy
curl -X POST http://localhost:8080/api/sessions/{session_id}/proxy \
  -d '{"protocol": "postgresql", "remote_host": "pg-server", "remote_port": 5432}'

# Connect through proxy
psql -h localhost -p 10001 -U username -d database
```

## Best Practices

### 1. Session Management
- Always create sessions before using proxies
- Terminate sessions when access is complete
- Monitor active sessions regularly

### 2. Security Monitoring
- Review high-risk commands daily
- Monitor security alerts in real-time
- Regularly audit session recordings

### 3. Network Security
- Use TLS for all API communications
- Restrict proxy access to authorized networks
- Implement network segmentation

### 4. Access Control
- Use ephemeral credentials when possible
- Implement time-based access restrictions
- Regular access reviews and cleanup

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   - Secretary automatically finds available ports
   - Check if another proxy is using the same port

2. **Connection Refused**
   - Verify target server is accessible
   - Check firewall rules
   - Ensure target service is running

3. **Command Blocking**
   - Review security policies
   - Check command risk analysis
   - Adjust blocking rules if needed

### Debug Commands

```bash
# Check active proxies
curl -X GET http://localhost:8080/api/proxies/active \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get proxy status
curl -X GET http://localhost:8080/api/sessions/{session_id}/proxy \
  -H "Authorization: Bearer YOUR_TOKEN"

# Stop proxy
curl -X POST http://localhost:8080/api/proxies/{proxy_id}/stop \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Integration Examples

### Python Client
```python
import requests
import paramiko

# Create session
session_response = requests.post(
    "http://localhost:8080/api/sessions",
    headers={"Authorization": "Bearer YOUR_TOKEN"},
    json={"user_id": "user123", "resource_id": "resource456"}
)
session_id = session_response.json()["data"]["id"]

# Create proxy
proxy_response = requests.post(
    f"http://localhost:8080/api/sessions/{session_id}/proxy",
    headers={"Authorization": "Bearer YOUR_TOKEN"},
    json={"protocol": "ssh", "remote_host": "target-server", "remote_port": 22}
)
proxy_data = proxy_response.json()["data"]

# Start proxy
requests.post(
    f"http://localhost:8080/api/proxies/{proxy_data['id']}/start",
    headers={"Authorization": "Bearer YOUR_TOKEN"}
)

# Connect through proxy
ssh = paramiko.SSHClient()
ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
ssh.connect('localhost', port=proxy_data['local_port'], username='user')
```

### Go Client
```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "golang.org/x/crypto/ssh"
)

type ProxyRequest struct {
    Protocol   string `json:"protocol"`
    RemoteHost string `json:"remote_host"`
    RemotePort int    `json:"remote_port"`
}

func main() {
    // Create session
    sessionReq := map[string]string{
        "user_id":     "user123",
        "resource_id": "resource456",
    }
    sessionJSON, _ := json.Marshal(sessionReq)
    
    resp, _ := http.Post(
        "http://localhost:8080/api/sessions",
        "application/json",
        bytes.NewBuffer(sessionJSON),
    )
    
    var sessionResp map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&sessionResp)
    sessionID := sessionResp["data"].(map[string]interface{})["id"].(string)
    
    // Create proxy
    proxyReq := ProxyRequest{
        Protocol:   "ssh",
        RemoteHost: "target-server",
        RemotePort: 22,
    }
    proxyJSON, _ := json.Marshal(proxyReq)
    
    resp, _ = http.Post(
        fmt.Sprintf("http://localhost:8080/api/sessions/%s/proxy", sessionID),
        "application/json",
        bytes.NewBuffer(proxyJSON),
    )
    
    var proxyResp map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&proxyResp)
    proxyData := proxyResp["data"].(map[string]interface{})
    
    // Connect through proxy
    config := &ssh.ClientConfig{
        User: "user",
        Auth: []ssh.AuthMethod{
            ssh.Password("password"),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    
    client, _ := ssh.Dial("tcp", fmt.Sprintf("localhost:%v", proxyData["local_port"]), config)
    defer client.Close()
}
```

## Security Considerations

1. **Network Security**: Always use TLS for API communications
2. **Access Control**: Implement proper authentication and authorization
3. **Monitoring**: Regularly review logs and alerts
4. **Updates**: Keep Secretary updated with latest security patches
5. **Backup**: Regularly backup session recordings and audit logs

## Compliance

The proxy functionality supports various compliance requirements:

- **SOC 2**: Access controls and monitoring
- **PCI DSS**: Secure access to cardholder data
- **HIPAA**: Protected health information access
- **SOX**: Financial data access controls
- **ISO 27001**: Information security management

## Support

For issues or questions about the proxy functionality:

1. Check the troubleshooting section
2. Review security alerts and logs
3. Consult the API documentation
4. Contact the development team

---

*This guide covers the TCP proxy functionality in Secretary. For additional information, refer to the main documentation and API specifications.* 