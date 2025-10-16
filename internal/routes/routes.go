package routes

import (
	"auto-forge/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 应用CORS中间件
	r.Use(middleware.CORSMiddleware())

	// 注册前端路由
	RegisterClientRoutes(r)

	// 注册文件访问路由（独立路由，不在 /api/v1 下）
	RegisterFileRoutes(r)

	// 注册OAuth2回调路由（独立路由，不在 /api/v1 下）
	RegisterOAuth2CallbackRoutes(r)

	// 注册预览路由（独立路由，不在 /api/v1 下）
	RegisterPreviewRoutes(r)

	prefix := r.Group("/api")
	version := prefix.Group("/v1")

	{
		// 配置相关路由（公开接口）
		RegisterConfigRoutes(version)

		// 用户相关路由
		userRoutes := version.Group("/user")
		RegisterUserRoutes(userRoutes)

		// 上传相关路由
		RegisterUploadRoutes(version)

		// 任务相关路由
		RegisterTaskRoutes(version)

		// 管理员相关路由
		adminRoutes := version.Group("/admin")
		RegisterAdminRoutes(adminRoutes)

		// OAuth2相关路由
		oauth2Routes := version.Group("/auth")
		RegisterOAuth2Routes(oauth2Routes)

		// 工具相关路由
		RegisterToolRoutes(version)

		// 工作流相关路由
		RegisterWorkflowRoutes(version)

		// 公开工作流调用路由（无需认证）
		RegisterPublicWorkflowRoutes(version)

		// 模板市场路由
		RegisterTemplateRoutes(version)

		// 在这里添加其他模块路由
		// 例如：
		// productRoutes := api.Group("/product")
		// RegisterProductRoutes(productRoutes)
	}

	// 静态文件服务 - 暂时注释掉，使用embed版本
	// r.Static("/static", "./internal/static")
}
