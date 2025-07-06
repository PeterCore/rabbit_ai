# Rabbit AI Backend

一个基于 Go 和 Gin 框架的现代化后端服务，提供用户认证、设备管理、缓存管理和 AI 聊天等功能。

## 功能特性

- 🔐 **用户认证**: JWT 认证、GitHub OAuth 登录
- 📱 **设备管理**: 设备标识、平台检测、设备绑定
- 🚀 **Redis 缓存**: 高性能缓存、用户信息缓存
- 🤖 **AI 聊天**: MiniMax AI 集成、流式响应、参数控制
- 🛡️ **中间件**: JWT 验证、设备识别、CORS 支持
- 📊 **监控**: 缓存统计、使用情况监控

## MiniMax AI 功能

### 支持的参数

- **temperature**: 温度参数，控制随机性 (0.0-2.0)
- **max_tokens**: 最大生成token数
- **top_p**: 核采样参数 (0.0-1.0)
- **stream**: 流式响应支持
- **tool_choices**: 工具选择
- **stop**: 停止词列表
- **user**: 用户标识
- **repetition_penalty**: 重复惩罚参数
- **presence_penalty**: 存在惩罚参数
- **frequency_penalty**: 频率惩罚参数

### 响应类型

- **普通响应**: 完整的AI回复和使用统计
- **流式响应**: Server-Sent Events (SSE) 实时流式输出

## 快速开始

### 1. 环境要求

- Go 1.21+
- PostgreSQL 12+
- Redis 6+
- Docker (可选)

### 2. 安装依赖

```bash
go mod download
```

### 3. 环境配置

复制环境变量文件并配置：

```bash
cp env.example .env
```

编辑 `.env` 文件：

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

### 4. 启动服务

#### 使用 Docker Compose (推荐)

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

#### 手动启动

```bash
# 启动数据库和Redis
make db-start
make redis-start

# 运行项目
make run
```

### 5. 验证服务

```bash
# 健康检查
curl http://localhost:8080/api/v1/health

# 测试环境变量
go run scripts/test_env.go
```

## API 使用示例

### MiniMax AI 聊天

#### 简单聊天

```bash
curl -X POST "http://localhost:8080/api/v1/ai/chat/simple" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好，请介绍一下自己",
    "temperature": 0.7,
    "max_tokens": 2048
  }'
```

#### 完整聊天（支持流式）

```bash
curl -X POST "http://localhost:8080/api/v1/ai/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "写一首关于春天的诗",
    "temperature": 0.8,
    "max_tokens": 500,
    "top_p": 0.9,
    "stream": false,
    "stop": ["END", "STOP"],
    "user": "poet"
  }'
```

#### 流式聊天

```bash
curl -X POST "http://localhost:8080/api/v1/ai/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "请写一个关于人工智能的短文",
    "stream": true,
    "temperature": 0.7,
    "max_tokens": 300
  }'
```

### 用户认证

```bash
# 用户注册
curl -X POST "http://localhost:8080/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com"
  }'

# 用户登录
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 设备管理

```bash
# 绑定设备
curl -X POST "http://localhost:8080/api/v1/device/bind" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "unique-device-identifier",
    "platform": "iOS"
  }'
```

## 开发

### 项目结构

```
rabbit_ai_be/
├── cmd/server/          # 主程序入口
├── internal/            # 内部包
│   ├── auth/           # 认证相关
│   ├── cache/          # 缓存管理
│   ├── device/         # 设备管理
│   ├── minimax/        # MiniMax AI 集成
│   ├── model/          # 数据模型
│   ├── repository/     # 数据访问层
│   └── user/           # 用户管理
├── config/             # 配置文件
├── docs/               # 文档
├── scripts/            # 脚本文件
└── examples/           # 使用示例
```

### 常用命令

```bash
# 运行测试
make test

# 构建项目
make build

# 运行项目
make run

# 清理构建文件
make clean

# 格式化代码
make fmt

# 代码检查
make lint

# 数据库迁移
make db-migrate

# 启动开发环境
make dev
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/minimax/...

# 运行测试并显示覆盖率
go test -cover ./...

# 运行基准测试
go test -bench=. ./...
```

## 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `DB_HOST` | 数据库主机 | localhost |
| `DB_PORT` | 数据库端口 | 5432 |
| `DB_USER` | 数据库用户 | postgres |
| `DB_PASSWORD` | 数据库密码 | - |
| `DB_NAME` | 数据库名称 | rabbit_ai |
| `REDIS_HOST` | Redis主机 | localhost |
| `REDIS_PORT` | Redis端口 | 6379 |
| `REDIS_PASSWORD` | Redis密码 | - |
| `REDIS_DB` | Redis数据库 | 0 |
| `JWT_SECRET` | JWT密钥 | - |
| `GITHUB_CLIENT_ID` | GitHub OAuth客户端ID | - |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth客户端密钥 | - |
| `MINIMAX_API_KEY` | MiniMax API密钥 | - |
| `MINIMAX_BASE_URL` | MiniMax API基础URL | https://api.minimaxi.com/v1 |
| `PORT` | 服务器端口 | 8080 |

### MiniMax AI 参数

| 参数 | 类型 | 范围 | 默认值 | 说明 |
|------|------|------|--------|------|
| `temperature` | float64 | 0.0-2.0 | 0.7 | 控制输出的随机性 |
| `max_tokens` | int | > 0 | 2048 | 最大生成token数 |
| `top_p` | float64 | 0.0-1.0 | 0.9 | 核采样参数 |
| `stream` | bool | - | false | 是否启用流式响应 |
| `tool_choices` | array | - | - | 工具选择列表 |
| `stop` | array | - | - | 停止词列表 |
| `user` | string | - | - | 用户标识 |
| `repetition_penalty` | float64 | > 0 | - | 重复惩罚参数 |
| `presence_penalty` | float64 | - | - | 存在惩罚参数 |
| `frequency_penalty` | float64 | - | - | 频率惩罚参数 |

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t rabbit-ai-backend .

# 运行容器
docker run -d \
  --name rabbit-ai-backend \
  -p 8080:8080 \
  --env-file .env \
  rabbit-ai-backend
```

### 生产环境

1. 配置生产环境变量
2. 使用反向代理 (Nginx)
3. 配置 SSL 证书
4. 设置监控和日志
5. 配置数据库备份

## 贡献

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 支持

如果您遇到问题或有建议，请：

1. 查看 [API 文档](docs/API.md)
2. 检查 [GitHub Issues](https://github.com/your-repo/rabbit_ai_be/issues)
3. 创建新的 Issue 或 Pull Request

## 更新日志

### v1.0.0
- 初始版本发布
- 用户认证和授权
- 设备管理功能
- Redis 缓存集成
- MiniMax AI 聊天功能
- GitHub OAuth 登录
- 完整的 API 文档