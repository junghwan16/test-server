package verification

import (
	"time"

	"github.com/google/uuid"
	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

// EmailVerification is an aggregate for email verification tokens
type EmailVerification struct {
	token     string
	userID    user.UserID
	expiresAt time.Time
	createdAt time.Time
}

// NewEmailVerification creates a new email verification token
func NewEmailVerification(userID user.UserID, ttl time.Duration) *EmailVerification {
	return &EmailVerification{
		token:     uuid.New().String(),
		userID:    userID,
		expiresAt: time.Now().Add(ttl),
		createdAt: time.Now(),
	}
}

// ReconstructEmailVerification reconstructs from persistence
func ReconstructEmailVerification(token string, userID user.UserID, expiresAt, createdAt time.Time) *EmailVerification {
	return &EmailVerification{
		token:     token,
		userID:    userID,
		expiresAt: expiresAt,
		createdAt: createdAt,
	}
}

// Getters
func (e *EmailVerification) Token() string        { return e.token }
func (e *EmailVerification) UserID() user.UserID  { return e.userID }
func (e *EmailVerification) ExpiresAt() time.Time { return e.expiresAt }
func (e *EmailVerification) CreatedAt() time.Time { return e.createdAt }

// IsExpired returns true if the token is expired
func (e *EmailVerification) IsExpired() bool {
	return time.Now().After(e.expiresAt)
}

// IsValid returns true if the token is still valid
func (e *EmailVerification) IsValid() bool {
	return !e.IsExpired()
}
