package application

import (
	"log/slog"
	"time"

	"github.com/junghwan16/test-server/internal/auth/domain"
)

type AuthService struct {
	userRepo      domain.UserRepository
	logger        *slog.Logger
	sessionExpiry time.Duration
}

func NewAuthService(userRepo domain.UserRepository, logger *slog.Logger, sessionExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		logger:        logger,
		sessionExpiry: sessionExpiry,
	}
}

func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	username, err := domain.NewUsername(req.Username)
	if err != nil {
		s.logger.Warn("invalid username format", "error", err)
		return nil, err
	}

	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		s.logger.Warn("user not found", "username", req.Username)
		return nil, domain.ErrInvalidCredentials
	}

	if err := user.Authenticate(req.Password); err != nil {
		s.logger.Warn("authentication failed", "username", req.Username, "error", err)
		return nil, domain.ErrInvalidCredentials
	}

	expiresAt := time.Now().Add(s.sessionExpiry)
	session := domain.NewSession(user.ID(), user.Username(), expiresAt)

	s.logger.Info("user authenticated successfully",
		"username", req.Username,
		"userID", user.ID().String())

	return &LoginResponse{Session: session}, nil
}

func (s *AuthService) Register(req RegisterRequest) (*RegisterResponse, error) {
	username, err := domain.NewUsername(req.Username)
	if err != nil {
		s.logger.Warn("invalid username", "error", err)
		return nil, err
	}

	exists, err := s.userRepo.ExistsByUsername(username)
	if err != nil {
		s.logger.Error("failed to check username existence", "error", err)
		return nil, err
	}
	if exists {
		s.logger.Warn("username already exists", "username", req.Username)
		return nil, domain.ErrUsernameAlreadyExists
	}

	userID := domain.NewUserID(0)

	user, err := domain.NewUser(userID, username, req.Password)
	if err != nil {
		s.logger.Warn("failed to create user", "error", err)
		return nil, err
	}

	if err := s.userRepo.Save(user); err != nil {
		s.logger.Error("failed to save user", "error", err)
		return nil, err
	}

	savedUser, err := s.userRepo.FindByUsername(username)
	if err != nil {
		s.logger.Error("failed to retrieve created user", "error", err)
		return nil, err
	}

	expiresAt := time.Now().Add(s.sessionExpiry)
	session := domain.NewSession(savedUser.ID(), savedUser.Username(), expiresAt)

	s.logger.Info("user registered successfully",
		"username", req.Username,
		"userID", savedUser.ID().String())

	return &RegisterResponse{
		UserID:    savedUser.ID().Value(),
		Username:  savedUser.Username().String(),
		Session:   session,
		ExpiresAt: expiresAt,
	}, nil
}
