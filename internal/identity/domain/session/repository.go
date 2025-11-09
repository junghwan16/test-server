package session

import (
	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

// Repository defines the interface for Session aggregate persistence
type Repository interface {
	Save(session *Session) error

	// FindByID retrieves a Session by ID
	FindByID(id SessionID) (*Session, error)

	Delete(id SessionID) error

	DeleteExpired() error

	DeleteByUserID(userID user.UserID) error
}
