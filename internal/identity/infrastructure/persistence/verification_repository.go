package persistence

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/junghwan16/test-server/internal/identity/domain/user"
	"github.com/junghwan16/test-server/internal/identity/domain/verification"
)

// EmailVerificationModel is the GORM model
type EmailVerificationModel struct {
	Token     string    `gorm:"primarykey"`
	UserID    uint      `gorm:"index;not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	CreatedAt time.Time
}

func (EmailVerificationModel) TableName() string {
	return "email_verifications"
}

// PasswordResetModel is the GORM model
type PasswordResetModel struct {
	Token     string    `gorm:"primarykey"`
	UserID    uint      `gorm:"index;not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	CreatedAt time.Time
}

func (PasswordResetModel) TableName() string {
	return "password_resets"
}

// EmailVerificationRepository implements the repository
type EmailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) *EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

func (r *EmailVerificationRepository) Save(v *verification.EmailVerification) error {
	model := EmailVerificationModel{
		Token:     v.Token(),
		UserID:    v.UserID().Value(),
		ExpiresAt: v.ExpiresAt(),
		CreatedAt: v.CreatedAt(),
	}
	return r.db.Create(&model).Error
}

func (r *EmailVerificationRepository) FindByToken(token string) (*verification.EmailVerification, error) {
	var model EmailVerificationModel
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("verification token not found or expired")
		}
		return nil, err
	}

	userID, _ := user.NewUserID(model.UserID)
	return verification.ReconstructEmailVerification(
		model.Token,
		userID,
		model.ExpiresAt,
		model.CreatedAt,
	), nil
}

func (r *EmailVerificationRepository) Delete(token string) error {
	return r.db.Delete(&EmailVerificationModel{}, "token = ?", token).Error
}

// PasswordResetRepository implements the repository
type PasswordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Save(p *verification.PasswordReset) error {
	model := PasswordResetModel{
		Token:     p.Token(),
		UserID:    p.UserID().Value(),
		ExpiresAt: p.ExpiresAt(),
		CreatedAt: p.CreatedAt(),
	}
	return r.db.Create(&model).Error
}

func (r *PasswordResetRepository) FindByToken(token string) (*verification.PasswordReset, error) {
	var model PasswordResetModel
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reset token not found or expired")
		}
		return nil, err
	}

	userID, _ := user.NewUserID(model.UserID)
	return verification.ReconstructPasswordReset(
		model.Token,
		userID,
		model.ExpiresAt,
		model.CreatedAt,
	), nil
}

func (r *PasswordResetRepository) Delete(token string) error {
	return r.db.Delete(&PasswordResetModel{}, "token = ?", token).Error
}

func (r *PasswordResetRepository) DeleteByUserID(userID user.UserID) error {
	return r.db.Where("user_id = ?", userID.Value()).Delete(&PasswordResetModel{}).Error
}
