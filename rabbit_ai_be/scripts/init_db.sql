-- 创建数据库（PostgreSQL 语法）
SELECT 'CREATE DATABASE rabbit_ai'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'rabbit_ai')\gexec

-- 使用数据库
\c rabbit_ai;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE,
    password VARCHAR(255), -- 加密后的密码，可以为空（阿里一键登录用户）
    nickname VARCHAR(100) NOT NULL,
    avatar TEXT,
    status INTEGER DEFAULT 1,
    github_id VARCHAR(100) UNIQUE, -- GitHub用户ID
    email VARCHAR(255) UNIQUE, -- 邮箱
    device_id VARCHAR(255) UNIQUE, -- 设备唯一标识
    platform VARCHAR(20), -- 终端平台: ios/android/browser
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_github_id ON users(github_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_device_id ON users(device_id);
CREATE INDEX IF NOT EXISTS idx_users_platform ON users(platform);

-- 创建对话表
CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL DEFAULT '新对话',
    status INTEGER DEFAULT 1, -- 1: 活跃, 0: 已删除
    message_count INTEGER DEFAULT 0,
    last_message_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建消息表
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL,
    role VARCHAR(20) NOT NULL, -- user/assistant
    content TEXT NOT NULL,
    tokens INTEGER DEFAULT 0,
    model VARCHAR(50) DEFAULT 'glm-4',
    finish_reason VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_conversations_status ON conversations(status);
CREATE INDEX IF NOT EXISTS idx_conversations_last_message_at ON conversations(last_message_at);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

-- 插入测试数据（可选）
INSERT INTO users (phone, nickname, avatar, status) 
VALUES ('13800138000', '测试用户', '', 1)
ON CONFLICT (phone) DO NOTHING; 