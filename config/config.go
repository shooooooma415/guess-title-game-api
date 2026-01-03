package config

import (
	"os"
	"strconv"
)

// Config represents application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		dbPort = 5432
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", getEnv("SERVER_PORT", "8080")),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "guess_title_game"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
