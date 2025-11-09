package application

import (
	"errors"

	"github.com/junghwan16/test-server/internal/identity/domain/user"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrCannotDeleteSelf = errors.New("cannot delete yourself")
)

// UserService handles user-related application logic
type UserService struct {
	userRepo user.Repository
}

// NewUserService creates a new UserService
func NewUserService(userRepo user.Repository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(email, password string) (*user.User, error) {
	emailVO, err := user.NewEmail(email)
	if err != nil {
		return nil, err
	}

	passwordVO, err := user.NewPassword(password)
	if err != nil {
		return nil, err
	}

	existing, _ := s.userRepo.FindByEmail(emailVO)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	id := s.userRepo.NextID()
	u, err := user.NewUser(id, emailVO, passwordVO)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Save(u); err != nil {
		return nil, err
	}

	return u, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id uint) (*user.User, error) {
	userID, err := user.NewUserID(id)
	if err != nil {
		return nil, err
	}

	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return u, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*user.User, error) {
	emailVO, err := user.NewEmail(email)
	if err != nil {
		return nil, err
	}

	u, err := s.userRepo.FindByEmail(emailVO)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return u, nil
}

// ListUsers retrieves all users with pagination
func (s *UserService) ListUsers(limit, offset int) ([]*user.User, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.userRepo.FindAll(limit, offset)
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(id uint, newPassword string) error {
	userID, err := user.NewUserID(id)
	if err != nil {
		return err
	}

	u, err := s.userRepo.FindByID(userID)
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

	return s.userRepo.Save(u)
}

// VerifyEmail marks a user's email as verified
func (s *UserService) VerifyEmail(id uint) error {
	userID, err := user.NewUserID(id)
	if err != nil {
		return err
	}

	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := u.VerifyEmail(); err != nil {
		return err
	}

	return s.userRepo.Save(u)
}

// ChangeRole changes a user's role (admin operation)
func (s *UserService) ChangeRole(id uint, roleName string) error {
	userID, err := user.NewUserID(id)
	if err != nil {
		return err
	}

	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	newRole, err := user.NewRole(roleName)
	if err != nil {
		return err
	}

	if err := u.ChangeRole(newRole); err != nil {
		return err
	}

	return s.userRepo.Save(u)
}

// SetActive sets a user's active status (admin operation)
func (s *UserService) SetActive(id uint, active bool) error {
	userID, err := user.NewUserID(id)
	if err != nil {
		return err
	}

	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if active {
		if err := u.Activate(); err != nil {
			return err
		}
	} else {
		if err := u.Deactivate(); err != nil {
			return err
		}
	}

	return s.userRepo.Save(u)
}

func (s *UserService) DeleteUser(id uint, currentUserID uint) error {
	if id == currentUserID {
		return ErrCannotDeleteSelf
	}

	userID, err := user.NewUserID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(userID)
}
