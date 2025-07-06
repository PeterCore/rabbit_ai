# Rabbit AI Backend API 文档

## 概述

Rabbit AI Backend 是一个基于 Go 和 Gin 框架的后端服务，提供用户认证、设备管理、缓存管理和 AI 聊天等功能。

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **Content-Type**: `application/json`

## 认证

大部分 API 需要在请求头中包含 JWT Token：

```
Authorization: Bearer <your-jwt-token>
```

## API 端点

### 用户管理

#### 1. 用户注册

```http
POST /auth/register
```

**请求体:**
```json
{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "User registered successfully",
  "data": {
    "user_id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 2. 用户登录

```http
POST /auth/login
```

**请求体:**
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    }
  }
}
```

#### 3. GitHub 登录

```http
POST /auth/github
```

**请求体:**
```json
{
  "code": "github_oauth_code"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "GitHub login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "github_user",
      "email": "user@github.com",
      "github_id": "123456",
      "github_email": "user@github.com"
    }
  }
}
```

### 设备管理

#### 1. 绑定设备

```http
POST /device/bind
```

**请求体:**
```json
{
  "device_id": "unique-device-identifier",
  "platform": "iOS"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Device bound successfully",
  "data": {
    "device_id": "unique-device-identifier",
    "platform": "iOS",
    "user_id": 1,
    "bound_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 2. 获取设备信息

```http
GET /device/info
```

**响应:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "device_id": "unique-device-identifier",
    "platform": "iOS",
    "user_id": 1,
    "bound_at": "2024-01-01T00:00:00Z"
  }
}
```

### 缓存管理

#### 1. 获取缓存统计

```http
GET /cache/stats
```

**响应:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "total_keys": 100,
    "memory_usage": "10MB",
    "hit_rate": 0.85
  }
}
```

#### 2. 清除缓存

```http
DELETE /cache/clear
```

**响应:**
```json
{
  "code": 200,
  "message": "Cache cleared successfully"
}
```

### MiniMax AI 聊天

#### 1. 完整聊天接口（支持流式响应）

```http
POST /ai/chat
```

**请求体:**
```json
{
  "message": "你好，请介绍一下自己",
  "temperature": 0.7,
  "max_tokens": 2048,
  "top_p": 0.9,
  "stream": false,
  "tool_choices": [
    {
      "type": "function",
      "function": {
        "name": "get_weather"
      }
    }
  ],
  "stop": ["END", "STOP"],
  "user": "user123",
  "repetition_penalty": 1.1,
  "presence_penalty": 0.0,
  "frequency_penalty": 0.0
}
```

**参数说明:**
- `message` (必需): 用户消息内容
- `temperature` (可选): 温度参数，控制随机性 (0.0-2.0)，默认 0.7
- `max_tokens` (可选): 最大生成token数，默认 2048
- `top_p` (可选): 核采样参数 (0.0-1.0)，默认 0.9
- `stream` (可选): 是否启用流式响应，默认 false
- `tool_choices` (可选): 工具选择列表
- `stop` (可选): 停止词列表
- `user` (可选): 用户标识
- `repetition_penalty` (可选): 重复惩罚参数
- `presence_penalty` (可选): 存在惩罚参数
- `frequency_penalty` (可选): 频率惩罚参数

**普通响应:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "content": "你好！我是 MiniMax AI，一个智能助手...",
    "usage": {
      "total_tokens": 150,
      "total_characters": 300,
      "prompt_tokens": 50,
      "completion_tokens": 100
    }
  }
}
```

**流式响应 (SSE):**
当 `stream: true` 时，响应为 Server-Sent Events 格式：

```
event: message
data: {"content": "你好", "index": 0}

event: message
data: {"content": "！我是", "index": 0}

event: message
data: {"content": " MiniMax AI", "index": 0}

event: done
data: {"usage": {"total_tokens": 150, "prompt_tokens": 50, "completion_tokens": 100}}
```

#### 2. 简单聊天接口

```http
POST /ai/chat/simple
```

**请求体:**
```json
{
  "message": "你好",
  "temperature": 0.7,
  "max_tokens": 2048
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "content": "你好！有什么可以帮助你的吗？"
  }
}
```

## 错误响应

### 通用错误格式

所有 API 在发生错误时都会返回以下格式：

```json
{
  "code": 400,
  "message": "Error description",
  "details": "Additional error details"
}
```

### MiniMax AI 错误码

MiniMax AI 接口可能返回以下错误码：

| 错误码 | 说明 | HTTP状态码 |
|--------|------|------------|
| 1000 | 未知错误 | 500 |
| 1001 | 请求超时 | 408 |
| 1002 | 触发RPM限流 | 429 |
| 1004 | 鉴权失败 | 401 |
| 1008 | 余额不足 | 402 |
| 1013 | 服务内部错误 | 500 |
| 1027 | 输出内容错误 | 400 |
| 1039 | Token限制 | 400 |
| 2013 | 参数错误 | 400 |

### 常见错误码

- `400`: 请求参数错误
- `401`: 未授权（需要登录）
- `402`: 余额不足
- `403`: 禁止访问
- `404`: 资源不存在
- `408`: 请求超时
- `429`: 请求频率限制
- `500`: 服务器内部错误

### 错误处理示例

#### 限流错误
```json
{
  "code": 1002,
  "message": "触发RPM限流",
  "details": "MiniMax API error: 1002 - 触发RPM限流"
}
```

#### 认证失败
```json
{
  "code": 1004,
  "message": "鉴权失败",
  "details": "MiniMax API error: 1004 - 鉴权失败"
}
```

#### 余额不足
```json
{
  "code": 1008,
  "message": "余额不足",
  "details": "MiniMax API error: 1008 - 余额不足"
}
```

## MiniMax AI 响应结构

### 完整响应格式

```json
{
  "id": "03d3f5bd571f85faa1d980d2f779630f",
  "choices": [
    {
      "finish_reason": "stop",
      "index": 0,
      "message": {
        "content": "你好！有什么我可以帮助你的吗？",
        "role": "assistant",
        "name": "MiniMax AI"
      }
    }
  ],
  "created": 1736753853,
  "model": "MiniMax-M1",
  "object": "chat.completion",
  "usage": {
    "total_tokens": 70,
    "total_characters": 0,
    "prompt_tokens": 62,
    "completion_tokens": 8
  },
  "input_sensitive": false,
  "output_sensitive": false,
  "input_sensitive_type": 0,
  "output_sensitive_type": 0,
  "output_sensitive_int": 0,
  "base_resp": {
    "status_code": 0,
    "status_msg": ""
  }
}
```

### 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string | 响应ID |
| `choices` | array | 选择列表 |
| `choices[].finish_reason` | string | 完成原因 (stop, length, content_filter) |
| `choices[].index` | int | 选择索引 |
| `choices[].message.content` | string | AI回复内容 |
| `choices[].message.role` | string | 消息角色 |
| `choices[].message.name` | string | 消息名称 |
| `created` | int64 | 创建时间戳 |
| `model` | string | 使用的模型 |
| `object` | string | 对象类型 |
| `usage.total_tokens` | int | 总token数 |
| `usage.total_characters` | int | 总字符数 |
| `usage.prompt_tokens` | int | 提示token数 |
| `usage.completion_tokens` | int | 完成token数 |
| `input_sensitive` | bool | 输入是否敏感 |
| `output_sensitive` | bool | 输出是否敏感 |
| `input_sensitive_type` | int | 输入敏感类型 |
| `output_sensitive_type` | int | 输出敏感类型 |
| `output_sensitive_int` | int | 输出敏感整数 |
| `base_resp.status_code` | int | 状态码 (0表示成功) |
| `base_resp.status_msg` | string | 状态消息 |

## 环境变量

项目使用以下环境变量：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=rabbit_ai

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT配置
JWT_SECRET=your-jwt-secret-key

# GitHub OAuth配置
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

# MiniMax AI配置
MINIMAX_API_KEY=your-minimax-api-key
MINIMAX_BASE_URL=https://api.minimaxi.com/v1

# 服务器配置
PORT=8080
```

## 快速开始

1. 克隆项目并安装依赖
2. 复制 `env.example` 为 `.env` 并配置环境变量
3. 启动数据库和Redis
4. 运行 `make run` 启动服务
5. 访问 `http://localhost:8080/api/v1` 测试API

## 开发

```bash
# 运行测试
make test

# 构建项目
make build

# 运行项目
make run

# 清理构建文件
make clean
```

## 使用示例

### 基本聊天

```bash
curl -X POST "http://localhost:8080/api/v1/ai/chat/simple" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}'
```

### 带参数的聊天

```bash
curl -X POST "http://localhost:8080/api/v1/ai/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "写一首诗",
    "temperature": 0.8,
    "max_tokens": 500,
    "stream": false
  }'
```

### 流式聊天

```bash
curl -X POST "http://localhost:8080/api/v1/ai/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "请写一个故事",
    "stream": true,
    "temperature": 0.7
  }'
```

### 错误处理

```bash
# 处理限流错误
if [ $? -eq 429 ]; then
    echo "请求被限流，请稍后重试"
fi

# 处理认证错误
if [ $? -eq 401 ]; then
    echo "认证失败，请检查API密钥"
fi
``` 