package config

import (
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Auth     AuthConfig
	Database DatabaseConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MaxHeaderBytes  int
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecretKey    string
	TokenExpiry     time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver string
	DSN    string
}

// NewDefaultConfig returns a default configuration
func NewDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:           ":8080",
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1 MB
		},
		Auth: AuthConfig{
			JWTSecretKey: "your-secret-key-change-in-production", // Change in production!
			TokenExpiry:  24 * time.Hour,
		},
		Database: DatabaseConfig{
			Driver: "memory", // Using in-memory for now, can be changed to postgres/mysql
			DSN:    "",
		},
	}
}
