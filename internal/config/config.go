package config

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"

	"secretary/alpha/pkg/utils"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Security SecurityConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	TLSCertPath  string
	TLSKeyPath   string
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Driver   string
	FilePath string
}

// SecurityConfig holds security-specific configuration
type SecurityConfig struct {
	JWTSecret     string
	JWTExpiration time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	// Security: Generate secure secrets if not provided
	secret := os.Getenv("SECRETARY_SESSION_SECRET")
	if secret == "" {
		// Generate a secure random secret
		bytes := make([]byte, 32)
		if _, err := rand.Read(bytes); err != nil {
			utils.Fatalf("Failed to generate session secret: %v", err)
		}
		secret = base64.URLEncoding.EncodeToString(bytes)
		utils.Warn("WARNING: No SECRETARY_SESSION_SECRET provided. Generated temporary secret. Set this in production!")
	}

	// Security: Validate secret length
	if len(secret) < 32 {
		utils.Fatalf("SECRETARY_SESSION_SECRET must be at least 32 characters long")
	}

	// Security: Generate CSRF secret if not provided
	csrfSecret := os.Getenv("SECRETARY_CSRF_SECRET")
	if csrfSecret == "" {
		bytes := make([]byte, 32)
		if _, err := rand.Read(bytes); err != nil {
			utils.Fatalf("Failed to generate CSRF secret: %v", err)
		}
		csrfSecret = base64.URLEncoding.EncodeToString(bytes)
		utils.Warn("WARNING: No SECRETARY_CSRF_SECRET provided. Generated temporary secret. Set this in production!")
	}

	// Security: Validate CSRF secret length
	if len(csrfSecret) < 32 {
		utils.Fatalf("SECRETARY_CSRF_SECRET must be at least 32 characters long")
	}

	// Parse timeouts with secure defaults
	readTimeout, err := time.ParseDuration(getEnv("SECRETARY_READ_TIMEOUT", "15s"))
	if err != nil {
		utils.Fatalf("Invalid SECRETARY_READ_TIMEOUT: %v", err)
	}

	writeTimeout, err := time.ParseDuration(getEnv("SECRETARY_WRITE_TIMEOUT", "15s"))
	if err != nil {
		utils.Fatalf("Invalid SECRETARY_WRITE_TIMEOUT: %v", err)
	}

	idleTimeout, err := time.ParseDuration(getEnv("SECRETARY_IDLE_TIMEOUT", "60s"))
	if err != nil {
		utils.Fatalf("Invalid SECRETARY_IDLE_TIMEOUT: %v", err)
	}

	jwtExpiration, err := time.ParseDuration(getEnv("SECRETARY_JWT_EXPIRATION", "24h"))
	if err != nil {
		utils.Fatalf("Invalid SECRETARY_JWT_EXPIRATION: %v", err)
	}

	// Security: TLS configuration
	tlsCertPath := os.Getenv("SECRETARY_TLS_CERT_PATH")
	tlsKeyPath := os.Getenv("SECRETARY_TLS_KEY_PATH")

	// Security: Check if running in production without TLS
	if os.Getenv("SECRETARY_ENVIRONMENT") == "production" {
		if tlsCertPath == "" || tlsKeyPath == "" {
			utils.Fatalf("TLS certificate and key paths are required in production environment")
		}
	}

	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SECRETARY_HOST", "localhost"),
			Port:         getEnv("SECRETARY_PORT", "6080"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
			TLSCertPath:  tlsCertPath,
			TLSKeyPath:   tlsKeyPath,
		},
		Database: DatabaseConfig{
			Driver:   getEnv("SECRETARY_DB_DRIVER", "sqlite3"),
			FilePath: getEnv("SECRETARY_DB_PATH", "./data/secretary.db"),
		},
		Security: SecurityConfig{
			JWTSecret:     secret,
			JWTExpiration: jwtExpiration,
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
