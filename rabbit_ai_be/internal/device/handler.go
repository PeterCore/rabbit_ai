package device

import (
	"net/http"

	"rabbit_ai/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Handler 设备处理器
type Handler struct {
	deviceService *DeviceService
}

// NewHandler 创建设备处理器实例
func NewHandler(deviceService *DeviceService) *Handler {
	return &Handler{
		deviceService: deviceService,
	}
}

// GetOrCreateUser 根据设备ID获取或创建用户
func (h *Handler) GetOrCreateUser(c *gin.Context) {
	deviceID, exists := middleware.GetDeviceIDFromContext(c)
	if !exists || deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Device ID is required",
		})
		return
	}

	platform, _ := middleware.GetPlatformFromContext(c)
	user, err := h.deviceService.GetOrCreateUserByDeviceID(deviceID, platform)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get or create user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data": gin.H{
			"device_id": deviceID,
			"platform":  platform,
			"user":      user,
			"is_new":    user.Nickname[:8] == "设备用户_",
		},
	})
}

// GetUserByDevice 根据设备ID获取用户信息
func (h *Handler) GetUserByDevice(c *gin.Context) {
	deviceID, exists := middleware.GetDeviceIDFromContext(c)
	if !exists || deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Device ID is required",
		})
		return
	}

	user, err := h.deviceService.GetUserByDeviceID(deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "User not found for this device",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data": gin.H{
			"device_id": deviceID,
			"user":      user,
		},
	})
}

// BindDevice 绑定设备到当前用户
func (h *Handler) BindDevice(c *gin.Context) {
	deviceID, exists := middleware.GetDeviceIDFromContext(c)
	if !exists || deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Device ID is required",
		})
		return
	}

	// 从JWT中获取用户ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "User not authenticated",
		})
		return
	}

	err := h.deviceService.BindDeviceToUser(deviceID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Failed to bind device: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Device bound successfully",
		"data": gin.H{
			"device_id": deviceID,
			"user_id":   userID,
		},
	})
}

// UnbindDevice 解绑设备
func (h *Handler) UnbindDevice(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "User not authenticated",
		})
		return
	}

	err := h.deviceService.UnbindDevice(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to unbind device: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Device unbound successfully",
		"data": gin.H{
			"user_id": userID,
		},
	})
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	device := rg.Group("/device")
	{
		device.GET("/user", h.GetOrCreateUser)      // 根据设备ID获取或创建用户
		device.GET("/user/info", h.GetUserByDevice) // 根据设备ID获取用户信息
		device.POST("/bind", h.BindDevice)          // 绑定设备到用户（需要JWT）
		device.DELETE("/unbind", h.UnbindDevice)    // 解绑设备（需要JWT）
	}
}
