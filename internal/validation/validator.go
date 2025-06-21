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
	MaxReasonLength   = 120
	MaxResourceName   = 50
	MaxHostLength     = 60
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
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
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

	// Check for reserved usernames
	reservedNames := []string{
		"admin", "root", "system", "user", "guest", "test", "temp", "tmp",
		"backup", "www", "mail", "ftp", "localhost", "127.0.0.1", "::1",
	}
	for _, reserved := range reservedNames {
		if strings.ToLower(username) == reserved {
			return ValidationError{
				Field:   "username",
				Message: fmt.Sprintf("'%s' is a reserved username", username),
			}
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

	if !hasUpper {
		return ValidationError{
			Field:   "password",
			Message: "must contain at least one uppercase letter",
		}
	}

	if !hasLower {
		return ValidationError{
			Field:   "password",
			Message: "must contain at least one lowercase letter",
		}
	}

	if !hasDigit {
		return ValidationError{
			Field:   "password",
			Message: "must contain at least one digit",
		}
	}

	if !hasSpecial {
		return ValidationError{
			Field:   "password",
			Message: "must contain at least one special character",
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

	if len(name) < 2 {
		return ValidationError{
			Field:   "resource_name",
			Message: "must be at least 2 characters long",
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

	// Check if starts with invalid character
	if len(name) > 0 && (name[0] >= '0' && name[0] <= '9' || name[0] == '_' || name[0] == '-') {
		return ValidationError{
			Field:   "resource_name",
			Message: "cannot start with a number, underscore, or hyphen",
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
		// Additional validation for IP addresses
		if ip.IsLoopback() && host != "localhost" && host != "::1" {
			return ValidationError{
				Field:   "host",
				Message: "loopback IP addresses are not allowed",
			}
		}
		return nil // Valid IP address
	}

	// If it looks like an IP (all numeric parts separated by dots or colons) but failed to parse, it's invalid
	ipLike := regexp.MustCompile(`^([0-9]+\.){3}[0-9]+$|^([0-9a-fA-F:]+)$`)
	if ipLike.MatchString(host) {
		return ValidationError{
			Field:   "host",
			Message: "invalid IP address",
		}
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

	// More specific patterns that are likely to be SQL injection
	patterns := []string{
		"';", "\";", "--", "/*", "*/", "@@", "#",
		"union select", "union all select", "union distinct select",
		"drop table", "drop database", "drop index", "drop view",
		"create table", "create database", "create index", "create view",
		"alter table", "alter database", "alter index", "alter view",
		"exec(", "execute(", "sp_executesql", "xp_", "sp_",
		"char(0x", "ascii(0x", "substring(0x", "length(0x", "version(0x",
		"database(0x", "user(0x", "system_user(0x", "session_user(0x",
		"0x", "0b", "\\x", "\\u",
		"or '1'='1", "or 1=1", "and 1=1", "or '1'='1'",
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
		"onreset=", "onselect=", "onunload=", "onresize=", "onscroll=", "onkeydown=",
		"onkeyup=", "onkeypress=",
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

	// Remove script content first (including content between script tags)
	scriptRegex := regexp.MustCompile(`<script[^>]*>.*?</script>`)
	input = scriptRegex.ReplaceAllString(input, "")

	// Remove HTML tags using regex
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	input = htmlTagRegex.ReplaceAllString(input, "")

	// Trim whitespace
	input = strings.TrimSpace(input)

	return input
}
