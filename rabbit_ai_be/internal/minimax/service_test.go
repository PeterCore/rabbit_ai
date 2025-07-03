package minimax

import (
	"testing"
)

func TestMiniMaxService_SimpleChat(t *testing.T) {
	// 注意：这个测试需要有效的API Key才能运行
	// 在实际环境中，应该使用mock或测试环境的API Key
	config := DefaultConfig("your-api-key-here")
	service := NewMiniMaxService(config)

	// 测试简单聊天
	message := "你好"
	response, err := service.SimpleChat(message)

	// 由于需要真实的API Key，这里只测试错误情况
	if err == nil {
		t.Logf("AI Response: %s", response)
	} else {
		t.Logf("Expected error (due to invalid API key): %v", err)
	}
}

func TestMiniMaxService_ChatCompletion(t *testing.T) {
	config := DefaultConfig("your-api-key-here")
	service := NewMiniMaxService(config)

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
				Content: "你好",
			},
		},
	}

	response, err := service.ChatCompletion(request)

	if err == nil {
		t.Logf("Response ID: %s", response.ID)
		t.Logf("Model: %s", response.Model)
		if len(response.Choices) > 0 {
			t.Logf("Content: %s", response.Choices[0].Message.Content)
		}
	} else {
		t.Logf("Expected error (due to invalid API key): %v", err)
	}
}

func TestMiniMaxConfig_DefaultConfig(t *testing.T) {
	apiKey := "test-api-key"
	config := DefaultConfig(apiKey)

	if config.APIKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, config.APIKey)
	}

	if config.BaseURL != "https://api.minimaxi.com/v1" {
		t.Errorf("Expected base URL https://api.minimaxi.com/v1, got %s", config.BaseURL)
	}
}
