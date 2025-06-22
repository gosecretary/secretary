# Secretary Security Scanning Pipeline

## Overview

The Secretary project includes a comprehensive security scanning pipeline that automatically analyzes code for security vulnerabilities, compliance issues, and best practice violations. This pipeline is designed to ensure that all code meets the project's security standards and follows industry best practices.

## Pipeline Components

### 1. Go Security Analysis
- **Tools**: `govulncheck`, `gosec`, `gocyclo`, `golangci-lint`
- **Purpose**: Static analysis of Go code for security vulnerabilities
- **Checks**:
  - Known vulnerabilities in dependencies
  - Security anti-patterns (hardcoded secrets, SQL injection, etc.)
  - Code complexity analysis
  - Comprehensive linting with security focus

### 2. Dependency Security Scanning
- **Tools**: `Trivy`, `go list`
- **Purpose**: Vulnerability scanning of dependencies
- **Checks**:
  - Known vulnerabilities in Go modules
  - Outdated dependencies
  - License compliance issues

### 3. Snyk Security Analysis
- **Tools**: Snyk CLI
- **Purpose**: Multi-layered security analysis
- **Checks**:
  - **SAST**: Static Application Security Testing
  - **SCA**: Software Composition Analysis
  - **Container**: Container image security
  - **IaC**: Infrastructure as Code security

### 4. Container Security
- **Tools**: Trivy, Snyk Container
- **Purpose**: Security analysis of Docker images
- **Checks**:
  - Base image vulnerabilities
  - Package vulnerabilities
  - Configuration security
  - Best practices compliance

### 5. CodeQL Analysis
- **Tools**: GitHub CodeQL
- **Purpose**: Advanced static analysis
- **Checks**:
  - Security vulnerabilities
  - Code quality issues
  - Custom security queries

### 6. License Compliance
- **Tools**: `go-licenses`
- **Purpose**: License compliance checking
- **Checks**:
  - License compatibility
  - License documentation
  - Compliance reporting

### 7. Security Policy Compliance
- **Tools**: Custom scripts
- **Purpose**: Policy enforcement
- **Checks**:
  - Security documentation presence
  - Security headers implementation
  - Hardcoded secrets detection
  - TLS configuration validation

### 8. Security Test Execution
- **Tools**: Go test
- **Purpose**: Security-focused testing
- **Checks**:
  - Security component tests
  - Test coverage for security code
  - Security validation tests

## Configuration

### Environment Variables

The pipeline requires the following environment variables:

```bash
# Required for Snyk
SNYK_TOKEN=your-snyk-api-token

# Optional for enhanced reporting
GITHUB_TOKEN=${{ secrets.GITHUB_TOKEN }}
```

### Security Configuration

The pipeline uses `.security/config.yml` for configuration:

```yaml
policies:
  severity_thresholds:
    critical: 0
    high: 0
    medium: 5
    low: 20
```

## Usage

### Automatic Execution

The pipeline runs automatically on:
- **Push to master/main**: Full security scan
- **Pull requests**: Security validation
- **Weekly schedule**: Comprehensive security audit

### Manual Execution

To run the security pipeline manually:

```bash
# Run all security checks locally
./scripts/security-scan.sh

# Run specific security tools
gosec ./...
govulncheck ./...
trivy fs .
```

### Local Development

For local development, install the required tools:

```bash
# Install Go security tools
go install golang.org/x/vuln/cmd/govulncheck@latest
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install Trivy
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin

# Install Snyk CLI
npm install -g snyk
```

## Security Standards

### Vulnerability Severity Levels

1. **Critical (0 allowed)**
   - Remote code execution
   - Authentication bypass
   - Data exfiltration
   - Must be fixed immediately

2. **High (0 allowed)**
   - SQL injection
   - XSS vulnerabilities
   - Privilege escalation
   - Must be fixed within 24 hours

3. **Medium (≤5 allowed)**
   - Information disclosure
   - Denial of service
   - Should be fixed within 7 days

4. **Low (≤20 allowed)**
   - Best practice violations
   - Code quality issues
   - Should be addressed within 30 days

### Code Quality Standards

- **Test Coverage**: Minimum 80% overall, 90% for security components
- **Cyclomatic Complexity**: Maximum 15 per function
- **Function Length**: Maximum 50 lines
- **File Length**: Maximum 500 lines

### Security Requirements

1. **Authentication & Authorization**
   - Secure password policies
   - Session management
   - Role-based access control

2. **Input Validation**
   - All user inputs validated
   - SQL injection protection
   - XSS protection

3. **Encryption & Security**
   - TLS 1.2+ required
   - Secure hashing (bcrypt)
   - No hardcoded secrets

4. **Security Headers**
   - Content Security Policy
   - X-Frame-Options
   - X-Content-Type-Options
   - X-XSS-Protection

## Compliance Standards

The pipeline ensures compliance with:

- **OWASP Top 10**: Web application security
- **NIST Cybersecurity Framework**: Security controls
- **SOC 2**: Security, availability, processing integrity
- **ISO 27001**: Information security management
- **PCI DSS**: Payment card industry security

## Reporting

### Security Reports

The pipeline generates comprehensive reports:

1. **SARIF Files**: Machine-readable security results
2. **GitHub Security Tab**: Integrated security findings
3. **PR Comments**: Automated security summaries
4. **Artifacts**: Detailed security reports

### Report Contents

- Vulnerability summary by severity
- Remediation recommendations
- Compliance status
- Test coverage metrics
- Security policy compliance

## Remediation Process

### Critical & High Severity

1. **Immediate Action Required**
   - Block merge until fixed
   - Create security issue
   - Assign to security team

2. **Remediation Steps**
   - Identify root cause
   - Implement fix
   - Add security tests
   - Update documentation

### Medium & Low Severity

1. **Scheduled Remediation**
   - Track in project management
   - Plan for next sprint
   - Monitor progress

2. **Documentation**
   - Update security documentation
   - Add prevention measures
   - Share lessons learned

## Best Practices

### Development Workflow

1. **Pre-commit Checks**
   ```bash
   # Run security checks before committing
   make security-check
   ```

2. **Security-First Development**
   - Write security tests first
   - Validate all inputs
   - Use secure defaults

3. **Regular Security Reviews**
   - Weekly security audits
   - Monthly dependency updates
   - Quarterly security assessments

### Security Testing

1. **Unit Tests**
   ```go
   func TestPasswordValidation(t *testing.T) {
       // Test password complexity requirements
   }
   ```

2. **Integration Tests**
   ```go
   func TestAuthenticationFlow(t *testing.T) {
       // Test complete authentication flow
   }
   ```

3. **Security Tests**
   ```go
   func TestSQLInjectionProtection(t *testing.T) {
       // Test SQL injection prevention
   }
   ```

## Troubleshooting

### Common Issues

1. **False Positives**
   - Review security configuration
   - Add exclusions if needed
   - Document false positive reasons

2. **Tool Failures**
   - Check tool versions
   - Verify configuration
   - Review error logs

3. **Performance Issues**
   - Optimize scan scope
   - Use caching
   - Parallel execution

### Getting Help

1. **Documentation**: Check this guide and project docs
2. **Issues**: Create GitHub issue with security label
3. **Security Team**: Contact for critical issues
4. **Community**: Discuss in project discussions

## Continuous Improvement

### Metrics Tracking

- Vulnerability count by severity
- Remediation time
- Test coverage trends
- Security compliance score

### Regular Reviews

- Monthly pipeline effectiveness review
- Quarterly security standard updates
- Annual compliance assessment

### Feedback Loop

- Developer feedback on false positives
- Security team input on new threats
- Compliance team guidance on requirements

---

## Quick Reference

### Commands

```bash
# Run all security checks
./scripts/security-scan.sh

# Run specific tools
gosec ./...
govulncheck ./...
trivy fs .
snyk test

# Check security configuration
golangci-lint run
```

### Configuration Files

- `.security/config.yml`: Security policies
- `.golangci-lint.yml`: Linting configuration
- `.github/workflows/security-scanning.yml`: Pipeline definition

### Important URLs

- **GitHub Security Tab**: Project security findings
- **Snyk Dashboard**: Dependency vulnerabilities
- **Trivy Reports**: Container security analysis
- **CodeQL Results**: Advanced static analysis

---

*This security pipeline is designed to ensure that Secretary maintains the highest security standards while providing comprehensive coverage of potential vulnerabilities and compliance requirements.* 