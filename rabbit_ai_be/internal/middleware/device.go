package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// DeviceConfig 设备配置
type DeviceConfig struct {
	HeaderNames []string // 可能的设备ID header名称
}

// DefaultDeviceConfig 默认设备配置
func DefaultDeviceConfig() DeviceConfig {
	return DeviceConfig{
		HeaderNames: []string{
			"X-Device-ID",
			"X-Device-Id",
			"Device-ID",
			"Device-Id",
			"X-Client-ID",
			"X-Client-Id",
			"Client-ID",
			"Client-Id",
			"X-User-Agent",
			"User-Agent",
		},
	}
}

// DeviceMiddleware 设备标识中间件
func DeviceMiddleware(config DeviceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := extractDeviceID(c, config)
		platform := extractPlatform(c)
		// 将设备ID和平台存储到上下文中
		if deviceID != "" {
			c.Set("device_id", deviceID)
		}
		if platform != "" {
			c.Set("platform", platform)
		}
		c.Next()
	}
}

// extractDeviceID 从请求中提取设备ID
func extractDeviceID(c *gin.Context, config DeviceConfig) string {
	// 按优先级尝试不同的header名称
	for _, headerName := range config.HeaderNames {
		if deviceID := c.GetHeader(headerName); deviceID != "" {
			return strings.TrimSpace(deviceID)
		}
	}

	// 如果没有找到设备ID，尝试从User-Agent生成一个简单的标识
	userAgent := c.GetHeader("User-Agent")
	if userAgent != "" {
		// 这里可以添加更复杂的User-Agent解析逻辑
		// 目前简单返回User-Agent的hash或截取
		return generateDeviceIDFromUserAgent(userAgent)
	}

	return ""
}

// generateDeviceIDFromUserAgent 从User-Agent生成设备ID
func generateDeviceIDFromUserAgent(userAgent string) string {
	// 简单的实现：取User-Agent的前50个字符作为设备ID
	// 实际项目中可以使用更复杂的算法，如MD5 hash等
	if len(userAgent) > 50 {
		return "ua_" + userAgent[:50]
	}
	return "ua_" + userAgent
}

// extractPlatform 从请求中提取平台类型
func extractPlatform(c *gin.Context) string {
	// 优先从常见header获取
	possibleHeaders := []string{"X-Platform", "Platform"}
	for _, h := range possibleHeaders {
		if v := c.GetHeader(h); v != "" {
			return normalizePlatform(v)
		}
	}
	// 可选：根据User-Agent简单判断
	ua := c.GetHeader("User-Agent")
	if ua != "" {
		uaLower := strings.ToLower(ua)
		if strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "ios") {
			return "ios"
		}
		if strings.Contains(uaLower, "android") {
			return "android"
		}
		if strings.Contains(uaLower, "windows") || strings.Contains(uaLower, "macintosh") || strings.Contains(uaLower, "linux") {
			return "browser"
		}
	}
	return ""
}

// normalizePlatform 规范化平台字符串
func normalizePlatform(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "ios" || v == "iphone" {
		return "ios"
	}
	if v == "android" {
		return "android"
	}
	if v == "browser" || v == "web" || v == "h5" {
		return "browser"
	}
	return v
}

// GetDeviceIDFromContext 从上下文中获取设备ID
func GetDeviceIDFromContext(c *gin.Context) (string, bool) {
	deviceID, exists := c.Get("device_id")
	if !exists {
		return "", false
	}

	if id, ok := deviceID.(string); ok {
		return id, true
	}

	return "", false
}

// GetPlatformFromContext 从上下文获取平台
func GetPlatformFromContext(c *gin.Context) (string, bool) {
	platform, exists := c.Get("platform")
	if !exists {
		return "", false
	}
	if p, ok := platform.(string); ok {
		return p, true
	}
	return "", false
}

// RequireDeviceID 要求设备ID的中间件
func RequireDeviceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID, exists := GetDeviceIDFromContext(c)
		if !exists || deviceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Device ID is required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
