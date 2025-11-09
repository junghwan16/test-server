package handler

import (
	"encoding/json"
	"net/http"

	"github.com/junghwan16/test-server/internal/identity/application"
)

type VerificationHandler struct {
	verifSvc *application.VerificationService
}

func NewVerificationHandler(verifSvc *application.VerificationService) *VerificationHandler {
	return &VerificationHandler{
		verifSvc: verifSvc,
	}
}

// RequestVerification requests a new email verification token
func (h *VerificationHandler) RequestVerification(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.EmailVerified() {
		http.Error(w, "Email already verified", http.StatusBadRequest)
		return
	}

	token, err := h.verifSvc.RequestEmailVerification(user.ID().Value())
	if err != nil {
		http.Error(w, "Failed to create verification token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Verification email sent",
		"token":   token, // In production, don't return token - only send via email
	})
}

// VerifyEmail verifies an email using a token
func (h *VerificationHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	if err := h.verifSvc.VerifyEmail(token); err != nil {
		if err == application.ErrInvalidToken {
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to verify email", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email verified successfully",
	})
}

// RequestPasswordReset requests a password reset
func (h *VerificationHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.verifSvc.RequestPasswordReset(req.Email)
	if err != nil {
		http.Error(w, "Failed to request password reset", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset email sent if account exists",
		"token":   token, // In production, don't return token - only send via email
	})
}

// ResetPassword resets a password using a token
func (h *VerificationHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		http.Error(w, "Token and new password required", http.StatusBadRequest)
		return
	}

	if err := h.verifSvc.ResetPassword(req.Token, req.NewPassword); err != nil {
		if err == application.ErrInvalidToken {
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset successfully",
	})
}
