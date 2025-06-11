package validation

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Constants for validation rules
const (
	MinUsernameLength = 3
	MaxUsernameLength = 32
	MinPasswordLength = 8
	MaxPasswordLength = 128
	MaxReasonLength   = 1000
	MaxResourceName   = 64
	MaxHostLength     = 253
	MinPortNumber     = 1
	MaxPortNumber     = 65535
)

// Regular expressions for validation
var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	hostnameRegex = regexp.MustCompile(`^[a-zA-Z0-9.-]+$`)
	resourceRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// ValidateUsername validates a username
func ValidateUsername(username string) error {
	if len(username) < MinUsernameLength {
		return ValidationError{
			Field:   "username",
			Message: fmt.Sprintf("must be at least %d characters long", MinUsernameLength),
		}
	}

	if len(username) > MaxUsernameLength {
		return ValidationError{
			Field:   "username",
			Message: fmt.Sprintf("must be no more than %d characters long", MaxUsernameLength),
		}
	}

	if !utf8.ValidString(username) {
		return ValidationError{
			Field:   "username",
			Message: "must be valid UTF-8",
		}
	}

	if !usernameRegex.MatchString(username) {
		return ValidationError{
			Field:   "username",
			Message: "can only contain letters, numbers, dots, hyphens, and underscores",
		}
	}

	// Check for SQL injection patterns
	if containsSQLInjectionPatterns(username) {
		return ValidationError{
			Field:   "username",
			Message: "contains invalid characters",
		}
	}

	return nil
}

// ValidatePassword validates a password
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ValidationError{
			Field:   "password",
			Message: fmt.Sprintf("must be at least %d characters long", MinPasswordLength),
		}
	}

	if len(password) > MaxPasswordLength {
		return ValidationError{
			Field:   "password",
			Message: fmt.Sprintf("must be no more than %d characters long", MaxPasswordLength),
		}
	}

	if !utf8.ValidString(password) {
		return ValidationError{
			Field:   "password",
			Message: "must be valid UTF-8",
		}
	}

	// Check password strength
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?~` + "`" + `]`).MatchString(password)

	strengthCount := 0
	if hasUpper {
		strengthCount++
	}
	if hasLower {
		strengthCount++
	}
	if hasDigit {
		strengthCount++
	}
	if hasSpecial {
		strengthCount++
	}

	if strengthCount < 3 {
		return ValidationError{
			Field:   "password",
			Message: "must contain at least 3 of: uppercase letter, lowercase letter, digit, special character",
		}
	}

	return nil
}

// ValidateResourceName validates a resource name
func ValidateResourceName(name string) error {
	if len(name) == 0 {
		return ValidationError{
			Field:   "resource_name",
			Message: "cannot be empty",
		}
	}

	if len(name) > MaxResourceName {
		return ValidationError{
			Field:   "resource_name",
			Message: fmt.Sprintf("must be no more than %d characters long", MaxResourceName),
		}
	}

	if !utf8.ValidString(name) {
		return ValidationError{
			Field:   "resource_name",
			Message: "must be valid UTF-8",
		}
	}

	if !resourceRegex.MatchString(name) {
		return ValidationError{
			Field:   "resource_name",
			Message: "can only contain letters, numbers, dots, hyphens, and underscores",
		}
	}

	if containsSQLInjectionPatterns(name) {
		return ValidationError{
			Field:   "resource_name",
			Message: "contains invalid characters",
		}
	}

	return nil
}

// ValidateHost validates a hostname or IP address
func ValidateHost(host string) error {
	if len(host) == 0 {
		return ValidationError{
			Field:   "host",
			Message: "cannot be empty",
		}
	}

	if len(host) > MaxHostLength {
		return ValidationError{
			Field:   "host",
			Message: fmt.Sprintf("must be no more than %d characters long", MaxHostLength),
		}
	}

	if !utf8.ValidString(host) {
		return ValidationError{
			Field:   "host",
			Message: "must be valid UTF-8",
		}
	}

	// Try to parse as IP address first
	if ip := net.ParseIP(host); ip != nil {
		return nil // Valid IP address
	}

	// Validate as hostname
	if !hostnameRegex.MatchString(host) {
		return ValidationError{
			Field:   "host",
			Message: "must be a valid hostname or IP address",
		}
	}

	if containsSQLInjectionPatterns(host) {
		return ValidationError{
			Field:   "host",
			Message: "contains invalid characters",
		}
	}

	return nil
}

// ValidatePort validates a port number
func ValidatePort(port int) error {
	if port < MinPortNumber || port > MaxPortNumber {
		return ValidationError{
			Field:   "port",
			Message: fmt.Sprintf("must be between %d and %d", MinPortNumber, MaxPortNumber),
		}
	}
	return nil
}

// ValidateReason validates a reason text (for requests)
func ValidateReason(reason string) error {
	if len(reason) == 0 {
		return ValidationError{
			Field:   "reason",
			Message: "cannot be empty",
		}
	}

	if len(reason) > MaxReasonLength {
		return ValidationError{
			Field:   "reason",
			Message: fmt.Sprintf("must be no more than %d characters long", MaxReasonLength),
		}
	}

	if !utf8.ValidString(reason) {
		return ValidationError{
			Field:   "reason",
			Message: "must be valid UTF-8",
		}
	}

	// Check for potential XSS or injection patterns
	if containsSQLInjectionPatterns(reason) || containsXSSPatterns(reason) {
		return ValidationError{
			Field:   "reason",
			Message: "contains invalid characters",
		}
	}

	return nil
}

// containsSQLInjectionPatterns checks for common SQL injection patterns
func containsSQLInjectionPatterns(input string) bool {
	input = strings.ToLower(input)

	patterns := []string{
		"'", "\"", ";", "--", "/*", "*/", "@@", "@",
		"union", "select", "insert", "update", "delete", "drop", "create", "alter",
		"exec", "execute", "sp_", "xp_", "sp_executesql",
		"char(", "ascii(", "substring(", "length(", "version(",
		"database(", "user(", "system_user", "session_user",
		"0x", "0b", "\\x", "\\u",
	}

	for _, pattern := range patterns {
		if strings.Contains(input, pattern) {
			return true
		}
	}

	return false
}

// containsXSSPatterns checks for common XSS patterns
func containsXSSPatterns(input string) bool {
	input = strings.ToLower(input)

	patterns := []string{
		"<script", "</script>", "javascript:", "vbscript:", "onload=", "onerror=",
		"onclick=", "onmouseover=", "onfocus=", "onblur=", "onchange=", "onsubmit=",
		"<iframe", "<object", "<embed", "<link", "<meta", "<style",
		"data:text/html", "data:text/javascript",
	}

	for _, pattern := range patterns {
		if strings.Contains(input, pattern) {
			return true
		}
	}

	return false
}

// SanitizeInput removes or escapes potentially dangerous characters
func SanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Trim whitespace
	input = strings.TrimSpace(input)

	return input
}
