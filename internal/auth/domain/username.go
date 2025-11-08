package domain

import (
	"errors"
	"strings"
)

// Username은 사용자명을 나타내는 Value Object입니다.
type Username struct {
	value string
}

// NewUsername은 새로운 Username을 생성하고 유효성을 검증합니다.
func NewUsername(value string) (Username, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Username{}, errors.New("username cannot be empty")
	}

	if len(value) < 3 {
		return Username{}, errors.New("username must be at least 3 characters long")
	}

	if len(value) > 50 {
		return Username{}, errors.New("username cannot exceed 50 characters")
	}

	return Username{value: value}, nil
}

// String은 Username의 문자열 값을 반환합니다.
func (u Username) String() string {
	return u.value
}

// Equals는 두 Username이 같은지 비교합니다.
func (u Username) Equals(other Username) bool {
	return u.value == other.value
}
