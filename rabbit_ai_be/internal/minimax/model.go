package minimax

// ChatCompletionRequest MiniMax聊天完成请求
type ChatCompletionRequest struct {
	Model             string        `json:"model"`
	Messages          []ChatMessage `json:"messages"`
	Stream            bool          `json:"stream,omitempty"`             // 是否流式响应
	Temperature       float64       `json:"temperature,omitempty"`        // 温度参数，控制随机性 (0.0-2.0)
	ToolChoices       []ToolChoice  `json:"tool_choices,omitempty"`       // 工具选择
	MaxTokens         int           `json:"max_tokens,omitempty"`         // 最大token数
	TopP              float64       `json:"top_p,omitempty"`              // 核采样参数 (0.0-1.0)
	TopK              int           `json:"top_k,omitempty"`              // Top-K采样
	RepetitionPenalty float64       `json:"repetition_penalty,omitempty"` // 重复惩罚参数
	Stop              []string      `json:"stop,omitempty"`               // 停止词
	PresencePenalty   float64       `json:"presence_penalty,omitempty"`   // 存在惩罚
	FrequencyPenalty  float64       `json:"frequency_penalty,omitempty"`  // 频率惩罚
	User              string        `json:"user,omitempty"`               // 用户标识
}

// ToolChoice 工具选择
type ToolChoice struct {
	Type     string `json:"type"` // "none", "auto", "function"
	Function *struct {
		Name string `json:"name"` // 函数名称
	} `json:"function,omitempty"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role      string     `json:"role"`                 // system, user, assistant, tool
	Name      string     `json:"name"`                 // 可选字段
	Content   string     `json:"content"`              // 消息内容
	ToolCalls []ToolCall `json:"tool_calls,omitempty"` // 工具调用
}

// ToolCall 工具调用
type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// BaseResponse 基础响应
type BaseResponse struct {
	StatusCode int    `json:"status_code"` // 状态码
	StatusMsg  string `json:"status_msg"`  // 错误详情
}

// ChatCompletionResponse MiniMax聊天完成响应
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
	// 敏感信息相关字段
	InputSensitive      bool         `json:"input_sensitive"`
	OutputSensitive     bool         `json:"output_sensitive"`
	InputSensitiveType  int          `json:"input_sensitive_type"`
	OutputSensitiveType int          `json:"output_sensitive_type"`
	OutputSensitiveInt  int          `json:"output_sensitive_int"`
	BaseResp            BaseResponse `json:"base_resp"`
}

// Choice 选择项
type Choice struct {
	FinishReason string       `json:"finish_reason"`
	Index        int          `json:"index"`
	Message      ChatMessage  `json:"message"`
	Delta        *ChatMessage `json:"delta,omitempty"` // 流式响应增量
}

// Usage 使用统计
type Usage struct {
	TotalTokens      int `json:"total_tokens"`
	TotalCharacters  int `json:"total_characters"`
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

// MiniMaxError MiniMax错误码
type MiniMaxError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 错误码常量
const (
	ErrorUnknown       = 1000 // 未知错误
	ErrorTimeout       = 1001 // 请求超时
	ErrorRateLimit     = 1002 // 触发RPM限流
	ErrorAuthFailed    = 1004 // 鉴权失败
	ErrorInsufficient  = 1008 // 余额不足
	ErrorInternal      = 1013 // 服务内部错误
	ErrorOutput        = 1027 // 输出内容错误
	ErrorTokenLimit    = 1039 // Token限制
	ErrorInvalidParams = 2013 // 参数错误
)

// GetErrorMessage 根据错误码获取错误信息
func GetErrorMessage(code int) string {
	switch code {
	case ErrorUnknown:
		return "未知错误"
	case ErrorTimeout:
		return "请求超时"
	case ErrorRateLimit:
		return "触发RPM限流"
	case ErrorAuthFailed:
		return "鉴权失败"
	case ErrorInsufficient:
		return "余额不足"
	case ErrorInternal:
		return "服务内部错误"
	case ErrorOutput:
		return "输出内容错误"
	case ErrorTokenLimit:
		return "Token限制"
	case ErrorInvalidParams:
		return "参数错误"
	default:
		return "未知错误"
	}
}

// IsSuccess 检查响应是否成功
func (r *ChatCompletionResponse) IsSuccess() bool {
	return r.BaseResp.StatusCode == 0
}

// GetError 获取错误信息
func (r *ChatCompletionResponse) GetError() *MiniMaxError {
	if r.IsSuccess() {
		return nil
	}
	return &MiniMaxError{
		Code:    r.BaseResp.StatusCode,
		Message: r.BaseResp.StatusMsg,
	}
}

// GetContent 获取回复内容
func (r *ChatCompletionResponse) GetContent() string {
	if len(r.Choices) == 0 {
		return ""
	}
	return r.Choices[0].Message.Content
}

// GetFinishReason 获取完成原因
func (r *ChatCompletionResponse) GetFinishReason() string {
	if len(r.Choices) == 0 {
		return ""
	}
	return r.Choices[0].FinishReason
}

// MiniMaxConfig MiniMax配置
type MiniMaxConfig struct {
	APIKey  string
	BaseURL string
}

// DefaultConfig 默认配置
func DefaultConfig(apiKey string) MiniMaxConfig {
	return MiniMaxConfig{
		APIKey:  apiKey,
		BaseURL: "https://api.minimaxi.com/v1",
	}
}

// NewChatCompletionRequest 创建聊天完成请求
func NewChatCompletionRequest(model string, messages []ChatMessage) *ChatCompletionRequest {
	return &ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 0.7,  // 默认温度
		MaxTokens:   2048, // 默认最大token数
		TopP:        0.9,  // 默认top_p
	}
}

// WithStream 设置流式响应
func (r *ChatCompletionRequest) WithStream(stream bool) *ChatCompletionRequest {
	r.Stream = stream
	return r
}

// WithTemperature 设置温度参数
func (r *ChatCompletionRequest) WithTemperature(temp float64) *ChatCompletionRequest {
	if temp >= 0.0 && temp <= 2.0 {
		r.Temperature = temp
	}
	return r
}

// WithMaxTokens 设置最大token数
func (r *ChatCompletionRequest) WithMaxTokens(maxTokens int) *ChatCompletionRequest {
	if maxTokens > 0 {
		r.MaxTokens = maxTokens
	}
	return r
}

// WithTopP 设置top_p参数
func (r *ChatCompletionRequest) WithTopP(topP float64) *ChatCompletionRequest {
	if topP >= 0.0 && topP <= 1.0 {
		r.TopP = topP
	}
	return r
}

// WithToolChoices 设置工具选择
func (r *ChatCompletionRequest) WithToolChoices(toolChoices []ToolChoice) *ChatCompletionRequest {
	r.ToolChoices = toolChoices
	return r
}

// WithStop 设置停止词
func (r *ChatCompletionRequest) WithStop(stop []string) *ChatCompletionRequest {
	r.Stop = stop
	return r
}

// WithUser 设置用户标识
func (r *ChatCompletionRequest) WithUser(user string) *ChatCompletionRequest {
	r.User = user
	return r
}
