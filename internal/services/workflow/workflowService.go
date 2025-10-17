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

type WorkflowService struct{}

var workflowChangeCallback func()

func NewWorkflowService() *WorkflowService {
	return &WorkflowService{}
}

func SetWorkflowChangeCallback(callback func()) {
	workflowChangeCallback = callback
}

func (s *WorkflowService) CreateWorkflow(userID string, req *request.CreateWorkflowRequest) (*models.Workflow, error) {
	db := database.GetDB()

	if err := s.ValidateWorkflowConfig(req.Nodes, req.Edges); err != nil {
		return nil, fmt.Errorf("工作流配置无效: %w", err)
	}

	var nextRunTime *int64
	if req.Enabled && req.ScheduleType != "" && req.ScheduleType != "manual" {
		t := s.CalculateNextRunTime(req.ScheduleType, req.ScheduleValue)
		nextRunTime = &t
	}

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
		Viewport:      req.Viewport,
		ScheduleType:  req.ScheduleType,
		ScheduleValue: req.ScheduleValue,
		Enabled:       req.Enabled,
		NextRunTime:   nextRunTime,
		APIParams:     apiParams,
	}

	if workflow.APIKey == "" {
		if key, err := utils.GenerateWorkflowAPIKey(); err == nil {
			workflow.APIKey = key
		}
	}

	if err := db.Create(workflow).Error; err != nil {
		log.Error("创建工作流失败: %v", err)
		return nil, err
	}

	log.Info("用户 %s 创建工作流: %s (ID: %s)", userID, workflow.Name, workflow.ID)

	if workflow.Enabled && workflow.ScheduleType != "" && workflow.ScheduleType != "manual" {
		s.reloadScheduler()
	}

	return workflow, nil
}

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

func (s *WorkflowService) GetWorkflowList(userID string, query *request.WorkflowListQuery) (*response.WorkflowListResponse, error) {
	db := database.GetDB()

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	queryDB := db.Model(&models.Workflow{}).Where("user_id = ?", userID)

	if query.Keyword != "" {
		queryDB = queryDB.Where("name LIKE ? OR description LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if query.Enabled != nil {
		queryDB = queryDB.Where("enabled = ?", *query.Enabled)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, err
	}

	var workflows []models.Workflow
	offset := (query.Page - 1) * query.PageSize
	if err := queryDB.Order("created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&workflows).Error; err != nil {
		return nil, err
	}

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

func (s *WorkflowService) UpdateWorkflow(workflowID, userID string, req *request.UpdateWorkflowRequest) (*models.Workflow, error) {
	db := database.GetDB()

	workflow, err := s.GetWorkflowByID(workflowID, userID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Nodes != nil {

		edges := workflow.Edges
		if req.Edges != nil {
			edges = *req.Edges
		}
		if err := s.ValidateWorkflowConfig(*req.Nodes, edges); err != nil {
			return nil, fmt.Errorf("工作流配置无效: %w", err)
		}
		updates["nodes"] = models.WorkflowNodes(*req.Nodes)

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
	if req.Viewport != nil {
		updates["viewport"] = req.Viewport
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

	if err := db.First(workflow, "id = ?", workflowID).Error; err != nil {
		return nil, err
	}

	log.Info("用户 %s 更新工作流: %s (ID: %s)", userID, workflow.Name, workflow.ID)

	s.reloadScheduler()

	return workflow, nil
}

func (s *WorkflowService) DeleteWorkflow(workflowID, userID string) error {
	db := database.GetDB()

	workflow, err := s.GetWorkflowByID(workflowID, userID)
	if err != nil {
		return err
	}

	if err := db.Delete(workflow).Error; err != nil {
		return err
	}

	log.Info("用户 %s 删除工作流: %s (ID: %s)", userID, workflow.Name, workflow.ID)

	s.reloadScheduler()

	return nil
}

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

	s.reloadScheduler()

	return workflow, nil
}

func (s *WorkflowService) ValidateWorkflowConfig(nodes []models.WorkflowNode, edges []models.WorkflowEdge) error {
	if len(nodes) == 0 {
		return errors.New("工作流至少需要一个节点")
	}

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

	for _, edge := range edges {
		if !nodeIDs[edge.Source] {
			return fmt.Errorf("连接线的源节点不存在: %s", edge.Source)
		}
		if !nodeIDs[edge.Target] {
			return fmt.Errorf("连接线的目标节点不存在: %s", edge.Target)
		}
	}

	return nil
}

func (s *WorkflowService) ExtractExternalTriggerParams(nodes []models.WorkflowNode, edges []models.WorkflowEdge) (models.WorkflowAPIParams, error) {

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

	var externalTriggerNode *models.WorkflowNode
	for _, node := range startNodes {
		if node.Type == "external_trigger" {
			externalTriggerNode = &node
			break
		}
	}

	if externalTriggerNode == nil {
		return models.WorkflowAPIParams{}, nil
	}

	params := models.WorkflowAPIParams{}

	if paramsInterface, ok := externalTriggerNode.Config["params"]; ok {
		if paramsList, ok := paramsInterface.([]interface{}); ok {
			for _, paramInterface := range paramsList {
				if paramMap, ok := paramInterface.(map[string]interface{}); ok {
					param := models.WorkflowAPIParam{
						Key:      getStringValue(paramMap, "key"),
						Type:     getStringValue(paramMap, "type"),
						Required: getBoolValue(paramMap, "required"),
					}

					if desc := getStringValue(paramMap, "description"); desc != "" {
						param.Description = desc
					}
					if defaultVal, ok := paramMap["defaultValue"]; ok {
						param.DefaultValue = defaultVal
					}
					if example, ok := paramMap["example"]; ok {
						param.Example = example
					}

					// 文件类型特有属性
					if param.Type == "file" {
						if accept := getStringValue(paramMap, "accept"); accept != "" {
							param.Accept = accept
						}
						if maxSize, ok := paramMap["maxSize"].(float64); ok {
							param.MaxSize = int(maxSize)
						}
					}

					params = append(params, param)
				}
			}
		}
	}

	return params, nil
}

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

func (s *WorkflowService) GetWorkflowStats(workflowID, userID string) (*response.WorkflowStatsResponse, error) {
	workflow, err := s.GetWorkflowByID(workflowID, userID)
	if err != nil {
		return nil, err
	}

	db := database.GetDB()

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

func (s *WorkflowService) ToWorkflowResponse(workflow *models.Workflow) response.WorkflowResponse {
	return s.toWorkflowResponse(workflow)
}

func (s *WorkflowService) toWorkflowResponse(workflow *models.Workflow) response.WorkflowResponse {
	return response.WorkflowResponse{
		ID:              workflow.GetID(),
		UserID:          workflow.UserID,
		Name:            workflow.Name,
		Description:     workflow.Description,
		Nodes:           workflow.Nodes,
		Edges:           workflow.Edges,
		EnvVars:         workflow.EnvVars,
		Viewport:        workflow.Viewport,
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

func GetCurrentTime() int64 {
	return time.Now().Unix()
}

func (s *WorkflowService) reloadScheduler() {
	if workflowChangeCallback != nil {
		go func() {
			workflowChangeCallback()
			log.Info("工作流调度器已触发重新加载")
		}()
	}
}

func (s *WorkflowService) CalculateNextRunTime(scheduleType, scheduleValue string) int64 {
	now := time.Now()

	switch scheduleType {
	case "daily":

		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 3 {
			return now.Unix()
		}
		hour, _ := strconv.Atoi(parts[0])
		minute, _ := strconv.Atoi(parts[1])
		second, _ := strconv.Atoi(parts[2])

		next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())

		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		return next.Unix()

	case "weekly":

		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 4 {
			return now.Unix()
		}
		dayStrs := strings.Split(parts[0], ",")
		hour, _ := strconv.Atoi(parts[1])
		minute, _ := strconv.Atoi(parts[2])
		second, _ := strconv.Atoi(parts[3])

		var weekdays []int
		for _, dayStr := range dayStrs {
			day, _ := strconv.Atoi(dayStr)
			weekdays = append(weekdays, day)
		}

		currentWeekday := int(now.Weekday())
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())

		minDaysToAdd := 7
		for _, targetWeekday := range weekdays {
			daysToAdd := (targetWeekday - currentWeekday + 7) % 7

			if daysToAdd == 0 && next.After(now) {
				minDaysToAdd = 0
				break
			}

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

		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 4 {
			return now.Unix()
		}
		day, _ := strconv.Atoi(parts[0])
		hour, _ := strconv.Atoi(parts[1])
		minute, _ := strconv.Atoi(parts[2])
		second, _ := strconv.Atoi(parts[3])

		next := time.Date(now.Year(), now.Month(), day, hour, minute, second, 0, now.Location())

		if next.Before(now) {
			next = next.AddDate(0, 1, 0)
		}

		for next.Day() != day {
			next = next.AddDate(0, 1, 0)
			next = time.Date(next.Year(), next.Month(), day, hour, minute, second, 0, next.Location())
		}

		return next.Unix()

	case "hourly":

		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 2 {
			return now.Unix()
		}
		minute, _ := strconv.Atoi(parts[0])
		second, _ := strconv.Atoi(parts[1])

		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), minute, second, 0, now.Location())

		if next.Before(now) {
			next = next.Add(1 * time.Hour)
		}
		return next.Unix()

	case "interval":

		seconds, _ := strconv.ParseInt(scheduleValue, 10, 64)
		return now.Add(time.Duration(seconds) * time.Second).Unix()

	case "cron":

		return now.Add(1 * time.Minute).Unix()

	default:
		return now.Unix()
	}
}

func (s *WorkflowService) EnableWorkflowAPI(workflowID, userID string) (string, error) {
	db := database.GetDB()

	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return "", fmt.Errorf("工作流不存在")
	}

	if workflow.APIEnabled && workflow.APIKey != "" {
		return workflow.APIKey, nil
	}

	apiKey, err := utils.GenerateWorkflowAPIKey()
	if err != nil {
		return "", fmt.Errorf("生成 API Key 失败: %w", err)
	}

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

func (s *WorkflowService) DisableWorkflowAPI(workflowID, userID string) error {
	db := database.GetDB()

	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return fmt.Errorf("工作流不存在")
	}

	updates := map[string]interface{}{
		"api_enabled": false,
	}

	if err := db.Model(&workflow).Updates(updates).Error; err != nil {
		return fmt.Errorf("禁用 API 失败: %w", err)
	}

	log.Info("工作流 API 已禁用: WorkflowID=%s", workflowID)
	return nil
}

func (s *WorkflowService) RegenerateAPIKey(workflowID, userID string) (string, error) {
	db := database.GetDB()

	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return "", fmt.Errorf("工作流不存在")
	}

	if !workflow.APIEnabled {
		return "", fmt.Errorf("工作流 API 未启用")
	}

	apiKey, err := utils.GenerateWorkflowAPIKey()
	if err != nil {
		return "", fmt.Errorf("生成 API Key 失败: %w", err)
	}

	if err := db.Model(&workflow).Update("api_key", apiKey).Error; err != nil {
		return "", fmt.Errorf("更新 API Key 失败: %w", err)
	}

	log.Info("工作流 API Key 已重新生成: WorkflowID=%s, APIKey=%s", workflowID, apiKey)
	return apiKey, nil
}

func (s *WorkflowService) UpdateAPIParams(workflowID, userID string, params models.WorkflowAPIParams) error {
	db := database.GetDB()

	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return fmt.Errorf("工作流不存在")
	}

	if err := db.Model(&workflow).Update("api_params", params).Error; err != nil {
		return fmt.Errorf("更新 API 参数配置失败: %w", err)
	}

	log.Info("工作流 API 参数已更新: WorkflowID=%s, ParamsCount=%d", workflowID, len(params))
	return nil
}

func (s *WorkflowService) UpdateAPITimeout(workflowID, userID string, timeout int) error {
	db := database.GetDB()

	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return fmt.Errorf("工作流不存在")
	}

	if err := db.Model(&workflow).Update("api_timeout", timeout).Error; err != nil {
		return fmt.Errorf("更新 API 超时时间失败: %w", err)
	}

	log.Info("工作流 API 超时时间已更新: WorkflowID=%s, Timeout=%d", workflowID, timeout)
	return nil
}

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

func (s *WorkflowService) ValidateAPIParams(workflow *models.Workflow, userParams map[string]interface{}) error {
	if len(workflow.APIParams) == 0 {
		return nil
	}

	for _, param := range workflow.APIParams {
		value, exists := userParams[param.Key]

		if !exists {
			// 如果参数不存在，尝试使用默认值
			if param.DefaultValue != nil {
				userParams[param.Key] = param.DefaultValue
				continue
			}

			// 如果没有默认值且是必填，报错
			if param.Required {
				return fmt.Errorf("缺少必填参数: %s", param.Key)
			}
			continue
		}

		if err := validateParamType(value, param.Type); err != nil {
			return fmt.Errorf("参数 %s 类型错误: %w", param.Key, err)
		}
	}

	return nil
}

func (s *WorkflowService) ApplyAPIParams(workflow *models.Workflow, userParams map[string]interface{}) error {
	if len(workflow.APIParams) == 0 {
		return nil
	}

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

	for _, param := range workflow.APIParams {
		value, exists := userParams[param.Key]

		if !exists {
			if param.Required {
				return fmt.Errorf("缺少必填参数: %s", param.Key)
			}

			if param.DefaultValue != nil {
				value = param.DefaultValue
			} else {
				continue
			}
		}

		if err := validateParamType(value, param.Type); err != nil {
			return fmt.Errorf("参数 %s 类型错误: %w", param.Key, err)
		}

		if err := utils.SetValueByPath(workflowData, param.MappingPath, value); err != nil {
			return fmt.Errorf("应用参数 %s 失败: %w", param.Key, err)
		}

		log.Info("应用参数: %s = %v -> %s", param.Key, value, param.MappingPath)
	}

	if nodesInterface, ok := workflowData["nodes"].([]interface{}); ok {
		for i, nodeInterface := range nodesInterface {
			if nodeMap, ok := nodeInterface.(map[string]interface{}); ok {
				workflow.Nodes[i].Config = nodeMap["config"].(map[string]interface{})
			}
		}
	}

	return nil
}

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

func (s *WorkflowService) ExecuteWorkflowSync(executionID, userID string, timeoutSeconds int, externalParams map[string]interface{}) (map[string]interface{}, error) {
	engineSvc := NewEngineService()
	executionSvc := NewExecutionService()

	done := make(chan error, 1)

	go func() {
		err := engineSvc.ExecuteWorkflow(executionID, nil, externalParams)
		done <- err
	}()

	timeout := time.Duration(timeoutSeconds) * time.Second
	select {
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("执行失败: %w", err)
		}

		execution, err := executionSvc.GetExecutionByID(executionID, userID)
		if err != nil {
			return nil, fmt.Errorf("查询执行结果失败: %w", err)
		}

		var finalOutput map[string]interface{}
		if len(execution.NodeLogs) > 0 {
			lastNode := execution.NodeLogs[len(execution.NodeLogs)-1]
			finalOutput = lastNode.Output
		}

		if finalOutput == nil {
			finalOutput = make(map[string]interface{})
		}

		return finalOutput, nil

	case <-time.After(timeout):
		return nil, fmt.Errorf("执行超时 (%d 秒)", timeoutSeconds)
	}
}

func (s *WorkflowService) IncrementAPICallCount(workflowID string) error {
	db := database.GetDB()
	now := time.Now().Unix()

	updates := map[string]interface{}{
		"api_call_count":     gorm.Expr("api_call_count + 1"),
		"api_last_called_at": now,
	}

	if err := db.Model(&models.Workflow{}).Where("id = ?", workflowID).Updates(updates).Error; err != nil {
		log.Error("更新 API 调用统计失败: %v", err)
		return err
	}

	return nil
}

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
