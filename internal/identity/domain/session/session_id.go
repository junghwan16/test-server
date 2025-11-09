package session

import (
	"errors"

	"github.com/google/uuid"
)

// SessionID is a value object representing a session's unique identifier
type SessionID struct {
	value string
}

var ErrInvalidSessionID = errors.New("invalid session ID")

// NewSessionID creates a new SessionID from string
func NewSessionID(value string) (SessionID, error) {
	if value == "" {
		return SessionID{}, ErrInvalidSessionID
	}
	return SessionID{value: value}, nil
}

// GenerateSessionID generates a new random SessionID
func GenerateSessionID() SessionID {
	return SessionID{value: uuid.New().String()}
}

// Value returns the underlying value
func (id SessionID) Value() string {
	return id.value
}

// Equals checks if two SessionIDs are equal
func (id SessionID) Equals(other SessionID) bool {
	return id.value == other.value
}

// String returns string representation
func (id SessionID) String() string {
	return id.value
}

// IsZero returns true if the ID is zero value
func (id SessionID) IsZero() bool {
	return id.value == ""
}
