package cache

import (
	"context"
	"testing"
	"time"

	"rabbit_ai/internal/model"
)

func TestRedisCache(t *testing.T) {
	// 注意：这个测试需要Redis服务器运行
	// 在实际项目中，应该使用测试容器或mock

	// 创建Redis缓存实例
	cache := NewRedisCache("localhost:6379", "", 0)
	defer cache.Close()

	ctx := context.Background()

	// 测试连接
	err := cache.Ping(ctx)
	if err != nil {
		t.Skipf("Redis not available, skipping test: %v", err)
	}

	// 创建测试用户
	testUser := &model.User{
		ID:        1,
		Phone:     "13800138000",
		Nickname:  "测试用户",
		Avatar:    "https://example.com/avatar.jpg",
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("SetUser", func(t *testing.T) {
		err := cache.SetUser(ctx, testUser)
		if err != nil {
			t.Errorf("Failed to set user cache: %v", err)
		}
	})

	t.Run("GetUser", func(t *testing.T) {
		// 先设置缓存
		err := cache.SetUser(ctx, testUser)
		if err != nil {
			t.Fatalf("Failed to set user cache: %v", err)
		}

		// 从缓存获取
		cachedUser, err := cache.GetUser(ctx, testUser.ID)
		if err != nil {
			t.Errorf("Failed to get user from cache: %v", err)
		}

		if cachedUser == nil {
			t.Error("Expected cached user, got nil")
			return
		}

		// 验证用户信息
		if cachedUser.ID != testUser.ID {
			t.Errorf("Expected user ID %d, got %d", testUser.ID, cachedUser.ID)
		}

		if cachedUser.Phone != testUser.Phone {
			t.Errorf("Expected phone %s, got %s", testUser.Phone, cachedUser.Phone)
		}

		if cachedUser.Nickname != testUser.Nickname {
			t.Errorf("Expected nickname %s, got %s", testUser.Nickname, cachedUser.Nickname)
		}
	})

	t.Run("GetUserNotFound", func(t *testing.T) {
		// 获取不存在的用户
		cachedUser, err := cache.GetUser(ctx, 99999)
		if err != nil {
			t.Errorf("Expected no error for non-existent user, got %v", err)
		}

		if cachedUser != nil {
			t.Error("Expected nil for non-existent user")
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		// 先设置缓存
		err := cache.SetUser(ctx, testUser)
		if err != nil {
			t.Fatalf("Failed to set user cache: %v", err)
		}

		// 删除缓存
		err = cache.DeleteUser(ctx, testUser.ID)
		if err != nil {
			t.Errorf("Failed to delete user cache: %v", err)
		}

		// 验证缓存已被删除
		cachedUser, err := cache.GetUser(ctx, testUser.ID)
		if err != nil {
			t.Errorf("Expected no error when getting deleted user, got %v", err)
		}

		if cachedUser != nil {
			t.Error("Expected nil for deleted user")
		}
	})

	t.Run("InvalidateUser", func(t *testing.T) {
		// 先设置缓存
		err := cache.SetUser(ctx, testUser)
		if err != nil {
			t.Fatalf("Failed to set user cache: %v", err)
		}

		// 使缓存失效
		err = cache.InvalidateUser(ctx, testUser.ID)
		if err != nil {
			t.Errorf("Failed to invalidate user cache: %v", err)
		}

		// 验证缓存已被删除
		cachedUser, err := cache.GetUser(ctx, testUser.ID)
		if err != nil {
			t.Errorf("Expected no error when getting invalidated user, got %v", err)
		}

		if cachedUser != nil {
			t.Error("Expected nil for invalidated user")
		}
	})
}
