package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rabbit_ai/internal/model"

	"github.com/redis/go-redis/v9"
)

// RedisCache Redis缓存服务
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存实例
func NewRedisCache(addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: client,
	}
}

// Close 关闭Redis连接
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// 缓存键常量
const (
	UserKeyPrefix = "user:"
	UserTTL       = 30 * time.Minute // 用户信息缓存30分钟
)

// getUserKey 生成用户缓存键
func getUserKey(userID int64) string {
	return fmt.Sprintf("%s%d", UserKeyPrefix, userID)
}

// SetUser 缓存用户信息
func (c *RedisCache) SetUser(ctx context.Context, user *model.User) error {
	key := getUserKey(user.ID)

	// 将用户信息序列化为JSON
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	// 设置缓存，带过期时间
	err = c.client.Set(ctx, key, userData, UserTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set user cache: %w", err)
	}

	return nil
}

// GetUser 从缓存获取用户信息
func (c *RedisCache) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	key := getUserKey(userID)

	// 从Redis获取数据
	userData, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("failed to get user from cache: %w", err)
	}

	// 反序列化用户信息
	var user model.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// DeleteUser 删除用户缓存
func (c *RedisCache) DeleteUser(ctx context.Context, userID int64) error {
	key := getUserKey(userID)

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user cache: %w", err)
	}

	return nil
}

// InvalidateUser 使用户缓存失效（删除缓存）
func (c *RedisCache) InvalidateUser(ctx context.Context, userID int64) error {
	return c.DeleteUser(ctx, userID)
}

// Ping 测试Redis连接
func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}
