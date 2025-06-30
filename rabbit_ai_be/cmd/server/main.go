package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"rabbit_ai/internal/auth"
	"rabbit_ai/internal/cache"
	"rabbit_ai/internal/middleware"
	"rabbit_ai/internal/model"
	"rabbit_ai/internal/repository"
	"rabbit_ai/internal/user"
)

// Config 配置结构
type Config struct {
	Server struct {
		Port string `yaml:"port"`
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
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		Password     string `yaml:"password"`
		DB           int    `yaml:"db"`
		PoolSize     int    `yaml:"pool_size"`
		MinIdleConns int    `yaml:"min_idle_conns"`
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
}

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 初始化配置
	config := loadConfig()

	// 设置Gin模式
	gin.SetMode(config.Server.Mode)

	// 连接数据库
	db, err := connectDatabase(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 初始化数据库表
	if err := initDatabase(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 连接Redis
	redisCache, err := connectRedis(config)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisCache.Close()

	// 测试Redis连接
	ctx := context.Background()
	if err := redisCache.Ping(ctx); err != nil {
		log.Printf("Warning: Redis ping failed: %v", err)
	} else {
		log.Println("Redis connection established successfully")
	}

	// 创建基础用户仓库
	baseUserRepo := model.NewUserRepository(db)

	// 创建带缓存的用户仓库
	userRepo := repository.NewCachedUserRepository(baseUserRepo, redisCache)

	// 创建缓存管理器
	cacheManager := cache.NewCacheManager(redisCache)

	// 创建JWT配置
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

	// 创建服务
	authService := auth.NewAuthService(userRepo, jwtConfig, aliyunConfig, githubOAuth)
	userService := user.NewUserService(userRepo)

	// 创建处理器
	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)
	cacheHandler := cache.NewHandler(cacheManager)

	// 创建路由
	r := gin.Default()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证路由（无需JWT验证）
		authHandler.RegisterRoutes(api)

		// 需要JWT验证的路由
		protected := api.Group("")
		protected.Use(middleware.JWTMiddleware(jwtConfig))
		{
			userHandler.RegisterRoutes(protected)
			cacheHandler.RegisterRoutes(protected)
		}
	}

	// 启动服务器
	port := config.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// loadConfig 加载配置
func loadConfig() *Config {
	config := &Config{}

	// 从环境变量加载配置
	config.Server.Port = getEnv("SERVER_PORT", "8080")
	config.Server.Mode = getEnv("SERVER_MODE", "debug")

	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = 5432
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "password")
	config.Database.DBName = getEnv("DB_NAME", "rabbit_ai")
	config.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// Redis配置
	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = 6379
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.DB = 0
	config.Redis.PoolSize = 10
	config.Redis.MinIdleConns = 5

	config.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key-here")
	config.JWT.ExpireHours = 24

	config.Aliyun.AccessKeyID = getEnv("ALIYUN_ACCESS_KEY_ID", "")
	config.Aliyun.AccessKeySecret = getEnv("ALIYUN_ACCESS_KEY_SECRET", "")
	config.Aliyun.Region = getEnv("ALIYUN_REGION", "cn-hangzhou")
	config.Aliyun.OneClickAppID = getEnv("ALIYUN_ONE_CLICK_APP_ID", "")

	config.GitHub.ClientID = getEnv("GITHUB_CLIENT_ID", "")
	config.GitHub.ClientSecret = getEnv("GITHUB_CLIENT_SECRET", "")
	config.GitHub.RedirectURL = getEnv("GITHUB_REDIRECT_URL", "")

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
func connectDatabase(config *Config) (*sql.DB, error) {
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

	return db, nil
}

// connectRedis 连接Redis
func connectRedis(config *Config) (*cache.RedisCache, error) {
	addr := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)

	redisCache := cache.NewRedisCache(addr, config.Redis.Password, config.Redis.DB)

	// 测试连接
	ctx := context.Background()
	if err := redisCache.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return redisCache, nil
}

// initDatabase 初始化数据库表
func initDatabase(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			phone VARCHAR(20) UNIQUE,
			password VARCHAR(255),
			nickname VARCHAR(100) NOT NULL,
			avatar TEXT,
			status INTEGER DEFAULT 1,
			github_id VARCHAR(100) UNIQUE,
			email VARCHAR(255) UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
		CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
		CREATE INDEX IF NOT EXISTS idx_users_github_id ON users(github_id);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	_, err := db.Exec(query)
	return err
}
