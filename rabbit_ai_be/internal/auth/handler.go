package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler 认证处理器
type Handler struct {
	authService *AuthService
}

// NewHandler 创建认证处理器实例
func NewHandler(authService *AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

// Login 用户登录（阿里一键登录）
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// 调用认证服务
	response, err := h.authService.Login(req.AuthCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Login failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Login successful",
		"data":    response,
	})
}

// PasswordLogin 密码登录
func (h *Handler) PasswordLogin(c *gin.Context) {
	var req PasswordLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// 调用认证服务
	response, err := h.authService.PasswordLogin(req.Phone, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Login failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Login successful",
		"data":    response,
	})
}

// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// 调用认证服务
	response, err := h.authService.Register(req.Phone, req.Password, req.Nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Registration failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Registration successful",
		"data":    response,
	})
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)                  // 阿里一键登录
		auth.POST("/login/password", h.PasswordLogin) // 密码登录
		auth.POST("/register", h.Register)            // 用户注册
	}
}
