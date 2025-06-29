-- 创建数据库
CREATE DATABASE IF NOT EXISTS rabbit_ai;

-- 使用数据库
\c rabbit_ai;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255), -- 加密后的密码，可以为空（阿里一键登录用户）
    nickname VARCHAR(100) NOT NULL,
    avatar TEXT,
    status INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- 插入测试数据（可选）
INSERT INTO users (phone, nickname, avatar, status) 
VALUES ('13800138000', '测试用户', '', 1)
ON CONFLICT (phone) DO NOTHING; 