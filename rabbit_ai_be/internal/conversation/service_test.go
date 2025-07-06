package conversation

import (
	"context"
	"testing"

	"rabbit_ai/internal/cache"
	"rabbit_ai/internal/minimax"
	"rabbit_ai/internal/model"
)

// MockConversationRepository 模拟对话仓库
type MockConversationRepository struct {
	conversations map[int64]*model.Conversation
	nextID        int64
}

func NewMockConversationRepository() *MockConversationRepository {
	return &MockConversationRepository{
		conversations: make(map[int64]*model.Conversation),
		nextID:        1,
	}
}

func (m *MockConversationRepository) Create(conversation *model.Conversation) error {
	conversation.ID = m.nextID
	m.nextID++
	m.conversations[conversation.ID] = conversation
	return nil
}

func (m *MockConversationRepository) GetByID(id int64) (*model.Conversation, error) {
	if conv, exists := m.conversations[id]; exists {
		return conv, nil
	}
	return nil, model.ErrConversationNotFound
}

func (m *MockConversationRepository) GetByUserID(userID int64, limit, offset int) ([]*model.Conversation, error) {
	var conversations []*model.Conversation
	for _, conv := range m.conversations {
		if conv.UserID == userID && conv.Status == 1 {
			conversations = append(conversations, conv)
		}
	}
	return conversations, nil
}

func (m *MockConversationRepository) Update(conversation *model.Conversation) error {
	if _, exists := m.conversations[conversation.ID]; !exists {
		return model.ErrConversationNotFound
	}
	m.conversations[conversation.ID] = conversation
	return nil
}

func (m *MockConversationRepository) Delete(id int64) error {
	if conv, exists := m.conversations[id]; exists {
		conv.Status = 0
		return nil
	}
	return model.ErrConversationNotFound
}

func (m *MockConversationRepository) GetUserConversationCount(userID int64) (int, error) {
	count := 0
	for _, conv := range m.conversations {
		if conv.UserID == userID && conv.Status == 1 {
			count++
		}
	}
	return count, nil
}

// MockMessageRepository 模拟消息仓库
type MockMessageRepository struct {
	messages map[int64]*model.Message
	nextID   int64
}

func NewMockMessageRepository() *MockMessageRepository {
	return &MockMessageRepository{
		messages: make(map[int64]*model.Message),
		nextID:   1,
	}
}

func (m *MockMessageRepository) Create(message *model.Message) error {
	message.ID = m.nextID
	m.nextID++
	m.messages[message.ID] = message
	return nil
}

func (m *MockMessageRepository) GetByID(id int64) (*model.Message, error) {
	if msg, exists := m.messages[id]; exists {
		return msg, nil
	}
	return nil, model.ErrMessageNotFound
}

func (m *MockMessageRepository) GetByConversationID(conversationID int64, limit, offset int) ([]*model.Message, error) {
	var messages []*model.Message
	for _, msg := range m.messages {
		if msg.ConversationID == conversationID {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

func (m *MockMessageRepository) GetConversationMessages(conversationID int64) ([]*model.Message, error) {
	var messages []*model.Message
	for _, msg := range m.messages {
		if msg.ConversationID == conversationID {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

func (m *MockMessageRepository) Update(message *model.Message) error {
	if _, exists := m.messages[message.ID]; !exists {
		return model.ErrMessageNotFound
	}
	m.messages[message.ID] = message
	return nil
}

func (m *MockMessageRepository) Delete(id int64) error {
	if _, exists := m.messages[id]; !exists {
		return model.ErrMessageNotFound
	}
	delete(m.messages, id)
	return nil
}

func (m *MockMessageRepository) GetConversationMessageCount(conversationID int64) (int, error) {
	count := 0
	for _, msg := range m.messages {
		if msg.ConversationID == conversationID {
			count++
		}
	}
	return count, nil
}

// MockUserRepository 模拟用户仓库
type MockUserRepository struct {
	users map[int64]*model.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[int64]*model.User),
	}
}

func (m *MockUserRepository) GetByID(id int64) (*model.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, model.ErrUserNotFound
}

func (m *MockUserRepository) Create(user *model.User) error {
	m.users[user.ID] = user
	return nil
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
	for _, user := range m.users {
		if user.DeviceID == deviceID {
			return user, nil
		}
	}
	return nil, model.ErrUserNotFound
}

func (m *MockUserRepository) Update(user *model.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return model.ErrUserNotFound
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(id int64) error {
	if _, exists := m.users[id]; !exists {
		return model.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) CreateWithPassword(user *model.User, password string) error {
	return m.Create(user)
}

func (m *MockUserRepository) VerifyPassword(phone, password string) (*model.User, error) {
	return m.GetByPhone(phone)
}

func (m *MockUserRepository) UpdatePassword(userID int64, newPassword string) error {
	return nil
}

// MockMiniMaxService 模拟MiniMax服务
type MockMiniMaxService struct{}

func NewMockMiniMaxService() *MockMiniMaxService {
	return &MockMiniMaxService{}
}

func (m *MockMiniMaxService) ChatCompletion(request minimax.ChatCompletionRequest) (*minimax.ChatCompletionResponse, error) {
	return &minimax.ChatCompletionResponse{
		ID: "test_id",
		Choices: []minimax.Choice{
			{
				FinishReason: "stop",
				Index:        0,
				Message: minimax.ChatMessage{
					Role:    "assistant",
					Content: "这是一个模拟的AI回复",
				},
			},
		},
		Model: "glm-4",
		Usage: minimax.Usage{
			TotalTokens: 100,
		},
		BaseResp: minimax.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		},
	}, nil
}

func (m *MockMiniMaxService) ChatCompletionStream(request minimax.ChatCompletionRequest) (<-chan minimax.ChatCompletionResponse, error) {
	ch := make(chan minimax.ChatCompletionResponse, 1)
	go func() {
		defer close(ch)
		ch <- minimax.ChatCompletionResponse{
			BaseResp: minimax.BaseResponse{
				StatusCode: 0,
				StatusMsg:  "success",
			},
		}
	}()
	return ch, nil
}

func (m *MockMiniMaxService) SimpleChat(userMessage string) (string, error) {
	return "模拟回复", nil
}

func (m *MockMiniMaxService) SimpleChatWithParams(userMessage string, temperature float64, maxTokens int) (string, error) {
	return "模拟回复", nil
}

func (m *MockMiniMaxService) GetResponseContent(response *minimax.ChatCompletionResponse) (string, error) {
	return response.GetContent(), nil
}

func (m *MockMiniMaxService) GetUsage(response *minimax.ChatCompletionResponse) *minimax.Usage {
	return &response.Usage
}

func (m *MockMiniMaxService) IsRateLimited(response *minimax.ChatCompletionResponse) bool {
	return false
}

func (m *MockMiniMaxService) IsAuthFailed(response *minimax.ChatCompletionResponse) bool {
	return false
}

func (m *MockMiniMaxService) IsInsufficientBalance(response *minimax.ChatCompletionResponse) bool {
	return false
}

func (m *MockMiniMaxService) IsTokenLimited(response *minimax.ChatCompletionResponse) bool {
	return false
}

// TestCreateConversation 测试创建对话
func TestCreateConversation(t *testing.T) {
	// 创建模拟依赖
	conversationRepo := NewMockConversationRepository()
	messageRepo := NewMockMessageRepository()
	userRepo := NewMockUserRepository()
	minimaxService := NewMockMiniMaxService()

	// 创建用户
	user := &model.User{
		ID:     1,
		Phone:  "13800138000",
		Status: 1,
	}
	userRepo.Create(user)

	// 创建对话缓存
	conversationCache := cache.NewConversationCache("localhost:6379", "", 0)

	// 创建服务
	service := NewService(conversationRepo, messageRepo, userRepo, conversationCache, minimaxService)

	// 测试创建对话
	req := &CreateConversationRequest{
		UserID: 1,
		Title:  "测试对话",
	}

	ctx := context.Background()
	response, err := service.CreateConversation(ctx, req)

	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	if response.Conversation == nil {
		t.Fatal("Expected conversation in response")
	}

	if response.Conversation.Title != "测试对话" {
		t.Errorf("Expected title '测试对话', got '%s'", response.Conversation.Title)
	}

	if response.Conversation.UserID != 1 {
		t.Errorf("Expected user ID 1, got %d", response.Conversation.UserID)
	}
}

// TestSendMessage 测试发送消息
func TestSendMessage(t *testing.T) {
	// 创建模拟依赖
	conversationRepo := NewMockConversationRepository()
	messageRepo := NewMockMessageRepository()
	userRepo := NewMockUserRepository()
	minimaxService := NewMockMiniMaxService()

	// 创建用户
	user := &model.User{
		ID:     1,
		Phone:  "13800138000",
		Status: 1,
	}
	userRepo.Create(user)

	// 创建对话
	conversation := &model.Conversation{
		ID:     1,
		UserID: 1,
		Title:  "测试对话",
		Status: 1,
	}
	conversationRepo.Create(conversation)

	// 创建对话缓存
	conversationCache := cache.NewConversationCache("localhost:6379", "", 0)

	// 创建服务
	service := NewService(conversationRepo, messageRepo, userRepo, conversationCache, minimaxService)

	// 测试发送消息
	req := &SendMessageRequest{
		ConversationID: 1,
		UserID:         1,
		Content:        "你好",
		Model:          "glm-4",
	}

	ctx := context.Background()
	response, err := service.SendMessage(ctx, req)

	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	if response.UserMessage == nil {
		t.Fatal("Expected user message in response")
	}

	if response.AssistantMessage == nil {
		t.Fatal("Expected assistant message in response")
	}

	if response.UserMessage.Content != "你好" {
		t.Errorf("Expected user message content '你好', got '%s'", response.UserMessage.Content)
	}

	if response.AssistantMessage.Content != "这是一个模拟的AI回复" {
		t.Errorf("Expected assistant message content '这是一个模拟的AI回复', got '%s'", response.AssistantMessage.Content)
	}
}
