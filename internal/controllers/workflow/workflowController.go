package workflow

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/services/workflow"
	"auto-forge/pkg/errors"
	log "auto-forge/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

var workflowService = workflow.NewWorkflowService()

// CreateWorkflow 创建工作流
func CreateWorkflow(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	var req request.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	wf, err := workflowService.CreateWorkflow(userID, &req)
	if err != nil {
		log.Error("创建工作流失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "创建工作流失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, workflowService.ToWorkflowResponse(wf), "创建工作流成功")
}

// GetWorkflowList 获取工作流列表
func GetWorkflowList(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	var query request.WorkflowListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	result, err := workflowService.GetWorkflowList(userID, &query)
	if err != nil {
		log.Error("获取工作流列表失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "获取工作流列表失败"))
		return
	}

	errors.ResponseSuccess(c, result, "获取工作流列表成功")
}

// GetWorkflowByID 获取工作流详情
func GetWorkflowByID(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "工作流ID不能为空"))
		return
	}

	wf, err := workflowService.GetWorkflowByID(workflowID, userID)
	if err != nil {
		log.Error("获取工作流失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeNotFound, "获取工作流失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, workflowService.ToWorkflowResponse(wf), "获取工作流详情成功")
}

// UpdateWorkflow 更新工作流
func UpdateWorkflow(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "工作流ID不能为空"))
		return
	}

	var req request.UpdateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	wf, err := workflowService.UpdateWorkflow(workflowID, userID, &req)
	if err != nil {
		log.Error("更新工作流失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "更新工作流失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, workflowService.ToWorkflowResponse(wf), "更新工作流成功")
}

// DeleteWorkflow 删除工作流
func DeleteWorkflow(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "工作流ID不能为空"))
		return
	}

	if err := workflowService.DeleteWorkflow(workflowID, userID); err != nil {
		log.Error("删除工作流失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "删除工作流失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, nil, "删除成功")
}

// ToggleEnabled 切换工作流启用状态
func ToggleEnabled(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "工作流ID不能为空"))
		return
	}

	var req request.ToggleEnabledRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	wf, err := workflowService.ToggleEnabled(workflowID, userID, req.Enabled)
	if err != nil {
		log.Error("切换工作流状态失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "操作失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, workflowService.ToWorkflowResponse(wf), "操作成功")
}

// GetWorkflowStats 获取工作流统计
func GetWorkflowStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "工作流ID不能为空"))
		return
	}

	stats, err := workflowService.GetWorkflowStats(workflowID, userID)
	if err != nil {
		log.Error("获取工作流统计失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "获取统计失败"))
		return
	}

	errors.ResponseSuccess(c, stats, "获取统计成功")
}

// ValidateWorkflow 验证工作流配置
func ValidateWorkflow(c *gin.Context) {
	var req request.ValidateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	if err := workflowService.ValidateWorkflowConfig(req.Nodes, req.Edges); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"valid":  false,
				"errors": []string{err.Error()},
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"valid":  true,
			"errors": nil,
		},
	})
}
