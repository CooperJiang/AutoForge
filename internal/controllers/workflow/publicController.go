package workflow

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/models"
	"auto-forge/internal/services/workflow"
	"auto-forge/pkg/errors"
	log "auto-forge/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func InvokeWorkflow(c *gin.Context) {

	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {

		auth := c.GetHeader("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			apiKey = strings.TrimPrefix(auth, "Bearer ")
		}
	}

	if apiKey == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "缺少 API Key"))
		return
	}

	// 判断 Content-Type
	contentType := c.GetHeader("Content-Type")
	var req request.InvokeWorkflowRequest

	if strings.Contains(contentType, "multipart/form-data") {
		// 处理 FormData（包含文件）
		if err := parseFormData(c, &req); err != nil {
			errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "解析 FormData 失败: "+err.Error()))
			return
		}
	} else {
		// 处理 JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
			return
		}
	}

	mode := c.Query("mode")
	if mode == "" {
		mode = "async"
	}

	svc := workflow.NewWorkflowService()

	wf, err := svc.GetWorkflowByAPIKey(apiKey)
	if err != nil {
		log.Error("API Key 验证失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "无效的 API Key"))
		return
	}

	if err := svc.ValidateAPIParams(wf, req.Params); err != nil {
		log.Error("验证 API 参数失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	svc.IncrementAPICallCount(wf.GetID())

	if mode == "sync" {

		handleSyncInvoke(c, wf, req.Params)
	} else {

		handleAsyncInvoke(c, wf, req.WebhookURL, req.Params)
	}
}

func handleSyncInvoke(c *gin.Context, wf *models.Workflow, externalParams map[string]interface{}) {
	svc := workflow.NewWorkflowService()
	executionSvc := workflow.NewExecutionService()

	execution, err := executionSvc.CreateExecution(wf.GetID(), wf.UserID, "api")
	if err != nil {
		log.Error("创建执行记录失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "执行失败"))
		return
	}

	result, err := svc.ExecuteWorkflowSync(execution.GetID(), wf.UserID, wf.APITimeout, externalParams)
	if err != nil {
		log.Error("同步执行工作流失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "执行失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

func handleAsyncInvoke(c *gin.Context, wf *models.Workflow, webhookURL string, externalParams map[string]interface{}) {
	executionSvc := workflow.NewExecutionService()
	engineSvc := workflow.NewEngineService()

	if webhookURL == "" {
		webhookURL = wf.APIWebhookURL
	}

	execution, err := executionSvc.CreateExecution(wf.GetID(), wf.UserID, "api")
	if err != nil {
		log.Error("创建执行记录失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "执行失败"))
		return
	}

	userID := wf.UserID
	execID := execution.GetID()
	go func(webhookURL, execID, userID string, externalParams map[string]interface{}) {
		if err := engineSvc.ExecuteWorkflow(execID, nil, externalParams); err != nil {
			log.Error("异步执行工作流失败: ExecutionID=%s, Error=%v", execID, err)
		} else {
			log.Info("异步执行工作流完成: ExecutionID=%s", execID)

			if webhookURL != "" {
				sendWebhookNotification(webhookURL, execID, userID)
			}
		}
	}(webhookURL, execID, userID, externalParams)

	errors.ResponseSuccess(c, gin.H{
		"execution_id": execution.GetID(),
		"status":       "running",
		"message":      "工作流已开始执行",
		"webhook_url":  webhookURL,
	}, "已接受")
}

func sendWebhookNotification(webhookURL string, executionID string, userID string) {
	executionSvc := workflow.NewExecutionService()

	execution, err := executionSvc.GetExecutionByID(executionID, userID)
	if err != nil {
		log.Error("查询执行记录失败: ExecutionID=%s, Error=%v", executionID, err)
		return
	}

	outputs := make(map[string]interface{})
	for _, nodeLog := range execution.NodeLogs {
		if len(nodeLog.Output) > 0 {
			outputs[nodeLog.NodeID] = nodeLog.Output
		}
	}

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

func parseFormData(c *gin.Context, req *request.InvokeWorkflowRequest) error {
	// 1. 获取 JSON 参数（从 params 字段）
	paramsJSON := c.PostForm("params")
	if paramsJSON != "" {
		if err := json.Unmarshal([]byte(paramsJSON), &req.Params); err != nil {
			return fmt.Errorf("解析 params 参数失败: %w", err)
		}
	}

	if req.Params == nil {
		req.Params = make(map[string]interface{})
	}

	// 2. 获取 webhook_url（可选）
	if webhookURL := c.PostForm("webhook_url"); webhookURL != "" {
		req.WebhookURL = webhookURL
	}

	// 3. 处理文件参数
	form, err := c.MultipartForm()
	if err != nil {
		// 如果没有文件也没关系
		return nil
	}

	if form.File != nil {
		for paramKey, fileHeaders := range form.File {
			// 跳过 params 和 webhook_url 字段
			if paramKey == "params" || paramKey == "webhook_url" {
				continue
			}

			if len(fileHeaders) > 0 {
				fileHeader := fileHeaders[0]

				// 生成临时执行 ID（后续会被替换）
				tempExecID := fmt.Sprintf("temp_%d", time.Now().UnixNano())

				// 保存文件并获取文件信息对象
				fileInfo, err := saveUploadedFile(c, tempExecID, paramKey, fileHeader)
				if err != nil {
					return fmt.Errorf("保存文件 %s 失败: %w", paramKey, err)
				}

				// 将文件信息注入到 params
				req.Params[paramKey] = fileInfo
			}
		}
	}

	return nil
}

func saveUploadedFile(c *gin.Context, executionID, paramKey string, fileHeader *multipart.FileHeader) (map[string]interface{}, error) {
	// 1. 创建临时目录（简化路径，执行完就删除）
	baseDir := "/tmp/workflow-files"
	dir := fmt.Sprintf("%s/%s", baseDir, executionID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 2. 安全处理文件名
	filename := filepath.Base(fileHeader.Filename)
	filename = strings.ReplaceAll(filename, "..", "")
	filePath := filepath.Join(dir, filename)

	// 3. 保存文件
	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 4. 检测 MIME 类型
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, _ := file.Read(buffer)
	mimeType := http.DetectContentType(buffer[:n])

	// 5. 返回文件信息对象（给工具使用，不包含 url）
	fileInfo := map[string]interface{}{
		"type":      "file",
		"path":      filePath,        // 文件路径，供工具读取
		"filename":  filename,        // 原始文件名
		"size":      fileHeader.Size, // 文件大小
		"mime_type": mimeType,        // MIME 类型
	}

	log.Info("临时文件已保存: %s, 大小: %d bytes, MIME: %s (执行完成后会自动删除)", filePath, fileHeader.Size, mimeType)

	return fileInfo, nil
}
