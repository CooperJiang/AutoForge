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
		avatarURL = "https:
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
