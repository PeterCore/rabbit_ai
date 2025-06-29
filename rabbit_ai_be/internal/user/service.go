package user

import (
	"fmt"

	"rabbit_ai/internal/model"
)

// UserService 用户服务
type UserService struct {
	userRepo model.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo model.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(userID int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID int64, nickname, avatar string) (*model.User, error) {
	// 先获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 更新用户信息
	if nickname != "" {
		user.Nickname = nickname
	}
	if avatar != "" {
		user.Avatar = avatar
	}

	// 保存到数据库
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// UpdatePassword 修改密码
func (s *UserService) UpdatePassword(userID int64, oldPassword, newPassword string) error {
	// 先获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 验证旧密码
	if user.Password != "" {
		_, err = s.userRepo.VerifyPassword(user.Phone, oldPassword)
		if err != nil {
			return fmt.Errorf("invalid old password: %w", err)
		}
	}

	// 更新密码
	err = s.userRepo.UpdatePassword(userID, newPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID int64) error {
	err := s.userRepo.Delete(userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
