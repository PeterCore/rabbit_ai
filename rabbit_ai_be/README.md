# Rabbit AI 登录注册系统

一个基于 Golang + Gin 框架的 AI 应用登录注册系统，支持阿里一键登录、JWT 认证、PostgreSQL 数据库和 Redis 缓存。

## 功能特性

- 🔐 **阿里一键登录**: 集成阿里云一键登录服务，用户可通过手机号快速登录
- 🐙 **GitHub登录**: 集成GitHub OAuth，用户可通过GitHub账号快速登录
- 🛡️ **JWT 认证**: 使用 JWT 进行用户身份验证和授权
- 👤 **用户管理**: 完整的用户 CRUD 操作
- 🗄️ **PostgreSQL**: 使用 PostgreSQL 作为主数据库
- ⚡ **Redis 缓存**: 使用 Redis 缓存用户信息，提升查询性能
- 🏗️ **分层架构**: 清晰的分层架构设计，易于维护和扩展
- 📚 **完整文档**: 提供详细的 API 文档和使用说明

## 技术栈

- **后端框架**: Gin
- **认证**: JWT (github.com/dgrijalva/jwt-go)
- **数据库**: PostgreSQL
- **缓存**: Redis (github.com/redis/go-redis/v9)
- **配置管理**: 环境变量 + godotenv
- **API 文档**: Markdown 格式

## 项目结构

```
rabbit_ai/
├── cmd/
│   └── server/
│       └── main.go              # 主程序入口
├── internal/
│   ├── auth/
│   │   ├── handler.go           # 认证处理器
│   │   └── service.go           # 认证服务
│   ├── user/
│   │   ├── handler.go           # 用户处理器
│   │   └── service.go           # 用户服务
│   ├── middleware/
│   │   └── jwt.go               # JWT 中间件
│   ├── model/
│   │   └── user.go              # 用户模型
│   ├── cache/
│   │   ├── redis.go             # Redis 缓存服务
│   │   └── redis_test.go        # Redis 缓存测试
│   └── repository/
│       └── user_cache.go        # 带缓存的用户仓库
├── config/
│   └── config.yaml              # 配置文件
├── scripts/
│   └── init_db.sql              # 数据库初始化脚本
├── docs/
│   └── API.md                   # API 文档
├── go.mod                       # Go 模块文件
├── Makefile                     # 构建脚本
├── env.example                  # 环境变量示例
└── README.md                    # 项目说明
```

## 快速开始

### 1. 环境要求

- Go 1.21+
- PostgreSQL 12+
- Redis 6+
- 阿里云账号（用于一键登录服务）
- GitHub账号（用于创建OAuth应用）

### 2. 克隆项目

```bash
git clone <repository-url>
cd rabbit_ai
```

### 3. 安装依赖

```bash
make deps
# 或者
go mod tidy
```

### 4. 配置环境变量

复制环境变量示例文件并修改配置：

```bash
cp env.example .env
```

编辑 `.env` 文件，配置以下参数：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=rabbit_ai
DB_SSLMODE=disable

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT 配置
JWT_SECRET=your-secret-key-here
JWT_EXPIRE_HOURS=24

# 阿里云配置
ALIYUN_ACCESS_KEY_ID=your-access-key-id
ALIYUN_ACCESS_KEY_SECRET=your-access-key-secret
ALIYUN_REGION=cn-hangzhou
ALIYUN_ONE_CLICK_APP_ID=your-one-click-app-id

# GitHub配置
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://localhost:8080/api/v1/auth/github/callback
```

### 5. 启动服务

#### 使用 Docker Compose（推荐）

```bash
# 启动所有服务（PostgreSQL + Redis + 应用）
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f app
```

#### 手动启动

```bash
# 启动 PostgreSQL
# 根据你的系统安装和启动 PostgreSQL

# 启动 Redis
redis-server

# 初始化数据库
make init-db

# 运行应用
make run
```

### 6. 验证服务

访问健康检查接口：

```bash
curl http://localhost:8080/health
```

## 缓存功能

### Redis 缓存特性

- **用户信息缓存**: 用户信息缓存30分钟，提升查询性能
- **缓存策略**: 采用 Cache-Aside 模式，先查缓存，缓存未命中则查数据库
- **数据同步**: 确保缓存与数据库数据一致性
- **自动失效**: 用户信息更新时自动使缓存失效

### 缓存操作

- **读取**: 优先从 Redis 缓存获取，缓存未命中则从数据库获取并缓存
- **写入**: 先写入数据库，再更新缓存
- **更新**: 先更新数据库，再更新缓存
- **删除**: 先删除数据库记录，再删除缓存
- **密码更新**: 密码更新时使缓存失效（安全考虑）

## API 使用

### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"auth_code": "your_auth_code_here"}'
```

### 获取用户信息

```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer your_token_here"
```

### 更新用户信息

```bash
curl -X PUT http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer your_token_here" \
  -H "Content-Type: application/json" \
  -d '{"nickname": "新昵称", "avatar": "https://example.com/avatar.jpg"}'
```

详细的 API 文档请参考 [docs/API.md](docs/API.md)。

## 开发指南

### 代码格式化

```bash
make fmt
```

### 代码检查

```bash
make lint
```

### 运行测试

```bash
# 运行所有测试
make test

# 运行缓存测试（需要 Redis 运行）
go test ./internal/cache/
```

### 热重载开发

```bash
# 安装 air
make install-air

# 启动热重载
make dev
```

## 部署

### 构建生产版本

```bash
make build
```

### Docker 部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

## 配置说明

### 数据库配置

- `DB_HOST`: 数据库主机地址
- `DB_PORT`: 数据库端口
- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称
- `DB_SSLMODE`: SSL 模式

### Redis 配置

- `REDIS_HOST`: Redis 主机地址
- `REDIS_PORT`: Redis 端口
- `REDIS_PASSWORD`: Redis 密码（可选）
- `REDIS_DB`: Redis 数据库编号

### JWT 配置

- `JWT_SECRET`: JWT 签名密钥
- `JWT_EXPIRE_HOURS`: Token 过期时间（小时）

### 阿里云配置

- `ALIYUN_ACCESS_KEY_ID`: 阿里云 Access Key ID
- `ALIYUN_ACCESS_KEY_SECRET`: 阿里云 Access Key Secret
- `ALIYUN_REGION`: 阿里云地域
- `ALIYUN_ONE_CLICK_APP_ID`: 一键登录应用 ID

### GitHub配置

- `GITHUB_CLIENT_ID`: GitHub Client ID
- `GITHUB_CLIENT_SECRET`: GitHub Client Secret
- `GITHUB_REDIRECT_URL`: GitHub Redirect URL

## 注意事项

1. **阿里云配置**: 需要先在阿里云控制台开通一键登录服务并获取相关配置
2. **数据库安全**: 生产环境中请使用强密码和 SSL 连接
3. **Redis 安全**: 生产环境中请设置 Redis 密码和访问控制
4. **JWT 密钥**: 生产环境中请使用足够复杂的密钥
5. **环境变量**: 敏感信息请通过环境变量配置，不要硬编码
6. **缓存一致性**: 确保缓存与数据库的数据一致性，避免数据不一致问题

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License