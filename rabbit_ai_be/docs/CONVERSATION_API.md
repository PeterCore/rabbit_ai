# 多轮对话 API 文档

## 概述

本文档描述了 Rabbit AI 多轮对话功能的 API 接口。该功能支持用户创建对话会话、发送消息、获取对话历史等功能。

## 功能特性

- ✅ 每个用户（通过手机号和设备）对应自己的对话列表
- ✅ 每个用户都有自己的对话列表
- ✅ PostgreSQL 存储对话列表，Redis 缓存对话列表
- ✅ PostgreSQL 和 Redis 数据同步
- ✅ 每个用户可以获取对应对话历史内容
- ✅ 支持多轮对话上下文
- ✅ 自动生成对话标题
- ✅ 软删除对话

## 认证

所有对话相关的 API 都需要 JWT 认证。请在请求头中包含：

```
Authorization: Bearer <your_jwt_token>
```

## API 接口

### 1. 创建对话

**POST** `/api/v1/conversations`

创建新的对话会话。

#### 请求参数

```json
{
  "title": "新对话标题"
}
```

#### 响应示例

```json
{
  "success": true,
  "data": {
    "conversation": {
      "id": 1,
      "user_id": 123,
      "title": "新对话标题",
      "status": 1,
      "message_count": 0,
      "last_message_at": "2024-01-01T12:00:00Z",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  }
}
```

### 2. 获取对话列表

**GET** `/api/v1/conversations?limit=20&offset=0`

获取当前用户的对话列表。

#### 查询参数

- `limit`: 每页数量（默认 20，最大 100）
- `offset`: 偏移量（默认 0）

#### 响应示例

```json
{
  "success": true,
  "data": {
    "conversations": [
      {
        "id": 1,
        "user_id": 123,
        "title": "关于Go语言的问题",
        "status": 1,
        "message_count": 6,
        "last_message_at": "2024-01-01T12:30:00Z",
        "created_at": "2024-01-01T12:00:00Z",
        "updated_at": "2024-01-01T12:30:00Z"
      }
    ],
    "total": 1
  }
}
```

### 3. 获取对话消息

**GET** `/api/v1/conversations/{conversation_id}/messages?limit=50&offset=0`

获取指定对话的消息历史。

#### 路径参数

- `conversation_id`: 对话ID

#### 查询参数

- `limit`: 每页数量（默认 50，最大 200）
- `offset`: 偏移量（默认 0）

#### 响应示例

```json
{
  "success": true,
  "data": {
    "messages": [
      {
        "id": 1,
        "conversation_id": 1,
        "role": "user",
        "content": "你好，请介绍一下Go语言",
        "tokens": 15,
        "model": "glm-4",
        "finish_reason": null,
        "created_at": "2024-01-01T12:00:00Z"
      },
      {
        "id": 2,
        "conversation_id": 1,
        "role": "assistant",
        "content": "Go语言是由Google开发的一种静态类型、编译型语言...",
        "tokens": 120,
        "model": "glm-4",
        "finish_reason": "stop",
        "created_at": "2024-01-01T12:00:05Z"
      }
    ],
    "total": 2
  }
}
```

### 4. 发送消息

**POST** `/api/v1/conversations/{conversation_id}/messages`

向指定对话发送消息并获取AI回复。

#### 路径参数

- `conversation_id`: 对话ID

#### 请求参数

```json
{
  "content": "请详细解释Go语言的并发特性",
  "model": "glm-4"
}
```

#### 响应示例

```json
{
  "success": true,
  "data": {
    "user_message": {
      "id": 3,
      "conversation_id": 1,
      "role": "user",
      "content": "请详细解释Go语言的并发特性",
      "tokens": 20,
      "model": "glm-4",
      "finish_reason": null,
      "created_at": "2024-01-01T12:35:00Z"
    },
    "assistant_message": {
      "id": 4,
      "conversation_id": 1,
      "role": "assistant",
      "content": "Go语言的并发特性主要通过goroutine和channel实现...",
      "tokens": 150,
      "model": "glm-4",
      "finish_reason": "stop",
      "created_at": "2024-01-01T12:35:05Z"
    },
    "conversation": {
      "id": 1,
      "user_id": 123,
      "title": "关于Go语言的问题",
      "status": 1,
      "message_count": 4,
      "last_message_at": "2024-01-01T12:35:05Z",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:35:05Z"
    }
  }
}
```

### 5. 删除对话

**DELETE** `/api/v1/conversations/{conversation_id}`

删除指定的对话（软删除）。

#### 路径参数

- `conversation_id`: 对话ID

#### 响应示例

```json
{
  "success": true,
  "message": "Conversation deleted successfully"
}
```

## 错误处理

### 错误响应格式

```json
{
  "error": "错误描述",
  "details": "详细错误信息"
}
```

### 常见错误码

- `400`: 请求参数错误
- `401`: 未认证
- `403`: 权限不足
- `404`: 资源不存在
- `500`: 服务器内部错误

### 错误示例

```json
{
  "error": "Failed to send message",
  "details": "conversation does not belong to user"
}
```

## 数据模型

### Conversation（对话）

```go
type Conversation struct {
    ID             int64     `json:"id"`
    UserID         int64     `json:"user_id"`
    Title          string    `json:"title"`
    Status         int       `json:"status"`         // 1: 活跃, 0: 已删除
    MessageCount   int       `json:"message_count"`
    LastMessageAt  time.Time `json:"last_message_at"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}
```

### Message（消息）

```go
type Message struct {
    ID             int64     `json:"id"`
    ConversationID int64     `json:"conversation_id"`
    Role           string    `json:"role"`           // user/assistant
    Content        string    `json:"content"`
    Tokens         int       `json:"tokens"`
    Model          string    `json:"model"`
    FinishReason   string    `json:"finish_reason"`
    CreatedAt      time.Time `json:"created_at"`
}
```

## 缓存策略

### Redis 缓存

- **对话信息**: 缓存 30 分钟
- **消息信息**: 缓存 1 小时
- **用户对话列表**: 缓存 15 分钟
- **对话消息列表**: 缓存 30 分钟

### 缓存失效

- 创建新对话时，使用户对话列表缓存失效
- 发送消息时，使对话相关缓存失效
- 删除对话时，使相关缓存失效

## 使用示例

### 完整对话流程

1. **创建对话**
```bash
curl -X POST http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"title": "技术咨询"}'
```

2. **发送消息**
```bash
curl -X POST http://localhost:8080/api/v1/conversations/1/messages \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"content": "请介绍一下微服务架构", "model": "glm-4"}'
```

3. **获取对话历史**
```bash
curl -X GET http://localhost:8080/api/v1/conversations/1/messages \
  -H "Authorization: Bearer <your_token>"
```

4. **获取对话列表**
```bash
curl -X GET http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer <your_token>"
```

## 注意事项

1. 所有对话相关的操作都需要用户认证
2. 用户只能访问自己的对话
3. 对话标题会自动根据第一条用户消息生成
4. 消息按时间顺序排列
5. 支持分页查询，避免数据量过大
6. 使用软删除，删除的对话不会真正从数据库中移除
7. 缓存会自动失效，确保数据一致性 