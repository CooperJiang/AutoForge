package workflow

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/models"
	"auto-forge/internal/services/workflow"
	"auto-forge/pkg/errors"
	log "auto-forge/pkg/logger"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// InvokeWorkflow 公开调用工作流 API
func InvokeWorkflow(c *gin.Context) {
	// 从 Header 获取 API Key
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		// 也支持从 Authorization Bearer 获取
		auth := c.GetHeader("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			apiKey = strings.TrimPrefix(auth, "Bearer ")
		}
	}

	if apiKey == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "缺少 API Key"))
		return
	}

	// 解析请求参数
	var req request.InvokeWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	// 获取执行模式（同步或异步）
	mode := c.Query("mode") // sync 或 async，默认 async
	if mode == "" {
		mode = "async"
	}

	svc := workflow.NewWorkflowService()

	// 验证 API Key 并获取工作流
	wf, err := svc.GetWorkflowByAPIKey(apiKey)
	if err != nil {
		log.Error("API Key 验证失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "无效的 API Key"))
		return
	}

	// 验证外部参数
	if err := svc.ValidateAPIParams(wf, req.Params); err != nil {
		log.Error("验证 API 参数失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	// 更新调用统计
	svc.IncrementAPICallCount(wf.GetID())

	if mode == "sync" {
		// 同步模式
		handleSyncInvoke(c, wf, req.Params)
	} else {
		// 异步模式
		handleAsyncInvoke(c, wf, req.WebhookURL, req.Params)
	}
}

// handleSyncInvoke 同步执行
func handleSyncInvoke(c *gin.Context, wf *models.Workflow, externalParams map[string]interface{}) {
	svc := workflow.NewWorkflowService()
	executionSvc := workflow.NewExecutionService()

	// 创建执行记录
	execution, err := executionSvc.CreateExecution(wf.GetID(), wf.UserID, "api")
	if err != nil {
		log.Error("创建执行记录失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "执行失败"))
		return
	}

	// 同步执行工作流（带超时），传入外部参数
	result, err := svc.ExecuteWorkflowSync(execution.GetID(), wf.UserID, wf.APITimeout, externalParams)
	if err != nil {
		log.Error("同步执行工作流失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "执行失败: "+err.Error()))
		return
	}

	// 返回结果
	errors.ResponseSuccess(c, result, "执行成功")
}

// handleAsyncInvoke 异步执行
func handleAsyncInvoke(c *gin.Context, wf *models.Workflow, webhookURL string, externalParams map[string]interface{}) {
	executionSvc := workflow.NewExecutionService()
	engineSvc := workflow.NewEngineService()

	// 使用请求中的 webhook 或配置的 webhook
	if webhookURL == "" {
		webhookURL = wf.APIWebhookURL
	}

	// 创建执行记录
	execution, err := executionSvc.CreateExecution(wf.GetID(), wf.UserID, "api")
	if err != nil {
		log.Error("创建执行记录失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "执行失败"))
		return
	}

	// 异步执行工作流，传入外部参数
	userID := wf.UserID
	execID := execution.GetID()
	go func(webhookURL, execID, userID string, externalParams map[string]interface{}) {
		if err := engineSvc.ExecuteWorkflow(execID, nil, externalParams); err != nil {
			log.Error("异步执行工作流失败: ExecutionID=%s, Error=%v", execID, err)
		} else {
			log.Info("异步执行工作流完成: ExecutionID=%s", execID)

			// 如果配置了 webhook，发送回调
			if webhookURL != "" {
				sendWebhookNotification(webhookURL, execID, userID)
			}
		}
	}(webhookURL, execID, userID, externalParams)

	// 立即返回 execution_id
	errors.ResponseSuccess(c, gin.H{
		"execution_id": execution.GetID(),
		"status":       "running",
		"message":      "工作流已开始执行",
		"webhook_url":  webhookURL,
	}, "已接受")
}

// sendWebhookNotification 发送 Webhook 通知
func sendWebhookNotification(webhookURL string, executionID string, userID string) {
	executionSvc := workflow.NewExecutionService()

	// 查询执行结果
	execution, err := executionSvc.GetExecutionByID(executionID, userID)
	if err != nil {
		log.Error("查询执行记录失败: ExecutionID=%s, Error=%v", executionID, err)
		return
	}

	// 提取输出结果
	outputs := make(map[string]interface{})
	for _, nodeLog := range execution.NodeLogs {
		if len(nodeLog.Output) > 0 {
			outputs[nodeLog.NodeID] = nodeLog.Output
		}
	}

	// 构建 Webhook 通知数据
	payload := gin.H{
		"execution_id": execution.GetID(),
		"workflow_id":  execution.WorkflowID,
		"status":       execution.Status,
		"start_time":   execution.StartTime,
		"end_time":     execution.EndTime,
		"duration_ms":  execution.DurationMs,
		"error":        execution.Error,
		"outputs":      outputs,
	}

	// 发送 POST 请求
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Error("序列化 Webhook 数据失败: Error=%v", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Error("发送 Webhook 通知失败: URL=%s, Error=%v", webhookURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Info("Webhook 通知发送成功: URL=%s, ExecutionID=%s, StatusCode=%d", webhookURL, executionID, resp.StatusCode)
	} else {
		log.Warn("Webhook 通知响应异常: URL=%s, ExecutionID=%s, StatusCode=%d", webhookURL, executionID, resp.StatusCode)
	}
}
