package application

import (
	"log/slog"

	"github.com/junghwan16/test-server/internal/auth/domain"
)

// UserService는 사용자 정보 관리를 담당합니다.
type UserService struct {
	userRepo domain.UserRepository
	logger   *slog.Logger
}

// NewUserService는 새로운 사용자 서비스를 생성합니다.
func NewUserService(userRepo domain.UserRepository, logger *slog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetProfile은 사용자 프로필을 조회합니다.
func (s *UserService) GetProfile(userID uint) (*UserProfileResponse, error) {
	id := domain.NewUserID(userID)
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Error("failed to find user", "user_id", userID, "error", err)
		return nil, err
	}

	return &UserProfileResponse{
		UserID:   user.ID().Value(),
		Username: user.Username().String(),
	}, nil
}

// UpdateProfile은 사용자 프로필을 수정합니다.
func (s *UserService) UpdateProfile(userID uint, req UpdateProfileRequest) error {
	id := domain.NewUserID(userID)
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Error("failed to find user", "user_id", userID, "error", err)
		return err
	}

	// Username 변경
	newUsername, err := domain.NewUsername(req.Username)
	if err != nil {
		return err
	}

	// 새 username이 이미 사용 중인지 확인
	existingUser, err := s.userRepo.FindByUsername(newUsername)
	if err == nil && existingUser.ID().Value() != userID {
		return domain.ErrUsernameTaken
	}

	user.ChangeUsername(newUsername)

	if err := s.userRepo.Update(user); err != nil {
		s.logger.Error("failed to update user", "user_id", userID, "error", err)
		return err
	}

	s.logger.Info("user profile updated", "user_id", userID)
	return nil
}

// ChangePassword는 사용자 비밀번호를 변경합니다.
func (s *UserService) ChangePassword(userID uint, req ChangePasswordRequest) error {
	id := domain.NewUserID(userID)
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Error("failed to find user", "user_id", userID, "error", err)
		return err
	}

	// 현재 비밀번호 확인
	if !user.VerifyPassword(req.CurrentPassword) {
		return domain.ErrInvalidPassword
	}

	// 새 비밀번호로 변경
	if err := user.ChangePassword(req.NewPassword); err != nil {
		return err
	}

	if err := s.userRepo.Update(user); err != nil {
		s.logger.Error("failed to update password", "user_id", userID, "error", err)
		return err
	}

	s.logger.Info("password changed", "user_id", userID)
	return nil
}

// DeleteAccount는 계정을 삭제합니다.
func (s *UserService) DeleteAccount(userID uint) error {
	id := domain.NewUserID(userID)

	if err := s.userRepo.Delete(id); err != nil {
		s.logger.Error("failed to delete user", "user_id", userID, "error", err)
		return err
	}

	s.logger.Info("user deleted", "user_id", userID)
	return nil
}
