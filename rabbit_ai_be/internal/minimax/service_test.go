package minimax

import (
	"fmt"
	"log"
	"testing"
)

func TestNewChatCompletionRequest(t *testing.T) {
	messages := []ChatMessage{
		{
			Role:    "user",
			Content: "你好",
		},
	}

	request := NewChatCompletionRequest("MiniMax-M1", messages)

	if request.Model != "MiniMax-M1" {
		t.Errorf("Expected model MiniMax-M1, got %s", request.Model)
	}

	if len(request.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(request.Messages))
	}

	if request.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", request.Temperature)
	}

	if request.MaxTokens != 2048 {
		t.Errorf("Expected max tokens 2048, got %d", request.MaxTokens)
	}

	if request.TopP != 0.9 {
		t.Errorf("Expected top_p 0.9, got %f", request.TopP)
	}
}

func TestChatCompletionRequest_WithMethods(t *testing.T) {
	messages := []ChatMessage{
		{
			Role:    "user",
			Content: "你好",
		},
	}

	request := NewChatCompletionRequest("MiniMax-M1", messages)

	// 测试链式调用
	request.WithTemperature(0.5).
		WithMaxTokens(1000).
		WithTopP(0.8).
		WithStream(true).
		WithUser("test-user").
		WithStop([]string{"END", "STOP"})

	if request.Temperature != 0.5 {
		t.Errorf("Expected temperature 0.5, got %f", request.Temperature)
	}

	if request.MaxTokens != 1000 {
		t.Errorf("Expected max tokens 1000, got %d", request.MaxTokens)
	}

	if request.TopP != 0.8 {
		t.Errorf("Expected top_p 0.8, got %f", request.TopP)
	}

	if !request.Stream {
		t.Error("Expected stream to be true")
	}

	if request.User != "test-user" {
		t.Errorf("Expected user test-user, got %s", request.User)
	}

	if len(request.Stop) != 2 {
		t.Errorf("Expected 2 stop words, got %d", len(request.Stop))
	}
}

func TestChatCompletionRequest_WithToolChoices(t *testing.T) {
	messages := []ChatMessage{
		{
			Role:    "user",
			Content: "你好",
		},
	}

	request := NewChatCompletionRequest("MiniMax-M1", messages)

	toolChoices := []ToolChoice{
		{
			Type: "function",
			Function: &struct {
				Name string `json:"name"`
			}{
				Name: "get_weather",
			},
		},
	}

	request.WithToolChoices(toolChoices)

	if len(request.ToolChoices) != 1 {
		t.Errorf("Expected 1 tool choice, got %d", len(request.ToolChoices))
	}

	if request.ToolChoices[0].Type != "function" {
		t.Errorf("Expected tool choice type function, got %s", request.ToolChoices[0].Type)
	}

	if request.ToolChoices[0].Function.Name != "get_weather" {
		t.Errorf("Expected function name get_weather, got %s", request.ToolChoices[0].Function.Name)
	}
}

func TestChatCompletionResponse_IsSuccess(t *testing.T) {
	// 测试成功响应
	successResponse := ChatCompletionResponse{
		BaseResp: BaseResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
	}

	if !successResponse.IsSuccess() {
		t.Error("Expected success response to be true")
	}

	// 测试失败响应
	failResponse := ChatCompletionResponse{
		BaseResp: BaseResponse{
			StatusCode: ErrorAuthFailed,
			StatusMsg:  "鉴权失败",
		},
	}

	if failResponse.IsSuccess() {
		t.Error("Expected fail response to be false")
	}
}

func TestChatCompletionResponse_GetError(t *testing.T) {
	// 测试成功响应
	successResponse := ChatCompletionResponse{
		BaseResp: BaseResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
	}

	err := successResponse.GetError()
	if err != nil {
		t.Errorf("Expected no error for success response, got %v", err)
	}

	// 测试失败响应
	failResponse := ChatCompletionResponse{
		BaseResp: BaseResponse{
			StatusCode: ErrorRateLimit,
			StatusMsg:  "触发RPM限流",
		},
	}

	err = failResponse.GetError()
	if err == nil {
		t.Error("Expected error for fail response")
		return
	}

	if err.Code != ErrorRateLimit {
		t.Errorf("Expected error code %d, got %d", ErrorRateLimit, err.Code)
	}

	if err.Message != "触发RPM限流" {
		t.Errorf("Expected error message '触发RPM限流', got %s", err.Message)
	}
}

func TestChatCompletionResponse_GetContent(t *testing.T) {
	// 测试有内容的响应
	response := ChatCompletionResponse{
		BaseResp: BaseResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		Choices: []Choice{
			{
				Message: ChatMessage{
					Content: "你好！有什么可以帮助你的吗？",
				},
			},
		},
	}

	content := response.GetContent()
	if content != "你好！有什么可以帮助你的吗？" {
		t.Errorf("Expected content '你好！有什么可以帮助你的吗？', got %s", content)
	}

	// 测试无内容的响应
	emptyResponse := ChatCompletionResponse{
		BaseResp: BaseResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		Choices: []Choice{},
	}

	content = emptyResponse.GetContent()
	if content != "" {
		t.Errorf("Expected empty content, got %s", content)
	}
}

func TestGetErrorMessage(t *testing.T) {
	testCases := []struct {
		code     int
		expected string
	}{
		{ErrorUnknown, "未知错误"},
		{ErrorTimeout, "请求超时"},
		{ErrorRateLimit, "触发RPM限流"},
		{ErrorAuthFailed, "鉴权失败"},
		{ErrorInsufficient, "余额不足"},
		{ErrorInternal, "服务内部错误"},
		{ErrorOutput, "输出内容错误"},
		{ErrorTokenLimit, "Token限制"},
		{ErrorInvalidParams, "参数错误"},
		{9999, "未知错误"}, // 未知错误码
	}

	for _, tc := range testCases {
		result := GetErrorMessage(tc.code)
		if result != tc.expected {
			t.Errorf("For code %d, expected '%s', got '%s'", tc.code, tc.expected, result)
		}
	}
}

func TestMiniMaxService_ErrorChecking(t *testing.T) {
	config := DefaultConfig("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJHcm91cE5hbWUiOiLlvKDmt7MiLCJVc2VyTmFtZSI6IuW8oOa3syIsIkFjY291bnQiOiIiLCJTdWJqZWN0SUQiOiIxOTM5NjUzODAyNzg5NDQ1NzQyIiwiUGhvbmUiOiIxMzc3NDY3OTAwMiIsIkdyb3VwSUQiOiIxOTM5NjUzODAyNzgxMDU3MTk1IiwiUGFnZU5hbWUiOiIiLCJNYWlsIjoiZnRleG1fcmVhbHRpbWVAMTYzLmNvbSIsIkNyZWF0ZVRpbWUiOiIyMDI1LTA3LTAxIDE1OjIwOjIyIiwiVG9rZW5UeXBlIjoxLCJpc3MiOiJtaW5pbWF4In0.Crc3SUXtT-pKMDoQ0YcIplJK537JZudp33L5BZn6WBBNYnx6F5LGsWYNo2KDhBXpb89VfEO5lBQwuRy0uDef5JBynKJWQ93KbjX1rC8FAPYIXBKJjXA1BFfxSjC_n1k5djvOyGzC9_42XIsjvpoiyI8YnnpwqGJMs2LTou2FV4ilZYq5McjIp-gBfGIyVtoyjkScX9TPxgWtFyoWyIvpNztjf3eUGNPkE5yKVNTEYu-DPBnp4MGdgoMAeubl1-Ayc8Txi3cp8DQDseNGLlCaPLsbObFD6lERPjNhGXjXPqtEEZiLTSHCoTwTgHfWEn8rvKWfbyTIEso-oFIwMNzKxA")
	service := NewMiniMaxService(config)

	// 测试限流检查
	rateLimitResponse := &ChatCompletionResponse{
		BaseResp: BaseResponse{
			StatusCode: ErrorRateLimit,
			StatusMsg:  "触发RPM限流",
		},
	}

	if !service.IsRateLimited(rateLimitResponse) {
		t.Error("Expected rate limit check to return true")
	}

	// 测试认证失败检查
	authFailResponse := &ChatCompletionResponse{
		BaseResp: BaseResponse{
			StatusCode: ErrorAuthFailed,
			StatusMsg:  "鉴权失败",
		},
	}

	if !service.IsAuthFailed(authFailResponse) {
		t.Error("Expected auth failed check to return true")
	}

	// 测试余额不足检查
	insufficientResponse := &ChatCompletionResponse{
		BaseResp: BaseResponse{
			StatusCode: ErrorInsufficient,
			StatusMsg:  "余额不足",
		},
	}

	if !service.IsInsufficientBalance(insufficientResponse) {
		t.Error("Expected insufficient balance check to return true")
	}

	// 测试Token限制检查
	tokenLimitResponse := &ChatCompletionResponse{
		BaseResp: BaseResponse{
			StatusCode: ErrorTokenLimit,
			StatusMsg:  "Token限制",
		},
	}

	if !service.IsTokenLimited(tokenLimitResponse) {
		t.Error("Expected token limit check to return true")
	}
}

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

func TestMiniMaxService_SimpleChatWithParams(t *testing.T) {
	config := DefaultConfig("your-api-key-here")
	service := NewMiniMaxService(config)

	// 测试带参数的聊天
	message := "你好"
	temperature := 0.5
	maxTokens := 300

	response, err := service.SimpleChatWithParams(message, temperature, maxTokens)

	if err == nil {
		t.Logf("AI Response with params: %s", response)
	} else {
		t.Logf("Expected error (due to invalid API key): %v", err)
	}
}

func TestMiniMaxService_ChatCompletion(t *testing.T) {
	config := DefaultConfig("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJHcm91cE5hbWUiOiLlvKDmt7MiLCJVc2VyTmFtZSI6IuW8oOa3syIsIkFjY291bnQiOiIiLCJTdWJqZWN0SUQiOiIxOTM5NjUzODAyNzg5NDQ1NzQyIiwiUGhvbmUiOiIxMzc3NDY3OTAwMiIsIkdyb3VwSUQiOiIxOTM5NjUzODAyNzgxMDU3MTk1IiwiUGFnZU5hbWUiOiIiLCJNYWlsIjoiZnRleG1fcmVhbHRpbWVAMTYzLmNvbSIsIkNyZWF0ZVRpbWUiOiIyMDI1LTA3LTAxIDE1OjIwOjIyIiwiVG9rZW5UeXBlIjoxLCJpc3MiOiJtaW5pbWF4In0.Crc3SUXtT-pKMDoQ0YcIplJK537JZudp33L5BZn6WBBNYnx6F5LGsWYNo2KDhBXpb89VfEO5lBQwuRy0uDef5JBynKJWQ93KbjX1rC8FAPYIXBKJjXA1BFfxSjC_n1k5djvOyGzC9_42XIsjvpoiyI8YnnpwqGJMs2LTou2FV4ilZYq5McjIp-gBfGIyVtoyjkScX9TPxgWtFyoWyIvpNztjf3eUGNPkE5yKVNTEYu-DPBnp4MGdgoMAeubl1-Ayc8Txi3cp8DQDseNGLlCaPLsbObFD6lERPjNhGXjXPqtEEZiLTSHCoTwTgHfWEn8rvKWfbyTIEso-oFIwMNzKxA")
	service := NewMiniMaxService(config)

	request := NewChatCompletionRequest("MiniMax-M1", []ChatMessage{
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
	}).WithTemperature(0.7).WithMaxTokens(500)

	response, err := service.ChatCompletion(*request)

	if err == nil {
		t.Logf("Response ID: %s", response.ID)
		t.Logf("Model: %s", response.Model)
		t.Logf("Content: %s", response.GetContent())
		t.Logf("Finish Reason: %s", response.GetFinishReason())
		t.Logf("Choices count: %d", len(response.Choices))
		if len(response.Choices) > 0 {
			t.Logf("First choice message: %+v", response.Choices[0].Message)
			t.Logf("First choice delta: %+v", response.Choices[0].Delta)
		}
		t.Logf("Usage: %+v", response.Usage)
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

func TestMiniMaxService_SimpleChat_Usage(t *testing.T) {
	// 创建MiniMax服务
	config := DefaultConfig("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJHcm91cE5hbWUiOiLlvKDmt7MiLCJVc2VyTmFtZSI6IuW8oOa3syIsIkFjY291bnQiOiIiLCJTdWJqZWN0SUQiOiIxOTM5NjUzODAyNzg5NDQ1NzQyIiwiUGhvbmUiOiIxMzc3NDY3OTAwMiIsIkdyb3VwSUQiOiIxOTM5NjUzODAyNzgxMDU3MTk1IiwiUGFnZU5hbWUiOiIiLCJNYWlsIjoiZnRleG1fcmVhbHRpbWVAMTYzLmNvbSIsIkNyZWF0ZVRpbWUiOiIyMDI1LTA3LTAxIDE1OjIwOjIyIiwiVG9rZW5UeXBlIjoxLCJpc3MiOiJtaW5pbWF4In0.Crc3SUXtT-pKMDoQ0YcIplJK537JZudp33L5BZn6WBBNYnx6F5LGsWYNo2KDhBXpb89VfEO5lBQwuRy0uDef5JBynKJWQ93KbjX1rC8FAPYIXBKJjXA1BFfxSjC_n1k5djvOyGzC9_42XIsjvpoiyI8YnnpwqGJMs2LTou2FV4ilZYq5McjIp-gBfGIyVtoyjkScX9TPxgWtFyoWyIvpNztjf3eUGNPkE5yKVNTEYu-DPBnp4MGdgoMAeubl1-Ayc8Txi3cp8DQDseNGLlCaPLsbObFD6lERPjNhGXjXPqtEEZiLTSHCoTwTgHfWEn8rvKWfbyTIEso-oFIwMNzKxA")
	service := NewMiniMaxService(config)

	// 简单聊天
	message := "你好"
	response, err := service.SimpleChat(message)
	if err != nil {
		log.Printf("聊天失败: %v", err)
	} else {
		fmt.Printf("AI回复: %s\n", response)
	}
}

func TestMiniMaxService_ChatCompletionStream(t *testing.T) {
	config := DefaultConfig("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJHcm91cE5hbWUiOiLlvKDmt7MiLCJVc2VyTmFtZSI6IuW8oOa3syIsIkFjY291bnQiOiIiLCJTdWJqZWN0SUQiOiIxOTM5NjUzODAyNzg5NDQ1NzQyIiwiUGhvbmUiOiIxMzc3NDY3OTAwMiIsIkdyb3VwSUQiOiIxOTM5NjUzODAyNzgxMDU3MTk1IiwiUGFnZU5hbWUiOiIiLCJNYWlsIjoiZnRleG1fcmVhbHRpbWVAMTYzLmNvbSIsIkNyZWF0ZVRpbWUiOiIyMDI1LTA3LTAxIDE1OjIwOjIyIiwiVG9rZW5UeXBlIjoxLCJpc3MiOiJtaW5pbWF4In0.Crc3SUXtT-pKMDoQ0YcIplJK537JZudp33L5BZn6WBBNYnx6F5LGsWYNo2KDhBXpb89VfEO5lBQwuRy0uDef5JBynKJWQ93KbjX1rC8FAPYIXBKJjXA1BFfxSjC_n1k5djvOyGzC9_42XIsjvpoiyI8YnnpwqGJMs2LTou2FV4ilZYq5McjIp-gBfGIyVtoyjkScX9TPxgWtFyoWyIvpNztjf3eUGNPkE5yKVNTEYu-DPBnp4MGdgoMAeubl1-Ayc8Txi3cp8DQDseNGLlCaPLsbObFD6lERPjNhGXjXPqtEEZiLTSHCoTwTgHfWEn8rvKWfbyTIEso-oFIwMNzKxA")
	service := NewMiniMaxService(config)

	messages := []ChatMessage{
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
	}

	toolChoices := []ToolChoice{
		{
			Type: "function",
			Function: &struct {
				Name string `json:"name"`
			}{
				Name: "get_weather",
			},
		},
	}

	request := NewChatCompletionRequest("MiniMax-M1", messages).
		WithTemperature(0.7).
		WithMaxTokens(2048).
		WithStream(true).
		WithToolChoices(toolChoices).
		WithStop([]string{"END", "STOP"}).
		WithUser("user123")

	responseChan, err := service.ChatCompletionStream(*request)
	if err != nil {
		t.Errorf("Expected error (due to invalid API key): %v", err)
		return
	}

	for response := range responseChan {
		// 处理流式数据
		t.Logf("收到流式响应: %+v", response)
	}
}
