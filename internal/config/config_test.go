package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Test with environment variables set
	t.Setenv("SECRETARY_HOST", "0.0.0.0")
	t.Setenv("SECRETARY_PORT", "8080")
	t.Setenv("SECRETARY_SESSION_SECRET", "test-secret-key-that-is-long-enough-for-validation")
	t.Setenv("SECRETARY_CSRF_SECRET", "test-csrf-secret-that-is-long-enough-for-validation")
	t.Setenv("SECRETARY_READ_TIMEOUT", "30s")
	t.Setenv("SECRETARY_WRITE_TIMEOUT", "30s")
	t.Setenv("SECRETARY_IDLE_TIMEOUT", "120s")
	t.Setenv("SECRETARY_JWT_EXPIRATION", "12h")
	t.Setenv("SECRETARY_TLS_CERT_PATH", "/path/to/cert.pem")
	t.Setenv("SECRETARY_TLS_KEY_PATH", "/path/to/key.pem")
	t.Setenv("SECRETARY_DB_DRIVER", "postgres")
	t.Setenv("SECRETARY_DB_PATH", "/path/to/db")

	cfg := Load()

	expected := &Config{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         "8080",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
			TLSCertPath:  "/path/to/cert.pem",
			TLSKeyPath:   "/path/to/key.pem",
		},
		Database: DatabaseConfig{
			Driver:   "postgres",
			FilePath: "/path/to/db",
		},
		Security: SecurityConfig{
			JWTSecret:     "test-secret-key-that-is-long-enough-for-validation",
			JWTExpiration: 12 * time.Hour,
		},
	}

	if cfg.Server.Host != expected.Server.Host {
		t.Errorf("Expected Host %s, got %s", expected.Server.Host, cfg.Server.Host)
	}

	if cfg.Server.Port != expected.Server.Port {
		t.Errorf("Expected Port %s, got %s", expected.Server.Port, cfg.Server.Port)
	}

	if cfg.Server.ReadTimeout != expected.Server.ReadTimeout {
		t.Errorf("Expected ReadTimeout %v, got %v", expected.Server.ReadTimeout, cfg.Server.ReadTimeout)
	}

	if cfg.Server.WriteTimeout != expected.Server.WriteTimeout {
		t.Errorf("Expected WriteTimeout %v, got %v", expected.Server.WriteTimeout, cfg.Server.WriteTimeout)
	}

	if cfg.Server.IdleTimeout != expected.Server.IdleTimeout {
		t.Errorf("Expected IdleTimeout %v, got %v", expected.Server.IdleTimeout, cfg.Server.IdleTimeout)
	}

	if cfg.Server.TLSCertPath != expected.Server.TLSCertPath {
		t.Errorf("Expected TLSCertPath %s, got %s", expected.Server.TLSCertPath, cfg.Server.TLSCertPath)
	}

	if cfg.Server.TLSKeyPath != expected.Server.TLSKeyPath {
		t.Errorf("Expected TLSKeyPath %s, got %s", expected.Server.TLSKeyPath, cfg.Server.TLSKeyPath)
	}

	if cfg.Database.Driver != expected.Database.Driver {
		t.Errorf("Expected Driver %s, got %s", expected.Database.Driver, cfg.Database.Driver)
	}

	if cfg.Database.FilePath != expected.Database.FilePath {
		t.Errorf("Expected FilePath %s, got %s", expected.Database.FilePath, cfg.Database.FilePath)
	}

	if cfg.Security.JWTSecret != expected.Security.JWTSecret {
		t.Errorf("Expected JWTSecret %s, got %s", expected.Security.JWTSecret, cfg.Security.JWTSecret)
	}

	if cfg.Security.JWTExpiration != expected.Security.JWTExpiration {
		t.Errorf("Expected JWTExpiration %v, got %v", expected.Security.JWTExpiration, cfg.Security.JWTExpiration)
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable exists",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "environment variable does not exist",
			key:          "NON_EXISTENT_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "empty environment variable",
			key:          "EMPTY_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			original := os.Getenv(tt.key)
			defer func() {
				if original == "" {
					os.Unsetenv(tt.key)
				} else {
					os.Setenv(tt.key, original)
				}
			}()

			// Set test environment
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}

			// Test
			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
