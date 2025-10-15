package routes

import (
	workflowController "auto-forge/internal/controllers/workflow"
	"auto-forge/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterWorkflowRoutes 注册工作流相关路由
func RegisterWorkflowRoutes(r *gin.RouterGroup) {
	// 需要认证的工作流路由
	workflows := r.Group("/workflows")
	workflows.Use(middleware.RequireAuth())
	{
		// 工作流 CRUD
		workflows.POST("", workflowController.CreateWorkflow)              // 创建工作流
		workflows.GET("", workflowController.GetWorkflowList)              // 获取工作流列表
		workflows.GET("/:id", workflowController.GetWorkflowByID)          // 获取工作流详情
		workflows.PUT("/:id", workflowController.UpdateWorkflow)           // 更新工作流
		workflows.DELETE("/:id", workflowController.DeleteWorkflow)        // 删除工作流
		workflows.PATCH("/:id/toggle", workflowController.ToggleEnabled)   // 切换启用状态
		workflows.GET("/:id/stats", workflowController.GetWorkflowStats)   // 获取统计信息

		// 工作流执行
		workflows.POST("/:id/execute", workflowController.ExecuteWorkflow)                         // 执行工作流
		workflows.GET("/:id/executions", workflowController.GetExecutionList)                      // 获取执行历史
		workflows.GET("/:id/executions/:executionId", workflowController.GetExecutionDetail)       // 获取执行详情
		workflows.DELETE("/:id/executions/:executionId", workflowController.DeleteExecution)       // 删除执行记录
		workflows.POST("/:id/executions/:executionId/stop", workflowController.StopExecution)      // 停止执行

		// 工作流验证
		workflows.POST("/validate", workflowController.ValidateWorkflow)   // 验证工作流配置

		// API 管理
		workflows.POST("/:id/api/enable", workflowController.EnableAPI)              // 启用 API
		workflows.POST("/:id/api/disable", workflowController.DisableAPI)            // 禁用 API
		workflows.POST("/:id/api/regenerate", workflowController.RegenerateAPIKey)   // 重新生成 API Key
		workflows.PUT("/:id/api/params", workflowController.UpdateAPIParams)         // 更新 API 参数配置
		workflows.PUT("/:id/api/timeout", workflowController.UpdateAPITimeout)       // 更新 API 超时时间
		workflows.PUT("/:id/api/webhook", workflowController.UpdateAPIWebhook)       // 更新 Webhook URL
	}
}
