package routes

import (
	workflowController "auto-forge/internal/controllers/workflow"

	"github.com/gin-gonic/gin"
)

// RegisterPublicWorkflowRoutes 注册公开工作流路由（无需认证）
func RegisterPublicWorkflowRoutes(rg *gin.RouterGroup) {
	public := rg.Group("/public/workflows")
	{
		// 工作流调用接口（通过 API Key 认证）
		public.POST("/invoke", workflowController.InvokeWorkflow)
	}
}
