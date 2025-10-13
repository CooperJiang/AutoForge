package workflow

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/dto/response"
	"auto-forge/internal/models"
	"auto-forge/pkg/database"
	log "auto-forge/pkg/logger"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ExecutionService 工作流执行服务
type ExecutionService struct {
	workflowService *WorkflowService
}

// NewExecutionService 创建执行服务实例
func NewExecutionService() *ExecutionService {
	return &ExecutionService{
		workflowService: NewWorkflowService(),
	}
}

// CreateExecution 创建执行记录
func (s *ExecutionService) CreateExecution(workflowID, userID, triggerType string) (*models.WorkflowExecution, error) {
	db := database.GetDB()

	// 检查工作流是否存在
	workflow, err := s.workflowService.GetWorkflowByID(workflowID, userID)
	if err != nil {
		return nil, err
	}

	if !workflow.Enabled && triggerType != "manual" {
		return nil, errors.New("工作流未启用")
	}

	startTime := time.Now().Unix()
	execution := &models.WorkflowExecution{
		WorkflowID:   workflowID,
		UserID:       userID,
		Status:       models.ExecutionStatusPending,
		TriggerType:  triggerType,
		StartTime:    &startTime,
		TotalNodes:   len(workflow.Nodes),
		SuccessNodes: 0,
		FailedNodes:  0,
		SkippedNodes: 0,
		NodeLogs:     models.NodeExecutionLogs{},
	}

	if err := db.Create(execution).Error; err != nil {
		return nil, err
	}

	log.Info("创建工作流执行记录: WorkflowID=%s, ExecutionID=%s, TriggerType=%s",
		workflowID, execution.ID, triggerType)

	return execution, nil
}

// GetExecutionByID 根据ID获取执行记录
func (s *ExecutionService) GetExecutionByID(executionID, userID string) (*models.WorkflowExecution, error) {
	db := database.GetDB()

	var execution models.WorkflowExecution
	if err := db.Where("id = ? AND user_id = ?", executionID, userID).First(&execution).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("执行记录不存在")
		}
		return nil, err
	}

	return &execution, nil
}

// GetExecutionList 获取执行历史列表
func (s *ExecutionService) GetExecutionList(workflowID, userID string, query *request.ExecutionListQuery) (*response.ExecutionListResponse, error) {
	db := database.GetDB()

	// 默认值
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	// 构建查询
	queryDB := db.Model(&models.WorkflowExecution{}).
		Where("workflow_id = ? AND user_id = ?", workflowID, userID)

	// 状态过滤
	if query.Status != "" {
		queryDB = queryDB.Where("status = ?", query.Status)
	}

	// 时间范围过滤
	if query.StartTime != nil {
		queryDB = queryDB.Where("start_time >= ?", *query.StartTime)
	}
	if query.EndTime != nil {
		queryDB = queryDB.Where("start_time <= ?", *query.EndTime)
	}

	// 计算总数
	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	var executions []models.WorkflowExecution
	offset := (query.Page - 1) * query.PageSize
	if err := queryDB.Order("created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&executions).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]response.WorkflowExecutionResponse, len(executions))
	for i, exec := range executions {
		items[i] = s.toExecutionResponse(&exec)
	}

	return &response.ExecutionListResponse{
		Items:    items,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

// UpdateExecutionStatus 更新执行状态
func (s *ExecutionService) UpdateExecutionStatus(executionID, status string, errorMsg string) error {
	db := database.GetDB()

	updates := map[string]interface{}{
		"status": status,
	}

	if status == models.ExecutionStatusSuccess || status == models.ExecutionStatusFailed || status == models.ExecutionStatusCancelled {
		endTime := time.Now().Unix()
		updates["end_time"] = endTime

		// 计算执行时长
		var execution models.WorkflowExecution
		if err := db.First(&execution, "id = ?", executionID).Error; err == nil {
			if execution.StartTime != nil {
				durationMs := (endTime - *execution.StartTime) * 1000
				updates["duration_ms"] = durationMs
			}
		}
	}

	if errorMsg != "" {
		updates["error"] = errorMsg
	}

	if err := db.Model(&models.WorkflowExecution{}).
		Where("id = ?", executionID).
		Updates(updates).Error; err != nil {
		return err
	}

	log.Info("更新执行状态: ExecutionID=%s, Status=%s", executionID, status)
	return nil
}

// AddNodeLog 添加节点执行日志
func (s *ExecutionService) AddNodeLog(executionID string, nodeLog models.NodeExecutionLog) error {
	db := database.GetDB()

	var execution models.WorkflowExecution
	if err := db.First(&execution, "id = ?", executionID).Error; err != nil {
		return err
	}

	// 添加日志
	execution.NodeLogs = append(execution.NodeLogs, nodeLog)

	// 更新统计
	switch nodeLog.Status {
	case models.ExecutionStatusSuccess:
		execution.SuccessNodes++
	case models.ExecutionStatusFailed:
		execution.FailedNodes++
	case "skipped":
		execution.SkippedNodes++
	}

	if err := db.Model(&execution).
		Updates(map[string]interface{}{
			"node_logs":     execution.NodeLogs,
			"success_nodes": execution.SuccessNodes,
			"failed_nodes":  execution.FailedNodes,
			"skipped_nodes": execution.SkippedNodes,
		}).Error; err != nil {
		return err
	}

	return nil
}

// UpdateWorkflowStats 更新工作流统计信息
func (s *ExecutionService) UpdateWorkflowStats(workflowID string, success bool) error {
	db := database.GetDB()

	now := time.Now().Unix()
	updates := map[string]interface{}{
		"total_executions": gorm.Expr("total_executions + 1"),
		"last_executed_at": now,
	}

	if success {
		updates["success_count"] = gorm.Expr("success_count + 1")
	} else {
		updates["failed_count"] = gorm.Expr("failed_count + 1")
	}

	if err := db.Model(&models.Workflow{}).
		Where("id = ?", workflowID).
		Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// DeleteExecution 删除执行记录
func (s *ExecutionService) DeleteExecution(executionID, userID string) error {
	db := database.GetDB()

	// 检查执行记录是否存在且属于该用户
	var execution models.WorkflowExecution
	if err := db.Where("id = ? AND user_id = ?", executionID, userID).First(&execution).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("执行记录不存在")
		}
		return err
	}

	// 删除执行记录
	if err := db.Delete(&execution).Error; err != nil {
		return fmt.Errorf("删除执行记录失败: %v", err)
	}

	log.Info("删除执行记录: ExecutionID=%s, UserID=%s", executionID, userID)
	return nil
}

// ToExecutionResponse 转换为响应格式（公开方法）
func (s *ExecutionService) ToExecutionResponse(execution *models.WorkflowExecution) response.WorkflowExecutionResponse {
	return s.toExecutionResponse(execution)
}

// toExecutionResponse 转换为响应格式（内部方法）
func (s *ExecutionService) toExecutionResponse(execution *models.WorkflowExecution) response.WorkflowExecutionResponse {
	return response.WorkflowExecutionResponse{
		ID:           execution.GetID(),
		WorkflowID:   execution.WorkflowID,
		UserID:       execution.UserID,
		Status:       execution.Status,
		TriggerType:  execution.TriggerType,
		StartTime:    execution.StartTime,
		EndTime:      execution.EndTime,
		DurationMs:   execution.DurationMs,
		TotalNodes:   execution.TotalNodes,
		SuccessNodes: execution.SuccessNodes,
		FailedNodes:  execution.FailedNodes,
		SkippedNodes: execution.SkippedNodes,
		NodeLogs:     execution.NodeLogs,
		Error:        execution.Error,
		CreatedAt:    execution.GetCreatedAt().Unix(),
		UpdatedAt:    execution.GetUpdatedAt().Unix(),
	}
}
