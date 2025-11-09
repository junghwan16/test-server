package verification

import (
	"time"

	"github.com/google/uuid"
	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

// PasswordReset is an aggregate for password reset tokens
type PasswordReset struct {
	token     string
	userID    user.UserID
	expiresAt time.Time
	createdAt time.Time
}

// NewPasswordReset creates a new password reset token
func NewPasswordReset(userID user.UserID, ttl time.Duration) *PasswordReset {
	return &PasswordReset{
		token:     uuid.New().String(),
		userID:    userID,
		expiresAt: time.Now().Add(ttl),
		createdAt: time.Now(),
	}
}

// ReconstructPasswordReset reconstructs from persistence
func ReconstructPasswordReset(token string, userID user.UserID, expiresAt, createdAt time.Time) *PasswordReset {
	return &PasswordReset{
		token:     token,
		userID:    userID,
		expiresAt: expiresAt,
		createdAt: createdAt,
	}
}

// Getters
func (p *PasswordReset) Token() string        { return p.token }
func (p *PasswordReset) UserID() user.UserID  { return p.userID }
func (p *PasswordReset) ExpiresAt() time.Time { return p.expiresAt }
func (p *PasswordReset) CreatedAt() time.Time { return p.createdAt }

// IsExpired returns true if the token is expired
func (p *PasswordReset) IsExpired() bool {
	return time.Now().After(p.expiresAt)
}

// IsValid returns true if the token is still valid
func (p *PasswordReset) IsValid() bool {
	return !p.IsExpired()
}
