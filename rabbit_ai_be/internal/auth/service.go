package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"rabbit_ai/internal/middleware"
	"rabbit_ai/internal/model"
)

// AliyunConfig 阿里云配置
type AliyunConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	Region          string
	OneClickAppID   string
}

// AuthService 认证服务
type AuthService struct {
	userRepo     model.UserRepository
	jwtConfig    middleware.JWTConfig
	aliyunConfig AliyunConfig
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo model.UserRepository, jwtConfig middleware.JWTConfig, aliyunConfig AliyunConfig) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		jwtConfig:    jwtConfig,
		aliyunConfig: aliyunConfig,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	AuthCode string `json:"auth_code" binding:"required"`
}

// PasswordLoginRequest 密码登录请求
type PasswordLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Login 用户登录（阿里一键登录）
func (s *AuthService) Login(authCode string) (*LoginResponse, error) {
	// 1. 调用阿里云接口获取手机号
	phone, err := s.getPhoneFromAliyun(authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get phone from aliyun: %w", err)
	}

	// 2. 查找或创建用户
	user, err := s.userRepo.GetByPhone(phone)
	if err != nil {
		// 用户不存在，创建新用户
		user = &model.User{
			Phone:    phone,
			Nickname: "用户" + phone[len(phone)-4:], // 使用手机号后4位作为默认昵称
			Avatar:   "",                          // 默认头像
			Status:   1,                           // 正常状态
		}

		err = s.userRepo.Create(user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// 3. 生成 JWT token
	token, err := middleware.GenerateToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// PasswordLogin 密码登录
func (s *AuthService) PasswordLogin(phone, password string) (*LoginResponse, error) {
	// 1. 验证密码
	user, err := s.userRepo.VerifyPassword(phone, password)
	if err != nil {
		return nil, fmt.Errorf("invalid phone or password: %w", err)
	}

	// 2. 检查用户状态
	if user.Status != 1 {
		return nil, fmt.Errorf("user account is disabled")
	}

	// 3. 生成 JWT token
	token, err := middleware.GenerateToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Register 用户注册
func (s *AuthService) Register(phone, password, nickname string) (*LoginResponse, error) {
	// 1. 检查用户是否已存在
	existingUser, err := s.userRepo.GetByPhone(phone)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with this phone number already exists")
	}

	// 2. 创建新用户
	user := &model.User{
		Phone:    phone,
		Nickname: nickname,
		Avatar:   "", // 默认头像
		Status:   1,  // 正常状态
	}

	err = s.userRepo.CreateWithPassword(user, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 3. 生成 JWT token
	token, err := middleware.GenerateToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// getPhoneFromAliyun 从阿里云获取手机号
func (s *AuthService) getPhoneFromAliyun(authCode string) (string, error) {
	// 这里需要根据阿里云一键登录的实际API文档来实现
	// 以下是示例实现，实际使用时需要替换为真实的阿里云API调用

	// 构建请求URL和参数
	apiURL := "https://dypnsapi.aliyuncs.com/"
	params := url.Values{}
	params.Set("Action", "GetMobile")
	params.Set("Version", "2017-05-25")
	params.Set("AccessKeyId", s.aliyunConfig.AccessKeyID)
	params.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	params.Set("SignatureMethod", "HMAC-SHA1")
	params.Set("SignatureVersion", "1.0")
	params.Set("SignatureNonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	params.Set("Format", "JSON")
	params.Set("AppId", s.aliyunConfig.OneClickAppID)
	params.Set("Token", authCode)

	// 发送HTTP请求
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 解析响应
	var result struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    struct {
			Mobile string `json:"Mobile"`
		} `json:"Data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Code != "OK" {
		return "", fmt.Errorf("aliyun API error: %s", result.Message)
	}

	return result.Data.Mobile, nil
}

// 注意：实际项目中，阿里云API调用需要正确的签名算法
// 这里提供的是简化版本，实际使用时需要参考阿里云官方SDK或文档
