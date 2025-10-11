package routes

import (
	"auto-forge/internal/controllers"
	"auto-forge/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes 注册管理员路由
func RegisterAdminRoutes(r *gin.RouterGroup) {
	adminController := controllers.NewAdminController()

	// 登录接口（无需认证）
	r.POST("/login", adminController.Login)

	// 需要认证的接口
	auth := r.Group("")
	auth.Use(middleware.AdminAuth())
	{
		// 登出
		auth.POST("/logout", adminController.Logout)

		// 任务管理
		auth.GET("/tasks", adminController.GetTasks)
		auth.PUT("/tasks/:id/status", adminController.UpdateTaskStatus)
		auth.DELETE("/tasks/:id", adminController.DeleteTask)
		auth.POST("/tasks/:id/execute", adminController.ExecuteTask)

		// 执行记录
		auth.GET("/executions", adminController.GetExecutions)

		// 统计数据
		auth.GET("/stats", adminController.GetStats)

		// 用户管理
		auth.GET("/users", adminController.GetUsers)
		auth.PUT("/users/:id/status", adminController.UpdateUserStatus)
	}
}
