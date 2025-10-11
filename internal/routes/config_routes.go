package routes

import (
	"auto-forge/internal/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterConfigRoutes 注册配置路由
func RegisterConfigRoutes(r *gin.RouterGroup) {
	configController := controllers.NewConfigController()

	// 获取公开配置（无需认证）
	r.GET("/config", configController.GetPublicConfig)
}
