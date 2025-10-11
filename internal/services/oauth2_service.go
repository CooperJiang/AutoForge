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

// OAuth2Service OAuth2服务
type OAuth2Service struct{}

// NewOAuth2Service 创建OAuth2服务实例
func NewOAuth2Service() *OAuth2Service {
	return &OAuth2Service{}
}

// LinuxDoUserInfo Linux.do用户信息结构
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

// LinuxDoTokenResponse Linux.do token响应
type LinuxDoTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// GetLinuxDoAuthURL 获取Linux.do授权URL
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

// ExchangeLinuxDoToken 使用授权码exchange访问令牌
func (s *OAuth2Service) ExchangeLinuxDoToken(code string) (*LinuxDoTokenResponse, error) {
	cfg := config.GetConfig()

	// 构造请求参数
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", cfg.OAuth2.LinuxDo.ClientID)
	data.Set("client_secret", cfg.OAuth2.LinuxDo.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", cfg.OAuth2.LinuxDo.RedirectURL)

	// 发送POST请求
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

	// 解析响应
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

// GetLinuxDoUserInfo 获取Linux.do用户信息
func (s *OAuth2Service) GetLinuxDoUserInfo(accessToken string) (*LinuxDoUserInfo, error) {
	cfg := config.GetConfig()

	// 发送GET请求
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

	// 解析响应
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

// FindOrCreateLinuxDoUser 查找或创建Linux.do用户
func (s *OAuth2Service) FindOrCreateLinuxDoUser(userInfo *LinuxDoUserInfo) (*models.User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New(errors.CodeDBConnectionFailed, "数据库连接失败")
	}

	// 构建头像URL（使用120x120尺寸）
	avatarURL := strings.Replace(userInfo.AvatarTemplate, "{size}", "120", 1)
	if !strings.HasPrefix(avatarURL, "http") {
		avatarURL = "https://linux.do" + avatarURL
	}

	externalID := fmt.Sprintf("%d", userInfo.ID)

	// 查找是否已存在此Linux.do用户
	var user models.User
	err := db.Where("provider = ? AND external_id = ?", "oauth2_linuxdo", externalID).First(&user).Error
	if err == nil {
		// 用户已存在，更新信息
		user.Username = userInfo.Username
		user.Avatar = avatarURL
		user.TrustLevel = userInfo.TrustLevel
		user.Status = 1 // 确保状态为正常
		if userInfo.Silenced {
			user.Status = 2 // 如果被禁言，设置为禁用
		}
		if err := db.Save(&user).Error; err != nil {
			return nil, errors.New(errors.CodeQueryFailed, "更新用户失败")
		}
		return &user, nil
	}

	// 用户不存在，创建新用户
	newUser := models.User{
		BaseModel: models.BaseModel{
			ID: common.NewUUID(),
		},
		Username:   userInfo.Username,
		Email:      fmt.Sprintf("%s@linuxdo.oauth2", userInfo.Username), // OAuth2用户使用虚拟邮箱
		Password:   "", // OAuth2用户无密码
		Avatar:     avatarURL,
		Provider:   "oauth2_linuxdo",
		ExternalID: externalID,
		TrustLevel: userInfo.TrustLevel,
		Status:     1,
		Role:       3, // 默认普通用户
	}

	if userInfo.Silenced {
		newUser.Status = 2 // 被禁言则设置为禁用
	}

	if err := db.Create(&newUser).Error; err != nil {
		return nil, errors.New(errors.CodeQueryFailed, "创建用户失败")
	}

	return &newUser, nil
}
