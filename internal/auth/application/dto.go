package application

import (
	"time"

	"github.com/junghwan16/test-server/internal/auth/domain"
)

type LoginRequest struct {
	Username  string
	Password  string
	IPAddress string // 감사 로깅을 위한 IP 주소
	UserAgent string // 감사 로깅을 위한 User-Agent
}

type LoginResponse struct {
	Session *domain.Session
}

type RegisterRequest struct {
	Username  string
	Password  string
	IPAddress string // 감사 로깅을 위한 IP 주소
}

type RegisterResponse struct {
	UserID    uint
	Username  string
	Session   *domain.Session
	ExpiresAt time.Time
}

type UserProfileResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
