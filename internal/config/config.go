package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Logger    LoggerConfig
	SMTP      SMTPConfig
	Session   SessionConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

type LoggerConfig struct {
	Environment string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type SessionConfig struct {
	TTL int // seconds
}

type RateLimitConfig struct {
	RequestsPerSecond float64
	Burst             int
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "postgres"),
			Port:     getEnv("DB_PORT", "5432"),
		},
		Logger: LoggerConfig{
			Environment: getEnv("ENV", "development"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     getEnv("SMTP_PORT", "1025"),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@example.com"),
		},
		Session: SessionConfig{
			TTL: getEnvInt("SESSION_TTL", 86400), // 24 hours
		},
		RateLimit: RateLimitConfig{
			RequestsPerSecond: getEnvFloat("RATE_LIMIT_RPS", 10),
			Burst:             getEnvInt("RATE_LIMIT_BURST", 20),
		},
	}

	return cfg, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as int with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := fmt.Sscanf(value, "%d", new(int)); err == nil && intVal == 1 {
			var result int
			fmt.Sscanf(value, "%d", &result)
			return result
		}
	}
	return defaultValue
}

// getEnvFloat gets an environment variable as float64 with a default value
func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := fmt.Sscanf(value, "%f", new(float64)); err == nil && floatVal == 1 {
			var result float64
			fmt.Sscanf(value, "%f", &result)
			return result
		}
	}
	return defaultValue
}

// DSN returns the database connection string
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		d.Host, d.User, d.Password, d.Name, d.Port)
}

// IsProduction returns true if the environment is production
func (l *LoggerConfig) IsProduction() bool {
	return l.Environment == "production"
}
