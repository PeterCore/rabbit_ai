.PHONY: build run test clean deps

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