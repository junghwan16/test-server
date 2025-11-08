package persistence

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/junghwan16/test-server/internal/auth/domain"
)

type UserModel struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Username     string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (UserModel) TableName() string {
	return "users"
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) (*GormUserRepository, error) {
	if err := db.AutoMigrate(&UserModel{}); err != nil {
		return nil, err
	}
	return &GormUserRepository{db: db}, nil
}

func (r *GormUserRepository) FindByID(id domain.UserID) (*domain.User, error) {
	var model UserModel
	err := r.db.First(&model, id.Value()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return r.toDomain(&model)
}

func (r *GormUserRepository) FindByUsername(username domain.Username) (*domain.User, error) {
	var model UserModel
	err := r.db.First(&model, "username = ?", username.String()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return r.toDomain(&model)
}

func (r *GormUserRepository) Save(user *domain.User) error {
	model := r.toModel(user)
	result := r.db.Save(model)
	return result.Error
}

func (r *GormUserRepository) Update(user *domain.User) error {
	model := r.toModel(user)
	result := r.db.Model(&UserModel{}).Where("id = ?", model.ID).Updates(map[string]interface{}{
		"username":      model.Username,
		"password_hash": model.PasswordHash,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *GormUserRepository) Delete(id domain.UserID) error {
	result := r.db.Delete(&UserModel{}, id.Value())
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *GormUserRepository) ExistsByUsername(username domain.Username) (bool, error) {
	var count int64
	err := r.db.Model(&UserModel{}).Where("username = ?", username.String()).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormUserRepository) toDomain(model *UserModel) (*domain.User, error) {
	userID := domain.NewUserID(model.ID)

	username, err := domain.NewUsername(model.Username)
	if err != nil {
		return nil, err
	}

	return domain.ReconstructUser(userID, username, model.PasswordHash)
}

func (r *GormUserRepository) toModel(user *domain.User) *UserModel {
	return &UserModel{
		ID:           user.ID().Value(),
		Username:     user.Username().String(),
		PasswordHash: user.PasswordHash(),
	}
}
