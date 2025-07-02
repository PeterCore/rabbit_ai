package model

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ErrUserNotFound 用户未找到错误
var ErrUserNotFound = errors.New("user not found")

// User 用户模型
type User struct {
	ID        int64     `json:"id" db:"id"`
	Phone     string    `json:"phone" db:"phone"`
	Password  string    `json:"-" db:"password"` // 密码不返回给前端
	Nickname  string    `json:"nickname" db:"nickname"`
	Avatar    string    `json:"avatar" db:"avatar"`
	Status    int       `json:"status" db:"status"`       // 1: 正常, 0: 禁用
	GitHubID  string    `json:"github_id" db:"github_id"` // GitHub用户ID
	Email     string    `json:"email" db:"email"`         // 邮箱
	DeviceID  string    `json:"device_id" db:"device_id"` // 设备唯一标识
	Platform  string    `json:"platform" db:"platform"`   // 终端平台: ios/android/browser
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserRepository 用户数据访问接口
type UserRepository interface {
	Create(user *User) error
	GetByID(id int64) (*User, error)
	GetByPhone(phone string) (*User, error)
	GetByGitHubID(githubID string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByDeviceID(deviceID string) (*User, error)
	Update(user *User) error
	Delete(id int64) error
	CreateWithPassword(user *User, password string) error
	VerifyPassword(phone, password string) (*User, error)
	UpdatePassword(userID int64, newPassword string) error
}

// UserRepositoryImpl 用户数据访问实现
type UserRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(user *User) error {
	query := `
		INSERT INTO users (phone, nickname, avatar, status, github_id, email, device_id, platform, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return r.db.QueryRow(
		query,
		user.Phone,
		user.Nickname,
		user.Avatar,
		user.Status,
		user.GitHubID,
		user.Email,
		user.DeviceID,
		user.Platform,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
}

// CreateWithPassword 创建带密码的用户
func (r *UserRepositoryImpl) CreateWithPassword(user *User, password string) error {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (phone, password, nickname, avatar, status, github_id, email, device_id, platform, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return r.db.QueryRow(
		query,
		user.Phone,
		string(hashedPassword),
		user.Nickname,
		user.Avatar,
		user.Status,
		user.GitHubID,
		user.Email,
		user.DeviceID,
		user.Platform,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
}

// GetByID 根据ID获取用户
func (r *UserRepositoryImpl) GetByID(id int64) (*User, error) {
	user := &User{}
	query := `
		SELECT id, phone, password, nickname, avatar, status, github_id, email, device_id, platform, created_at, updated_at
		FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Phone,
		&user.Password,
		&user.Nickname,
		&user.Avatar,
		&user.Status,
		&user.GitHubID,
		&user.Email,
		&user.DeviceID,
		&user.Platform,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByPhone 根据手机号获取用户
func (r *UserRepositoryImpl) GetByPhone(phone string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, phone, password, nickname, avatar, status, github_id, email, device_id, platform, created_at, updated_at
		FROM users WHERE phone = $1`

	err := r.db.QueryRow(query, phone).Scan(
		&user.ID,
		&user.Phone,
		&user.Password,
		&user.Nickname,
		&user.Avatar,
		&user.Status,
		&user.GitHubID,
		&user.Email,
		&user.DeviceID,
		&user.Platform,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByGitHubID 根据GitHub ID获取用户
func (r *UserRepositoryImpl) GetByGitHubID(githubID string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, phone, password, nickname, avatar, status, github_id, email, device_id, platform, created_at, updated_at
		FROM users WHERE github_id = $1`

	err := r.db.QueryRow(query, githubID).Scan(
		&user.ID,
		&user.Phone,
		&user.Password,
		&user.Nickname,
		&user.Avatar,
		&user.Status,
		&user.GitHubID,
		&user.Email,
		&user.DeviceID,
		&user.Platform,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepositoryImpl) GetByEmail(email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, phone, password, nickname, avatar, status, github_id, email, device_id, platform, created_at, updated_at
		FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Phone,
		&user.Password,
		&user.Nickname,
		&user.Avatar,
		&user.Status,
		&user.GitHubID,
		&user.Email,
		&user.DeviceID,
		&user.Platform,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByDeviceID 根据设备ID获取用户
func (r *UserRepositoryImpl) GetByDeviceID(deviceID string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, phone, password, nickname, avatar, status, github_id, email, device_id, platform, created_at, updated_at
		FROM users WHERE device_id = $1`

	err := r.db.QueryRow(query, deviceID).Scan(
		&user.ID,
		&user.Phone,
		&user.Password,
		&user.Nickname,
		&user.Avatar,
		&user.Status,
		&user.GitHubID,
		&user.Email,
		&user.DeviceID,
		&user.Platform,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// VerifyPassword 验证密码
func (r *UserRepositoryImpl) VerifyPassword(phone, password string) (*User, error) {
	user, err := r.GetByPhone(phone)
	if err != nil {
		return nil, err
	}

	// 检查用户是否有密码
	if user.Password == "" {
		return nil, sql.ErrNoRows
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update 更新用户信息
func (r *UserRepositoryImpl) Update(user *User) error {
	query := `
		UPDATE users 
		SET nickname = $1, avatar = $2, status = $3, github_id = $4, email = $5, device_id = $6, platform = $7, updated_at = $8
		WHERE id = $9`

	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(
		query,
		user.Nickname,
		user.Avatar,
		user.Status,
		user.GitHubID,
		user.Email,
		user.DeviceID,
		user.Platform,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// UpdatePassword 更新密码
func (r *UserRepositoryImpl) UpdatePassword(userID int64, newPassword string) error {
	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		UPDATE users 
		SET password = $1, updated_at = $2
		WHERE id = $3`

	_, err = r.db.Exec(
		query,
		string(hashedPassword),
		time.Now(),
		userID,
	)

	return err
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
