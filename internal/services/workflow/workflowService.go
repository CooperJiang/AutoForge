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

// WorkflowService 工作流服务
type WorkflowService struct{}

var workflowChangeCallback func()

// NewWorkflowService 创建工作流服务实例
func NewWorkflowService() *WorkflowService {
	return &WorkflowService{}
}

// SetWorkflowChangeCallback 设置工作流变更回调
func SetWorkflowChangeCallback(callback func()) {
	workflowChangeCallback = callback
}

// CreateWorkflow 创建工作流
func (s *WorkflowService) CreateWorkflow(userID string, req *request.CreateWorkflowRequest) (*models.Workflow, error) {
	db := database.GetDB()

	// 验证工作流配置
	if err := s.ValidateWorkflowConfig(req.Nodes, req.Edges); err != nil {
		return nil, fmt.Errorf("工作流配置无效: %w", err)
	}

	workflow := &models.Workflow{
		UserID:        userID,
		Name:          req.Name,
		Description:   req.Description,
		Nodes:         models.WorkflowNodes(req.Nodes),
		Edges:         models.WorkflowEdges(req.Edges),
		EnvVars:       models.WorkflowEnvVars(req.EnvVars),
		ScheduleType:  req.ScheduleType,
		ScheduleValue: req.ScheduleValue,
		Enabled:       req.Enabled,
	}

	if err := db.Create(workflow).Error; err != nil {
		log.Error("创建工作流失败: %v", err)
		return nil, err
	}

	log.Info("用户 %s 创建工作流: %s (ID: %s)", userID, workflow.Name, workflow.ID)

	// 如果启用了且有调度配置，触发调度器重载
	if workflow.Enabled && workflow.ScheduleType != "" && workflow.ScheduleType != "manual" {
		s.reloadScheduler()
	}

	return workflow, nil
}

// GetWorkflowByID 根据ID获取工作流
func (s *WorkflowService) GetWorkflowByID(workflowID, userID string) (*models.Workflow, error) {
	db := database.GetDB()

	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("工作流不存在")
		}
		return nil, err
	}

	return &workflow, nil
}

// GetWorkflowList 获取工作流列表
func (s *WorkflowService) GetWorkflowList(userID string, query *request.WorkflowListQuery) (*response.WorkflowListResponse, error) {
	db := database.GetDB()

	// 默认值
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	// 构建查询
	queryDB := db.Model(&models.Workflow{}).Where("user_id = ?", userID)

	// 关键词搜索
	if query.Keyword != "" {
		queryDB = queryDB.Where("name LIKE ? OR description LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 启用状态过滤
	if query.Enabled != nil {
		queryDB = queryDB.Where("enabled = ?", *query.Enabled)
	}

	// 计算总数
	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	var workflows []models.Workflow
	offset := (query.Page - 1) * query.PageSize
	if err := queryDB.Order("created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&workflows).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]response.WorkflowResponse, len(workflows))
	for i, wf := range workflows {
		items[i] = s.toWorkflowResponse(&wf)
	}

	return &response.WorkflowListResponse{
		Items:    items,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

// UpdateWorkflow 更新工作流
func (s *WorkflowService) UpdateWorkflow(workflowID, userID string, req *request.UpdateWorkflowRequest) (*models.Workflow, error) {
	db := database.GetDB()

	// 查询工作流
	workflow, err := s.GetWorkflowByID(workflowID, userID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Nodes != nil {
		// 验证节点配置
		edges := workflow.Edges
		if req.Edges != nil {
			edges = *req.Edges
		}
		if err := s.ValidateWorkflowConfig(*req.Nodes, edges); err != nil {
			return nil, fmt.Errorf("工作流配置无效: %w", err)
		}
		updates["nodes"] = models.WorkflowNodes(*req.Nodes)
	}
	if req.Edges != nil {
		nodes := workflow.Nodes
		if req.Nodes != nil {
			nodes = *req.Nodes
		}
		if err := s.ValidateWorkflowConfig(nodes, *req.Edges); err != nil {
			return nil, fmt.Errorf("工作流配置无效: %w", err)
		}
		updates["edges"] = models.WorkflowEdges(*req.Edges)
	}
	if req.EnvVars != nil {
		updates["env_vars"] = models.WorkflowEnvVars(*req.EnvVars)
	}
	if req.ScheduleType != nil {
		updates["schedule_type"] = *req.ScheduleType
	}
	if req.ScheduleValue != nil {
		updates["schedule_value"] = *req.ScheduleValue
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if err := db.Model(workflow).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新加载
	if err := db.First(workflow, "id = ?", workflowID).Error; err != nil {
		return nil, err
	}

	log.Info("用户 %s 更新工作流: %s (ID: %s)", userID, workflow.Name, workflow.ID)

	// 触发调度器重载（任何更新都重载，因为可能影响调度配置）
	s.reloadScheduler()

	return workflow, nil
}

// DeleteWorkflow 删除工作流
func (s *WorkflowService) DeleteWorkflow(workflowID, userID string) error {
	db := database.GetDB()

	// 检查工作流是否存在
	workflow, err := s.GetWorkflowByID(workflowID, userID)
	if err != nil {
		return err
	}

	// 删除工作流
	if err := db.Delete(workflow).Error; err != nil {
		return err
	}

	log.Info("用户 %s 删除工作流: %s (ID: %s)", userID, workflow.Name, workflow.ID)

	// 触发调度器重载
	s.reloadScheduler()

	return nil
}

// ToggleEnabled 切换工作流启用状态
func (s *WorkflowService) ToggleEnabled(workflowID, userID string, enabled bool) (*models.Workflow, error) {
	db := database.GetDB()

	workflow, err := s.GetWorkflowByID(workflowID, userID)
	if err != nil {
		return nil, err
	}

	workflow.Enabled = enabled
	if err := db.Model(workflow).Update("enabled", enabled).Error; err != nil {
		return nil, err
	}

	log.Info("用户 %s %s工作流: %s (ID: %s)",
		userID,
		map[bool]string{true: "启用", false: "禁用"}[enabled],
		workflow.Name,
		workflow.ID,
	)

	// 触发调度器重载
	s.reloadScheduler()

	return workflow, nil
}

// ValidateWorkflowConfig 验证工作流配置
func (s *WorkflowService) ValidateWorkflowConfig(nodes []models.WorkflowNode, edges []models.WorkflowEdge) error {
	if len(nodes) == 0 {
		return errors.New("工作流至少需要一个节点")
	}

	// 创建节点ID集合
	nodeIDs := make(map[string]bool)
	for _, node := range nodes {
		if node.ID == "" {
			return errors.New("节点ID不能为空")
		}
		if nodeIDs[node.ID] {
			return fmt.Errorf("节点ID重复: %s", node.ID)
		}
		nodeIDs[node.ID] = true
	}

	// 验证连接线
	for _, edge := range edges {
		if !nodeIDs[edge.Source] {
			return fmt.Errorf("连接线的源节点不存在: %s", edge.Source)
		}
		if !nodeIDs[edge.Target] {
			return fmt.Errorf("连接线的目标节点不存在: %s", edge.Target)
		}
	}

	// TODO: 检测循环依赖
	// TODO: 验证节点配置完整性

	return nil
}

// GetWorkflowStats 获取工作流统计信息
func (s *WorkflowService) GetWorkflowStats(workflowID, userID string) (*response.WorkflowStatsResponse, error) {
	workflow, err := s.GetWorkflowByID(workflowID, userID)
	if err != nil {
		return nil, err
	}

	db := database.GetDB()

	// 计算平均执行时间
	var avgDuration *int64
	db.Model(&models.WorkflowExecution{}).
		Where("workflow_id = ? AND status = ?", workflowID, models.ExecutionStatusSuccess).
		Select("AVG(duration_ms)").
		Scan(&avgDuration)

	avgDurationMs := int64(0)
	if avgDuration != nil {
		avgDurationMs = *avgDuration
	}

	return &response.WorkflowStatsResponse{
		TotalExecutions: workflow.TotalExecutions,
		SuccessCount:    workflow.SuccessCount,
		FailedCount:     workflow.FailedCount,
		AvgDurationMs:   avgDurationMs,
		LastExecutionAt: workflow.LastExecutedAt,
	}, nil
}

// ToWorkflowResponse 转换为响应格式（公开方法）
func (s *WorkflowService) ToWorkflowResponse(workflow *models.Workflow) response.WorkflowResponse {
	return s.toWorkflowResponse(workflow)
}

// toWorkflowResponse 转换为响应格式（内部方法）
func (s *WorkflowService) toWorkflowResponse(workflow *models.Workflow) response.WorkflowResponse {
	return response.WorkflowResponse{
		ID:              workflow.GetID(),
		UserID:          workflow.UserID,
		Name:            workflow.Name,
		Description:     workflow.Description,
		Nodes:           workflow.Nodes,
		Edges:           workflow.Edges,
		EnvVars:         workflow.EnvVars,
		ScheduleType:    workflow.ScheduleType,
		ScheduleValue:   workflow.ScheduleValue,
		Enabled:         workflow.Enabled,
		NextRunTime:     workflow.NextRunTime,
		TotalExecutions: workflow.TotalExecutions,
		SuccessCount:    workflow.SuccessCount,
		FailedCount:     workflow.FailedCount,
		LastExecutedAt:  workflow.LastExecutedAt,
		CreatedAt:       workflow.GetCreatedAt().Unix(),
		UpdatedAt:       workflow.GetUpdatedAt().Unix(),
	}
}

// GetCurrentTime 获取当前时间戳
func GetCurrentTime() int64 {
	return time.Now().Unix()
}

// reloadScheduler 重新加载工作流调度器
func (s *WorkflowService) reloadScheduler() {
	if workflowChangeCallback != nil {
		go func() {
			workflowChangeCallback()
			log.Info("工作流调度器已触发重新加载")
		}()
	}
}
