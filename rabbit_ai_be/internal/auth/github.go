package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHubConfig GitHub OAuth配置
type GitHubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// GitHubUser GitHub用户信息
type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

// GitHubOAuth GitHub OAuth服务
type GitHubOAuth struct {
	config *oauth2.Config
}

// NewGitHubOAuth 创建GitHub OAuth实例
func NewGitHubOAuth(clientID, clientSecret, redirectURL string) *GitHubOAuth {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}

	return &GitHubOAuth{
		config: config,
	}
}

// GetAuthURL 获取GitHub授权URL
func (g *GitHubOAuth) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode 使用授权码交换访问令牌
func (g *GitHubOAuth) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}

// GetUserInfo 获取GitHub用户信息
func (g *GitHubOAuth) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GitHubUser, error) {
	client := g.config.Client(ctx, token)

	// 获取用户基本信息
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// 如果用户信息中没有邮箱，尝试获取邮箱列表
	if user.Email == "" {
		emails, err := g.getUserEmails(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to get user emails: %w", err)
		}

		// 使用主邮箱或第一个邮箱
		for _, email := range emails {
			if email.Primary {
				user.Email = email.Email
				break
			}
		}

		// 如果没有主邮箱，使用第一个邮箱
		if user.Email == "" && len(emails) > 0 {
			user.Email = emails[0].Email
		}
	}

	return &user, nil
}

// GitHubEmail GitHub邮箱信息
type GitHubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

// getUserEmails 获取用户邮箱列表
func (g *GitHubOAuth) getUserEmails(ctx context.Context, client *http.Client) ([]GitHubEmail, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return nil, fmt.Errorf("failed to get user emails: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var emails []GitHubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, fmt.Errorf("failed to decode emails: %w", err)
	}

	return emails, nil
}
