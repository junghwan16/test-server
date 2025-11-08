package domain

import "time"

// Session은 인증된 사용자의 활성 세션을 나타내는 도메인 개념입니다.
type Session struct {
	userID    UserID
	username  Username
	createdAt time.Time
	expiresAt time.Time
}

func NewSession(userID UserID, username Username, expiresAt time.Time) *Session {
	return &Session{
		userID:    userID,
		username:  username,
		createdAt: time.Now(),
		expiresAt: expiresAt,
	}
}

func (s *Session) UserID() UserID {
	return s.userID
}

func (s *Session) Username() Username {
	return s.username
}

func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

func (s *Session) IsValid() bool {
	return !s.IsExpired()
}
