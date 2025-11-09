package application

import (
	"errors"

	"github.com/junghwan16/test-server/internal/identity/domain/session"
	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

// AuthService handles authentication logic
type AuthService struct {
	userRepo    user.Repository
	sessionRepo session.Repository
	sessionTTL  int
}

// NewAuthService creates a new AuthService
func NewAuthService(
	userRepo user.Repository,
	sessionRepo session.Repository,
	sessionTTL int,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  sessionTTL,
	}
}

// Login authenticates a user and creates a session
func (s *AuthService) Login(email, password string) (*session.Session, *user.User, error) {
	emailVO, err := user.NewEmail(email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	u, err := s.userRepo.FindByEmail(emailVO)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if !u.Authenticate(password) {
		return nil, nil, ErrInvalidCredentials
	}

	// Create session
	sess := session.NewSession(
		session.GenerateSessionID(),
		u.ID(),
		s.sessionTTL,
	)

	if err := s.sessionRepo.Save(sess); err != nil {
		return nil, nil, err
	}

	return sess, u, nil
}

// ValidateSession validates a session
func (s *AuthService) ValidateSession(sessionID string) (*session.Session, *user.User, error) {
	sid, err := session.NewSessionID(sessionID)
	if err != nil {
		return nil, nil, errors.New("invalid session")
	}

	sess, err := s.sessionRepo.FindByID(sid)
	if err != nil {
		return nil, nil, errors.New("invalid session")
	}

	if sess.IsExpired() {
		return nil, nil, errors.New("session expired")
	}

	u, err := s.userRepo.FindByID(sess.UserID())
	if err != nil {
		return nil, nil, err
	}

	return sess, u, nil
}

// Logout destroys a session
func (s *AuthService) Logout(sessionID string) error {
	sid, _ := session.NewSessionID(sessionID)
	return s.sessionRepo.Delete(sid)
}
