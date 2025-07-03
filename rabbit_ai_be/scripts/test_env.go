package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// 测试MiniMax配置
	minimaxAPIKey := os.Getenv("MINIMAX_API_KEY")
	minimaxBaseURL := os.Getenv("MINIMAX_BASE_URL")

	fmt.Println("=== MiniMax配置测试 ===")
	fmt.Printf("API Key: %s\n", minimaxAPIKey)
	fmt.Printf("Base URL: %s\n", minimaxBaseURL)

	if minimaxAPIKey == "" || minimaxAPIKey == "your-minimax-api-key" {
		fmt.Println("❌ MiniMax API Key未配置或使用默认值")
	} else {
		fmt.Println("✅ MiniMax API Key已配置")
	}

	if minimaxBaseURL == "" {
		fmt.Println("⚠️  MiniMax Base URL未配置，将使用默认值")
	} else {
		fmt.Println("✅ MiniMax Base URL已配置")
	}

	// 测试其他重要配置
	fmt.Println("\n=== 其他配置测试 ===")

	dbHost := os.Getenv("DB_HOST")
	fmt.Printf("数据库主机: %s\n", dbHost)

	redisHost := os.Getenv("REDIS_HOST")
	fmt.Printf("Redis主机: %s\n", redisHost)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" || jwtSecret == "your-secret-key-here" {
		fmt.Println("❌ JWT Secret未配置或使用默认值")
	} else {
		fmt.Println("✅ JWT Secret已配置")
	}
}
