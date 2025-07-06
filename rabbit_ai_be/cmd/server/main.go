package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"rabbit_ai/internal/auth"
	"rabbit_ai/internal/cache"
	"rabbit_ai/internal/conversation"
	"rabbit_ai/internal/device"
	"rabbit_ai/internal/middleware"
	"rabbit_ai/internal/minimax"
	"rabbit_ai/internal/model"
	"rabbit_ai/internal/repository"
	"rabbit_ai/internal/user"
)

// Config 配置结构
type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	JWT struct {
		Secret      string `yaml:"secret"`
		ExpireHours int    `yaml:"expire_hours"`
	} `yaml:"jwt"`
	Aliyun struct {
		AccessKeyID     string `yaml:"access_key_id"`
		AccessKeySecret string `yaml:"access_key_secret"`
		Region          string `yaml:"region"`
		OneClickAppID   string `yaml:"one_click_app_id"`
	} `yaml:"aliyun"`
	GitHub struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		RedirectURL  string `yaml:"redirect_url"`
	} `yaml:"github"`
	MiniMax struct {
		APIKey  string `yaml:"api_key"`
		BaseURL string `yaml:"base_url"`
	} `yaml:"minimax"`
}

func main() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// 加载配置
	config := loadConfig()

	// 设置Gin模式
	gin.SetMode(config.Server.Mode)

	// 初始化数据库连接
	db, err := connectDatabase(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 初始化Redis连接
	redisClient, err := connectRedis(config)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// 初始化用户仓库（带缓存）
	userRepo := repository.NewCachedUserRepository(
		model.NewUserRepository(db),
		redisClient,
	)

	// 初始化对话和消息仓库
	conversationRepo := model.NewConversationRepository(db)
	messageRepo := model.NewMessageRepository(db)

	// 初始化对话缓存
	conversationCache := cache.NewConversationCache(
		fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		config.Redis.Password,
		config.Redis.DB+1, // 使用不同的数据库避免冲突
	)

	// 初始化JWT配置
	jwtConfig := middleware.JWTConfig{
		Secret:     config.JWT.Secret,
		ExpireTime: time.Duration(config.JWT.ExpireHours) * time.Hour,
	}

	// 创建阿里云配置
	aliyunConfig := auth.AliyunConfig{
		AccessKeyID:     config.Aliyun.AccessKeyID,
		AccessKeySecret: config.Aliyun.AccessKeySecret,
		Region:          config.Aliyun.Region,
		OneClickAppID:   config.Aliyun.OneClickAppID,
	}

	// 创建GitHub OAuth配置
	githubOAuth := auth.NewGitHubOAuth(
		config.GitHub.ClientID,
		config.GitHub.ClientSecret,
		config.GitHub.RedirectURL,
	)

	// 初始化服务
	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(userRepo, jwtConfig, aliyunConfig, githubOAuth)
	deviceService := device.NewDeviceService(userRepo)

	// 初始化MiniMax AI服务
	minimaxConfig := minimax.MiniMaxConfig{
		APIKey:  config.MiniMax.APIKey,
		BaseURL: config.MiniMax.BaseURL,
	}
	if minimaxConfig.BaseURL == "" {
		minimaxConfig.BaseURL = "https://api.minimaxi.com/v1"
	}
	minimaxService := minimax.NewMiniMaxService(minimaxConfig)

	// 初始化对话服务
	conversationService := conversation.NewService(
		conversationRepo,
		messageRepo,
		userRepo,
		conversationCache,
		minimaxService,
	)

	// 初始化处理器
	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService)
	deviceHandler := device.NewHandler(deviceService)
	minimaxHandler := minimax.NewHandler(minimaxService)
	conversationHandler := conversation.NewHandler(conversationService)

	// 初始化设备中间件配置
	deviceConfig := middleware.DefaultDeviceConfig()

	// 创建路由
	r := gin.Default()

	// 添加设备中间件（全局）
	r.Use(middleware.DeviceMiddleware(deviceConfig))

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Device-ID, X-Client-ID, X-Platform, Platform")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 用户相关路由
		userHandler.RegisterRoutes(api)

		// 认证相关路由
		authHandler.RegisterRoutes(api)

		// 设备相关路由
		deviceHandler.RegisterRoutes(api)

		// AI相关路由
		minimaxHandler.RegisterRoutes(api)

		// 需要JWT认证的路由组
		authorized := api.Group("/")
		authorized.Use(middleware.JWTMiddleware(jwtConfig))
		{
			// 对话相关路由（需要认证）
			conversationHandler.RegisterRoutes(authorized)

			// 这里可以添加需要认证的路由
			authorized.GET("/profile", func(c *gin.Context) {
				userID, _ := middleware.GetUserIDFromContext(c)
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"message": "Profile endpoint",
					"data":    gin.H{"user_id": userID},
				})
			})
		}
	}

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.Server.Port)
	log.Printf("Server starting on port %d", config.Server.Port)
	log.Fatal(r.Run(addr))
}

// loadConfig 从环境变量加载配置
func loadConfig() Config {
	var config Config

	// 从环境变量加载配置
	config.Server.Port = 8080 // 默认端口
	if portStr := getEnv("SERVER_PORT", ""); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Server.Port = port
		}
	}
	config.Server.Mode = getEnv("SERVER_MODE", "debug")

	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = 5432
	if portStr := getEnv("DB_PORT", ""); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Database.Port = port
		}
	}
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "password")
	config.Database.DBName = getEnv("DB_NAME", "rabbit_ai")
	config.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = 6379
	if portStr := getEnv("REDIS_PORT", ""); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Redis.Port = port
		}
	}
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.DB = 0

	config.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key-here")
	config.JWT.ExpireHours = 24
	if hoursStr := getEnv("JWT_EXPIRE_HOURS", ""); hoursStr != "" {
		if hours, err := strconv.Atoi(hoursStr); err == nil {
			config.JWT.ExpireHours = hours
		}
	}

	config.Aliyun.AccessKeyID = getEnv("ALIYUN_ACCESS_KEY_ID", "")
	config.Aliyun.AccessKeySecret = getEnv("ALIYUN_ACCESS_KEY_SECRET", "")
	config.Aliyun.Region = getEnv("ALIYUN_REGION", "cn-hangzhou")
	config.Aliyun.OneClickAppID = getEnv("ALIYUN_ONE_CLICK_APP_ID", "")

	config.GitHub.ClientID = getEnv("GITHUB_CLIENT_ID", "")
	config.GitHub.ClientSecret = getEnv("GITHUB_CLIENT_SECRET", "")
	config.GitHub.RedirectURL = getEnv("GITHUB_REDIRECT_URL", "")

	// MiniMax配置从.env文件获取
	config.MiniMax.APIKey = getEnv("MINIMAX_API_KEY", "")
	config.MiniMax.BaseURL = getEnv("MINIMAX_BASE_URL", "https://api.minimaxi.com/v1")

	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// connectDatabase 连接数据库
func connectDatabase(config Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 创建用户表
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			phone VARCHAR(15) UNIQUE,
			status VARCHAR(20),
			github_id VARCHAR(100) UNIQUE,
			email VARCHAR(255) UNIQUE,
			device_id VARCHAR(255) UNIQUE,
			platform VARCHAR(20),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
		CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
		CREATE INDEX IF NOT EXISTS idx_users_github_id ON users(github_id);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_device_id ON users(device_id);
		CREATE INDEX IF NOT EXISTS idx_users_platform ON users(platform);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	return db, nil
}

// connectRedis 连接Redis
func connectRedis(config Config) (*cache.RedisCache, error) {
	addr := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)

	redisCache := cache.NewRedisCache(addr, config.Redis.Password, config.Redis.DB)

	// 测试连接
	ctx := context.Background()
	if err := redisCache.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return redisCache, nil
}
