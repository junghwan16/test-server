package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/junghwan16/test-server/internal/auth/application"
)

// UserHandler는 사용자 관련 HTTP 요청을 처리합니다.
type UserHandler struct {
	userService *application.UserService
	logger      *slog.Logger
}

// NewUserHandler는 새로운 사용자 핸들러를 생성합니다.
func NewUserHandler(userService *application.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetProfile은 본인 프로필을 조회합니다.
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		h.logger.Error("failed to get profile", "user_id", userID, "error", err)
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfile은 본인 프로필을 수정합니다.
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req application.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdateProfile(userID, req); err != nil {
		h.logger.Error("failed to update profile", "user_id", userID, "error", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ChangePassword는 비밀번호를 변경합니다.
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req application.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.userService.ChangePassword(userID, req); err != nil {
		h.logger.Error("failed to change password", "user_id", userID, "error", err)
		http.Error(w, "Failed to change password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteAccount는 계정을 삭제합니다.
func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.userService.DeleteAccount(userID); err != nil {
		h.logger.Error("failed to delete account", "user_id", userID, "error", err)
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
