package config

import (
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Auth     AuthConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

type AuthConfig struct {
	JWTSecretKey string
	TokenExpiry  time.Duration
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:           ":8080",
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		Auth: AuthConfig{
			JWTSecretKey: getEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production"),
			TokenExpiry:  24 * time.Hour,
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "memory"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "stocktrack"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
