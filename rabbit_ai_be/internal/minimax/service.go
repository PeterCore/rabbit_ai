package minimax

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MiniMaxService MiniMax AI服务
type MiniMaxService struct {
	config MiniMaxConfig
	client *http.Client
}

// NewMiniMaxService 创建MiniMax服务实例
func NewMiniMaxService(config MiniMaxConfig) *MiniMaxService {
	return &MiniMaxService{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ChatCompletion 聊天完成
func (s *MiniMaxService) ChatCompletion(request ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// 构建请求URL
	url := fmt.Sprintf("%s/text/chatcompletion_v2", s.config.BaseURL)

	// 序列化请求体
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.APIKey))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 检查MiniMax API错误
	if !response.IsSuccess() {
		err := response.GetError()
		return &response, fmt.Errorf("MiniMax API error: %d - %s", err.Code, err.Message)
	}

	return &response, nil
}

// ChatCompletionStream 流式聊天完成
func (s *MiniMaxService) ChatCompletionStream(request ChatCompletionRequest) (<-chan ChatCompletionResponse, error) {
	// 确保启用流式响应
	request.Stream = true

	// 构建请求URL
	url := fmt.Sprintf("%s/text/chatcompletion_v2", s.config.BaseURL)

	// 序列化请求体
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.APIKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 创建响应通道
	responseChan := make(chan ChatCompletionResponse, 10)

	// 在goroutine中处理流式响应
	go func() {
		defer resp.Body.Close()
		defer close(responseChan)

		reader := bufio.NewReader(resp.Body)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				// 发送错误响应
				responseChan <- ChatCompletionResponse{
					BaseResp: BaseResponse{
						StatusCode: -1,
						StatusMsg:  fmt.Sprintf("Stream read error: %v", err),
					},
				}
				return
			}

			// 处理SSE格式的数据
			if len(line) > 6 && line[:6] == "data: " {
				data := line[6:]
				if data == "[DONE]" {
					break
				}

				// 解析JSON响应
				var response ChatCompletionResponse
				if err := json.Unmarshal([]byte(data), &response); err != nil {
					continue // 跳过无效的JSON
				}

				responseChan <- response
			}
		}
	}()

	return responseChan, nil
}

// SimpleChat 简单聊天（便捷方法）
func (s *MiniMaxService) SimpleChat(userMessage string) (string, error) {
	request := NewChatCompletionRequest("MiniMax-M1", []ChatMessage{
		{
			Role:    "system",
			Name:    "MiniMax AI",
			Content: "",
		},
		{
			Role:    "user",
			Name:    "用户",
			Content: userMessage,
		},
	})

	response, err := s.ChatCompletion(*request)
	if err != nil {
		return "", err
	}

	// 返回第一个选择的回复内容
	return response.GetContent(), nil
}

// SimpleChatWithParams 带参数的简单聊天
func (s *MiniMaxService) SimpleChatWithParams(userMessage string, temperature float64, maxTokens int) (string, error) {
	request := NewChatCompletionRequest("MiniMax-M1", []ChatMessage{
		{
			Role:    "system",
			Name:    "MiniMax AI",
			Content: "",
		},
		{
			Role:    "user",
			Name:    "用户",
			Content: userMessage,
		},
	}).WithTemperature(temperature).WithMaxTokens(maxTokens)

	response, err := s.ChatCompletion(*request)
	if err != nil {
		return "", err
	}

	// 返回第一个选择的回复内容
	return response.GetContent(), nil
}

// GetResponseContent 获取响应内容
func (s *MiniMaxService) GetResponseContent(response *ChatCompletionResponse) (string, error) {
	if !response.IsSuccess() {
		err := response.GetError()
		return "", fmt.Errorf("MiniMax API error: %d - %s", err.Code, err.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	return response.Choices[0].Message.Content, nil
}

// GetUsage 获取使用统计
func (s *MiniMaxService) GetUsage(response *ChatCompletionResponse) *Usage {
	return &response.Usage
}

// IsRateLimited 检查是否被限流
func (s *MiniMaxService) IsRateLimited(response *ChatCompletionResponse) bool {
	return response.BaseResp.StatusCode == ErrorRateLimit
}

// IsAuthFailed 检查是否认证失败
func (s *MiniMaxService) IsAuthFailed(response *ChatCompletionResponse) bool {
	return response.BaseResp.StatusCode == ErrorAuthFailed
}

// IsInsufficientBalance 检查是否余额不足
func (s *MiniMaxService) IsInsufficientBalance(response *ChatCompletionResponse) bool {
	return response.BaseResp.StatusCode == ErrorInsufficient
}

// IsTokenLimited 检查是否Token限制
func (s *MiniMaxService) IsTokenLimited(response *ChatCompletionResponse) bool {
	return response.BaseResp.StatusCode == ErrorTokenLimit
}
