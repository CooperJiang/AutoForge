package controllers

import (
	"auto-forge/internal/services/tool"
	"auto-forge/pkg/errors"

	"github.com/gin-gonic/gin"
)

type ToolController struct {
	toolService *tool.ToolService
}

func NewToolController() *ToolController {
	return &ToolController{
		toolService: tool.GetToolService(),
	}
}

func (c *ToolController) ListTools(ctx *gin.Context) {
	tools, err := c.toolService.ListTools()
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "获取工具列表失败"))
		return
	}

	errors.ResponseSuccess(ctx, tools, "获取成功")
}

func (c *ToolController) GetToolDetail(ctx *gin.Context) {
	code := ctx.Param("code")

	tool, err := c.toolService.GetToolByCode(code)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeNotFound, "工具不存在"))
		return
	}

	errors.ResponseSuccess(ctx, tool, "获取成功")
}

func (c *ToolController) DescribeToolOutput(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "工具代码不能为空"))
		return
	}

	var cfg map[string]interface{}
	if err := ctx.ShouldBindJSON(&cfg); err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	schema, err := c.toolService.DescribeOutput(code, cfg)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "获取动态输出结构失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(ctx, schema, "获取成功")
}

// ListToolsGroupedByCategory 获取按分类分组的工具列表
func (c *ToolController) ListToolsGroupedByCategory(ctx *gin.Context) {
	categories, err := c.toolService.GetToolsGroupedByCategory()
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "获取工具列表失败"))
		return
	}

	errors.ResponseSuccess(ctx, categories, "获取成功")
}
