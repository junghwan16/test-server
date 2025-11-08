package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/junghwan16/test-server/internal/auth/application"
	"github.com/junghwan16/test-server/internal/auth/domain"
)

type AuthHandler struct {
	authService *application.AuthService
	jwtEncoder  *JWTEncoder
	logger      *slog.Logger
}

func NewAuthHandler(authService *application.AuthService, jwtEncoder *JWTEncoder, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtEncoder:  jwtEncoder,
		logger:      logger,
	}
}

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body loginRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.logger.Error("failed to decode login request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req := application.LoginRequest{
		Username:  body.Username,
		Password:  body.Password,
		IPAddress: getIPAddress(r),
		UserAgent: r.UserAgent(),
	}

	res, err := h.authService.Login(req)
	if err != nil {
		h.logger.Warn("login failed", "error", err, "username", body.Username)

		// Prevent username enumeration attack
		if errors.Is(err, domain.ErrInvalidCredentials) || errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.jwtEncoder.EncodeSession(res.Session)
	if err != nil {
		h.logger.Error("failed to encode session to JWT", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"token":      token,
		"expires_at": res.Session.ExpiresAt(),
	})
}

type registerRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body registerRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.logger.Error("failed to decode register request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req := application.RegisterRequest{
		Username:  body.Username,
		Password:  body.Password,
		IPAddress: getIPAddress(r),
	}

	res, err := h.authService.Register(req)
	if err != nil {
		h.logger.Warn("registration failed", "error", err, "username", body.Username)

		switch {
		case errors.Is(err, domain.ErrUsernameAlreadyExists):
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	token, err := h.jwtEncoder.EncodeSession(res.Session)
	if err != nil {
		h.logger.Error("failed to encode session to JWT", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"user_id":    res.UserID,
		"username":   res.Username,
		"token":      token,
		"expires_at": res.ExpiresAt,
	})
}

func getIPAddress(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}
