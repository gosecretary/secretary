package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Security SecurityConfig
	Audit    AuditConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	TLSCertPath  string
	TLSKeyPath   string
}

type DatabaseConfig struct {
	Driver   string
	FilePath string
}

type SecurityConfig struct {
	SessionSecret     []byte
	SessionMaxAge     int
	BcryptCost        int
	RateLimitRequests int
	RateLimitWindow   time.Duration
	CSRFSecret        []byte
	SecureCookies     bool
	TrustedProxies    []string
}

type AuditConfig struct {
	Enabled     bool
	Directory   string
	LogToStdout bool
	LogToDB     bool
}

var GlobalConfig *Config

func generateRandomSecret(length int) []byte {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Failed to generate random secret:", err)
	}
	return bytes
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvSecret(key string, length int) []byte {
	if value := os.Getenv(key); value != "" {
		if decoded, err := hex.DecodeString(value); err == nil && len(decoded) >= length {
			return decoded[:length]
		}
		// If env var exists but is invalid, log warning and generate new
		log.Printf("Warning: Invalid %s in environment, generating new secret", key)
	}
	return generateRandomSecret(length)
}

func Load() *Config {
	config := &Config{
		Server: ServerConfig{
			Host:         getEnv("SECRETARY_HOST", "0.0.0.0"),
			Port:         getEnv("SECRETARY_PORT", "6080"),
			ReadTimeout:  getEnvDuration("SECRETARY_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvDuration("SECRETARY_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvDuration("SECRETARY_IDLE_TIMEOUT", 60*time.Second),
			TLSCertPath:  getEnv("SECRETARY_TLS_CERT_PATH", ""),
			TLSKeyPath:   getEnv("SECRETARY_TLS_KEY_PATH", ""),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("SECRETARY_DB_DRIVER", "sqlite3"),
			FilePath: getEnv("SECRETARY_DB_PATH", "./data/secretary.db"),
		},
		Security: SecurityConfig{
			SessionSecret:     getEnvSecret("SECRETARY_SESSION_SECRET", 32),
			SessionMaxAge:     getEnvInt("SECRETARY_SESSION_MAX_AGE", 3600), // 1 hour default
			BcryptCost:        getEnvInt("SECRETARY_BCRYPT_COST", 12),
			RateLimitRequests: getEnvInt("SECRETARY_RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:   getEnvDuration("SECRETARY_RATE_LIMIT_WINDOW", time.Hour),
			CSRFSecret:        getEnvSecret("SECRETARY_CSRF_SECRET", 32),
			SecureCookies:     getEnvBool("SECRETARY_SECURE_COOKIES", true),
			TrustedProxies:    []string{}, // TODO: Parse from env
		},
		Audit: AuditConfig{
			Enabled:     getEnvBool("SECRETARY_AUDIT_ENABLED", true),
			Directory:   getEnv("SECRETARY_AUDIT_DIR", "./data/audit/"),
			LogToStdout: getEnvBool("SECRETARY_AUDIT_STDOUT", false),
			LogToDB:     getEnvBool("SECRETARY_AUDIT_DB", true),
		},
	}

	// Validate TLS configuration
	if config.Server.TLSCertPath != "" && config.Server.TLSKeyPath == "" {
		log.Fatal("TLS certificate path provided but key path is missing")
	}
	if config.Server.TLSKeyPath != "" && config.Server.TLSCertPath == "" {
		log.Fatal("TLS key path provided but certificate path is missing")
	}

	GlobalConfig = config
	return config
}
