package workflow

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/dto/response"
	"auto-forge/internal/services/workflow"
	"auto-forge/pkg/errors"
	log "auto-forge/pkg/logger"
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

var executionService = workflow.NewExecutionService()
var engineService = workflow.NewEngineService()

func ExecuteWorkflow(c *gin.Context) {
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

	// 判断 Content-Type
	contentType := c.GetHeader("Content-Type")
	var req request.ExecuteWorkflowRequest

	log.Info("执行工作流 - Content-Type: %s", contentType)

	if strings.Contains(contentType, "multipart/form-data") {
		// 处理 FormData（包含文件）
		log.Info("检测到 FormData 请求，开始解析文件...")
		if err := parseExecuteFormData(c, &req); err != nil {
			errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "解析 FormData 失败: "+err.Error()))
			return
		}
		log.Info("FormData 解析完成，参数: %+v", req.Params)
	} else {
		// 处理 JSON
		log.Info("使用 JSON 格式解析请求")
		_ = c.ShouldBindJSON(&req)
	}

	// 创建执行记录
	execution, err := executionService.CreateExecution(workflowID, userID, "manual")
	if err != nil {
		log.Error("创建执行记录失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "执行失败: "+err.Error()))
		return
	}

	// 如果有文件，需要先保存并更新 executionID
	if req.Params != nil {
		updateExecutionIDInParams(req.Params, execution.GetID())
	}

	// 异步执行工作流
	go func() {
		if err := engineService.ExecuteWorkflow(execution.GetID(), req.EnvVars, req.Params); err != nil {
			log.Error("工作流执行失败: ExecutionID=%s, Error=%v", execution.GetID(), err)
		}
	}()

	errors.ResponseSuccess(c, response.ExecuteWorkflowResponse{
		ExecutionID: execution.GetID(),
		Status:      execution.Status,
		Message:     "工作流已开始执行",
	}, "工作流已开始执行")
}

func GetExecutionList(c *gin.Context) {
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

	var query request.ExecutionListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	result, err := executionService.GetExecutionList(workflowID, userID, &query)
	if err != nil {
		log.Error("获取执行历史失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "获取执行历史失败"))
		return
	}

	errors.ResponseSuccess(c, result, "获取执行历史成功")
}

func GetExecutionDetail(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	executionID := c.Param("executionId")
	if executionID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "执行ID不能为空"))
		return
	}

	execution, err := executionService.GetExecutionByID(executionID, userID)
	if err != nil {
		log.Error("获取执行详情失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeNotFound, "获取执行详情失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, executionService.ToExecutionResponse(execution), "获取执行详情成功")
}

func StopExecution(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	executionID := c.Param("executionId")
	if executionID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "执行ID不能为空"))
		return
	}

	execution, err := executionService.GetExecutionByID(executionID, userID)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeNotFound, "执行记录不存在"))
		return
	}

	if execution.Status != "running" && execution.Status != "pending" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "只能停止运行中的执行"))
		return
	}

	if err := executionService.UpdateExecutionStatus(executionID, "cancelled", "用户手动取消"); err != nil {
		log.Error("停止执行失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "停止失败"))
		return
	}

	errors.ResponseSuccess(c, nil, "已停止执行")
}

func DeleteExecution(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	executionID := c.Param("executionId")
	if executionID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "执行ID不能为空"))
		return
	}

	if err := executionService.DeleteExecution(executionID, userID); err != nil {
		log.Error("删除执行记录失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, nil, "删除成功")
}

// parseExecuteFormData 解析执行工作流的 FormData
func parseExecuteFormData(c *gin.Context, req *request.ExecuteWorkflowRequest) error {
	// 1. 获取 JSON 参数（从 params 字段）
	paramsJSON := c.PostForm("params")
	log.Info("FormData - params 字段: %s", paramsJSON)

	if paramsJSON != "" {
		if err := json.Unmarshal([]byte(paramsJSON), &req.Params); err != nil {
			return fmt.Errorf("解析 params 参数失败: %w", err)
		}
	}

	if req.Params == nil {
		req.Params = make(map[string]interface{})
	}

	// 2. 处理文件参数（使用临时 executionID，稍后会更新）
	form, err := c.MultipartForm()
	if err != nil {
		log.Warn("获取 MultipartForm 失败: %v", err)
		return nil
	}

	if form.File != nil {
		log.Info("FormData - 文件字段数量: %d", len(form.File))

		// 使用临时 ID
		tempExecID := fmt.Sprintf("temp_%d", time.Now().UnixNano())

		for paramKey, fileHeaders := range form.File {
			log.Info("FormData - 处理文件字段: %s, 文件数量: %d", paramKey, len(fileHeaders))

			// 跳过 params 字段
			if paramKey == "params" {
				continue
			}

			if len(fileHeaders) > 0 {
				fileHeader := fileHeaders[0]
				log.Info("FormData - 文件名: %s, 大小: %d bytes", fileHeader.Filename, fileHeader.Size)

				// 保存文件并获取文件信息对象
				fileInfo, err := saveExecutionFile(c, tempExecID, paramKey, fileHeader)
				if err != nil {
					return fmt.Errorf("保存文件 %s 失败: %w", paramKey, err)
				}

				// 将文件信息注入到 params
				req.Params[paramKey] = fileInfo
				log.Info("FormData - 文件 %s 已添加到 params: %+v", paramKey, fileInfo)
			}
		}
	} else {
		log.Warn("FormData - form.File 为 nil")
	}

	return nil
}

// saveExecutionFile 保存执行工作流的上传文件
func saveExecutionFile(c *gin.Context, executionID, paramKey string, fileHeader *multipart.FileHeader) (map[string]interface{}, error) {
	// 1. 创建临时目录
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

	// 5. 返回文件信息对象
	fileInfo := map[string]interface{}{
		"type":      "file",
		"path":      filePath,
		"filename":  filename,
		"size":      fileHeader.Size,
		"mime_type": mimeType,
	}

	log.Info("执行工作流的临时文件已保存: %s, 大小: %d bytes, MIME: %s", filePath, fileHeader.Size, mimeType)

	return fileInfo, nil
}

// updateExecutionIDInParams 更新文件路径中的 executionID
func updateExecutionIDInParams(params map[string]interface{}, newExecutionID string) {
	baseDir := "/tmp/workflow-files"

	for _, value := range params {
		if fileObj, ok := value.(map[string]interface{}); ok {
			if fileType, exists := fileObj["type"]; exists && fileType == "file" {
				// 更新文件路径中的 executionID
				if oldPath, pathExists := fileObj["path"].(string); pathExists {
					// 获取文件名
					filename := filepath.Base(oldPath)

					// 构造新路径
					newDir := filepath.Join(baseDir, newExecutionID)
					newPath := filepath.Join(newDir, filename)

					// 获取旧目录
					oldDir := filepath.Dir(oldPath)

					// 重命名目录
					if oldDir != newDir {
						if err := os.Rename(oldDir, newDir); err != nil {
							log.Error("重命名文件目录失败: OldDir=%s, NewDir=%s, Error=%v", oldDir, newDir, err)
						} else {
							fileObj["path"] = newPath
							log.Info("已更新文件路径: %s -> %s", oldPath, newPath)
						}
					}
				}
			}
		}
	}
}
