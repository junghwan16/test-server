package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashedPassword는 해시된 비밀번호를 나타내는 Value Object입니다.
type HashedPassword struct {
	hash string
}

// NewHashedPasswordFromPlain은 평문 비밀번호로부터 해시된 비밀번호를 생성합니다.
func NewHashedPasswordFromPlain(plainPassword string) (HashedPassword, error) {
	if len(plainPassword) < 6 {
		return HashedPassword{}, errors.New("password must be at least 6 characters long")
	}

	if len(plainPassword) > 72 {
		// bcrypt has a maximum password length of 72 bytes
		return HashedPassword{}, errors.New("password cannot exceed 72 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return HashedPassword{}, err
	}

	return HashedPassword{hash: string(hash)}, nil
}

// NewHashedPasswordFromHash는 이미 해시된 값으로부터 HashedPassword를 생성합니다.
func NewHashedPasswordFromHash(hash string) (HashedPassword, error) {
	if hash == "" {
		return HashedPassword{}, errors.New("password hash cannot be empty")
	}
	return HashedPassword{hash: hash}, nil
}

// Verify는 평문 비밀번호가 해시된 비밀번호와 일치하는지 검증합니다.
func (p HashedPassword) Verify(plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plainPassword))
}

// Hash는 저장을 위한 해시 문자열을 반환합니다.
func (p HashedPassword) Hash() string {
	return p.hash
}
