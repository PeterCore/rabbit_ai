package minimax

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler MiniMax AI处理器
type Handler struct {
	service *MiniMaxService
}

// NewHandler 创建MiniMax处理器实例
func NewHandler(service *MiniMaxService) *Handler {
	return &Handler{
		service: service,
	}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Content string `json:"content"`
	Usage   *Usage `json:"usage,omitempty"`
}

// Chat 聊天接口
func (h *Handler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	// 调用MiniMax服务
	response, err := h.service.ChatCompletion(ChatCompletionRequest{
		Model: "MiniMax-M1",
		Messages: []ChatMessage{
			{
				Role:    "system",
				Name:    "MiniMax AI",
				Content: "",
			},
			{
				Role:    "user",
				Name:    "用户",
				Content: req.Message,
			},
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get AI response: " + err.Error(),
		})
		return
	}

	// 获取回复内容
	content, err := h.service.GetResponseContent(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get response content: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data": ChatResponse{
			Content: content,
			Usage:   h.service.GetUsage(response),
		},
	})
}

// SimpleChat 简单聊天接口
func (h *Handler) SimpleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	// 调用简单聊天方法
	content, err := h.service.SimpleChat(req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get AI response: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data": gin.H{
			"content": content,
		},
	})
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	ai := rg.Group("/ai")
	{
		ai.POST("/chat", h.Chat)              // 完整聊天接口
		ai.POST("/chat/simple", h.SimpleChat) // 简单聊天接口
	}
}
