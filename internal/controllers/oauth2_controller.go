package controllers

import (
	"auto-forge/internal/services"
	"auto-forge/pkg/common"
	"auto-forge/pkg/config"
	"auto-forge/pkg/errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)


type OAuth2Controller struct {
	oauth2Service *services.OAuth2Service
}


func NewOAuth2Controller() *OAuth2Controller {
	return &OAuth2Controller{
		oauth2Service: services.NewOAuth2Service(),
	}
}


func (c *OAuth2Controller) LinuxDoLogin(ctx *gin.Context) {

	state := common.GenerateRandomString(32)


	ctx.SetCookie("oauth2_state", state, 600, "/", "", false, true)


	authURL := c.oauth2Service.GetLinuxDoAuthURL(state)


	ctx.Redirect(302, authURL)
}


func (c *OAuth2Controller) LinuxDoCallback(ctx *gin.Context) {

	var req struct {
		Code  string `json:"code" binding:"required"`
		State string `json:"state"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}


	cfg := config.GetConfig()
	if cfg.App.Mode != "debug" && req.State != "" {
		savedState, err := ctx.Cookie("oauth2_state")
		if err != nil || savedState != req.State {
			ctx.JSON(400, gin.H{
				"code":    400,
				"message": "state验证失败",
			})
			return
		}
	}


	ctx.SetCookie("oauth2_state", "", -1, "/", "", false, true)


	tokenResp, err := c.oauth2Service.ExchangeLinuxDoToken(req.Code)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "获取访问令牌失败",
		})
		return
	}


	userInfo, err := c.oauth2Service.GetLinuxDoUserInfo(tokenResp.AccessToken)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}


	user, err := c.oauth2Service.FindOrCreateLinuxDoUser(userInfo)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "创建用户失败",
		})
		return
	}


	if user.Status != 1 {
		ctx.JSON(403, gin.H{
			"code":    403,
			"message": "您的账号已被禁用",
		})
		return
	}


	jwtToken, expiresIn, err := generateJWTToken(user.ID.String())
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "生成令牌失败",
		})
		return
	}


	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token":      jwtToken,
			"expires_in": expiresIn,
			"user": gin.H{
				"id":       user.ID.String(),
				"username": user.Username,
				"email":    user.Email,
				"avatar":   user.Avatar,
				"status":   user.Status,
				"role":     user.Role,
			},
		},
	})
}


func generateJWTToken(userID string) (string, string, error) {
	cfg := config.GetConfig()

	now := time.Now()
	expiresAt := now.Add(time.Duration(cfg.JWT.ExpiresIn) * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
		"iat":     now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		return "", "", errors.New(errors.CodeInternal, "生成令牌失败")
	}

	return tokenString, expiresAt.Format(time.RFC3339), nil
}

// GitHubLogin 发起 GitHub OAuth 登录
func (c *OAuth2Controller) GitHubLogin(ctx *gin.Context) {
	// 生成随机 state，用于防止 CSRF 攻击
	state := common.GenerateRandomString(32)

	// 将 state 保存到 cookie 中
	ctx.SetCookie("oauth2_state", state, 600, "/", "", false, true)

	// 生成 GitHub 授权链接
	authURL := c.oauth2Service.GetGitHubAuthURL(state)

	// 重定向到 GitHub 授权页面
	ctx.Redirect(302, authURL)
}

// GitHubCallback 处理 GitHub OAuth 回调
func (c *OAuth2Controller) GitHubCallback(ctx *gin.Context) {
	// 获取回调参数
	var req struct {
		Code  string `json:"code" binding:"required"`
		State string `json:"state"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 验证 state，防止 CSRF 攻击
	cfg := config.GetConfig()
	if cfg.App.Mode != "debug" && req.State != "" {
		savedState, err := ctx.Cookie("oauth2_state")
		if err != nil || savedState != req.State {
			ctx.JSON(400, gin.H{
				"code":    400,
				"message": "state验证失败",
			})
			return
		}
	}

	// 清除 state cookie
	ctx.SetCookie("oauth2_state", "", -1, "/", "", false, true)

	// 用授权码换取访问令牌
	tokenResp, err := c.oauth2Service.ExchangeGitHubToken(req.Code)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "获取访问令牌失败",
		})
		return
	}

	// 获取用户信息
	userInfo, err := c.oauth2Service.GetGitHubUserInfo(tokenResp.AccessToken)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}

	// 查找或创建用户
	user, err := c.oauth2Service.FindOrCreateGitHubUser(userInfo)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "创建用户失败",
		})
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		ctx.JSON(403, gin.H{
			"code":    403,
			"message": "您的账号已被禁用",
		})
		return
	}

	// 生成 JWT token
	jwtToken, expiresIn, err := generateJWTToken(user.ID.String())
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "生成令牌失败",
		})
		return
	}

	// 返回登录成功信息
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token":      jwtToken,
			"expires_in": expiresIn,
			"user": gin.H{
				"id":       user.ID.String(),
				"username": user.Username,
				"email":    user.Email,
				"avatar":   user.Avatar,
				"status":   user.Status,
				"role":     user.Role,
			},
		},
	})
}
