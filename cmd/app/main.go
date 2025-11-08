package main

import (
	"log/slog"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/junghwan16/test-server/internal/config"
	"github.com/junghwan16/test-server/internal/health"
	"github.com/junghwan16/test-server/internal/server"

	auth_app "github.com/junghwan16/test-server/internal/auth/application"
	auth_infra_persistence "github.com/junghwan16/test-server/internal/auth/infrastructure/persistence"
	auth_http "github.com/junghwan16/test-server/internal/auth/interfaces/http"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger := initLogger(&cfg.Logger)
	logger.Info("starting application", "port", cfg.Server.Port)

	db, err := initDB(&cfg.Database, logger)
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

	authAppService := auth_app.NewAuthService(userRepository, logger, cfg.Auth.SessionExpiry)
	userAppService := auth_app.NewUserService(userRepository, logger)

	srv := server.New(logger)

	// [HEALTH]
	healthHandler := health.NewHandler(db, logger)
	srv.Router().HandleFunc("GET /health/live", healthHandler.Live)
	srv.Router().HandleFunc("GET /health/ready", healthHandler.Ready)

	// [AUTH]
	jwtEncoder := auth_http.NewJWTEncoder(cfg.Auth.JWTSecret)
	auth_http.RegisterRoutes(srv.Router(), authAppService, userAppService, jwtEncoder, logger)

	addr := ":" + cfg.Server.Port
	logger.Info("server listening", "address", addr)

	if err := http.ListenAndServe(addr, srv); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func initDB(cfg *config.DatabaseConfig, logger *slog.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	logger.Info("database connected", "driver", "postgres")
	return db, nil
}

// TODO: zap logger가 더 효율적이라면 변경을 고려.
func initLogger(cfg *config.LoggerConfig) *slog.Logger {
	var handler slog.Handler
	if cfg.IsProduction() {
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
