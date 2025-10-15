package workflow

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/models"
	"auto-forge/internal/services/workflow"
	"auto-forge/pkg/errors"
	log "auto-forge/pkg/logger"

	"github.com/gin-gonic/gin"
)


func EnableAPI(c *gin.Context) {
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

	svc := workflow.NewWorkflowService()
	apiKey, err := svc.EnableWorkflowAPI(workflowID, userID)
	if err != nil {
		log.Error("启用工作流 API 失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "启用失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"api_key": apiKey,
	}, "API 已启用")
}


func DisableAPI(c *gin.Context) {
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

	svc := workflow.NewWorkflowService()
	if err := svc.DisableWorkflowAPI(workflowID, userID); err != nil {
		log.Error("禁用工作流 API 失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "禁用失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, nil, "API 已禁用")
}


func RegenerateAPIKey(c *gin.Context) {
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

	svc := workflow.NewWorkflowService()
	apiKey, err := svc.RegenerateAPIKey(workflowID, userID)
	if err != nil {
		log.Error("重新生成 API Key 失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "重新生成失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"api_key": apiKey,
	}, "API Key 已重新生成")
}


func UpdateAPIParams(c *gin.Context) {
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

	var req request.UpdateAPIParamsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}


	var params models.WorkflowAPIParams
	for _, p := range req.Params {
		params = append(params, models.WorkflowAPIParam{
			Key:          p.Key,
			Type:         p.Type,
			Required:     p.Required,
			DefaultValue: p.DefaultValue,
			Description:  p.Description,
			Example:      p.Example,
			MappingPath:  p.MappingPath,
		})
	}

	svc := workflow.NewWorkflowService()
	if err := svc.UpdateAPIParams(workflowID, userID, params); err != nil{
		log.Error("更新 API 参数失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "更新失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, nil, "API 参数已更新")
}


func UpdateAPITimeout(c *gin.Context) {
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

	var req request.UpdateAPITimeoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	svc := workflow.NewWorkflowService()
	if err := svc.UpdateAPITimeout(workflowID, userID, req.Timeout); err != nil {
		log.Error("更新 API 超时时间失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "更新失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, nil, "API 超时时间已更新")
}


func UpdateAPIWebhook(c *gin.Context) {
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

	var req request.UpdateAPIWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	svc := workflow.NewWorkflowService()
	if err := svc.UpdateAPIWebhook(workflowID, userID, req.WebhookURL); err != nil {
		log.Error("更新 Webhook URL 失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "更新失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, nil, "Webhook URL 已更新")
}
