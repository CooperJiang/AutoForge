package controllers

import (
	"auto-forge/internal/services/tool"
	"auto-forge/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ToolController 工具控制器
type ToolController struct {
	toolService *tool.ToolService
}

// NewToolController 创建工具控制器
func NewToolController() *ToolController {
	return &ToolController{
		toolService: tool.GetToolService(),
	}
}

// ListTools 获取工具列表
func (c *ToolController) ListTools(ctx *gin.Context) {
	tools, err := c.toolService.ListTools()
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "获取工具列表失败"))
		return
	}

	errors.ResponseSuccess(ctx, tools, "获取成功")
}

// GetToolDetail 获取工具详情
func (c *ToolController) GetToolDetail(ctx *gin.Context) {
	code := ctx.Param("code")

	tool, err := c.toolService.GetToolByCode(code)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeNotFound, "工具不存在"))
		return
	}

	errors.ResponseSuccess(ctx, tool, "获取成功")
}
