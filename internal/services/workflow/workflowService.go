package workflow

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/dto/response"
	"auto-forge/internal/models"
	"auto-forge/pkg/database"
	log "auto-forge/pkg/logger"
	"auto-forge/pkg/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

	// 计算下次执行时间
	var nextRunTime *int64
	if req.Enabled && req.ScheduleType != "" && req.ScheduleType != "manual" {
		t := s.CalculateNextRunTime(req.ScheduleType, req.ScheduleValue)
		nextRunTime = &t
	}

	// 提取外部触发器参数配置
	apiParams, err := s.ExtractExternalTriggerParams(req.Nodes, req.Edges)
	if err != nil {
		return nil, fmt.Errorf("提取外部触发器参数失败: %w", err)
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
		NextRunTime:   nextRunTime,
		APIParams:     apiParams,
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

		// 提取并更新外部触发器参数配置
		apiParams, err := s.ExtractExternalTriggerParams(*req.Nodes, edges)
		if err != nil {
			return nil, fmt.Errorf("提取外部触发器参数失败: %w", err)
		}
		updates["api_params"] = apiParams
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

		// 如果边更新了但节点没更新，仍需重新提取参数
		if req.Nodes == nil {
			apiParams, err := s.ExtractExternalTriggerParams(nodes, *req.Edges)
			if err != nil {
				return nil, fmt.Errorf("提取外部触发器参数失败: %w", err)
			}
			updates["api_params"] = apiParams
		}
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

	// 如果调度配置有变化，重新计算下次执行时间
	scheduleType := workflow.ScheduleType
	scheduleValue := workflow.ScheduleValue
	enabled := workflow.Enabled

	if req.ScheduleType != nil {
		scheduleType = *req.ScheduleType
	}
	if req.ScheduleValue != nil {
		scheduleValue = *req.ScheduleValue
	}
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	if enabled && scheduleType != "" && scheduleType != "manual" {
		nextRunTime := s.CalculateNextRunTime(scheduleType, scheduleValue)
		updates["next_run_time"] = nextRunTime
	} else {
		updates["next_run_time"] = nil
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

// ExtractExternalTriggerParams 从节点中提取外部触发器参数配置
func (s *WorkflowService) ExtractExternalTriggerParams(nodes []models.WorkflowNode, edges []models.WorkflowEdge) (models.WorkflowAPIParams, error) {
	// 找到没有入边的节点（起始节点）
	targetNodes := make(map[string]bool)
	for _, edge := range edges {
		targetNodes[edge.Target] = true
	}

	var startNodes []models.WorkflowNode
	for _, node := range nodes {
		if !targetNodes[node.ID] {
			startNodes = append(startNodes, node)
		}
	}

	// 查找 external_trigger 类型的起始节点
	var externalTriggerNode *models.WorkflowNode
	for _, node := range startNodes {
		if node.Type == "external_trigger" {
			externalTriggerNode = &node
			break
		}
	}

	// 如果没有找到 external_trigger 节点，返回空参数列表
	if externalTriggerNode == nil {
		return models.WorkflowAPIParams{}, nil
	}

	// 从节点配置中提取参数
	params := models.WorkflowAPIParams{}

	// 获取节点配置中的 params 数组
	if paramsInterface, ok := externalTriggerNode.Config["params"]; ok {
		if paramsList, ok := paramsInterface.([]interface{}); ok {
			for _, paramInterface := range paramsList {
				if paramMap, ok := paramInterface.(map[string]interface{}); ok {
					param := models.WorkflowAPIParam{
						Key:      getStringValue(paramMap, "key"),
						Type:     getStringValue(paramMap, "type"),
						Required: getBoolValue(paramMap, "required"),
					}

					// 可选字段
					if desc := getStringValue(paramMap, "description"); desc != "" {
						param.Description = desc
					}
					if defaultVal, ok := paramMap["defaultValue"]; ok {
						param.DefaultValue = defaultVal
					}
					if example, ok := paramMap["example"]; ok {
						param.Example = example
					}

					params = append(params, param)
				}
			}
		}
	}

	return params, nil
}

// Helper functions for safe type conversion
func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getBoolValue(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
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
		APIEnabled:      workflow.APIEnabled,
		APIKey:          workflow.APIKey,
		APIParams:       workflow.APIParams,
		APITimeout:      workflow.APITimeout,
		APIWebhookURL:   workflow.APIWebhookURL,
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

// CalculateNextRunTime 计算下次执行时间（公开方法）
func (s *WorkflowService) CalculateNextRunTime(scheduleType, scheduleValue string) int64 {
	now := time.Now()

	switch scheduleType {
	case "daily":
		// 解析时间 HH:MM:SS
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 3 {
			return now.Unix()
		}
		hour, _ := strconv.Atoi(parts[0])
		minute, _ := strconv.Atoi(parts[1])
		second, _ := strconv.Atoi(parts[2])

		// 设置为今天的指定时间
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())

		// 如果已经过了今天的时间，设置为明天
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		return next.Unix()

	case "weekly":
		// 解析 day1,day2,...:HH:MM:SS
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 4 {
			return now.Unix()
		}
		dayStrs := strings.Split(parts[0], ",")
		hour, _ := strconv.Atoi(parts[1])
		minute, _ := strconv.Atoi(parts[2])
		second, _ := strconv.Atoi(parts[3])

		// 转换星期几为整数
		var weekdays []int
		for _, dayStr := range dayStrs {
			day, _ := strconv.Atoi(dayStr)
			weekdays = append(weekdays, day)
		}

		// 找到最近的执行日
		currentWeekday := int(now.Weekday())
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())

		// 查找最近的星期几
		minDaysToAdd := 7
		for _, targetWeekday := range weekdays {
			daysToAdd := (targetWeekday - currentWeekday + 7) % 7
			// 如果是今天且时间未过，则为0天
			if daysToAdd == 0 && next.After(now) {
				minDaysToAdd = 0
				break
			}
			// 如果是今天但时间已过，需要等到下周
			if daysToAdd == 0 {
				daysToAdd = 7
			}
			if daysToAdd < minDaysToAdd {
				minDaysToAdd = daysToAdd
			}
		}

		next = next.Add(time.Duration(minDaysToAdd) * 24 * time.Hour)
		return next.Unix()

	case "monthly":
		// 解析 day:HH:MM:SS
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 4 {
			return now.Unix()
		}
		day, _ := strconv.Atoi(parts[0])
		hour, _ := strconv.Atoi(parts[1])
		minute, _ := strconv.Atoi(parts[2])
		second, _ := strconv.Atoi(parts[3])

		// 设置为本月的指定日期和时间
		next := time.Date(now.Year(), now.Month(), day, hour, minute, second, 0, now.Location())

		// 如果已经过了本月的时间，设置为下个月
		if next.Before(now) {
			next = next.AddDate(0, 1, 0)
		}

		// 处理月末不存在的日期（如2月30日）
		for next.Day() != day {
			next = next.AddDate(0, 1, 0)
			next = time.Date(next.Year(), next.Month(), day, hour, minute, second, 0, next.Location())
		}

		return next.Unix()

	case "hourly":
		// 解析 MM:SS
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 2 {
			return now.Unix()
		}
		minute, _ := strconv.Atoi(parts[0])
		second, _ := strconv.Atoi(parts[1])

		// 设置为本小时的指定时间
		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), minute, second, 0, now.Location())

		// 如果已经过了本小时的时间，设置为下一小时
		if next.Before(now) {
			next = next.Add(1 * time.Hour)
		}
		return next.Unix()

	case "interval":
		// 间隔秒数
		seconds, _ := strconv.ParseInt(scheduleValue, 10, 64)
		return now.Add(time.Duration(seconds) * time.Second).Unix()

	case "cron":
		// cron 表达式的处理需要使用 cron 库
		// 这里暂时返回一分钟后
		return now.Add(1 * time.Minute).Unix()

	default:
		return now.Unix()
	}
}

// EnableWorkflowAPI 启用工作流 API 并生成 API Key
func (s *WorkflowService) EnableWorkflowAPI(workflowID, userID string) (string, error) {
	db := database.GetDB()

	// 查询工作流
	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return "", fmt.Errorf("工作流不存在")
	}

	// 如果已启用，返回现有 API Key
	if workflow.APIEnabled && workflow.APIKey != "" {
		return workflow.APIKey, nil
	}

	// 生成新的 API Key
	apiKey, err := utils.GenerateWorkflowAPIKey()
	if err != nil {
		return "", fmt.Errorf("生成 API Key 失败: %w", err)
	}

	// 更新数据库
	updates := map[string]interface{}{
		"api_enabled": true,
		"api_key":     apiKey,
	}

	if err := db.Model(&workflow).Updates(updates).Error; err != nil {
		return "", fmt.Errorf("启用 API 失败: %w", err)
	}

	log.Info("工作流 API 已启用: WorkflowID=%s, APIKey=%s", workflowID, apiKey)
	return apiKey, nil
}

// DisableWorkflowAPI 禁用工作流 API
func (s *WorkflowService) DisableWorkflowAPI(workflowID, userID string) error {
	db := database.GetDB()

	// 查询工作流
	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return fmt.Errorf("工作流不存在")
	}

	// 更新数据库
	updates := map[string]interface{}{
		"api_enabled": false,
	}

	if err := db.Model(&workflow).Updates(updates).Error; err != nil {
		return fmt.Errorf("禁用 API 失败: %w", err)
	}

	log.Info("工作流 API 已禁用: WorkflowID=%s", workflowID)
	return nil
}

// RegenerateAPIKey 重新生成 API Key
func (s *WorkflowService) RegenerateAPIKey(workflowID, userID string) (string, error) {
	db := database.GetDB()

	// 查询工作流
	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return "", fmt.Errorf("工作流不存在")
	}

	if !workflow.APIEnabled {
		return "", fmt.Errorf("工作流 API 未启用")
	}

	// 生成新的 API Key
	apiKey, err := utils.GenerateWorkflowAPIKey()
	if err != nil {
		return "", fmt.Errorf("生成 API Key 失败: %w", err)
	}

	// 更新数据库
	if err := db.Model(&workflow).Update("api_key", apiKey).Error; err != nil {
		return "", fmt.Errorf("更新 API Key 失败: %w", err)
	}

	log.Info("工作流 API Key 已重新生成: WorkflowID=%s, APIKey=%s", workflowID, apiKey)
	return apiKey, nil
}

// UpdateAPIParams 更新工作流 API 参数配置
func (s *WorkflowService) UpdateAPIParams(workflowID, userID string, params models.WorkflowAPIParams) error {
	db := database.GetDB()

	// 查询工作流
	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return fmt.Errorf("工作流不存在")
	}

	// 更新 API 参数配置
	if err := db.Model(&workflow).Update("api_params", params).Error; err != nil {
		return fmt.Errorf("更新 API 参数配置失败: %w", err)
	}

	log.Info("工作流 API 参数已更新: WorkflowID=%s, ParamsCount=%d", workflowID, len(params))
	return nil
}

// UpdateAPITimeout 更新工作流 API 超时时间
func (s *WorkflowService) UpdateAPITimeout(workflowID, userID string, timeout int) error {
	db := database.GetDB()

	// 查询工作流
	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return fmt.Errorf("工作流不存在")
	}

	// 更新超时时间
	if err := db.Model(&workflow).Update("api_timeout", timeout).Error; err != nil {
		return fmt.Errorf("更新 API 超时时间失败: %w", err)
	}

	log.Info("工作流 API 超时时间已更新: WorkflowID=%s, Timeout=%d", workflowID, timeout)
	return nil
}

// GetWorkflowByAPIKey 通过 API Key 获取工作流
func (s *WorkflowService) GetWorkflowByAPIKey(apiKey string) (*models.Workflow, error) {
	db := database.GetDB()

	var workflow models.Workflow
	if err := db.Where("api_key = ? AND api_enabled = ?", apiKey, true).First(&workflow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("无效的 API Key")
		}
		return nil, err
	}

	return &workflow, nil
}

// ValidateAPIParams 验证 API 参数
func (s *WorkflowService) ValidateAPIParams(workflow *models.Workflow, userParams map[string]interface{}) error {
	if len(workflow.APIParams) == 0 {
		return nil
	}

	// 验证每个参数
	for _, param := range workflow.APIParams {
		value, exists := userParams[param.Key]

		// 处理必填参数
		if !exists {
			if param.Required {
				return fmt.Errorf("缺少必填参数: %s", param.Key)
			}
			continue
		}

		// 类型校验
		if err := validateParamType(value, param.Type); err != nil {
			return fmt.Errorf("参数 %s 类型错误: %w", param.Key, err)
		}
	}

	return nil
}

// ApplyAPIParams 应用 API 参数到工作流 (已废弃，保留用于兼容)
func (s *WorkflowService) ApplyAPIParams(workflow *models.Workflow, userParams map[string]interface{}) error {
	if len(workflow.APIParams) == 0 {
		return nil
	}

	// 将 Nodes 转换为 []interface{} 以便使用 path 工具
	nodesData := make([]interface{}, len(workflow.Nodes))
	for i, node := range workflow.Nodes {
		nodesData[i] = map[string]interface{}{
			"id":       node.ID,
			"type":     node.Type,
			"toolCode": node.ToolCode,
			"name":     node.Name,
			"config":   node.Config,
			"retry":    node.Retry,
			"position": node.Position,
			"data":     node.Data,
		}
	}

	workflowData := map[string]interface{}{
		"nodes": nodesData,
	}

	// 应用每个参数
	for _, param := range workflow.APIParams {
		value, exists := userParams[param.Key]

		// 处理必填参数
		if !exists {
			if param.Required {
				return fmt.Errorf("缺少必填参数: %s", param.Key)
			}
			// 使用默认值
			if param.DefaultValue != nil {
				value = param.DefaultValue
			} else {
				continue
			}
		}

		// 类型校验（简单校验）
		if err := validateParamType(value, param.Type); err != nil {
			return fmt.Errorf("参数 %s 类型错误: %w", param.Key, err)
		}

		// 应用参数到指定路径
		if err := utils.SetValueByPath(workflowData, param.MappingPath, value); err != nil {
			return fmt.Errorf("应用参数 %s 失败: %w", param.Key, err)
		}

		log.Info("应用参数: %s = %v -> %s", param.Key, value, param.MappingPath)
	}

	// 将修改后的数据转换回 workflow.Nodes
	if nodesInterface, ok := workflowData["nodes"].([]interface{}); ok {
		for i, nodeInterface := range nodesInterface {
			if nodeMap, ok := nodeInterface.(map[string]interface{}); ok {
				workflow.Nodes[i].Config = nodeMap["config"].(map[string]interface{})
			}
		}
	}

	return nil
}

// validateParamType 验证参数类型
func validateParamType(value interface{}, expectedType string) error {
	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("期望类型 string, 实际类型 %T", value)
		}
	case "number":
		switch value.(type) {
		case int, int64, float64, float32:
			return nil
		default:
			return fmt.Errorf("期望类型 number, 实际类型 %T", value)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("期望类型 boolean, 实际类型 %T", value)
		}
	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("期望类型 object, 实际类型 %T", value)
		}
	case "array":
		if _, ok := value.([]interface{}); !ok {
			return fmt.Errorf("期望类型 array, 实际类型 %T", value)
		}
	}
	return nil
}

// ExecuteWorkflowSync 同步执行工作流（等待完成）
func (s *WorkflowService) ExecuteWorkflowSync(executionID, userID string, timeoutSeconds int, externalParams map[string]interface{}) (map[string]interface{}, error) {
	engineSvc := NewEngineService()
	executionSvc := NewExecutionService()

	// 创建通道用于接收结果
	done := make(chan error, 1)

	// 异步执行
	go func() {
		err := engineSvc.ExecuteWorkflow(executionID, nil, externalParams)
		done <- err
	}()

	// 等待执行完成或超时
	timeout := time.Duration(timeoutSeconds) * time.Second
	select {
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("执行失败: %w", err)
		}

		// 查询执行结果
		execution, err := executionSvc.GetExecutionByID(executionID, userID)
		if err != nil {
			return nil, fmt.Errorf("查询执行结果失败: %w", err)
		}

		// 提取输出结果
		outputs := make(map[string]interface{})
		for _, nodeLog := range execution.NodeLogs {
			if len(nodeLog.Output) > 0 {
				outputs[nodeLog.NodeID] = nodeLog.Output
			}
		}

		// 构建返回结果
		result := map[string]interface{}{
			"execution_id":  execution.GetID(),
			"status":        execution.Status,
			"start_time":    execution.StartTime,
			"end_time":      execution.EndTime,
			"duration_ms":   execution.DurationMs,
			"error":         execution.Error,
			"outputs":       outputs,
		}

		return result, nil

	case <-time.After(timeout):
		return nil, fmt.Errorf("执行超时 (%d 秒)", timeoutSeconds)
	}
}

// IncrementAPICallCount 增加 API 调用次数
func (s *WorkflowService) IncrementAPICallCount(workflowID string) error {
	db := database.GetDB()
	now := time.Now().Unix()

	updates := map[string]interface{}{
		"api_call_count":    gorm.Expr("api_call_count + 1"),
		"api_last_called_at": now,
	}

	if err := db.Model(&models.Workflow{}).Where("id = ?", workflowID).Updates(updates).Error; err != nil {
		log.Error("更新 API 调用统计失败: %v", err)
		return err
	}

	return nil
}

// UpdateAPIWebhook 更新 API Webhook URL
func (s *WorkflowService) UpdateAPIWebhook(workflowID, userID, webhookURL string) error {
	db := database.GetDB()

	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return fmt.Errorf("工作流不存在")
	}

	if err := db.Model(&workflow).Update("api_webhook_url", webhookURL).Error; err != nil {
		return fmt.Errorf("更新 Webhook URL 失败: %w", err)
	}

	log.Info("工作流 API Webhook 已更新: WorkflowID=%s, WebhookURL=%s", workflowID, webhookURL)
	return nil
}
