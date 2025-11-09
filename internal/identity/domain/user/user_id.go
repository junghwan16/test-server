package user

import (
	"errors"
)

// UserID is a value object representing a user's unique identifier
type UserID struct {
	value uint
}

var ErrInvalidUserID = errors.New("invalid user ID")

// NewUserID creates a new UserID
func NewUserID(value uint) (UserID, error) {
	if value == 0 {
		return UserID{}, ErrInvalidUserID
	}
	return UserID{value: value}, nil
}

// MustNewUserID creates a UserID or panics
func MustNewUserID(value uint) UserID {
	id, err := NewUserID(value)
	if err != nil {
		panic(err)
	}
	return id
}

// Value returns the underlying value
func (id UserID) Value() uint {
	return id.value
}

// Equals checks if two UserIDs are equal
func (id UserID) Equals(other UserID) bool {
	return id.value == other.value
}

// String returns string representation
func (id UserID) String() string {
	return string(rune(id.value))
}

// IsZero returns true if the ID is zero value
func (id UserID) IsZero() bool {
	return id.value == 0
}
