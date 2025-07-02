# Rabbit AI Backend API 文档

## 基础信息

- 基础URL: `http://localhost:8080/api/v1`
- 认证方式: JWT Bearer Token
- 设备标识: 通过HTTP Header传递设备ID

### 设备标识Header

支持以下Header名称（按优先级排序）：
- `X-Device-ID`
- `X-Device-Id`
- `Device-ID`
- `Device-Id`
- `X-Client-ID`
- `X-Client-Id`
- `Client-ID`
- `Client-Id`
- `X-User-Agent`
- `User-Agent`

如果未提供设备ID Header，系统会尝试从User-Agent生成一个简单的设备标识。

### 平台标识Header

支持以下Header名称（按优先级排序）：
- `X-Platform`
- `Platform`

如果未提供平台Header，系统会尝试从User-Agent判断平台类型：
- iOS设备：返回 "ios"
- Android设备：返回 "android"
- 桌面浏览器：返回 "browser"

支持的平台值：
- `ios`：iOS设备
- `android`：Android设备
- `browser`：浏览器

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

### 4. GitHub登录

#### 4.1 获取GitHub授权URL

**接口地址:** `GET /auth/github/auth-url`

**请求参数:**
```
state: 可选的状态参数，用于防止CSRF攻击
```

**响应示例:**
```json
{
  "code": 200,
  "message": "GitHub auth URL generated",
  "data": {
    "auth_url": "https://github.com/login/oauth/authorize?client_id=...&redirect_uri=...&scope=user:email,read:user&state=random_state"
  }
}
```

#### 4.2 GitHub登录回调

**接口地址:** `POST /auth/github/login`

**请求参数:**
```json
{
  "code": "GitHub授权码",
  "state": "状态参数"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "GitHub login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "phone": "",
      "nickname": "GitHub用户名",
      "avatar": "https://avatars.githubusercontent.com/u/123456?v=4",
      "status": 1,
      "github_id": "123456",
      "email": "user@example.com",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

## 设备管理API

### 1. 根据设备ID获取或创建用户

**接口地址:** `GET /device/user`

**请求头:**
```
X-Device-ID: your-device-id
X-Platform: ios
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "device_id": "your-device-id",
    "platform": "ios",
    "user": {
      "id": 1,
      "phone": "",
      "nickname": "设备用户_12345678",
      "avatar": "",
      "status": 1,
      "github_id": "",
      "email": "",
      "device_id": "your-device-id",
      "platform": "ios",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "is_new": true
  }
}
```

### 2. 根据设备ID获取用户信息

**接口地址:** `GET /device/user/info`

**请求头:**
```
X-Device-ID: your-device-id
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "device_id": "your-device-id",
    "user": {
      "id": 1,
      "phone": "13800138000",
      "nickname": "用户昵称",
      "avatar": "",
      "status": 1,
      "github_id": "",
      "email": "",
      "device_id": "your-device-id",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 3. 绑定设备到用户（需要JWT认证）

**接口地址:** `POST /device/bind`

**请求头:**
```
Authorization: Bearer your-jwt-token
X-Device-ID: your-device-id
X-Platform: ios
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Device bound successfully",
  "data": {
    "device_id": "your-device-id",
    "user_id": 1
  }
}
```

### 4. 解绑设备（需要JWT认证）

**接口地址:** `DELETE /device/unbind`

**请求头:**
```
Authorization: Bearer your-jwt-token
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Device unbound successfully",
  "data": {
    "user_id": 1
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

### 1. 设备登录流程

```bash
# 1. 根据设备ID获取或创建用户
curl -X GET "http://localhost:8080/api/v1/device/user" \
  -H "X-Device-ID: your-device-id" \
  -H "X-Platform: ios"

# 2. 如果需要绑定到已有账户，先进行认证
curl -X POST "http://localhost:8080/api/v1/auth/password-login" \
  -H "Content-Type: application/json" \
  -H "X-Platform: ios" \
  -d '{"phone": "13800138000", "password": "password"}'

# 3. 绑定设备到用户
curl -X POST "http://localhost:8080/api/v1/device/bind" \
  -H "Authorization: Bearer your-jwt-token" \
  -H "X-Device-ID: your-device-id" \
  -H "X-Platform: ios"
```

### 2. 获取用户信息

```bash
curl -X GET "http://localhost:8080/api/v1/users/1" \
  -H "Authorization: Bearer your-jwt-token"
```

### 3. 更新用户信息

```bash
curl -X PUT "http://localhost:8080/api/v1/users/1" \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{"nickname": "新昵称", "avatar": "新头像URL"}'
``` 