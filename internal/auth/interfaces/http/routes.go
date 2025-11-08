package http

import (
	"log/slog"
	"net/http"

	"github.com/junghwan16/test-server/internal/auth/application"
)

func RegisterRoutes(router *http.ServeMux, authService *application.AuthService, userService *application.UserService, jwtEncoder *JWTEncoder, logger *slog.Logger) {
	authHandler := NewAuthHandler(authService, jwtEncoder, logger)
	userHandler := NewUserHandler(userService, logger)
	authMiddleware := NewAuthMiddleware(jwtEncoder, logger)

	// Auth routes (public)
	router.HandleFunc("POST /auth/login", authHandler.Login)
	router.HandleFunc("POST /auth/register", authHandler.Register)

	// User routes (authenticated)
	router.Handle("GET /users/me", authMiddleware.Middleware(http.HandlerFunc(userHandler.GetProfile)))
	router.Handle("PATCH /users/me", authMiddleware.Middleware(http.HandlerFunc(userHandler.UpdateProfile)))
	router.Handle("PATCH /users/me/password", authMiddleware.Middleware(http.HandlerFunc(userHandler.ChangePassword)))
	router.Handle("DELETE /users/me", authMiddleware.Middleware(http.HandlerFunc(userHandler.DeleteAccount)))
}
