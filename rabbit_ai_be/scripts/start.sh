#!/bin/bash

# Rabbit AI Backend 快速启动脚本

echo "🚀 Rabbit AI Backend 启动脚本"
echo "================================"

# 检查.env文件是否存在
if [ ! -f .env ]; then
    echo "⚠️  .env 文件不存在，正在创建..."
    cp env.example .env
    echo "✅ 已创建 .env 文件，请编辑配置信息后重新运行"
    echo "📝 主要需要配置的项："
    echo "   - MINIMAX_API_KEY: MiniMax AI API密钥"
    echo "   - JWT_SECRET: JWT密钥"
    echo "   - DB_PASSWORD: 数据库密码"
    exit 1
fi

# 测试环境配置
echo "🔍 测试环境配置..."
go run scripts/test_env.go

echo ""
echo "📦 构建项目..."
make build

echo ""
echo "🌐 启动服务器..."
echo "服务器将在 http://localhost:8080 启动"
echo "API文档: http://localhost:8080/docs"
echo ""
echo "按 Ctrl+C 停止服务器"
echo "================================"

# 启动服务器
./bin/server 