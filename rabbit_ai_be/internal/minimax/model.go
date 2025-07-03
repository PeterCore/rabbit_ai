package minimax

// ChatCompletionRequest MiniMax聊天完成请求
type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Name    string `json:"name"`    // 可选字段
	Content string `json:"content"` // 消息内容
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
	InputSensitive      bool `json:"input_sensitive"`
	OutputSensitive     bool `json:"output_sensitive"`
	InputSensitiveType  int  `json:"input_sensitive_type"`
	OutputSensitiveType int  `json:"output_sensitive_type"`
	OutputSensitiveInt  int  `json:"output_sensitive_int"`
	BaseResp            struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp"`
}

// Choice 选择项
type Choice struct {
	FinishReason string      `json:"finish_reason"`
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
}

// Usage 使用统计
type Usage struct {
	TotalTokens      int `json:"total_tokens"`
	TotalCharacters  int `json:"total_characters"`
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
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
