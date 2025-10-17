package routes

import (
	"auto-forge/internal/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterOAuth2Routes 注册OAuth2路由
func RegisterOAuth2Routes(r *gin.RouterGroup) {
	oauth2Controller := controllers.NewOAuth2Controller()

	// Linux.do OAuth2登录路由（跳转到授权页）
	r.GET("/linuxdo", oauth2Controller.LinuxDoLogin)
	// Linux.do OAuth2回调处理路由（前端通过POST传递code）
	r.POST("/linuxdo/callback", oauth2Controller.LinuxDoCallback)

	// GitHub OAuth2登录路由（跳转到授权页）
	r.GET("/github", oauth2Controller.GitHubLogin)
	// GitHub OAuth2回调处理路由（前端通过POST传递code）
	r.POST("/github/callback", oauth2Controller.GitHubCallback)
}

// RegisterOAuth2CallbackRoutes 注册OAuth2回调路由（已废弃）
func RegisterOAuth2CallbackRoutes(r *gin.Engine) {
	// 不再需要顶层回调路由
}
