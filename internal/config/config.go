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
	// Generate random secret if not provided
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		randomBytes := make([]byte, 32)
		if _, err := rand.Read(randomBytes); err != nil {
			utils.Fatal("Failed to generate random secret: " + err.Error())
		}
		secret = base64.URLEncoding.EncodeToString(randomBytes)
		utils.Warnf("Warning: Invalid JWT_SECRET in environment, generating new secret")
	}

	// Parse JWT expiration
	jwtExpiration := 24 * time.Hour // Default to 24 hours
	if exp := os.Getenv("JWT_EXPIRATION"); exp != "" {
		if duration, err := time.ParseDuration(exp); err == nil {
			jwtExpiration = duration
		} else {
			utils.Warnf("Warning: Invalid JWT_EXPIRATION in environment, using default")
		}
	}

	// Parse timeouts
	readTimeout := 15 * time.Second
	if rt := os.Getenv("SERVER_READ_TIMEOUT"); rt != "" {
		if duration, err := time.ParseDuration(rt); err == nil {
			readTimeout = duration
		}
	}

	writeTimeout := 15 * time.Second
	if wt := os.Getenv("SERVER_WRITE_TIMEOUT"); wt != "" {
		if duration, err := time.ParseDuration(wt); err == nil {
			writeTimeout = duration
		}
	}

	idleTimeout := 60 * time.Second
	if it := os.Getenv("SERVER_IDLE_TIMEOUT"); it != "" {
		if duration, err := time.ParseDuration(it); err == nil {
			idleTimeout = duration
		}
	}

	// Validate TLS configuration
	tlsCertPath := os.Getenv("TLS_CERT_PATH")
	tlsKeyPath := os.Getenv("TLS_KEY_PATH")
	if (tlsCertPath != "" && tlsKeyPath == "") || (tlsCertPath == "" && tlsKeyPath != "") {
		utils.Fatal("TLS certificate and key paths must be provided together")
	}

	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "localhost"),
			Port:         getEnv("SERVER_PORT", "6080"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
			TLSCertPath:  tlsCertPath,
			TLSKeyPath:   tlsKeyPath,
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite3"),
			FilePath: getEnv("DB_FILE_PATH", "secretary.db"),
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
