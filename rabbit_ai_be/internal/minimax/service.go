package minimax

import (
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

	return &response, nil
}

// SimpleChat 简单聊天（便捷方法）
func (s *MiniMaxService) SimpleChat(userMessage string) (string, error) {
	request := ChatCompletionRequest{
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
				Content: userMessage,
			},
		},
	}

	response, err := s.ChatCompletion(request)
	if err != nil {
		return "", err
	}

	// 返回第一个选择的回复内容
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from AI")
}

// GetResponseContent 获取响应内容
func (s *MiniMaxService) GetResponseContent(response *ChatCompletionResponse) (string, error) {
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	return response.Choices[0].Message.Content, nil
}

// GetUsage 获取使用统计
func (s *MiniMaxService) GetUsage(response *ChatCompletionResponse) *Usage {
	return &response.Usage
}
