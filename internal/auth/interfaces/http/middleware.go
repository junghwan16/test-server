package http

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type userContextKey string

const UserIDKey = userContextKey("userID")

func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

type AuthMiddleware struct {
	jwtEncoder *JWTEncoder
	logger     *slog.Logger
}

func NewAuthMiddleware(jwtEncoder *JWTEncoder, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtEncoder: jwtEncoder,
		logger:     logger,
	}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Warn("authorization header missing")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			m.logger.Warn("invalid authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]
		userIDStr, err := m.jwtEncoder.DecodeToken(tokenString)
		if err != nil {
			m.logger.Warn("token validation failed", "error", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			m.logger.Warn("invalid user ID in token", "error", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, uint(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
