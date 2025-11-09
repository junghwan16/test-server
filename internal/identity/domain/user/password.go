package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12
const minPasswordLength = 8

// Password is a value object representing a hashed password
type Password struct {
	hash string
}

var ErrPasswordTooShort = errors.New("password must be at least 8 characters")
var ErrInvalidPassword = errors.New("invalid password")

// NewPassword creates a new Password from plaintext
func NewPassword(plaintext string) (Password, error) {
	if len(plaintext) < minPasswordLength {
		return Password{}, ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcryptCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: string(hash)}, nil
}

// NewPasswordFromHash creates a Password from existing hash
func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

// Hash returns the bcrypt hash
func (p Password) Hash() string {
	return p.hash
}

// Matches checks if the plaintext matches the hashed password
func (p Password) Matches(plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plaintext))
	return err == nil
}

// Change creates a new Password with different plaintext
func (p Password) Change(newPlaintext string) (Password, error) {
	return NewPassword(newPlaintext)
}
