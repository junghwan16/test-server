package user

import (
	"errors"
	"strings"
)

// Email is a value object representing an email address
type Email struct {
	value string
}

var ErrInvalidEmail = errors.New("invalid email address")

// NewEmail creates a new Email value object
func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(strings.ToLower(value))

	if value == "" {
		return Email{}, ErrInvalidEmail
	}

	// Basic validation
	if !strings.Contains(value, "@") || !strings.Contains(value, ".") {
		return Email{}, ErrInvalidEmail
	}

	parts := strings.Split(value, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: value}, nil
}

// MustNewEmail creates an Email or panics
func MustNewEmail(value string) Email {
	email, err := NewEmail(value)
	if err != nil {
		panic(err)
	}
	return email
}

// Value returns the email address
func (e Email) Value() string {
	return e.value
}

// Equals checks if two emails are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// String returns string representation
func (e Email) String() string {
	return e.value
}
