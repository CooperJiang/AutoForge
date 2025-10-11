package routes

import (
	taskController "auto-forge/internal/controllers/task"
	"auto-forge/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterTaskRoutes 注册任务相关路由
func RegisterTaskRoutes(r *gin.RouterGroup) {
	// 所有任务路由都需要认证
	tasks := r.Group("/tasks")
	tasks.Use(middleware.RequireAuth())
	{
		// 任务管理
		tasks.POST("", taskController.CreateTask)                        // 创建任务
		tasks.GET("", taskController.GetTaskList)                        // 获取任务列表
		tasks.GET("/:id", taskController.GetTask)                        // 获取任务详情
		tasks.PUT("/:id", taskController.UpdateTask)                     // 更新任务
		tasks.DELETE("/:id", taskController.DeleteTask)                  // 删除任务
		tasks.POST("/:id/enable", taskController.EnableTask)             // 启用任务
		tasks.POST("/:id/disable", taskController.DisableTask)           // 禁用任务
		tasks.POST("/:id/trigger", taskController.TriggerTask)           // 手动触发任务
		tasks.GET("/:id/executions", taskController.GetTaskExecutions)   // 获取任务执行记录
		tasks.DELETE("/:id/executions", taskController.DeleteAllExecutions) // 删除任务的所有执行记录
		tasks.POST("/test", taskController.TestTask)                     // 测试任务配置
	}

	// 执行记录路由也需要认证
	executions := r.Group("/executions")
	executions.Use(middleware.RequireAuth())
	{
		executions.GET("/:id", taskController.GetExecution)       // 获取执行记录详情
		executions.DELETE("/:id", taskController.DeleteExecution) // 删除执行记录
	}
}
