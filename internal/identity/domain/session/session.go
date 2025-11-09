package session

import (
	"time"

	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

// Session is the aggregate root for user sessions
type Session struct {
	id        SessionID
	userID    user.UserID
	expiresAt time.Time
	createdAt time.Time
}

// NewSession creates a new Session aggregate
func NewSession(id SessionID, userID user.UserID, ttl int) *Session {
	return &Session{
		id:        id,
		userID:    userID,
		expiresAt: time.Now().Add(time.Duration(ttl) * time.Second),
		createdAt: time.Now(),
	}
}

// ReconstructSession reconstructs a Session from persistence
func ReconstructSession(
	id SessionID,
	userID user.UserID,
	expiresAt, createdAt time.Time,
) *Session {
	return &Session{
		id:        id,
		userID:    userID,
		expiresAt: expiresAt,
		createdAt: createdAt,
	}
}

// Getters

func (s *Session) ID() SessionID        { return s.id }
func (s *Session) UserID() user.UserID  { return s.userID }
func (s *Session) ExpiresAt() time.Time { return s.expiresAt }
func (s *Session) CreatedAt() time.Time { return s.createdAt }

// Business methods

// IsExpired returns true if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

// IsValid returns true if the session is still valid
func (s *Session) IsValid() bool {
	return !s.IsExpired()
}
