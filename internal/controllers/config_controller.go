package controllers

import (
	"auto-forge/pkg/config"

	"github.com/gin-gonic/gin"
)


type ConfigController struct{}


func NewConfigController() *ConfigController {
	return &ConfigController{}
}


func (c *ConfigController) GetPublicConfig(ctx *gin.Context) {
	cfg := config.GetConfig()

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "获取配置成功",
		"data": gin.H{
			"oauth2": gin.H{
				"linuxdo": gin.H{
					"enabled": cfg.OAuth2.LinuxDo.Enabled,
				},
				"github": gin.H{
					"enabled": cfg.OAuth2.GitHub.Enabled,
				},
			},
			"app": gin.H{
				"name": cfg.App.Name,
			},
		},
	})
}
