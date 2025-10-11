package controllers

import (
	"auto-forge/pkg/config"

	"github.com/gin-gonic/gin"
)

// ConfigController 配置控制器
type ConfigController struct{}

// NewConfigController 创建配置控制器实例
func NewConfigController() *ConfigController {
	return &ConfigController{}
}

// GetPublicConfig 获取公开配置（前端使用）
func (c *ConfigController) GetPublicConfig(ctx *gin.Context) {
	cfg := config.GetConfig()

	// 返回前端需要的配置
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "获取配置成功",
		"data": gin.H{
			"oauth2": gin.H{
				"linuxdo": gin.H{
					"enabled": cfg.OAuth2.LinuxDo.Enabled,
				},
			},
			"app": gin.H{
				"name": cfg.App.Name,
			},
		},
	})
}
