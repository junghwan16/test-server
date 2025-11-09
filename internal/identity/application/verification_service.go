package application

import (
	"errors"
	"time"

	"github.com/junghwan16/test-server/internal/identity/domain/user"
	"github.com/junghwan16/test-server/internal/identity/domain/verification"
)

var (
	ErrInvalidToken    = errors.New("invalid or expired token")
	ErrAlreadyVerified = errors.New("email already verified")
)

// VerificationService handles email verification and password reset
type VerificationService struct {
	userRepo          user.Repository
	emailVerifRepo    verification.EmailVerificationRepository
	passwordResetRepo verification.PasswordResetRepository
	verificationTTL   time.Duration
	passwordResetTTL  time.Duration
}

// NewVerificationService creates a new VerificationService
func NewVerificationService(
	userRepo user.Repository,
	emailVerifRepo verification.EmailVerificationRepository,
	passwordResetRepo verification.PasswordResetRepository,
	verificationTTL time.Duration,
	passwordResetTTL time.Duration,
) *VerificationService {
	return &VerificationService{
		userRepo:          userRepo,
		emailVerifRepo:    emailVerifRepo,
		passwordResetRepo: passwordResetRepo,
		verificationTTL:   verificationTTL,
		passwordResetTTL:  passwordResetTTL,
	}
}

// RequestEmailVerification creates a verification token for a user
func (s *VerificationService) RequestEmailVerification(userID uint) (string, error) {
	uid, err := user.NewUserID(userID)
	if err != nil {
		return "", err
	}

	u, err := s.userRepo.FindByID(uid)
	if err != nil {
		return "", ErrUserNotFound
	}

	if u.EmailVerified() {
		return "", ErrAlreadyVerified
	}

	verif := verification.NewEmailVerification(u.ID(), s.verificationTTL)

	if err := s.emailVerifRepo.Save(verif); err != nil {
		return "", err
	}

	return verif.Token(), nil
}

// VerifyEmail verifies an email using a token
func (s *VerificationService) VerifyEmail(token string) error {
	verif, err := s.emailVerifRepo.FindByToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	if verif.IsExpired() {
		return ErrInvalidToken
	}

	u, err := s.userRepo.FindByID(verif.UserID())
	if err != nil {
		return ErrUserNotFound
	}

	if err := u.VerifyEmail(); err != nil {
		return err
	}

	if err := s.userRepo.Save(u); err != nil {
		return err
	}

	return s.emailVerifRepo.Delete(token)
}

// RequestPasswordReset creates a password reset token
func (s *VerificationService) RequestPasswordReset(email string) (string, error) {
	emailVO, err := user.NewEmail(email)
	if err != nil {
		return "", err
	}

	u, err := s.userRepo.FindByEmail(emailVO)
	if err != nil {
		// Don't reveal if email exists or not
		return "", nil
	}

	if !u.Active() {
		return "", nil
	}

	s.passwordResetRepo.DeleteByUserID(u.ID())

	reset := verification.NewPasswordReset(u.ID(), s.passwordResetTTL)

	if err := s.passwordResetRepo.Save(reset); err != nil {
		return "", err
	}

	return reset.Token(), nil
}

// ResetPassword resets a password using a token
func (s *VerificationService) ResetPassword(token, newPassword string) error {
	reset, err := s.passwordResetRepo.FindByToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	if reset.IsExpired() {
		return ErrInvalidToken
	}

	u, err := s.userRepo.FindByID(reset.UserID())
	if err != nil {
		return ErrUserNotFound
	}

	newPass, err := user.NewPassword(newPassword)
	if err != nil {
		return err
	}

	if err := u.ChangePassword(newPass); err != nil {
		return err
	}

	if err := s.userRepo.Save(u); err != nil {
		return err
	}

	s.passwordResetRepo.Delete(token)

	return s.passwordResetRepo.DeleteByUserID(u.ID())
}
