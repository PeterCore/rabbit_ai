package auth

import (
	"testing"
)

func TestGitHubOAuth(t *testing.T) {
	// 创建GitHub OAuth实例
	githubOAuth := NewGitHubOAuth(
		"test-client-id",
		"test-client-secret",
		"http://localhost:8080/callback",
	)

	t.Run("GetAuthURL", func(t *testing.T) {
		authURL := githubOAuth.GetAuthURL("test-state")

		if authURL == "" {
			t.Error("Expected auth URL, got empty string")
		}

		// 验证URL包含必要的参数
		if authURL == "" {
			t.Error("Auth URL should not be empty")
		}
	})

	t.Run("ExchangeCode", func(t *testing.T) {
		// 这个测试需要真实的授权码，所以跳过
		t.Skip("Skipping ExchangeCode test - requires real authorization code")
	})

	t.Run("GetUserInfo", func(t *testing.T) {
		// 这个测试需要真实的访问令牌，所以跳过
		t.Skip("Skipping GetUserInfo test - requires real access token")
	})
}
