package routes

import (
	agentController "auto-forge/internal/controllers/agent"
	"auto-forge/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAgentRoutes 注册 Agent 路由
func RegisterAgentRoutes(r *gin.Engine) {
	agentGroup := r.Group("/api/v1/agent")
	agentGroup.Use(middleware.RequireAuth()) // 需要认证

	{
		// 对话管理
		agentGroup.POST("/conversations", agentController.CreateConversation)       // 创建对话
		agentGroup.GET("/conversations", agentController.GetConversations)          // 获取对话列表
		agentGroup.GET("/conversations/:id", agentController.GetConversationByID)   // 获取对话详情
		agentGroup.PUT("/conversations/:id", agentController.UpdateConversation)    // 更新对话
		agentGroup.DELETE("/conversations/:id", agentController.DeleteConversation) // 删除对话

		// 消息管理
		agentGroup.GET("/conversations/:id/messages", agentController.GetMessages)  // 获取消息列表
		agentGroup.POST("/conversations/:id/messages", agentController.SendMessage) // 发送消息（支持流式）
	}
}
