package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rabbit_ai/internal/model"

	"github.com/redis/go-redis/v9"
)

// ConversationCache 对话缓存服务
type ConversationCache struct {
	client *redis.Client
}

// NewConversationCache 创建对话缓存实例
func NewConversationCache(addr, password string, db int) *ConversationCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &ConversationCache{
		client: client,
	}
}

// Close 关闭Redis连接
func (c *ConversationCache) Close() error {
	return c.client.Close()
}

// 缓存键常量和TTL
const (
	ConversationKeyPrefix         = "conversation:"
	MessageKeyPrefix              = "message:"
	UserConversationsKeyPrefix    = "user_conversations:"
	ConversationMessagesKeyPrefix = "conversation_messages:"

	ConversationTTL         = 30 * time.Minute // 对话信息缓存30分钟
	MessageTTL              = 1 * time.Hour    // 消息缓存1小时
	UserConversationsTTL    = 15 * time.Minute // 用户对话列表缓存15分钟
	ConversationMessagesTTL = 30 * time.Minute // 对话消息列表缓存30分钟
)

// 生成缓存键
func getConversationKey(conversationID int64) string {
	return fmt.Sprintf("%s%d", ConversationKeyPrefix, conversationID)
}

func getMessageKey(messageID int64) string {
	return fmt.Sprintf("%s%d", MessageKeyPrefix, messageID)
}

func getUserConversationsKey(userID int64) string {
	return fmt.Sprintf("%s%d", UserConversationsKeyPrefix, userID)
}

func getConversationMessagesKey(conversationID int64) string {
	return fmt.Sprintf("%s%d", ConversationMessagesKeyPrefix, conversationID)
}

// SetConversation 缓存对话信息
func (c *ConversationCache) SetConversation(ctx context.Context, conversation *model.Conversation) error {
	key := getConversationKey(conversation.ID)

	conversationData, err := json.Marshal(conversation)
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %w", err)
	}

	err = c.client.Set(ctx, key, conversationData, ConversationTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set conversation cache: %w", err)
	}

	return nil
}

// GetConversation 从缓存获取对话信息
func (c *ConversationCache) GetConversation(ctx context.Context, conversationID int64) (*model.Conversation, error) {
	key := getConversationKey(conversationID)

	conversationData, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("failed to get conversation from cache: %w", err)
	}

	var conversation model.Conversation
	err = json.Unmarshal([]byte(conversationData), &conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversation: %w", err)
	}

	return &conversation, nil
}

// DeleteConversation 删除对话缓存
func (c *ConversationCache) DeleteConversation(ctx context.Context, conversationID int64) error {
	key := getConversationKey(conversationID)
	return c.client.Del(ctx, key).Err()
}

// SetMessage 缓存消息信息
func (c *ConversationCache) SetMessage(ctx context.Context, message *model.Message) error {
	key := getMessageKey(message.ID)

	messageData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = c.client.Set(ctx, key, messageData, MessageTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set message cache: %w", err)
	}

	return nil
}

// GetMessage 从缓存获取消息信息
func (c *ConversationCache) GetMessage(ctx context.Context, messageID int64) (*model.Message, error) {
	key := getMessageKey(messageID)

	messageData, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("failed to get message from cache: %w", err)
	}

	var message model.Message
	err = json.Unmarshal([]byte(messageData), &message)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &message, nil
}

// DeleteMessage 删除消息缓存
func (c *ConversationCache) DeleteMessage(ctx context.Context, messageID int64) error {
	key := getMessageKey(messageID)
	return c.client.Del(ctx, key).Err()
}

// SetUserConversations 缓存用户对话列表
func (c *ConversationCache) SetUserConversations(ctx context.Context, userID int64, conversations []*model.Conversation) error {
	key := getUserConversationsKey(userID)

	conversationsData, err := json.Marshal(conversations)
	if err != nil {
		return fmt.Errorf("failed to marshal conversations: %w", err)
	}

	err = c.client.Set(ctx, key, conversationsData, UserConversationsTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set user conversations cache: %w", err)
	}

	return nil
}

// GetUserConversations 从缓存获取用户对话列表
func (c *ConversationCache) GetUserConversations(ctx context.Context, userID int64) ([]*model.Conversation, error) {
	key := getUserConversationsKey(userID)

	conversationsData, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("failed to get user conversations from cache: %w", err)
	}

	var conversations []*model.Conversation
	err = json.Unmarshal([]byte(conversationsData), &conversations)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversations: %w", err)
	}

	return conversations, nil
}

// DeleteUserConversations 删除用户对话列表缓存
func (c *ConversationCache) DeleteUserConversations(ctx context.Context, userID int64) error {
	key := getUserConversationsKey(userID)
	return c.client.Del(ctx, key).Err()
}

// SetConversationMessages 缓存对话消息列表
func (c *ConversationCache) SetConversationMessages(ctx context.Context, conversationID int64, messages []*model.Message) error {
	key := getConversationMessagesKey(conversationID)

	messagesData, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal messages: %w", err)
	}

	err = c.client.Set(ctx, key, messagesData, ConversationMessagesTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set conversation messages cache: %w", err)
	}

	return nil
}

// GetConversationMessages 从缓存获取对话消息列表
func (c *ConversationCache) GetConversationMessages(ctx context.Context, conversationID int64) ([]*model.Message, error) {
	key := getConversationMessagesKey(conversationID)

	messagesData, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("failed to get conversation messages from cache: %w", err)
	}

	var messages []*model.Message
	err = json.Unmarshal([]byte(messagesData), &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages: %w", err)
	}

	return messages, nil
}

// DeleteConversationMessages 删除对话消息列表缓存
func (c *ConversationCache) DeleteConversationMessages(ctx context.Context, conversationID int64) error {
	key := getConversationMessagesKey(conversationID)
	return c.client.Del(ctx, key).Err()
}

// InvalidateConversationCache 使对话相关缓存失效
func (c *ConversationCache) InvalidateConversationCache(ctx context.Context, conversationID int64) error {
	// 删除对话缓存
	if err := c.DeleteConversation(ctx, conversationID); err != nil {
		return err
	}

	// 删除对话消息列表缓存
	if err := c.DeleteConversationMessages(ctx, conversationID); err != nil {
		return err
	}

	return nil
}

// InvalidateUserCache 使用户相关缓存失效
func (c *ConversationCache) InvalidateUserCache(ctx context.Context, userID int64) error {
	// 删除用户对话列表缓存
	return c.DeleteUserConversations(ctx, userID)
}

// Ping 测试Redis连接
func (c *ConversationCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}
