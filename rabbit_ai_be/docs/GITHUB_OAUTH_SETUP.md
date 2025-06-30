# GitHub OAuth 应用配置指南

## 1. 创建GitHub OAuth应用

### 步骤1：访问GitHub开发者设置
1. 登录GitHub账号
2. 点击右上角头像，选择 "Settings"
3. 在左侧菜单中点击 "Developer settings"
4. 点击 "OAuth Apps"
5. 点击 "New OAuth App"

### 步骤2：填写应用信息
```
Application name: Rabbit AI Login
Homepage URL: http://localhost:8080
Application description: Rabbit AI 登录注册系统
Authorization callback URL: http://localhost:8080/api/v1/auth/github/callback
```

### 步骤3：注册应用
点击 "Register application" 完成注册

## 2. 获取Client ID和Client Secret

注册完成后，你会看到：
- **Client ID**: 自动生成，复制保存
- **Client Secret**: 点击 "Generate a new client secret" 生成，复制保存

## 3. 配置环境变量

将获取到的信息配置到 `.env` 文件中：

```env
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://localhost:8080/api/v1/auth/github/callback
```

## 4. 测试GitHub登录

### 4.1 获取授权URL
```bash
curl -X GET "http://localhost:8080/api/v1/auth/github/auth-url?state=test123"
```

### 4.2 完成OAuth流程
1. 在浏览器中打开返回的 `auth_url`
2. 授权GitHub应用访问你的信息
3. 获取授权码（code参数）
4. 使用授权码调用登录接口

### 4.3 调用登录接口
```bash
curl -X POST "http://localhost:8080/api/v1/auth/github/login" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "your_authorization_code",
    "state": "test123"
  }'
```

## 5. 生产环境配置

在生产环境中，需要：

1. **更新回调URL**：
   - 将 `http://localhost:8080` 替换为你的生产域名
   - 例如：`https://yourdomain.com/api/v1/auth/github/callback`

2. **安全配置**：
   - 使用HTTPS
   - 设置合适的state参数防止CSRF攻击
   - 定期轮换Client Secret

3. **权限范围**：
   - 当前配置的权限范围：`user:email, read:user`
   - 可根据需要调整权限范围

## 6. 常见问题

### Q: 回调URL不匹配
A: 确保GitHub OAuth应用中的回调URL与代码中的配置完全一致

### Q: 授权失败
A: 检查Client ID和Client Secret是否正确配置

### Q: 获取不到邮箱信息
A: 确保GitHub账号的邮箱是公开的，或者应用有足够的权限

## 7. 安全注意事项

1. **保护Client Secret**：不要将Client Secret提交到代码仓库
2. **使用环境变量**：通过环境变量配置敏感信息
3. **验证State参数**：防止CSRF攻击
4. **HTTPS**：生产环境必须使用HTTPS
5. **权限最小化**：只申请必要的权限范围 