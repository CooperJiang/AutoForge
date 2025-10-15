package routes

import (
	"auto-forge/internal/controllers"
	"auto-forge/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterToolRoutes 注册工具相关路由
func RegisterToolRoutes(r *gin.RouterGroup) {
    controller := controllers.NewToolController()

    // 工具路由 (需要认证)
    toolGroup := r.Group("/tools")
    toolGroup.Use(middleware.RequireAuth())
    {
        toolGroup.GET("", controller.ListTools)           // 获取工具列表
        toolGroup.GET("/:code", controller.GetToolDetail) // 获取工具详情
        toolGroup.POST("/:code/describe-output", controller.DescribeToolOutput) // 动态输出结构
    }
}
