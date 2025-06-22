# Security Guide for Secretary

## Overview

Secretary is a security-focused privileged access management (PAM) application similar to HashiCorp Boundary and Teleport. This document outlines the security measures implemented and best practices for secure deployment.

## Security Improvements Implemented

### 1. **SQL Injection Prevention**
- ✅ All database queries use parameterized statements
- ✅ Input validation prevents malicious SQL patterns
- ✅ Database layer updated to support secure parameterized queries

### 2. **Authentication & Session Security**
- ✅ Secure session management with configurable expiration
- ✅ Session secrets generated from cryptographically secure random sources
- ✅ HttpOnly, Secure, and SameSite cookie flags
- ✅ Session rotation and proper cleanup
- ✅ Configurable bcrypt cost for password hashing (default: 12)

### 3. **CSRF Protection**
- ✅ CSRF tokens for all state-changing operations
- ✅ Token-based validation with expiration
- ✅ Automatic token cleanup to prevent memory leaks

### 4. **Rate Limiting**
- ✅ Per-IP rate limiting to prevent brute force attacks
- ✅ Configurable request limits and time windows
- ✅ Automatic cleanup of expired rate limit data

### 5. **Input Validation**
- ✅ Comprehensive input validation for all user inputs
- ✅ Username, password, hostname, and other field validation
- ✅ XSS and injection pattern detection
- ✅ UTF-8 validation and sanitization

### 6. **TLS/HTTPS Configuration**
- ✅ Modern TLS 1.2+ with secure cipher suites
- ✅ Configurable certificate paths
- ✅ HSTS headers when TLS is enabled
- ✅ HTTP/2 disabled for additional security

### 7. **Security Headers**
- ✅ Content Security Policy (CSP)
- ✅ X-Frame-Options: DENY
- ✅ X-Content-Type-Options: nosniff
- ✅ X-XSS-Protection
- ✅ Referrer-Policy: strict-origin-when-cross-origin

### 8. **Container Security**
- ✅ Non-root user execution
- ✅ Minimal Alpine Linux base image
- ✅ Security updates applied
- ✅ Non-privileged port (8080)
- ✅ Health checks implemented

### 9. **Audit & Logging**
- ✅ Comprehensive security event logging
- ✅ Failed authentication attempts logged
- ✅ Rate limit violations logged
- ✅ CSRF attacks logged
- ✅ Session anomalies logged

### 10. **Configuration Security**
- ✅ Environment variable-based configuration
- ✅ No hardcoded secrets
- ✅ Secure defaults
- ✅ Configuration validation

## Deployment Security Checklist

### Production Environment

#### 1. **TLS Configuration (CRITICAL)**
```bash
# Generate certificates (replace with your domain)
openssl req -x509 -newkey rsa:4096 -keyout secretary-key.pem -out secretary-cert.pem -days 365 -nodes -subj "/CN=your-domain.com"

# Set environment variables
export SECRETARY_TLS_CERT_PATH=/etc/ssl/certs/secretary-cert.pem
export SECRETARY_TLS_KEY_PATH=/etc/ssl/private/secretary-key.pem
export SECRETARY_SECURE_COOKIES=true
```

#### 2. **Generate Secure Secrets**
```bash
# Generate session secret
export SECRETARY_SESSION_SECRET=$(openssl rand -hex 32)

# Generate CSRF secret
export SECRETARY_CSRF_SECRET=$(openssl rand -hex 32)
```

#### 3. **Security Configuration**
```bash
export SECRETARY_BCRYPT_COST=14
export SECRETARY_RATE_LIMIT_REQUESTS=50
export SECRETARY_RATE_LIMIT_WINDOW=1h
export SECRETARY_SESSION_MAX_AGE=1800  # 30 minutes
```

#### 4. **Database Security**
```bash
# Ensure database file has proper permissions
chmod 600 /path/to/secretary.db
chown secretary:secretary /path/to/secretary.db
```

#### 5. **Network Security**
- Use a reverse proxy (nginx/Apache) with additional security headers
- Implement IP whitelisting if possible
- Use VPN or private networks for access
- Enable firewall rules to restrict access

### Docker Security

```bash
# Build with security
docker build -t secretary:secure .

# Run with security constraints
docker run -d \
  --name secretary \
  --user 1000:1000 \
  --read-only \
  --tmpfs /tmp \
  --tmpfs /app/data \
  --security-opt no-new-privileges \
  --cap-drop ALL \
  -e SECRETARY_TLS_CERT_PATH=/certs/cert.pem \
  -e SECRETARY_TLS_KEY_PATH=/certs/key.pem \
  -e SECRETARY_SESSION_SECRET=your-secure-session-secret \
  -e SECRETARY_CSRF_SECRET=your-secure-csrf-secret \
  -e SECRETARY_SECURE_COOKIES=true \
  -v /path/to/certs:/certs:ro \
  -v /path/to/data:/app/data \
  -p 443:8080 \
  secretary:secure
```

## Security Monitoring

### Key Metrics to Monitor
1. Failed authentication attempts
2. Rate limit violations
3. CSRF token validation failures
4. Session anomalies
5. Database connection errors
6. TLS handshake failures

### Log Analysis
```bash
# Monitor failed logins
grep "unauthorized" /var/log/secretary/audit.log

# Monitor rate limiting
grep "rate_limit_exceeded" /var/log/secretary/audit.log

# Monitor CSRF attacks
grep "csrf_invalid\|csrf_missing" /var/log/secretary/audit.log
```

## Vulnerability Reporting

If you discover a security vulnerability, please:

1. **DO NOT** create a public GitHub issue
2. Email security concerns to: security@gosecretary.com
3. Include detailed reproduction steps
4. Allow reasonable time for response before disclosure

## Regular Security Tasks

### Weekly
- [ ] Review audit logs for anomalies
- [ ] Check for failed authentication patterns
- [ ] Verify TLS certificate expiration dates

### Monthly
- [ ] Update dependencies (`go mod tidy && go mod verify`)
- [ ] Review rate limiting effectiveness
- [ ] Rotate session secrets in production
- [ ] Update base Docker images

### Quarterly
- [ ] Security penetration testing
- [ ] Review and update security policies
- [ ] Audit user permissions and access
- [ ] Review and rotate TLS certificates

## Security Best Practices

### For Administrators
1. **Always use HTTPS in production**
2. **Rotate secrets regularly**
3. **Monitor audit logs continuously**
4. **Keep the application updated**
5. **Use strong passwords and enforce password policies**
6. **Implement network segmentation**
7. **Regular security assessments**

### For Developers
1. **Never commit secrets to version control**
2. **Always validate user input**
3. **Use parameterized queries only**
4. **Follow the principle of least privilege**
5. **Audit third-party dependencies**
6. **Write security tests**

## Compliance Considerations

This implementation addresses common security frameworks:

- **OWASP Top 10**: SQL Injection, XSS, CSRF, Security Misconfiguration
- **NIST Cybersecurity Framework**: Authentication, Access Control, Logging
- **SOC 2**: Access Controls, Change Management, Monitoring
- **ISO 27001**: Information Security Management

## Known Limitations

1. SQLite is used by default (consider PostgreSQL for production)
2. In-memory rate limiting (consider Redis for distributed deployments)
3. File-based audit logs (consider centralized logging solutions)
4. Basic RBAC implementation (consider more sophisticated authorization)

## Future Security Enhancements

- [ ] Multi-factor authentication (MFA)
- [ ] Hardware security module (HSM) integration
- [ ] Certificate-based authentication
- [ ] Integration with LDAP/Active Directory
- [ ] Advanced threat detection
- [ ] Zero-trust network access features 
