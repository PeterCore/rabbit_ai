# Rabbit AI 登录注册系统 API 文档

## 基础信息

- 基础URL: `http://localhost:8080/api/v1`
- 认证方式: JWT Bearer Token
- 响应格式: JSON

## 通用响应格式

```json
{
  "code": 200,
  "message": "Success",
  "data": {}
}
```

## 认证相关接口

### 1. 阿里一键登录

**接口地址:** `POST /auth/login`

**请求参数:**
```json
{
  "auth_code": "阿里一键登录返回的auth_code"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "phone": "13800138000",
      "nickname": "用户8000",
      "avatar": "",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 2. 密码登录

**接口地址:** `POST /auth/login/password`

**请求参数:**
```json
{
  "phone": "13800138000",
  "password": "your_password"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "phone": "13800138000",
      "nickname": "用户8000",
      "avatar": "",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 3. 用户注册

**接口地址:** `POST /auth/register`

**请求参数:**
```json
{
  "phone": "13800138000",
  "password": "your_password",
  "nickname": "用户昵称"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Registration successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "phone": "13800138000",
      "nickname": "用户昵称",
      "avatar": "",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

## 用户相关接口（需要认证）

### 1. 获取用户信息

**接口地址:** `GET /users/profile`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": 1,
    "phone": "13800138000",
    "nickname": "用户8000",
    "avatar": "",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. 更新用户信息

**接口地址:** `PUT /users/profile`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "nickname": "新昵称",
  "avatar": "新头像URL"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "phone": "13800138000",
    "nickname": "新昵称",
    "avatar": "新头像URL",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

### 3. 修改密码

**接口地址:** `PUT /users/password`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "old_password": "旧密码",
  "new_password": "新密码"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Password updated successfully"
}
```

### 4. 删除用户

**接口地址:** `DELETE /users/profile`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "User deleted successfully"
}
```

### 5. 根据ID获取用户信息

**接口地址:** `GET /users/{id}`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": 1,
    "phone": "13800138000",
    "nickname": "用户8000",
    "avatar": "",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

## 健康检查

### 1. 服务健康检查

**接口地址:** `GET /health`

**响应示例:**
```json
{
  "status": "ok",
  "time": "2024-01-01T00:00:00Z"
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或认证失败 |
| 500 | 服务器内部错误 |

## 使用示例

### 1. 用户注册流程

```bash
# 1. 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "password": "your_password",
    "nickname": "用户昵称"
  }'

# 2. 使用返回的token访问受保护的接口
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer your_token_here"
```

### 2. 密码登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login/password \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "password": "your_password"
  }'
```

### 3. 修改密码

```bash
curl -X PUT http://localhost:8080/api/v1/users/password \
  -H "Authorization: Bearer your_token_here" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "旧密码",
    "new_password": "新密码"
  }'
```

### 4. 更新用户信息

```bash
curl -X PUT http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer your_token_here" \
  -H "Content-Type: application/json" \
  -d '{"nickname": "新昵称", "avatar": "https://example.com/avatar.jpg"}'
``` 