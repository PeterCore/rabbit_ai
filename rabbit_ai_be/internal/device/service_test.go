package device

import (
	"strings"
	"testing"

	"rabbit_ai/internal/model"
)

// MockUserRepository 模拟用户仓库
type MockUserRepository struct {
	users      map[int64]*model.User
	byDeviceID map[string]*model.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:      make(map[int64]*model.User),
		byDeviceID: make(map[string]*model.User),
	}
}

func (m *MockUserRepository) Create(user *model.User) error {
	user.ID = int64(len(m.users) + 1)
	m.users[user.ID] = user
	if user.DeviceID != "" {
		m.byDeviceID[user.DeviceID] = user
	}
	return nil
}

func (m *MockUserRepository) GetByID(id int64) (*model.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, model.ErrUserNotFound
}

func (m *MockUserRepository) GetByPhone(phone string) (*model.User, error) {
	for _, user := range m.users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return nil, model.ErrUserNotFound
}

func (m *MockUserRepository) GetByGitHubID(githubID string) (*model.User, error) {
	for _, user := range m.users {
		if user.GitHubID == githubID {
			return user, nil
		}
	}
	return nil, model.ErrUserNotFound
}

func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, model.ErrUserNotFound
}

func (m *MockUserRepository) GetByDeviceID(deviceID string) (*model.User, error) {
	if user, exists := m.byDeviceID[deviceID]; exists {
		return user, nil
	}
	return nil, model.ErrUserNotFound
}

func (m *MockUserRepository) Update(user *model.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return model.ErrUserNotFound
	}
	m.users[user.ID] = user
	if user.DeviceID != "" {
		m.byDeviceID[user.DeviceID] = user
	}
	return nil
}

func (m *MockUserRepository) Delete(id int64) error {
	if user, exists := m.users[id]; exists {
		delete(m.users, id)
		if user.DeviceID != "" {
			delete(m.byDeviceID, user.DeviceID)
		}
		return nil
	}
	return model.ErrUserNotFound
}

func (m *MockUserRepository) CreateWithPassword(user *model.User, password string) error {
	return m.Create(user)
}

func (m *MockUserRepository) VerifyPassword(phone, password string) (*model.User, error) {
	return m.GetByPhone(phone)
}

func (m *MockUserRepository) UpdatePassword(userID int64, newPassword string) error {
	if _, exists := m.users[userID]; !exists {
		return model.ErrUserNotFound
	}
	return nil
}

func TestDeviceService_GetOrCreateUserByDeviceID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewDeviceService(mockRepo)

	deviceID := "test-device-123"
	platform := "ios"

	// 测试创建新用户
	user, err := service.GetOrCreateUserByDeviceID(deviceID, platform)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("Expected user to be created")
	}

	if user.DeviceID != deviceID {
		t.Errorf("Expected device ID %s, got %s", deviceID, user.DeviceID)
	}

	if user.Platform != platform {
		t.Errorf("Expected platform %s, got %s", platform, user.Platform)
	}

	// 检查昵称是否以"设备用户_"开头
	if !strings.HasPrefix(user.Nickname, "设备用户_") {
		t.Errorf("Expected nickname to start with '设备用户_', got %s", user.Nickname)
	}

	// 测试获取已存在的用户
	existingUser, err := service.GetOrCreateUserByDeviceID(deviceID, platform)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if existingUser.ID != user.ID {
		t.Errorf("Expected same user ID, got %d vs %d", existingUser.ID, user.ID)
	}

	// 检查是否为新建用户
	if !strings.HasPrefix(user.Nickname, "设备用户_") {
		t.Errorf("Expected nickname to start with '设备用户_', got %s", user.Nickname)
	}
}

func TestDeviceService_GetUserByDeviceID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewDeviceService(mockRepo)

	deviceID := "test-device-456"

	// 测试获取不存在的用户
	_, err := service.GetUserByDeviceID(deviceID)
	if err == nil {
		t.Fatal("Expected error for non-existent user")
	}

	// 创建用户
	user := &model.User{
		DeviceID: deviceID,
		Nickname: "Test User",
	}
	mockRepo.Create(user)

	// 测试获取存在的用户
	foundUser, err := service.GetUserByDeviceID(deviceID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if foundUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, foundUser.ID)
	}
}

func TestDeviceService_BindDeviceToUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewDeviceService(mockRepo)

	// 创建用户
	user := &model.User{
		Nickname: "Test User",
	}
	mockRepo.Create(user)

	deviceID := "test-device-789"

	// 测试绑定设备
	err := service.BindDeviceToUser(deviceID, user.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 验证设备已绑定
	updatedUser, err := mockRepo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedUser.DeviceID != deviceID {
		t.Errorf("Expected device ID %s, got %s", deviceID, updatedUser.DeviceID)
	}
}

func TestDeviceService_UnbindDevice(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewDeviceService(mockRepo)

	// 创建带设备ID的用户
	user := &model.User{
		DeviceID: "test-device-999",
		Nickname: "Test User",
	}
	mockRepo.Create(user)

	// 测试解绑设备
	err := service.UnbindDevice(user.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 验证设备已解绑
	updatedUser, err := mockRepo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedUser.DeviceID != "" {
		t.Errorf("Expected empty device ID, got %s", updatedUser.DeviceID)
	}
}
