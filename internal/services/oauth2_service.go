package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"auto-forge/internal/models"
	"auto-forge/pkg/common"
	"auto-forge/pkg/config"
	"auto-forge/pkg/database"
	"auto-forge/pkg/errors"
)


type OAuth2Service struct{}


func NewOAuth2Service() *OAuth2Service {
	return &OAuth2Service{}
}


type LinuxDoUserInfo struct {
	ID             int                    `json:"id"`
	Username       string                 `json:"username"`
	Name           string                 `json:"name"`
	AvatarTemplate string                 `json:"avatar_template"`
	Active         bool                   `json:"active"`
	TrustLevel     int                    `json:"trust_level"`
	Silenced       bool                   `json:"silenced"`
	ExternalIDs    map[string]interface{} `json:"external_ids"`
	APIKey         string                 `json:"api_key"`
}


type LinuxDoTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}


func (s *OAuth2Service) GetLinuxDoAuthURL(state string) string {
	cfg := config.GetConfig()
	params := url.Values{}
	params.Add("client_id", cfg.OAuth2.LinuxDo.ClientID)
	params.Add("redirect_uri", cfg.OAuth2.LinuxDo.RedirectURL)
	params.Add("response_type", "code")
	params.Add("scope", strings.Join(cfg.OAuth2.LinuxDo.Scopes, " "))
	params.Add("state", state)

	return cfg.OAuth2.LinuxDo.AuthorizeURL + "?" + params.Encode()
}


func (s *OAuth2Service) ExchangeLinuxDoToken(code string) (*LinuxDoTokenResponse, error) {
	cfg := config.GetConfig()


	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", cfg.OAuth2.LinuxDo.ClientID)
	data.Set("client_secret", cfg.OAuth2.LinuxDo.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", cfg.OAuth2.LinuxDo.RedirectURL)


	req, err := http.NewRequest("POST", cfg.OAuth2.LinuxDo.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "创建请求失败")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "请求token失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(errors.CodeInternal, fmt.Sprintf("获取token失败: %s", string(body)))
	}


	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "读取响应失败")
	}

	var tokenResp LinuxDoTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, errors.New(errors.CodeInternal, "解析token响应失败")
	}

	return &tokenResp, nil
}


func (s *OAuth2Service) GetLinuxDoUserInfo(accessToken string) (*LinuxDoUserInfo, error) {
	cfg := config.GetConfig()


	req, err := http.NewRequest("GET", cfg.OAuth2.LinuxDo.UserInfoURL, nil)
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "创建请求失败")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "请求用户信息失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(errors.CodeInternal, fmt.Sprintf("获取用户信息失败: %s", string(body)))
	}


	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "读取响应失败")
	}

	var userInfo LinuxDoUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, errors.New(errors.CodeInternal, "解析用户信息失败")
	}

	return &userInfo, nil
}


func (s *OAuth2Service) FindOrCreateLinuxDoUser(userInfo *LinuxDoUserInfo) (*models.User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New(errors.CodeDBConnectionFailed, "数据库连接失败")
	}


	avatarURL := strings.Replace(userInfo.AvatarTemplate, "{size}", "120", 1)
	if !strings.HasPrefix(avatarURL, "http") {
		avatarURL = "https://linux.do" + avatarURL
	}

	externalID := fmt.Sprintf("%d", userInfo.ID)


	var user models.User
	err := db.Where("provider = ? AND external_id = ?", "oauth2_linuxdo", externalID).First(&user).Error
	if err == nil {

		user.Username = userInfo.Username
		user.Avatar = avatarURL
		user.TrustLevel = userInfo.TrustLevel
		user.Status = 1
		if userInfo.Silenced {
			user.Status = 2
		}
		if err := db.Save(&user).Error; err != nil {
			return nil, errors.New(errors.CodeQueryFailed, "更新用户失败")
		}
		return &user, nil
	}


	newUser := models.User{
		BaseModel: models.BaseModel{
			ID: common.NewUUID(),
		},
		Username:   userInfo.Username,
		Email:      fmt.Sprintf("%s@linuxdo.oauth2", userInfo.Username),
		Password:   "",
		Avatar:     avatarURL,
		Provider:   "oauth2_linuxdo",
		ExternalID: externalID,
		TrustLevel: userInfo.TrustLevel,
		Status:     1,
		Role:       3,
	}

	if userInfo.Silenced {
		newUser.Status = 2
	}

	if err := db.Create(&newUser).Error; err != nil {
		return nil, errors.New(errors.CodeQueryFailed, "创建用户失败")
	}

	return &newUser, nil
}

// GitHub OAuth2 相关结构体和方法

type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

type GitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GitHub OAuth2 固定配置
const (
	GitHubAuthorizeURL = "https://github.com/login/oauth/authorize"
	GitHubTokenURL     = "https://github.com/login/oauth/access_token"
	GitHubUserInfoURL  = "https://api.github.com/user"
	GitHubEmailsURL    = "https://api.github.com/user/emails"
	GitHubScopes       = "read:user user:email"
)

// GetGitHubAuthURL 生成 GitHub OAuth 授权链接
func (s *OAuth2Service) GetGitHubAuthURL(state string) string {
	cfg := config.GetConfig()
	params := url.Values{}
	params.Add("client_id", cfg.OAuth2.GitHub.ClientID)
	params.Add("redirect_uri", cfg.OAuth2.GitHub.RedirectURL)
	params.Add("response_type", "code")
	params.Add("scope", GitHubScopes)
	params.Add("state", state)

	return GitHubAuthorizeURL + "?" + params.Encode()
}

// ExchangeGitHubToken 用授权码换取 access token
func (s *OAuth2Service) ExchangeGitHubToken(code string) (*GitHubTokenResponse, error) {
	cfg := config.GetConfig()

	data := url.Values{}
	data.Set("client_id", cfg.OAuth2.GitHub.ClientID)
	data.Set("client_secret", cfg.OAuth2.GitHub.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", cfg.OAuth2.GitHub.RedirectURL)

	req, err := http.NewRequest("POST", GitHubTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "创建请求失败")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "请求token失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(errors.CodeInternal, fmt.Sprintf("获取token失败: %s", string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "读取响应失败")
	}

	var tokenResp GitHubTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, errors.New(errors.CodeInternal, "解析token响应失败")
	}

	return &tokenResp, nil
}

// GetGitHubUserInfo 获取 GitHub 用户信息
func (s *OAuth2Service) GetGitHubUserInfo(accessToken string) (*GitHubUserInfo, error) {
	req, err := http.NewRequest("GET", GitHubUserInfoURL, nil)
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "创建请求失败")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "请求用户信息失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(errors.CodeInternal, fmt.Sprintf("获取用户信息失败: %s", string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "读取响应失败")
	}

	var userInfo GitHubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, errors.New(errors.CodeInternal, "解析用户信息失败")
	}

	// 如果用户信息中没有 email，需要单独获取
	if userInfo.Email == "" {
		email, err := s.getGitHubPrimaryEmail(accessToken)
		if err == nil {
			userInfo.Email = email
		}
	}

	return &userInfo, nil
}

// getGitHubPrimaryEmail 获取 GitHub 用户的主邮箱
func (s *OAuth2Service) getGitHubPrimaryEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", GitHubEmailsURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取邮箱失败: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	// 查找主邮箱且已验证的
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	// 如果没有主邮箱，返回第一个已验证的
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("未找到已验证的邮箱")
}

// FindOrCreateGitHubUser 查找或创建 GitHub 用户
func (s *OAuth2Service) FindOrCreateGitHubUser(userInfo *GitHubUserInfo) (*models.User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New(errors.CodeDBConnectionFailed, "数据库连接失败")
	}

	externalID := fmt.Sprintf("%d", userInfo.ID)

	// 先查找是否已存在该 GitHub 用户
	var user models.User
	err := db.Where("provider = ? AND external_id = ?", "oauth2_github", externalID).First(&user).Error
	if err == nil {
		// 用户已存在，更新信息
		user.Username = userInfo.Login
		user.Avatar = userInfo.AvatarURL
		if userInfo.Email != "" {
			user.Email = userInfo.Email
		}
		if err := db.Save(&user).Error; err != nil {
			return nil, errors.New(errors.CodeQueryFailed, "更新用户失败")
		}
		return &user, nil
	}

	// 用户不存在，创建新用户
	email := userInfo.Email
	if email == "" {
		email = fmt.Sprintf("%s@github.oauth2", userInfo.Login)
	}

	newUser := models.User{
		BaseModel: models.BaseModel{
			ID: common.NewUUID(),
		},
		Username:   userInfo.Login,
		Email:      email,
		Password:   "",
		Avatar:     userInfo.AvatarURL,
		Provider:   "oauth2_github",
		ExternalID: externalID,
		TrustLevel: 0,
		Status:     1,
		Role:       3,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return nil, errors.New(errors.CodeQueryFailed, "创建用户失败")
	}

	return &newUser, nil
}
