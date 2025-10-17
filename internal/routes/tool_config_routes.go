package routes

import (
	"auto-forge/internal/controllers"
	"auto-forge/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterToolConfigRoutes 注册工具配置路由
func RegisterToolConfigRoutes(router *gin.Engine) {
	// 公开路由（获取工具分类）
	router.GET("/api/v1/tool-categories", controllers.GetToolCategories)

	// 管理端路由（工具配置管理）
	admin := router.Group("/api/v1/admin/tool-configs")
	admin.Use(middleware.RequireAdmin())
	{
		admin.GET("", controllers.GetAllTools)                         // 获取所有工具配置
		admin.GET("/:code", controllers.GetToolDetail)                 // 获取工具配置详情
		admin.PUT("/:code", controllers.UpdateToolConfig)              // 更新工具配置
		admin.PATCH("/:code/settings", controllers.UpdateToolSettings) // 更新工具设置
		admin.DELETE("/:id", controllers.DeleteTool)                   // 删除工具配置
		admin.POST("/sync", controllers.SyncTools)                     // 同步工具定义
	}
}
