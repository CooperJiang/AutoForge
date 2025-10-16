package routes

import (
	templateController "auto-forge/internal/controllers/template"
	"auto-forge/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterTemplateRoutes 注册模板市场路由
func RegisterTemplateRoutes(r *gin.RouterGroup) {
	templates := r.Group("/templates")
	{
		// 公开接口 - 浏览和查看模板
		templates.GET("", templateController.GetTemplateList)         // 获取模板列表
		templates.GET("/:id", templateController.GetTemplateDetail)   // 获取模板详情

		// 需要认证的接口
		authenticated := templates.Group("")
		authenticated.Use(middleware.RequireAuth())
		{
			// 用户操作
			authenticated.POST("/install", templateController.InstallTemplate)         // 安装模板
			authenticated.GET("/installs", templateController.GetInstallHistory)       // 获取安装历史

			// 管理员操作
			authenticated.POST("", templateController.CreateTemplate)                  // 创建模板
			authenticated.PUT("/:id", templateController.UpdateTemplate)               // 更新模板
			authenticated.DELETE("/:id", templateController.DeleteTemplate)            // 删除模板
		}
	}

	// 模板分类路由
	categories := r.Group("/template-categories")
	{
		// 公开接口 - 获取分类列表
		categories.GET("", templateController.GetCategoryList)             // 获取分类列表
		categories.GET("/:id", templateController.GetCategoryDetail)       // 获取分类详情

		// 管理员操作 - 分类管理
		authenticated := categories.Group("")
		authenticated.Use(middleware.RequireAuth())
		authenticated.Use(middleware.RequireAdmin())
		{
			authenticated.POST("", templateController.CreateCategory)        // 创建分类
			authenticated.PUT("/:id", templateController.UpdateCategory)     // 更新分类
			authenticated.DELETE("/:id", templateController.DeleteCategory)  // 删除分类
		}
	}
}
