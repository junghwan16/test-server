package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Logger   LoggerConfig
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
	// Connection pool settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type AuthConfig struct {
	JWTSecret     string
	SessionExpiry time.Duration
}

type LoggerConfig struct {
	Environment string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			Name:            getEnv("DB_NAME", "postgres"),
			Port:            getEnv("DB_PORT", "5432"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Auth: AuthConfig{
			JWTSecret:     os.Getenv("JWT_SECRET"),
			SessionExpiry: getEnvAsDuration("SESSION_EXPIRY", 1*time.Hour),
		},
		Logger: LoggerConfig{
			Environment: getEnv("ENV", "development"),
		},
	}

	// Validate required fields
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks if all required configuration fields are set
func (c *Config) validate() error {
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsDuration gets an environment variable as a duration with a default value
// The environment variable should be in the format accepted by time.ParseDuration (e.g., "1h", "30m")
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
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
