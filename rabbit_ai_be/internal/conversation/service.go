package conversation

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"rabbit_ai/internal/cache"
	"rabbit_ai/internal/minimax"
	"rabbit_ai/internal/model"
)

// MiniMaxServiceInterface MiniMax服务接口
type MiniMaxServiceInterface interface {
	ChatCompletion(request minimax.ChatCompletionRequest) (*minimax.ChatCompletionResponse, error)
	ChatCompletionStream(request minimax.ChatCompletionRequest) (<-chan minimax.ChatCompletionResponse, error)
	SimpleChat(userMessage string) (string, error)
	SimpleChatWithParams(userMessage string, temperature float64, maxTokens int) (string, error)
	GetResponseContent(response *minimax.ChatCompletionResponse) (string, error)
	GetUsage(response *minimax.ChatCompletionResponse) *minimax.Usage
	IsRateLimited(response *minimax.ChatCompletionResponse) bool
	IsAuthFailed(response *minimax.ChatCompletionResponse) bool
	IsInsufficientBalance(response *minimax.ChatCompletionResponse) bool
	IsTokenLimited(response *minimax.ChatCompletionResponse) bool
}

// Service 对话服务
type Service struct {
	conversationRepo  model.ConversationRepository
	messageRepo       model.MessageRepository
	userRepo          model.UserRepository
	conversationCache *cache.ConversationCache
	minimaxService    MiniMaxServiceInterface
}

// NewService 创建对话服务实例
func NewService(
	conversationRepo model.ConversationRepository,
	messageRepo model.MessageRepository,
	userRepo model.UserRepository,
	conversationCache *cache.ConversationCache,
	minimaxService MiniMaxServiceInterface,
) *Service {
	return &Service{
		conversationRepo:  conversationRepo,
		messageRepo:       messageRepo,
		userRepo:          userRepo,
		conversationCache: conversationCache,
		minimaxService:    minimaxService,
	}
}

// CreateConversationRequest 创建对话请求
type CreateConversationRequest struct {
	UserID int64  `json:"user_id" binding:"required"`
	Title  string `json:"title" binding:"required"`
}

// CreateConversationResponse 创建对话响应
type CreateConversationResponse struct {
	Conversation *model.Conversation `json:"conversation"`
}

// GetConversationsRequest 获取对话列表请求
type GetConversationsRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

// GetConversationsResponse 获取对话列表响应
type GetConversationsResponse struct {
	Conversations []*model.Conversation `json:"conversations"`
	Total         int                   `json:"total"`
}

// GetConversationMessagesRequest 获取对话消息请求
type GetConversationMessagesRequest struct {
	ConversationID int64 `json:"conversation_id" binding:"required"`
	Limit          int   `json:"limit"`
	Offset         int   `json:"offset"`
}

// GetConversationMessagesResponse 获取对话消息响应
type GetConversationMessagesResponse struct {
	Messages []*model.Message `json:"messages"`
	Total    int              `json:"total"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ConversationID int64  `json:"conversation_id" binding:"required"`
	UserID         int64  `json:"user_id" binding:"required"`
	Content        string `json:"content" binding:"required"`
	Model          string `json:"model"`
}

// SendMessageResponse 发送消息响应
type SendMessageResponse struct {
	UserMessage      *model.Message      `json:"user_message"`
	AssistantMessage *model.Message      `json:"assistant_message"`
	Conversation     *model.Conversation `json:"conversation"`
}

// DeleteConversationRequest 删除对话请求
type DeleteConversationRequest struct {
	ConversationID int64 `json:"conversation_id" binding:"required"`
	UserID         int64 `json:"user_id" binding:"required"`
}

// CreateConversation 创建新对话
func (s *Service) CreateConversation(ctx context.Context, req *CreateConversationRequest) (*CreateConversationResponse, error) {
	// 验证用户是否存在
	_, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 创建对话
	conversation := &model.Conversation{
		UserID:       req.UserID,
		Title:        req.Title,
		Status:       1,
		MessageCount: 0,
	}

	err = s.conversationRepo.Create(conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// 缓存对话信息
	err = s.conversationCache.SetConversation(ctx, conversation)
	if err != nil {
		// 缓存失败不影响主流程，只记录日志
		fmt.Printf("failed to cache conversation: %v\n", err)
	}

	// 使该用户的对话列表缓存失效
	err = s.conversationCache.InvalidateUserCache(ctx, req.UserID)
	if err != nil {
		fmt.Printf("failed to invalidate user cache: %v\n", err)
	}

	return &CreateConversationResponse{
		Conversation: conversation,
	}, nil
}

// GetConversations 获取用户对话列表
func (s *Service) GetConversations(ctx context.Context, req *GetConversationsRequest) (*GetConversationsResponse, error) {
	// 验证用户是否存在
	_, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 设置默认分页参数
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// 先从缓存获取
	conversations, err := s.conversationCache.GetUserConversations(ctx, req.UserID)
	if err != nil {
		fmt.Printf("failed to get conversations from cache: %v\n", err)
	}

	// 缓存未命中，从数据库获取
	if conversations == nil {
		conversations, err = s.conversationRepo.GetByUserID(req.UserID, req.Limit, req.Offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get conversations: %w", err)
		}

		// 缓存对话列表
		err = s.conversationCache.SetUserConversations(ctx, req.UserID, conversations)
		if err != nil {
			fmt.Printf("failed to cache user conversations: %v\n", err)
		}
	}

	// 获取总数
	total, err := s.conversationRepo.GetUserConversationCount(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation count: %w", err)
	}

	return &GetConversationsResponse{
		Conversations: conversations,
		Total:         total,
	}, nil
}

// GetConversationMessages 获取对话消息
func (s *Service) GetConversationMessages(ctx context.Context, req *GetConversationMessagesRequest) (*GetConversationMessagesResponse, error) {
	// 验证对话是否存在
	_, err := s.conversationRepo.GetByID(req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}

	// 设置默认分页参数
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// 先从缓存获取
	messages, err := s.conversationCache.GetConversationMessages(ctx, req.ConversationID)
	if err != nil {
		fmt.Printf("failed to get messages from cache: %v\n", err)
	}

	// 缓存未命中，从数据库获取
	if messages == nil {
		messages, err = s.messageRepo.GetByConversationID(req.ConversationID, req.Limit, req.Offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}

		// 缓存消息列表
		err = s.conversationCache.SetConversationMessages(ctx, req.ConversationID, messages)
		if err != nil {
			fmt.Printf("failed to cache conversation messages: %v\n", err)
		}
	}

	// 获取总数
	total, err := s.messageRepo.GetConversationMessageCount(req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message count: %w", err)
	}

	return &GetConversationMessagesResponse{
		Messages: messages,
		Total:    total,
	}, nil
}

// SendMessage 发送消息并获取AI回复
func (s *Service) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	// 验证用户是否存在
	_, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 验证对话是否存在
	conversation, err := s.conversationRepo.GetByID(req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}

	// 验证对话是否属于该用户
	if conversation.UserID != req.UserID {
		return nil, errors.New("conversation does not belong to user")
	}

	// 设置默认模型
	if req.Model == "" {
		req.Model = "glm-4"
	}

	// 创建用户消息
	userMessage := &model.Message{
		ConversationID: req.ConversationID,
		Role:           "user",
		Content:        req.Content,
		Model:          req.Model,
	}

	err = s.messageRepo.Create(userMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to create user message: %w", err)
	}

	// 缓存用户消息
	err = s.conversationCache.SetMessage(ctx, userMessage)
	if err != nil {
		fmt.Printf("failed to cache user message: %v\n", err)
	}

	// 获取对话历史消息
	historyMessages, err := s.messageRepo.GetConversationMessages(req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	// 构建MiniMax请求消息
	var minimaxMessages []minimax.ChatMessage
	for _, msg := range historyMessages {
		minimaxMessages = append(minimaxMessages, minimax.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 调用MiniMax API
	minimaxReq := minimax.NewChatCompletionRequest(req.Model, minimaxMessages).
		WithMaxTokens(2048).
		WithTemperature(0.7).
		WithUser(fmt.Sprintf("user_%d", req.UserID))

	minimaxResp, err := s.minimaxService.ChatCompletion(*minimaxReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	// 检查MiniMax响应
	if !minimaxResp.IsSuccess() {
		err := minimaxResp.GetError()
		return nil, fmt.Errorf("minimax API error: %d - %s", err.Code, err.Message)
	}

	content := minimaxResp.GetContent()
	if content == "" {
		return nil, errors.New("empty response from AI")
	}

	// 创建AI回复消息
	assistantMessage := &model.Message{
		ConversationID: req.ConversationID,
		Role:           "assistant",
		Content:        content,
		Model:          req.Model,
		FinishReason:   minimaxResp.GetFinishReason(),
		Tokens:         minimaxResp.Usage.TotalTokens,
	}

	err = s.messageRepo.Create(assistantMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to create assistant message: %w", err)
	}

	// 缓存AI消息
	err = s.conversationCache.SetMessage(ctx, assistantMessage)
	if err != nil {
		fmt.Printf("failed to cache assistant message: %v\n", err)
	}

	// 更新对话信息
	conversation.MessageCount += 2 // 用户消息 + AI回复
	conversation.LastMessageAt = time.Now()

	// 如果对话标题为空或为默认标题，使用用户消息的前20个字符作为标题
	if conversation.Title == "" || conversation.Title == "新对话" {
		title := strings.TrimSpace(req.Content)
		if len(title) > 20 {
			title = title[:20] + "..."
		}
		conversation.Title = title
	}

	err = s.conversationRepo.Update(conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to update conversation: %w", err)
	}

	// 缓存更新后的对话信息
	err = s.conversationCache.SetConversation(ctx, conversation)
	if err != nil {
		fmt.Printf("failed to cache updated conversation: %v\n", err)
	}

	// 使相关缓存失效
	err = s.conversationCache.InvalidateConversationCache(ctx, req.ConversationID)
	if err != nil {
		fmt.Printf("failed to invalidate conversation cache: %v\n", err)
	}

	err = s.conversationCache.InvalidateUserCache(ctx, req.UserID)
	if err != nil {
		fmt.Printf("failed to invalidate user cache: %v\n", err)
	}

	return &SendMessageResponse{
		UserMessage:      userMessage,
		AssistantMessage: assistantMessage,
		Conversation:     conversation,
	}, nil
}

// DeleteConversation 删除对话
func (s *Service) DeleteConversation(ctx context.Context, req *DeleteConversationRequest) error {
	// 验证用户是否存在
	_, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 验证对话是否存在
	conversation, err := s.conversationRepo.GetByID(req.ConversationID)
	if err != nil {
		return fmt.Errorf("conversation not found: %w", err)
	}

	// 验证对话是否属于该用户
	if conversation.UserID != req.UserID {
		return errors.New("conversation does not belong to user")
	}

	// 删除对话（软删除）
	err = s.conversationRepo.Delete(req.ConversationID)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	// 使相关缓存失效
	err = s.conversationCache.InvalidateConversationCache(ctx, req.ConversationID)
	if err != nil {
		fmt.Printf("failed to invalidate conversation cache: %v\n", err)
	}

	err = s.conversationCache.InvalidateUserCache(ctx, req.UserID)
	if err != nil {
		fmt.Printf("failed to invalidate user cache: %v\n", err)
	}

	return nil
}
