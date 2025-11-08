package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/junghwan16/test-server/internal/health"
	"github.com/junghwan16/test-server/internal/server"

	auth_app "github.com/junghwan16/test-server/internal/auth/application"
	auth_infra_persistence "github.com/junghwan16/test-server/internal/auth/infrastructure/persistence"
	auth_http "github.com/junghwan16/test-server/internal/auth/interfaces/http"
)

func main() {
	// Load .env file first
	_ = godotenv.Load()

	logger := initLogger()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Error("JWT_SECRET environment variable is required")
		os.Exit(1)
	}

	serverPort := getEnv("SERVER_PORT", "8080")
	sessionExpiry := 1 * time.Hour

	logger.Info("starting application", "port", serverPort)

	db, err := initDB(logger)
	if err != nil {
		logger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get database instance", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	userRepository, err := auth_infra_persistence.NewGormUserRepository(db)
	if err != nil {
		logger.Error("failed to initialize user repository", "error", err)
		os.Exit(1)
	}

	authAppService := auth_app.NewAuthService(userRepository, logger, sessionExpiry)
	userAppService := auth_app.NewUserService(userRepository, logger)

	// Setup health checks
	healthService := health.NewService(db, logger)
	healthHandler := health.NewHandler(healthService)

	// Create server with health handler
	srv := server.New(logger, healthHandler)

	// Mark application as ready after all initialization
	healthService.MarkReady()

	jwtEncoder := auth_http.NewJWTEncoder(jwtSecret)

	auth_http.RegisterRoutes(srv.Router(), authAppService, userAppService, jwtEncoder, logger)

	addr := ":" + serverPort
	logger.Info("server listening", "address", addr)

	if err := http.ListenAndServe(addr, srv); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func initDB(logger *slog.Logger) (*gorm.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "postgres")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	logger.Info("database connected", "driver", "postgres")
	return db, nil
}

func initLogger() *slog.Logger {
	env := getEnv("ENV", "development")

	var handler slog.Handler
	if env == "production" {
		// Production: JSON format
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Development: Text format with debug level
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}
