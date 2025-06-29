package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"rabbit_ai/internal/model"
)

// CacheManager 缓存管理器
type CacheManager struct {
	cache *RedisCache
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(cache *RedisCache) *CacheManager {
	return &CacheManager{
		cache: cache,
	}
}

// CacheStats 缓存统计信息
type CacheStats struct {
	TotalKeys   int64
	MemoryUsage int64
	HitRate     float64
	LastUpdated time.Time
}

// GetStats 获取缓存统计信息
func (m *CacheManager) GetStats(ctx context.Context) (*CacheStats, error) {
	// 这里可以添加更详细的统计信息
	// 实际项目中可以使用 Redis INFO 命令获取更多信息
	return &CacheStats{
		TotalKeys:   0, // 需要实现获取键数量的方法
		MemoryUsage: 0, // 需要实现获取内存使用的方法
		HitRate:     0, // 需要实现命中率统计
		LastUpdated: time.Now(),
	}, nil
}

// WarmUpCache 预热缓存（批量加载用户数据）
func (m *CacheManager) WarmUpCache(ctx context.Context, users []*model.User) error {
	log.Printf("Starting cache warm-up for %d users", len(users))

	for _, user := range users {
		err := m.cache.SetUser(ctx, user)
		if err != nil {
			log.Printf("Warning: failed to warm up cache for user %d: %v", user.ID, err)
			continue
		}
	}

	log.Printf("Cache warm-up completed for %d users", len(users))
	return nil
}

// ClearAllUserCache 清除所有用户缓存
func (m *CacheManager) ClearAllUserCache(ctx context.Context) error {
	// 注意：这是一个简化的实现
	// 实际项目中应该使用 Redis SCAN 命令来安全地删除所有用户缓存
	log.Println("Clearing all user cache...")

	// 这里可以实现批量删除逻辑
	// 例如：使用 SCAN 命令找到所有 user:* 键并删除

	return nil
}

// RefreshUserCache 刷新指定用户的缓存
func (m *CacheManager) RefreshUserCache(ctx context.Context, userID int64, user *model.User) error {
	// 先删除旧缓存
	err := m.cache.DeleteUser(ctx, userID)
	if err != nil {
		log.Printf("Warning: failed to delete old cache for user %d: %v", userID, err)
	}

	// 设置新缓存
	err = m.cache.SetUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to refresh cache for user %d: %w", userID, err)
	}

	log.Printf("Successfully refreshed cache for user %d", userID)
	return nil
}

// BatchSetUsers 批量设置用户缓存
func (m *CacheManager) BatchSetUsers(ctx context.Context, users []*model.User) error {
	log.Printf("Batch setting cache for %d users", len(users))

	for _, user := range users {
		err := m.cache.SetUser(ctx, user)
		if err != nil {
			log.Printf("Warning: failed to set cache for user %d: %v", user.ID, err)
			continue
		}
	}

	return nil
}

// BatchDeleteUsers 批量删除用户缓存
func (m *CacheManager) BatchDeleteUsers(ctx context.Context, userIDs []int64) error {
	log.Printf("Batch deleting cache for %d users", len(userIDs))

	for _, userID := range userIDs {
		err := m.cache.DeleteUser(ctx, userID)
		if err != nil {
			log.Printf("Warning: failed to delete cache for user %d: %v", userID, err)
			continue
		}
	}

	return nil
}

// HealthCheck 缓存健康检查
func (m *CacheManager) HealthCheck(ctx context.Context) error {
	return m.cache.Ping(ctx)
}
