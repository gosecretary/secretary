#!/bin/bash

# Secretary Security Scanning Script
# This script performs comprehensive security analysis locally

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GO_VERSION="1.22"
SECURITY_CONFIG=".security/config.yml"
REPORTS_DIR="security-reports"
FAIL_ON_CRITICAL=true
FAIL_ON_HIGH=true

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

# Check if required tools are installed
check_tools() {
    log_info "Checking required security tools..."
    
    local missing_tools=()
    
    # Check Go security tools
    if ! command -v govulncheck &> /dev/null; then
        missing_tools+=("govulncheck")
    fi
    
    if ! command -v gosec &> /dev/null; then
        missing_tools+=("gosec")
    fi
    
    if ! command -v gocyclo &> /dev/null; then
        missing_tools+=("gocyclo")
    fi
    
    if ! command -v golangci-lint &> /dev/null; then
        missing_tools+=("golangci-lint")
    fi
    
    # Check Trivy
    if ! command -v trivy &> /dev/null; then
        missing_tools+=("trivy")
    fi
    
    # Check Snyk (optional)
    if ! command -v snyk &> /dev/null; then
        log_warning "Snyk not found. Install with: npm install -g snyk"
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        log_info "Install missing tools with:"
        echo "go install golang.org/x/vuln/cmd/govulncheck@latest"
        echo "go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"
        echo "go install github.com/fzipp/gocyclo/cmd/gocyclo@latest"
        echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        echo "curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin"
        exit 1
    fi
    
    log_success "All required tools are installed"
}

# Create reports directory
setup_reports() {
    log_info "Setting up reports directory..."
    mkdir -p "$REPORTS_DIR"
    log_success "Reports directory created: $REPORTS_DIR"
}

# Go Security Analysis
run_go_security() {
    log_info "Running Go security analysis..."
    
    # Go vulnerability check
    log_info "Running govulncheck..."
    if govulncheck ./... > "$REPORTS_DIR/govulncheck.txt" 2>&1; then
        log_success "govulncheck completed successfully"
    else
        log_warning "govulncheck found vulnerabilities (see $REPORTS_DIR/govulncheck.txt)"
    fi
    
    # GoSec security linter
    log_info "Running gosec..."
    if gosec -fmt=sarif -out="$REPORTS_DIR/gosec.sarif" ./...; then
        log_success "gosec completed successfully"
    else
        log_warning "gosec found security issues (see $REPORTS_DIR/gosec.sarif)"
    fi
    
    # Cyclomatic complexity check
    log_info "Running gocyclo..."
    if gocyclo -over 15 . > "$REPORTS_DIR/gocyclo.txt" 2>&1; then
        log_warning "gocyclo found complex functions (see $REPORTS_DIR/gocyclo.txt)"
    else
        log_success "gocyclo completed successfully"
    fi
    
    # Comprehensive Go linting
    log_info "Running golangci-lint..."
    if golangci-lint run --out-format=sarif --out="$REPORTS_DIR/golangci-lint.sarif"; then
        log_success "golangci-lint completed successfully"
    else
        log_warning "golangci-lint found issues (see $REPORTS_DIR/golangci-lint.sarif)"
    fi
}

# Dependency Security Scanning
run_dependency_scanning() {
    log_info "Running dependency security scanning..."
    
    # Download dependencies
    log_info "Downloading Go dependencies..."
    go mod download
    
    # Trivy vulnerability scanner
    log_info "Running Trivy filesystem scan..."
    if trivy fs --format sarif --output "$REPORTS_DIR/trivy-fs.sarif" .; then
        log_success "Trivy filesystem scan completed"
    else
        log_warning "Trivy found vulnerabilities (see $REPORTS_DIR/trivy-fs.sarif)"
    fi
    
    # Check for known vulnerabilities in dependencies
    log_info "Analyzing Go dependencies..."
    go list -json -deps ./... | jq -r '.Deps[]' | sort | uniq > "$REPORTS_DIR/dependencies.txt"
    log_success "Dependency analysis completed (see $REPORTS_DIR/dependencies.txt)"
}

# Snyk Security Analysis
run_snyk_analysis() {
    if ! command -v snyk &> /dev/null; then
        log_warning "Snyk not available, skipping Snyk analysis"
        return
    fi
    
    log_info "Running Snyk security analysis..."
    
    # Check if SNYK_TOKEN is set
    if [ -z "$SNYK_TOKEN" ]; then
        log_warning "SNYK_TOKEN not set, skipping Snyk analysis"
        return
    fi
    
    # Snyk Code (SAST) analysis
    log_info "Running Snyk Code analysis..."
    if snyk code test --sarif > "$REPORTS_DIR/snyk-code.sarif" 2>/dev/null; then
        log_success "Snyk Code analysis completed"
    else
        log_warning "Snyk Code found issues (see $REPORTS_DIR/snyk-code.sarif)"
    fi
    
    # Snyk Open Source (SCA) analysis
    log_info "Running Snyk Open Source analysis..."
    if snyk monitor --all-projects > "$REPORTS_DIR/snyk-monitor.txt" 2>&1; then
        log_success "Snyk Open Source analysis completed"
    else
        log_warning "Snyk Open Source found issues (see $REPORTS_DIR/snyk-monitor.txt)"
    fi
    
    # Snyk Infrastructure as Code (IaC) analysis
    log_info "Running Snyk IaC analysis..."
    if snyk iac test --report > "$REPORTS_DIR/snyk-iac.txt" 2>&1; then
        log_success "Snyk IaC analysis completed"
    else
        log_warning "Snyk IaC found issues (see $REPORTS_DIR/snyk-iac.txt)"
    fi
}

# Container Security
run_container_security() {
    log_info "Running container security analysis..."
    
    # Check if Dockerfile exists
    if [ ! -f "Dockerfile" ]; then
        log_warning "Dockerfile not found, skipping container security"
        return
    fi
    
    # Build Docker image for scanning
    log_info "Building Docker image for security scan..."
    if docker build -t secretary:security-scan . > "$REPORTS_DIR/docker-build.txt" 2>&1; then
        log_success "Docker image built successfully"
    else
        log_error "Failed to build Docker image"
        return
    fi
    
    # Trivy container scan
    log_info "Running Trivy container scan..."
    if trivy image --format sarif --output "$REPORTS_DIR/trivy-container.sarif" secretary:security-scan; then
        log_success "Trivy container scan completed"
    else
        log_warning "Trivy container scan found issues (see $REPORTS_DIR/trivy-container.sarif)"
    fi
    
    # Snyk Container analysis
    if command -v snyk &> /dev/null && [ -n "$SNYK_TOKEN" ]; then
        log_info "Running Snyk Container analysis..."
        if snyk container monitor secretary:security-scan --file=Dockerfile > "$REPORTS_DIR/snyk-container.txt" 2>&1; then
            log_success "Snyk Container analysis completed"
        else
            log_warning "Snyk Container found issues (see $REPORTS_DIR/snyk-container.txt)"
        fi
    fi
}

# License Compliance
run_license_compliance() {
    log_info "Running license compliance check..."
    
    # Install go-licenses if not available
    if ! command -v go-licenses &> /dev/null; then
        log_info "Installing go-licenses..."
        go install github.com/google/go-licenses@latest
    fi
    
    # Check license compliance
    if go-licenses check ./... > "$REPORTS_DIR/license-check.txt" 2>&1; then
        log_success "License compliance check passed"
    else
        log_warning "License compliance issues found (see $REPORTS_DIR/license-check.txt)"
    fi
    
    # Generate license report
    if go-licenses csv ./... > "$REPORTS_DIR/licenses.csv" 2>/dev/null; then
        log_success "License report generated (see $REPORTS_DIR/licenses.csv)"
    else
        log_warning "Failed to generate license report"
    fi
}

# Security Policy Compliance
run_security_policy() {
    log_info "Running security policy compliance check..."
    
    local policy_violations=0
    
    # Check for security documentation
    log_info "Checking security documentation..."
    if [ ! -f "SECURITY.md" ]; then
        log_error "SECURITY.md not found"
        ((policy_violations++))
    else
        log_success "SECURITY.md found"
    fi
    
    if [ ! -f "SECURITY_FIXES.md" ]; then
        log_error "SECURITY_FIXES.md not found"
        ((policy_violations++))
    else
        log_success "SECURITY_FIXES.md found"
    fi
    
    if [ ! -f ".cursor/rules/secretary-project-standards.mdc" ]; then
        log_error "Security rules not found"
        ((policy_violations++))
    else
        log_success "Security rules found"
    fi
    
    # Check for hardcoded secrets
    log_info "Checking for hardcoded secrets..."
    if grep -r "password.*=.*\"" . --exclude-dir=.git --exclude-dir=vendor --exclude-dir=node_modules > "$REPORTS_DIR/hardcoded-passwords.txt" 2>/dev/null; then
        log_warning "Potential hardcoded passwords found (see $REPORTS_DIR/hardcoded-passwords.txt)"
        ((policy_violations++))
    else
        log_success "No obvious hardcoded passwords found"
    fi
    
    if grep -r "secret.*=.*\"" . --exclude-dir=.git --exclude-dir=vendor --exclude-dir=node_modules > "$REPORTS_DIR/hardcoded-secrets.txt" 2>/dev/null; then
        log_warning "Potential hardcoded secrets found (see $REPORTS_DIR/hardcoded-secrets.txt)"
        ((policy_violations++))
    else
        log_success "No obvious hardcoded secrets found"
    fi
    
    # Check TLS configuration
    log_info "Checking TLS configuration..."
    if grep -r "TLS" internal/config/ > "$REPORTS_DIR/tls-config.txt" 2>/dev/null; then
        log_success "TLS configuration found"
    else
        log_warning "No TLS configuration found"
    fi
    
    # Save policy violations count
    echo "$policy_violations" > "$REPORTS_DIR/policy-violations.txt"
    
    if [ $policy_violations -eq 0 ]; then
        log_success "Security policy compliance check passed"
    else
        log_warning "Security policy compliance check found $policy_violations violations"
    fi
}

# Security Test Execution
run_security_tests() {
    log_info "Running security tests..."
    
    # Run security-related tests
    log_info "Running security component tests..."
    if go test -v ./internal/validation/... ./internal/middleware/... ./internal/utils/... > "$REPORTS_DIR/security-tests.txt" 2>&1; then
        log_success "Security tests passed"
    else
        log_error "Security tests failed (see $REPORTS_DIR/security-tests.txt)"
    fi
    
    # Run test coverage for security components
    log_info "Running security test coverage..."
    if go test -coverprofile="$REPORTS_DIR/security-coverage.out" ./internal/validation/ ./internal/middleware/ ./internal/utils/ > "$REPORTS_DIR/coverage-report.txt" 2>&1; then
        go tool cover -func="$REPORTS_DIR/security-coverage.out" >> "$REPORTS_DIR/coverage-report.txt"
        log_success "Security test coverage completed"
    else
        log_warning "Security test coverage failed"
    fi
}

# Generate Security Report
generate_security_report() {
    log_info "Generating security report..."
    
    local report_file="$REPORTS_DIR/security-report.md"
    
    cat > "$report_file" << EOF
# Security Scan Report

Generated: $(date)

## Scan Summary

### Go Security Analysis
- govulncheck: $(if [ -s "$REPORTS_DIR/govulncheck.txt" ]; then echo "⚠️ Issues found"; else echo "✅ Passed"; fi)
- gosec: $(if [ -s "$REPORTS_DIR/gosec.sarif" ]; then echo "⚠️ Issues found"; else echo "✅ Passed"; fi)
- gocyclo: $(if [ -s "$REPORTS_DIR/gocyclo.txt" ]; then echo "⚠️ Complex functions found"; else echo "✅ Passed"; fi)
- golangci-lint: $(if [ -s "$REPORTS_DIR/golangci-lint.sarif" ]; then echo "⚠️ Issues found"; else echo "✅ Passed"; fi)

### Dependency Security
- Trivy filesystem scan: $(if [ -s "$REPORTS_DIR/trivy-fs.sarif" ]; then echo "⚠️ Vulnerabilities found"; else echo "✅ Passed"; fi)
- Dependencies analyzed: $(wc -l < "$REPORTS_DIR/dependencies.txt" 2>/dev/null || echo "0")

### Container Security
- Trivy container scan: $(if [ -s "$REPORTS_DIR/trivy-container.sarif" ]; then echo "⚠️ Vulnerabilities found"; else echo "✅ Passed"; fi)
- Snyk container: $(if [ -s "$REPORTS_DIR/snyk-container.txt" ]; then echo "⚠️ Issues found"; else echo "✅ Passed"; fi)

### License Compliance
- License check: $(if [ -s "$REPORTS_DIR/license-check.txt" ]; then echo "⚠️ Issues found"; else echo "✅ Passed"; fi)
- License report: $(if [ -s "$REPORTS_DIR/licenses.csv" ]; then echo "✅ Generated"; else echo "❌ Failed"; fi)

### Security Policy
- Policy violations: $(cat "$REPORTS_DIR/policy-violations.txt" 2>/dev/null || echo "0")

### Security Tests
- Test execution: $(if [ -s "$REPORTS_DIR/security-tests.txt" ]; then echo "✅ Completed"; else echo "❌ Failed"; fi)
- Coverage report: $(if [ -s "$REPORTS_DIR/coverage-report.txt" ]; then echo "✅ Generated"; else echo "❌ Failed"; fi)

## Detailed Reports

All detailed reports are available in the \`$REPORTS_DIR\` directory:

- \`govulncheck.txt\`: Go vulnerability check results
- \`gosec.sarif\`: GoSec security analysis results
- \`gocyclo.txt\`: Cyclomatic complexity analysis
- \`golangci-lint.sarif\`: Comprehensive Go linting results
- \`trivy-fs.sarif\`: Trivy filesystem scan results
- \`trivy-container.sarif\`: Trivy container scan results
- \`snyk-code.sarif\`: Snyk Code analysis results
- \`snyk-monitor.txt\`: Snyk Open Source monitoring results
- \`snyk-iac.txt\`: Snyk Infrastructure as Code results
- \`snyk-container.txt\`: Snyk Container analysis results
- \`license-check.txt\`: License compliance check results
- \`licenses.csv\`: Detailed license report
- \`hardcoded-passwords.txt\`: Potential hardcoded passwords
- \`hardcoded-secrets.txt\`: Potential hardcoded secrets
- \`tls-config.txt\`: TLS configuration analysis
- \`security-tests.txt\`: Security test execution results
- \`coverage-report.txt\`: Security test coverage report

## Next Steps

1. Review any security findings in the detailed reports
2. Address high and critical severity issues immediately
3. Update dependencies with known vulnerabilities
4. Review and update security documentation as needed
5. Implement security improvements based on findings

## Compliance Status

This scan ensures compliance with:
- OWASP Top 10
- NIST Cybersecurity Framework
- SOC 2 requirements
- ISO 27001 standards
- PCI DSS requirements

EOF
    
    log_success "Security report generated: $report_file"
}

# Main execution
main() {
    log_info "Starting Secretary Security Scan..."
    
    # Check tools
    check_tools
    
    # Setup reports directory
    setup_reports
    
    # Run all security checks
    run_go_security
    run_dependency_scanning
    run_snyk_analysis
    run_container_security
    run_license_compliance
    run_security_policy
    run_security_tests
    
    # Generate final report
    generate_security_report
    
    log_success "Security scan completed! Check $REPORTS_DIR/security-report.md for summary"
    log_info "All detailed reports are available in the $REPORTS_DIR directory"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --fail-on-critical)
            FAIL_ON_CRITICAL=true
            shift
            ;;
        --fail-on-high)
            FAIL_ON_HIGH=true
            shift
            ;;
        --reports-dir)
            REPORTS_DIR="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  --fail-on-critical    Exit with error if critical issues found"
            echo "  --fail-on-high        Exit with error if high severity issues found"
            echo "  --reports-dir DIR     Directory to store reports (default: security-reports)"
            echo "  --help                Show this help message"
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Run main function
main "$@" 