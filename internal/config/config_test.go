package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expected: &Config{
				Server: ServerConfig{
					Host:         "localhost",
					Port:         "6080",
					ReadTimeout:  15 * time.Second,
					WriteTimeout: 15 * time.Second,
					IdleTimeout:  60 * time.Second,
					TLSCertPath:  "",
					TLSKeyPath:   "",
				},
				Database: DatabaseConfig{
					Driver:   "sqlite3",
					FilePath: "secretary.db",
				},
				Security: SecurityConfig{
					JWTSecret:     "", // Will be generated randomly
					JWTExpiration: 24 * time.Hour,
				},
			},
		},
		{
			name: "custom values from environment",
			envVars: map[string]string{
				"SERVER_HOST":          "0.0.0.0",
				"SERVER_PORT":          "9090",
				"SERVER_READ_TIMEOUT":  "45s",
				"SERVER_WRITE_TIMEOUT": "45s",
				"SERVER_IDLE_TIMEOUT":  "90s",
				"TLS_CERT_PATH":        "/path/to/cert.pem",
				"TLS_KEY_PATH":         "/path/to/key.pem",
				"DB_DRIVER":            "postgres",
				"DB_FILE_PATH":         "/data/secretary.db",
				"JWT_SECRET":           "custom-secret",
				"JWT_EXPIRATION":       "12h",
			},
			expected: &Config{
				Server: ServerConfig{
					Host:         "0.0.0.0",
					Port:         "9090",
					ReadTimeout:  45 * time.Second,
					WriteTimeout: 45 * time.Second,
					IdleTimeout:  90 * time.Second,
					TLSCertPath:  "/path/to/cert.pem",
					TLSKeyPath:   "/path/to/key.pem",
				},
				Database: DatabaseConfig{
					Driver:   "postgres",
					FilePath: "/data/secretary.db",
				},
				Security: SecurityConfig{
					JWTSecret:     "custom-secret",
					JWTExpiration: 12 * time.Hour,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			originalEnv := make(map[string]string)
			for key := range tt.envVars {
				originalEnv[key] = os.Getenv(key)
				os.Unsetenv(key)
			}

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Load config
			config := Load()

			// Verify
			assert.Equal(t, tt.expected.Server.Host, config.Server.Host)
			assert.Equal(t, tt.expected.Server.Port, config.Server.Port)
			assert.Equal(t, tt.expected.Server.ReadTimeout, config.Server.ReadTimeout)
			assert.Equal(t, tt.expected.Server.WriteTimeout, config.Server.WriteTimeout)
			assert.Equal(t, tt.expected.Server.IdleTimeout, config.Server.IdleTimeout)
			assert.Equal(t, tt.expected.Server.TLSCertPath, config.Server.TLSCertPath)
			assert.Equal(t, tt.expected.Server.TLSKeyPath, config.Server.TLSKeyPath)
			assert.Equal(t, tt.expected.Database.Driver, config.Database.Driver)
			assert.Equal(t, tt.expected.Database.FilePath, config.Database.FilePath)

			// For JWT secret, if expected is empty (randomly generated), just check it's not empty
			if tt.expected.Security.JWTSecret == "" {
				assert.NotEmpty(t, config.Security.JWTSecret)
			} else {
				assert.Equal(t, tt.expected.Security.JWTSecret, config.Security.JWTSecret)
			}
			assert.Equal(t, tt.expected.Security.JWTExpiration, config.Security.JWTExpiration)

			// Restore original environment
			for key, value := range originalEnv {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}
		})
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
