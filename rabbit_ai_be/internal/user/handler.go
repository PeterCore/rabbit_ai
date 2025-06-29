package user

import (
	"net/http"
	"strconv"

	"rabbit_ai/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Handler 用户处理器
type Handler struct {
	userService *UserService
}

// NewHandler 创建用户处理器实例
func NewHandler(userService *UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

// GetProfile 获取用户信息
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "User not authenticated",
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get user profile: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    user,
	})
}

// UpdateProfile 更新用户信息
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "User not authenticated",
		})
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters: " + err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(userID, req.Nickname, req.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update user profile: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Profile updated successfully",
		"data":    user,
	})
}

// DeleteUser 删除用户
func (h *Handler) DeleteUser(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "User not authenticated",
		})
		return
	}

	err := h.userService.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "User deleted successfully",
	})
}

// GetUserByID 根据ID获取用户信息（管理员接口）
func (h *Handler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    user,
	})
}

// UpdatePassword 修改密码
func (h *Handler) UpdatePassword(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "User not authenticated",
		})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters: " + err.Error(),
		})
		return
	}

	err := h.userService.UpdatePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Failed to update password: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Password updated successfully",
	})
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("/profile", h.GetProfile)
		users.PUT("/profile", h.UpdateProfile)
		users.PUT("/password", h.UpdatePassword)
		users.DELETE("/profile", h.DeleteUser)
		users.GET("/:id", h.GetUserByID)
	}
}
