package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/junghwan16/test-server/internal/config"
	"github.com/junghwan16/test-server/internal/identity/application"
	"github.com/junghwan16/test-server/internal/identity/handler"
	"github.com/junghwan16/test-server/internal/identity/infrastructure/persistence"
	"github.com/junghwan16/test-server/internal/server"
	"github.com/junghwan16/test-server/internal/shared/domain"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := newLogger(cfg)

	db, err := connectDB(cfg, logger)
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	if err := db.AutoMigrate(
		&persistence.UserModel{},
		&persistence.EmailVerificationModel{},
		&persistence.PasswordResetModel{},
	); err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}
	logger.Info("database migrated")

	rdb, err := connectRedis(cfg, logger)
	if err != nil {
		logger.Error("failed to connect redis", "error", err)
		os.Exit(1)
	}

	eventBus := domain.NewSimpleEventBus()

	userRepo := persistence.NewUserRepository(db, eventBus)
	sessionRepo := persistence.NewRedisSessionRepository(rdb)
	emailVerifRepo := persistence.NewEmailVerificationRepository(db)
	passwordResetRepo := persistence.NewPasswordResetRepository(db)

	userSvc := application.NewUserService(userRepo)
	authSvc := application.NewAuthService(userRepo, sessionRepo, cfg.Session.TTL)
	verifSvc := application.NewVerificationService(
		userRepo,
		emailVerifRepo,
		passwordResetRepo,
		24*time.Hour,
		1*time.Hour,
	)

	authHandler := handler.NewAuthHandler(userSvc, authSvc, cfg.Session.TTL)
	usersHandler := handler.NewUsersHandler(userSvc)
	verifHandler := handler.NewVerificationHandler(verifSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health/live", server.HandleLive)
	mux.HandleFunc("GET /health/ready", server.HandleReady(db, logger))
	mux.HandleFunc("GET /health", server.HandleHealth(db, logger))

	mux.HandleFunc("POST /auth/signup", authHandler.Signup)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)
	mux.Handle("GET /me", server.RequireAuth(authSvc)(http.HandlerFunc(authHandler.Me)))

	mux.Handle("POST /verification/request", server.RequireAuth(authSvc)(http.HandlerFunc(verifHandler.RequestVerification)))
	mux.HandleFunc("GET /verification/verify", verifHandler.VerifyEmail)

	mux.HandleFunc("POST /password/reset/request", verifHandler.RequestPasswordReset)
	mux.HandleFunc("POST /password/reset/confirm", verifHandler.ResetPassword)

	mux.Handle("GET /admin/users", server.RequireAdmin(authSvc)(http.HandlerFunc(usersHandler.ListUsers)))
	mux.Handle("GET /admin/users/{id}", server.RequireAdmin(authSvc)(http.HandlerFunc(usersHandler.GetUser)))
	mux.Handle("PATCH /admin/users/{id}", server.RequireAdmin(authSvc)(http.HandlerFunc(usersHandler.UpdateUser)))
	mux.Handle("DELETE /admin/users/{id}", server.RequireAdmin(authSvc)(http.HandlerFunc(usersHandler.DeleteUser)))

	handler := server.Logging(logger)(server.RateLimit(cfg.RateLimit.RequestsPerSecond, cfg.RateLimit.Burst)(mux))

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: handler,
	}

	go func() {
		logger.Info("server listening", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
	}

	if sqlDB, _ := db.DB(); sqlDB != nil {
		logger.Info("closing database")
		sqlDB.Close()
	}

	if rdb != nil {
		logger.Info("closing redis")
		rdb.Close()
	}

	logger.Info("shutdown complete")
}

func newLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler
	if cfg.Logger.IsProduction() {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	return slog.New(handler)
}

func connectDB(cfg *config.Config, logger *slog.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	logger.Info("database connected")
	return db, nil
}

func connectRedis(cfg *config.Config, logger *slog.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	logger.Info("redis connected")
	return client, nil
}
