package repository

import (
	"context"
	"fmt"
	"log"

	"rabbit_ai/internal/cache"
	"rabbit_ai/internal/model"
)

// CachedUserRepository 带缓存的用户仓库
type CachedUserRepository struct {
	userRepo model.UserRepository
	cache    *cache.RedisCache
}

// NewCachedUserRepository 创建带缓存的用户仓库实例
func NewCachedUserRepository(userRepo model.UserRepository, cache *cache.RedisCache) *CachedUserRepository {
	return &CachedUserRepository{
		userRepo: userRepo,
		cache:    cache,
	}
}

// Create 创建用户（同时缓存）
func (r *CachedUserRepository) Create(user *model.User) error {
	// 先创建用户到数据库
	err := r.userRepo.Create(user)
	if err != nil {
		return fmt.Errorf("failed to create user in database: %w", err)
	}

	// 缓存用户信息
	ctx := context.Background()
	err = r.cache.SetUser(ctx, user)
	if err != nil {
		log.Printf("Warning: failed to cache user %d: %v", user.ID, err)
		// 不返回错误，因为数据库操作已经成功
	}

	return nil
}

// CreateWithPassword 创建带密码的用户（同时缓存）
func (r *CachedUserRepository) CreateWithPassword(user *model.User, password string) error {
	// 先创建用户到数据库
	err := r.userRepo.CreateWithPassword(user, password)
	if err != nil {
		return fmt.Errorf("failed to create user with password in database: %w", err)
	}

	// 缓存用户信息（不包含密码）
	ctx := context.Background()
	err = r.cache.SetUser(ctx, user)
	if err != nil {
		log.Printf("Warning: failed to cache user %d: %v", user.ID, err)
		// 不返回错误，因为数据库操作已经成功
	}

	return nil
}

// GetByID 根据ID获取用户（优先从缓存获取）
func (r *CachedUserRepository) GetByID(id int64) (*model.User, error) {
	ctx := context.Background()

	// 先从缓存获取
	cachedUser, err := r.cache.GetUser(ctx, id)
	if err != nil {
		log.Printf("Warning: failed to get user %d from cache: %v", id, err)
	}

	// 如果缓存命中，直接返回
	if cachedUser != nil {
		return cachedUser, nil
	}

	// 缓存未命中，从数据库获取
	user, err := r.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %w", err)
	}

	// 将用户信息缓存
	err = r.cache.SetUser(ctx, user)
	if err != nil {
		log.Printf("Warning: failed to cache user %d: %v", id, err)
		// 不返回错误，因为数据库操作已经成功
	}

	return user, nil
}

// GetByPhone 根据手机号获取用户（不缓存，因为手机号查询较少）
func (r *CachedUserRepository) GetByPhone(phone string) (*model.User, error) {
	return r.userRepo.GetByPhone(phone)
}

// Update 更新用户信息（同时更新缓存）
func (r *CachedUserRepository) Update(user *model.User) error {
	// 先更新数据库
	err := r.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user in database: %w", err)
	}

	// 更新缓存
	ctx := context.Background()
	err = r.cache.SetUser(ctx, user)
	if err != nil {
		log.Printf("Warning: failed to update user cache %d: %v", user.ID, err)
		// 不返回错误，因为数据库操作已经成功
	}

	return nil
}

// UpdatePassword 更新密码（同时使缓存失效）
func (r *CachedUserRepository) UpdatePassword(userID int64, newPassword string) error {
	// 先更新数据库
	err := r.userRepo.UpdatePassword(userID, newPassword)
	if err != nil {
		return fmt.Errorf("failed to update password in database: %w", err)
	}

	// 使缓存失效，因为密码已更改
	ctx := context.Background()
	err = r.cache.InvalidateUser(ctx, userID)
	if err != nil {
		log.Printf("Warning: failed to invalidate user cache %d: %v", userID, err)
		// 不返回错误，因为数据库操作已经成功
	}

	return nil
}

// Delete 删除用户（同时删除缓存）
func (r *CachedUserRepository) Delete(id int64) error {
	// 先删除数据库中的用户
	err := r.userRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete user from database: %w", err)
	}

	// 删除缓存
	ctx := context.Background()
	err = r.cache.DeleteUser(ctx, id)
	if err != nil {
		log.Printf("Warning: failed to delete user cache %d: %v", id, err)
		// 不返回错误，因为数据库操作已经成功
	}

	return nil
}

// VerifyPassword 验证密码（不缓存，因为涉及密码验证）
func (r *CachedUserRepository) VerifyPassword(phone, password string) (*model.User, error) {
	return r.userRepo.VerifyPassword(phone, password)
}
