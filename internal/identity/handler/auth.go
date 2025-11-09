package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/junghwan16/test-server/internal/identity/application"
	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

type AuthHandler struct {
	userSvc    *application.UserService
	authSvc    *application.AuthService
	sessionTTL int
}

func NewAuthHandler(userSvc *application.UserService, authSvc *application.AuthService, sessionTTL int) *AuthHandler {
	return &AuthHandler{
		userSvc:    userSvc,
		authSvc:    authSvc,
		sessionTTL: sessionTTL,
	}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Use IDDD UserService
	u, err := h.userSvc.RegisterUser(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, user.ErrInvalidEmail) {
			http.Error(w, "Invalid email", http.StatusBadRequest)
		} else if errors.Is(err, user.ErrPasswordTooShort) {
			http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		} else {
			http.Error(w, "Email already registered", http.StatusConflict)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "User created successfully",
		"user":    userToDTO(u),
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Use IDDD AuthService
	session, u, err := h.authSvc.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.ID().Value(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   h.sessionTTL,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"user": userToDTO(u),
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		_ = h.authSvc.Logout(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out",
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	u := GetUserFromContext(r.Context())
	if u == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userToDTO(u))
}
