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
	Message           string       `json:"message" binding:"required"`
	Temperature       float64      `json:"temperature,omitempty"`
	MaxTokens         int          `json:"max_tokens,omitempty"`
	TopP              float64      `json:"top_p,omitempty"`
	Stream            bool         `json:"stream,omitempty"`
	ToolChoices       []ToolChoice `json:"tool_choices,omitempty"`
	Stop              []string     `json:"stop,omitempty"`
	User              string       `json:"user,omitempty"`
	RepetitionPenalty float64      `json:"repetition_penalty,omitempty"`
	PresencePenalty   float64      `json:"presence_penalty,omitempty"`
	FrequencyPenalty  float64      `json:"frequency_penalty,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Content string `json:"content"`
	Usage   *Usage `json:"usage,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Chat 聊天接口
func (h *Handler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// 构建请求
	request := NewChatCompletionRequest("MiniMax-M1", []ChatMessage{
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
	})

	// 应用可选参数
	if req.Temperature > 0 {
		request.WithTemperature(req.Temperature)
	}
	if req.MaxTokens > 0 {
		request.WithMaxTokens(req.MaxTokens)
	}
	if req.TopP > 0 {
		request.WithTopP(req.TopP)
	}
	if req.Stream {
		request.WithStream(true)
	}
	if len(req.ToolChoices) > 0 {
		request.WithToolChoices(req.ToolChoices)
	}
	if len(req.Stop) > 0 {
		request.WithStop(req.Stop)
	}
	if req.User != "" {
		request.WithUser(req.User)
	}
	if req.RepetitionPenalty > 0 {
		request.RepetitionPenalty = req.RepetitionPenalty
	}
	if req.PresencePenalty != 0 {
		request.PresencePenalty = req.PresencePenalty
	}
	if req.FrequencyPenalty != 0 {
		request.FrequencyPenalty = req.FrequencyPenalty
	}

	// 检查是否为流式请求
	if req.Stream {
		h.handleStreamChat(c, *request)
		return
	}

	// 调用MiniMax服务
	response, err := h.service.ChatCompletion(*request)
	if err != nil {
		// 根据错误类型返回不同的状态码
		statusCode := http.StatusInternalServerError
		if h.service.IsRateLimited(response) {
			statusCode = http.StatusTooManyRequests
		} else if h.service.IsAuthFailed(response) {
			statusCode = http.StatusUnauthorized
		} else if h.service.IsInsufficientBalance(response) {
			statusCode = http.StatusPaymentRequired
		} else if h.service.IsTokenLimited(response) {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ErrorResponse{
			Code:    response.BaseResp.StatusCode,
			Message: response.BaseResp.StatusMsg,
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data": ChatResponse{
			Content: response.GetContent(),
			Usage:   h.service.GetUsage(response),
		},
	})
}

// handleStreamChat 处理流式聊天
func (h *Handler) handleStreamChat(c *gin.Context, request ChatCompletionRequest) {
	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// 获取流式响应
	responseChan, err := h.service.ChatCompletionStream(request)
	if err != nil {
		c.SSEvent("error", ErrorResponse{
			Code:    500,
			Message: "Failed to start stream",
			Details: err.Error(),
		})
		return
	}

	// 发送流式响应
	for response := range responseChan {
		// 检查是否有错误
		if !response.IsSuccess() {
			err := response.GetError()
			c.SSEvent("error", ErrorResponse{
				Code:    err.Code,
				Message: err.Message,
			})
			break
		}

		// 发送数据
		if len(response.Choices) > 0 {
			choice := response.Choices[0]
			if choice.Delta != nil && choice.Delta.Content != "" {
				c.SSEvent("message", gin.H{
					"content": choice.Delta.Content,
					"index":   choice.Index,
				})
			}
		}

		// 检查是否完成
		if len(response.Choices) > 0 && response.Choices[0].FinishReason == "stop" {
			c.SSEvent("done", gin.H{
				"usage": response.Usage,
			})
			break
		}
	}
}

// SimpleChat 简单聊天接口
func (h *Handler) SimpleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// 调用简单聊天方法
	var content string
	var err error

	if req.Temperature > 0 || req.MaxTokens > 0 {
		// 使用带参数的聊天
		temp := req.Temperature
		if temp == 0 {
			temp = 0.7 // 默认温度
		}
		maxTokens := req.MaxTokens
		if maxTokens == 0 {
			maxTokens = 2048 // 默认最大token数
		}
		content, err = h.service.SimpleChatWithParams(req.Message, temp, maxTokens)
	} else {
		// 使用简单聊天
		content, err = h.service.SimpleChat(req.Message)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "Failed to get AI response",
			Details: err.Error(),
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
		ai.POST("/chat", h.Chat)              // 完整聊天接口（支持流式）
		ai.POST("/chat/simple", h.SimpleChat) // 简单聊天接口
	}
}
