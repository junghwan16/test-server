package verification

import (
	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

// EmailVerificationRepository defines the interface for email verification persistence
type EmailVerificationRepository interface {
	Save(verification *EmailVerification) error
	FindByToken(token string) (*EmailVerification, error)
	Delete(token string) error
}

// PasswordResetRepository defines the interface for password reset persistence
type PasswordResetRepository interface {
	Save(reset *PasswordReset) error
	FindByToken(token string) (*PasswordReset, error)
	Delete(token string) error
	DeleteByUserID(userID user.UserID) error
}
