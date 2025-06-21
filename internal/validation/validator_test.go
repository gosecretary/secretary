package validation

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"valid username", "john_doe", false},
		{"valid username with numbers", "user123", false},
		{"valid username with dots", "john.doe", false},
		{"valid username with hyphens", "john-doe", false},
		{"too short", "ab", true},
		{"too long", "verylongusernameexceedingthelimit", true},
		{"contains invalid chars", "john@doe", true},
		{"contains spaces", "john doe", true},
		{"empty string", "", true},
		{"reserved name admin", "admin", true},
		{"reserved name root", "root", true},
		{"reserved name system", "system", true},
		{"reserved name user", "user", true},
		{"reserved name guest", "guest", true},
		{"reserved name test", "test", true},
		{"reserved name temp", "temp", true},
		{"reserved name tmp", "tmp", true},
		{"reserved name backup", "backup", true},
		{"reserved name www", "www", true},
		{"reserved name mail", "mail", true},
		{"reserved name ftp", "ftp", true},
		{"reserved name localhost", "localhost", true},
		{"reserved name 127.0.0.1", "127.0.0.1", true},
		{"reserved name ::1", "::1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "MySecurePass123!", false},
		{"valid password with special chars", "P@ssw0rd#123", false},
		{"too short", "short", true},
		{"too long", "verylongpasswordexceedingthemaximumlengthallowedbythesystem", true},
		{"no uppercase", "mypassword123!", true},
		{"no lowercase", "MYPASSWORD123!", true},
		{"no numbers", "MySecurePass!", true},
		{"no special chars", "MySecurePass123", true},
		{"empty string", "", true},
		{"only uppercase", "PASSWORD", true},
		{"only lowercase", "password", true},
		{"only numbers", "12345678", true},
		{"only special chars", "!@#$%^&*", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateResourceName(t *testing.T) {
	tests := []struct {
		testName string
		name     string
		wantErr  bool
	}{
		{"valid name", "production-db", false},
		{"valid name with numbers", "db-01", false},
		{"valid name with underscores", "test_server", false},
		{"too short", "a", true},
		{"too long", "verylongresourcenamethatexceedsthemaximumlengthallowed", true},
		{"contains spaces", "production db", true},
		{"contains special chars", "db@server", true},
		{"empty string", "", true},
		{"starts with number", "1db", true},
		{"starts with underscore", "_db", true},
		{"starts with hyphen", "-db", true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := ValidateResourceName(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateResourceName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantErr bool
	}{
		{"valid hostname", "example.com", false},
		{"valid IP", "192.168.1.1", false},
		{"valid localhost", "localhost", false},
		{"valid subdomain", "db.example.com", false},
		{"valid IPv6", "::1", false},
		{"valid IPv6 expanded", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", false},
		{"empty string", "", true},
		{"invalid chars", "host@name", true},
		{"invalid IP", "256.256.256.256", true},
		{"invalid format", "host:name", true},
		{"too long", "verylonghostnamethatexceedsthemaximumlengthallowedbythesystem", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHost(tt.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"valid port", 8080, false},
		{"valid port low", 1, false},
		{"valid port high", 65535, false},
		{"invalid port zero", 0, true},
		{"invalid port negative", -1, true},
		{"invalid port too high", 65536, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateReason(t *testing.T) {
	tests := []struct {
		name    string
		reason  string
		wantErr bool
	}{
		{"valid reason", "Need access for database maintenance", false},
		{"valid reason short", "Debug", false},
		{"valid reason long", "This is a very long reason that explains why access is needed for the specific resource and what tasks will be performed", false},
		{"empty string", "", true},
		{"too long", "This is an extremely long reason that exceeds the maximum allowed length for the reason field in the access request system", true},
		{"contains SQL injection", "'; DROP TABLE users; --", true},
		{"contains XSS", "<script>alert('xss')</script>", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReason(tt.reason)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateReason() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainsSQLInjectionPatterns(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"no injection", "SELECT * FROM users", false},
		{"SQL injection 1", "'; DROP TABLE users; --", true},
		{"SQL injection 2", "1' OR '1'='1", true},
		{"SQL injection 3", "admin'--", true},
		{"SQL injection 4", "1; DROP TABLE users", true},
		{"SQL injection 5", "1 UNION SELECT * FROM users", true},
		{"SQL injection 6", "1' UNION SELECT * FROM users--", true},
		{"SQL injection 7", "1' AND 1=1--", true},
		{"SQL injection 8", "1' OR 1=1#", true},
		{"SQL injection 9", "1' OR 1=1/*", true},
		{"SQL injection 10", "1' OR 1=1--", true},
		{"normal text", "This is normal text", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsSQLInjectionPatterns(tt.input)
			if got != tt.want {
				t.Errorf("containsSQLInjectionPatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsXSSPatterns(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"no XSS", "Hello world", false},
		{"XSS script tag", "<script>alert('xss')</script>", true},
		{"XSS img tag", "<img src=x onerror=alert('xss')>", true},
		{"XSS iframe tag", "<iframe src=javascript:alert('xss')>", true},
		{"XSS onload", "javascript:alert('xss')", true},
		{"XSS onmouseover", "onmouseover=alert('xss')", true},
		{"XSS onclick", "onclick=alert('xss')", true},
		{"XSS onerror", "onerror=alert('xss')", true},
		{"XSS onfocus", "onfocus=alert('xss')", true},
		{"XSS onblur", "onblur=alert('xss')", true},
		{"XSS onchange", "onchange=alert('xss')", true},
		{"XSS onsubmit", "onsubmit=alert('xss')", true},
		{"XSS onreset", "onreset=alert('xss')", true},
		{"XSS onselect", "onselect=alert('xss')", true},
		{"XSS onunload", "onunload=alert('xss')", true},
		{"XSS onresize", "onresize=alert('xss')", true},
		{"XSS onscroll", "onscroll=alert('xss')", true},
		{"XSS onkeydown", "onkeydown=alert('xss')", true},
		{"XSS onkeyup", "onkeyup=alert('xss')", true},
		{"XSS onkeypress", "onkeypress=alert('xss')", true},
		{"normal text", "This is normal text", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsXSSPatterns(tt.input)
			if got != tt.want {
				t.Errorf("containsXSSPatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"normal text", "Hello world", "Hello world"},
		{"with script tags", "<script>alert('xss')</script>Hello", "Hello"},
		{"with HTML tags", "<p>Hello</p> world", "Hello world"},
		{"with multiple tags", "<div><span>Hello</span> <b>world</b></div>", "Hello world"},
		{"with special chars", "Hello & world", "Hello & world"},
		{"empty string", "", ""},
		{"only tags", "<script></script><div></div>", ""},
		{"mixed content", "Hello<script>alert('xss')</script>world", "Helloworld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeInput(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Field:   "username",
		Message: "Invalid username format",
	}

	expected := "username: Invalid username format"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), expected)
	}
}
