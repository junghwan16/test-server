package persistence

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/junghwan16/test-server/internal/identity/domain/user"
	"github.com/junghwan16/test-server/internal/shared/domain"
)

// UserModel is the GORM model for User aggregate
type UserModel struct {
	ID            uint   `gorm:"primarykey"`
	Email         string `gorm:"uniqueIndex;not null"`
	PasswordHash  string `gorm:"not null"`
	Role          string `gorm:"not null;default:user"`
	EmailVerified bool   `gorm:"not null;default:false"`
	Active        bool   `gorm:"not null;default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName specifies the table name
func (UserModel) TableName() string {
	return "users"
}

// UserRepository implements user.Repository using GORM
type UserRepository struct {
	db       *gorm.DB
	eventBus domain.EventBus
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB, eventBus domain.EventBus) *UserRepository {
	return &UserRepository{
		db:       db,
		eventBus: eventBus,
	}
}

// NextID generates a new UserID
func (r *UserRepository) NextID() user.UserID {
	// In auto-increment scenario, return zero and let DB generate
	return user.UserID{}
}

func (r *UserRepository) Save(u *user.User) error {
	model := r.toModel(u)

	var err error
	if model.ID == 0 {
		err = r.db.Create(&model).Error
		if err == nil && !u.ID().IsZero() {
			// Update aggregate with generated ID
			// Note: This is a compromise in DDD - ideally ID should be known before persistence
		}
	} else {
		err = r.db.Save(&model).Error
	}

	if err != nil {
		return err
	}

	// Publish domain events
	for _, event := range u.DomainEvents() {
		_ = r.eventBus.Publish(event)
	}
	u.ClearEvents()

	return nil
}

// FindByID retrieves a User by ID
func (r *UserRepository) FindByID(id user.UserID) (*user.User, error) {
	var model UserModel
	err := r.db.First(&model, id.Value()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return r.toDomain(&model), nil
}

// FindByEmail retrieves a User by email
func (r *UserRepository) FindByEmail(email user.Email) (*user.User, error) {
	var model UserModel
	err := r.db.Where("email = ?", email.Value()).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return r.toDomain(&model), nil
}

// FindAll retrieves all users with pagination
func (r *UserRepository) FindAll(limit, offset int) ([]*user.User, int64, error) {
	var models []UserModel
	var total int64

	if err := r.db.Model(&UserModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	users := make([]*user.User, len(models))
	for i, model := range models {
		users[i] = r.toDomain(&model)
	}

	return users, total, nil
}

func (r *UserRepository) Delete(id user.UserID) error {
	return r.db.Delete(&UserModel{}, id.Value()).Error
}

// Mapping functions

func (r *UserRepository) toModel(u *user.User) UserModel {
	return UserModel{
		ID:            u.ID().Value(),
		Email:         u.Email().Value(),
		PasswordHash:  u.Password().Hash(),
		Role:          u.Role().Value(),
		EmailVerified: u.EmailVerified(),
		Active:        u.Active(),
		CreatedAt:     u.CreatedAt(),
		UpdatedAt:     u.UpdatedAt(),
	}
}

func (r *UserRepository) toDomain(m *UserModel) *user.User {
	id, _ := user.NewUserID(m.ID)
	email, _ := user.NewEmail(m.Email)
	password := user.NewPasswordFromHash(m.PasswordHash)
	role, _ := user.NewRole(m.Role)

	return user.ReconstructUser(
		id,
		email,
		password,
		role,
		m.EmailVerified,
		m.Active,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
