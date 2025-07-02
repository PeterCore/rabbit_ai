package device

import (
	"fmt"
	"log"

	"rabbit_ai/internal/model"
)

// DeviceService 设备服务
type DeviceService struct {
	userRepo model.UserRepository
}

// NewDeviceService 创建设备服务实例
func NewDeviceService(userRepo model.UserRepository) *DeviceService {
	return &DeviceService{
		userRepo: userRepo,
	}
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	DeviceID string      `json:"device_id"`
	User     *model.User `json:"user,omitempty"`
	Exists   bool        `json:"exists"`
}

// GetOrCreateUserByDeviceID 根据设备ID和平台获取或创建用户
func (s *DeviceService) GetOrCreateUserByDeviceID(deviceID, platform string) (*model.User, error) {
	// 先尝试根据设备ID查找用户
	user, err := s.userRepo.GetByDeviceID(deviceID)
	if err == nil && user != nil {
		// 用户存在，返回用户信息
		return user, nil
	}

	// 用户不存在，创建新用户
	user = &model.User{
		DeviceID: deviceID,
		Platform: platform,
		Nickname: "设备用户_" + deviceID[:8], // 使用设备ID前8位作为默认昵称
		Avatar:   "",                     // 默认头像
		Status:   1,                      // 正常状态
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user with device ID: %w", err)
	}

	log.Printf("Created new user with device ID: %s, platform: %s, user ID: %d", deviceID, platform, user.ID)
	return user, nil
}

// GetUserByDeviceID 根据设备ID获取用户
func (s *DeviceService) GetUserByDeviceID(deviceID string) (*model.User, error) {
	user, err := s.userRepo.GetByDeviceID(deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by device ID: %w", err)
	}
	return user, nil
}

// UpdateUserDeviceID 更新用户的设备ID
func (s *DeviceService) UpdateUserDeviceID(userID int64, deviceID string) error {
	// 先获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 更新设备ID
	user.DeviceID = deviceID

	// 保存到数据库
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user device ID: %w", err)
	}

	log.Printf("Updated device ID for user %d: %s", userID, deviceID)
	return nil
}

// BindDeviceToUser 将设备绑定到指定用户
func (s *DeviceService) BindDeviceToUser(deviceID string, userID int64) error {
	// 检查设备ID是否已被其他用户使用
	existingUser, err := s.userRepo.GetByDeviceID(deviceID)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		return fmt.Errorf("device ID %s is already bound to another user", deviceID)
	}

	// 更新用户的设备ID
	return s.UpdateUserDeviceID(userID, deviceID)
}

// UnbindDevice 解绑设备
func (s *DeviceService) UnbindDevice(userID int64) error {
	// 先获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 清空设备ID
	user.DeviceID = ""

	// 保存到数据库
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to unbind device: %w", err)
	}

	log.Printf("Unbound device from user %d", userID)
	return nil
}
