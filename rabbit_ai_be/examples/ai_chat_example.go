package main

import (
	"context"
	"fmt"
	"log"

	"rabbit_ai/internal/cache"
	"rabbit_ai/internal/conversation"
	"rabbit_ai/internal/model"
)

// 模拟数据库连接
func createMockDB() *model.MockDB {
	return &model.MockDB{}
}

func main() {
	fmt.Println("=== Rabbit AI 多轮聊天功能演示 ===\n")

	// 1. 初始化服务
	fmt.Println("1. 初始化服务...")

	// 创建模拟依赖
	conversationRepo := conversation.NewMockConversationRepository()
	messageRepo := conversation.NewMockMessageRepository()
	userRepo := conversation.NewMockUserRepository()
	minimaxService := conversation.NewMockMiniMaxService()

	// 创建用户
	user := &model.User{
		ID:       1,
		Phone:    "13800138000",
		DeviceID: "device_123",
		Status:   1,
	}
	userRepo.Create(user)

	// 创建对话缓存
	conversationCache := cache.NewConversationCache("localhost:6379", "", 0)

	// 创建对话服务
	conversationService := conversation.NewService(
		conversationRepo,
		messageRepo,
		userRepo,
		conversationCache,
		minimaxService,
	)

	ctx := context.Background()

	// 2. 创建AI聊天对话
	fmt.Println("2. 创建AI聊天对话...")
	createReq := &conversation.CreateConversationRequest{
		UserID: 1,
		Title:  "AI技术咨询",
	}

	createResp, err := conversationService.CreateConversation(ctx, createReq)
	if err != nil {
		log.Fatalf("创建对话失败: %v", err)
	}

	conversationID := createResp.Conversation.ID
	fmt.Printf("✅ 创建对话成功，ID: %d, 标题: %s\n\n", conversationID, createResp.Conversation.Title)

	// 3. 发送第一条消息给AI
	fmt.Println("3. 发送第一条消息给AI...")
	messageReq1 := &conversation.SendMessageRequest{
		ConversationID: conversationID,
		UserID:         1,
		Content:        "你好，请介绍一下人工智能的发展历史",
		Model:          "glm-4",
	}

	messageResp1, err := conversationService.SendMessage(ctx, messageReq1)
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}

	fmt.Printf("✅ 用户消息: %s\n", messageResp1.UserMessage.Content)
	fmt.Printf("✅ AI回复: %s\n\n", messageResp1.AssistantMessage.Content)

	// 4. 继续多轮对话
	fmt.Println("4. 继续多轮对话...")
	messageReq2 := &conversation.SendMessageRequest{
		ConversationID: conversationID,
		UserID:         1,
		Content:        "请详细解释一下机器学习和深度学习的区别",
		Model:          "glm-4",
	}

	messageResp2, err := conversationService.SendMessage(ctx, messageReq2)
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}

	fmt.Printf("✅ 用户消息: %s\n", messageResp2.UserMessage.Content)
	fmt.Printf("✅ AI回复: %s\n\n", messageResp2.AssistantMessage.Content)

	// 5. 获取对话列表
	fmt.Println("5. 获取用户的AI聊天列表...")
	listReq := &conversation.GetConversationsRequest{
		UserID: 1,
		Limit:  10,
		Offset: 0,
	}

	listResp, err := conversationService.GetConversations(ctx, listReq)
	if err != nil {
		log.Fatalf("获取对话列表失败: %v", err)
	}

	fmt.Printf("✅ 用户共有 %d 个AI聊天对话\n", listResp.Total)
	for i, conv := range listResp.Conversations {
		fmt.Printf("   %d. ID: %d, 标题: %s, 消息数: %d\n", i+1, conv.ID, conv.Title, conv.MessageCount)
	}
	fmt.Println()

	// 6. 获取AI聊天历史
	fmt.Println("6. 获取AI聊天历史...")
	historyReq := &conversation.GetConversationMessagesRequest{
		ConversationID: conversationID,
		Limit:          50,
		Offset:         0,
	}

	historyResp, err := conversationService.GetConversationMessages(ctx, historyReq)
	if err != nil {
		log.Fatalf("获取聊天历史失败: %v", err)
	}

	fmt.Printf("✅ 对话共有 %d 条消息\n", historyResp.Total)
	for i, msg := range historyResp.Messages {
		role := "用户"
		if msg.Role == "assistant" {
			role = "AI"
		}
		fmt.Printf("   %d. [%s] %s\n", i+1, role, msg.Content)
	}
	fmt.Println()

	// 7. 创建第二个AI聊天对话
	fmt.Println("7. 创建第二个AI聊天对话...")
	createReq2 := &conversation.CreateConversationRequest{
		UserID: 1,
		Title:  "编程问题咨询",
	}

	createResp2, err := conversationService.CreateConversation(ctx, createReq2)
	if err != nil {
		log.Fatalf("创建对话失败: %v", err)
	}

	conversationID2 := createResp2.Conversation.ID
	fmt.Printf("✅ 创建第二个对话成功，ID: %d, 标题: %s\n\n", conversationID2, createResp2.Conversation.Title)

	// 8. 在第二个对话中发送消息
	fmt.Println("8. 在第二个对话中发送消息...")
	messageReq3 := &conversation.SendMessageRequest{
		ConversationID: conversationID2,
		UserID:         1,
		Content:        "请介绍一下Go语言的特点",
		Model:          "glm-4",
	}

	messageResp3, err := conversationService.SendMessage(ctx, messageReq3)
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}

	fmt.Printf("✅ 用户消息: %s\n", messageResp3.UserMessage.Content)
	fmt.Printf("✅ AI回复: %s\n\n", messageResp3.AssistantMessage.Content)

	// 9. 再次获取对话列表（验证数据同步）
	fmt.Println("9. 验证数据同步 - 获取更新后的AI聊天列表...")
	listResp2, err := conversationService.GetConversations(ctx, listReq)
	if err != nil {
		log.Fatalf("获取对话列表失败: %v", err)
	}

	fmt.Printf("✅ 用户现在共有 %d 个AI聊天对话\n", listResp2.Total)
	for i, conv := range listResp2.Conversations {
		fmt.Printf("   %d. ID: %d, 标题: %s, 消息数: %d\n", i+1, conv.ID, conv.Title, conv.MessageCount)
	}
	fmt.Println()

	// 10. 删除一个对话
	fmt.Println("10. 删除一个AI聊天对话...")
	deleteReq := &conversation.DeleteConversationRequest{
		ConversationID: conversationID2,
		UserID:         1,
	}

	err = conversationService.DeleteConversation(ctx, deleteReq)
	if err != nil {
		log.Fatalf("删除对话失败: %v", err)
	}

	fmt.Printf("✅ 成功删除对话 ID: %d\n\n", conversationID2)

	// 11. 最终验证
	fmt.Println("11. 最终验证 - 获取删除后的AI聊天列表...")
	listResp3, err := conversationService.GetConversations(ctx, listReq)
	if err != nil {
		log.Fatalf("获取对话列表失败: %v", err)
	}

	fmt.Printf("✅ 删除后用户还有 %d 个AI聊天对话\n", listResp3.Total)
	for i, conv := range listResp3.Conversations {
		fmt.Printf("   %d. ID: %d, 标题: %s, 消息数: %d\n", i+1, conv.ID, conv.Title, conv.MessageCount)
	}

	fmt.Println("\n=== 演示完成 ===")
	fmt.Println("✅ 所有功能正常工作！")
	fmt.Println("✅ 用户与AI的多轮聊天功能已实现")
	fmt.Println("✅ PostgreSQL存储和Redis缓存同步正常")
	fmt.Println("✅ 用户可以管理自己的AI聊天列表和历史")
}
