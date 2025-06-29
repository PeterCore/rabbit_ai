package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Claims JWT 声明结构
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

// JWTMiddleware JWT 中间件
func JWTMiddleware(config JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// 解析 JWT token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 验证 token 是否有效
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// 检查 token 是否过期
			if time.Now().Unix() > claims.ExpiresAt {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "Token has expired",
				})
				c.Abort()
				return
			}

			// 将用户ID存储到上下文中
			c.Set("user_id", claims.UserID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid token claims",
			})
			c.Abort()
			return
		}
	}
}

// GenerateToken 生成 JWT token
func GenerateToken(userID int64, config JWTConfig) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.ExpireTime).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "rabbit_ai",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Secret))
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	if id, ok := userID.(int64); ok {
		return id, true
	}

	return 0, false
}
