.PHONY: help build run test clean docker-build docker-run setup-env test-env deps

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  setup-env    - 设置环境变量文件"
	@echo "  test-env     - 测试环境变量配置"
	@echo "  build        - 构建项目"
	@echo "  run          - 运行项目"
	@echo "  test         - 运行测试"
	@echo "  clean        - 清理构建文件"
	@echo "  docker-build - 构建Docker镜像"
	@echo "  docker-run   - 运行Docker容器"

# 设置环境变量文件
setup-env:
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "✅ 已创建 .env 文件，请编辑配置信息"; \
	else \
		echo "⚠️  .env 文件已存在"; \
	fi

# 测试环境变量配置
test-env:
	@go run scripts/test_env.go

# 构建项目
build:
	go build -o bin/server cmd/server/main.go

# 运行项目
run:
	go run cmd/server/main.go

# 运行测试
test:
	go test ./...

# 清理构建文件
clean:
	rm -rf bin/
	rm -rf tmp/

# 下载依赖
deps:
	go mod tidy
	go mod download

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 生成go.sum
sum:
	go mod tidy

# 数据库初始化（需要PostgreSQL运行）
init-db:
	psql -U postgres -f scripts/init_db.sql

# 开发模式运行（带热重载）
dev:
	air

# 安装air（热重载工具）
install-air:
	go install github.com/cosmtrek/air@latest

# 构建Docker镜像
docker-build:
	docker build -t rabbit-ai-backend .

# 运行Docker容器
docker-run:
	docker-compose up -d 