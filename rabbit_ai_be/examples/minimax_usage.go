package main

import (
	"fmt"
	"log"
	"os"

	"rabbit_ai/internal/minimax"
)

func main() {
	// 从环境变量获取API Key
	apiKey := os.Getenv("MINIMAX_API_KEY")
	if apiKey == "" {
		log.Fatal("MINIMAX_API_KEY environment variable is required")
	}

	// 创建MiniMax服务
	config := minimax.DefaultConfig(apiKey)
	service := minimax.NewMiniMaxService(config)

	fmt.Println("=== MiniMax AI 使用示例 ===\n")

	// 示例1: 简单聊天
	fmt.Println("1. 简单聊天:")
	response, err := service.SimpleChat("你好，请介绍一下自己")
	if err != nil {
		log.Printf("简单聊天失败: %v", err)
	} else {
		fmt.Printf("AI回复: %s\n\n", response)
	}

	// 示例2: 带参数的聊天
	fmt.Println("2. 带参数的聊天 (低温度，短回复):")
	response, err = service.SimpleChatWithParams("写一首关于春天的短诗", 0.3, 100)
	if err != nil {
		log.Printf("带参数聊天失败: %v", err)
	} else {
		fmt.Printf("AI回复: %s\n\n", response)
	}

	// 示例3: 完整聊天请求
	fmt.Println("3. 完整聊天请求:")
	messages := []minimax.ChatMessage{
		{
			Role:    "system",
			Content: "你是一个专业的Go语言开发者，请用简洁的语言回答问题。",
		},
		{
			Role:    "user",
			Content: "什么是goroutine？",
		},
	}

	request := minimax.NewChatCompletionRequest("MiniMax-M1", messages).
		WithTemperature(0.5).
		WithMaxTokens(500).
		WithTopP(0.8).
		WithUser("go-developer")

	responseObj, err := service.ChatCompletion(*request)
	if err != nil {
		log.Printf("完整聊天失败: %v", err)
	} else {
		// 使用新的响应方法
		content := responseObj.GetContent()
		finishReason := responseObj.GetFinishReason()

		fmt.Printf("AI回复: %s\n", content)
		fmt.Printf("完成原因: %s\n", finishReason)
		fmt.Printf("使用统计: 总token=%d, 提示token=%d, 完成token=%d\n",
			responseObj.Usage.TotalTokens,
			responseObj.Usage.PromptTokens,
			responseObj.Usage.CompletionTokens)

		// 检查响应状态
		if responseObj.IsSuccess() {
			fmt.Println("响应状态: 成功")
		} else {
			err := responseObj.GetError()
			fmt.Printf("响应状态: 失败 - %d: %s\n", err.Code, err.Message)
		}
		fmt.Println()
	}

	// 示例4: 流式聊天
	fmt.Println("4. 流式聊天:")
	streamRequest := minimax.NewChatCompletionRequest("MiniMax-M1", []minimax.ChatMessage{
		{
			Role:    "user",
			Content: "请写一个关于人工智能的短文",
		},
	}).WithStream(true).WithTemperature(0.7).WithMaxTokens(300)

	fmt.Print("AI回复 (流式): ")
	responseChan, err := service.ChatCompletionStream(*streamRequest)
	if err != nil {
		log.Printf("流式聊天失败: %v", err)
	} else {
		for response := range responseChan {
			// 检查错误
			if !response.IsSuccess() {
				err := response.GetError()
				fmt.Printf("\n流式响应错误: %d - %s\n", err.Code, err.Message)
				break
			}

			if len(response.Choices) > 0 {
				choice := response.Choices[0]
				if choice.Delta != nil && choice.Delta.Content != "" {
					fmt.Print(choice.Delta.Content)
				}
			}
		}
		fmt.Println("\n")
	}

	// 示例5: 使用工具选择
	fmt.Println("5. 使用工具选择:")
	toolRequest := minimax.NewChatCompletionRequest("MiniMax-M1", []minimax.ChatMessage{
		{
			Role:    "user",
			Content: "今天天气怎么样？",
		},
	}).WithToolChoices([]minimax.ToolChoice{
		{
			Type: "function",
			Function: &struct {
				Name string `json:"name"`
			}{
				Name: "get_weather",
			},
		},
	}).WithTemperature(0.3)

	responseObj, err = service.ChatCompletion(*toolRequest)
	if err != nil {
		log.Printf("工具选择聊天失败: %v", err)
	} else {
		content := responseObj.GetContent()
		fmt.Printf("AI回复: %s\n", content)

		// 检查敏感信息
		if responseObj.InputSensitive {
			fmt.Println("输入内容敏感")
		}
		if responseObj.OutputSensitive {
			fmt.Println("输出内容敏感")
		}
		fmt.Println()
	}

	// 示例6: 使用停止词
	fmt.Println("6. 使用停止词:")
	stopRequest := minimax.NewChatCompletionRequest("MiniMax-M1", []minimax.ChatMessage{
		{
			Role:    "user",
			Content: "请列出5个编程语言，每个语言用一句话描述",
		},
	}).WithStop([]string{"END", "STOP", "完成"}).
		WithTemperature(0.6).
		WithMaxTokens(200)

	responseObj, err = service.ChatCompletion(*stopRequest)
	if err != nil {
		log.Printf("停止词聊天失败: %v", err)
	} else {
		content := responseObj.GetContent()
		fmt.Printf("AI回复: %s\n", content)
		fmt.Printf("完成原因: %s\n", responseObj.GetFinishReason())
		fmt.Println()
	}

	// 示例7: 错误处理演示
	fmt.Println("7. 错误处理演示:")

	// 模拟各种错误情况
	errorResponses := []minimax.ChatCompletionResponse{
		{
			BaseResp: minimax.BaseResponse{
				StatusCode: minimax.ErrorRateLimit,
				StatusMsg:  "触发RPM限流",
			},
		},
		{
			BaseResp: minimax.BaseResponse{
				StatusCode: minimax.ErrorAuthFailed,
				StatusMsg:  "鉴权失败",
			},
		},
		{
			BaseResp: minimax.BaseResponse{
				StatusCode: minimax.ErrorInsufficient,
				StatusMsg:  "余额不足",
			},
		},
		{
			BaseResp: minimax.BaseResponse{
				StatusCode: minimax.ErrorTokenLimit,
				StatusMsg:  "Token限制",
			},
		},
	}

	for i, resp := range errorResponses {
		fmt.Printf("错误示例 %d:\n", i+1)
		if !resp.IsSuccess() {
			err := resp.GetError()
			fmt.Printf("  错误码: %d\n", err.Code)
			fmt.Printf("  错误信息: %s\n", err.Message)
			fmt.Printf("  错误描述: %s\n", minimax.GetErrorMessage(err.Code))

			// 检查具体错误类型
			if service.IsRateLimited(&resp) {
				fmt.Println("  类型: 限流错误")
			} else if service.IsAuthFailed(&resp) {
				fmt.Println("  类型: 认证错误")
			} else if service.IsInsufficientBalance(&resp) {
				fmt.Println("  类型: 余额不足")
			} else if service.IsTokenLimited(&resp) {
				fmt.Println("  类型: Token限制")
			}
		}
		fmt.Println()
	}

	fmt.Println("=== 示例完成 ===")
}
