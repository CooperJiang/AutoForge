package workflow

import (
	"auto-forge/internal/models"
	"auto-forge/pkg/database"
	log "auto-forge/pkg/logger"
	"auto-forge/pkg/utools"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// EngineService 工作流执行引擎
type EngineService struct {
	executionService *ExecutionService
}

// NewEngineService 创建引擎服务实例
func NewEngineService() *EngineService {
	return &EngineService{
		executionService: NewExecutionService(),
	}
}

// ExecuteWorkflow 执行工作流
func (s *EngineService) ExecuteWorkflow(executionID string, envVars map[string]string) error {
	db := database.GetDB()

	// 获取执行记录
	var execution models.WorkflowExecution
	if err := db.First(&execution, "id = ?", executionID).Error; err != nil {
		return fmt.Errorf("执行记录不存在: %w", err)
	}

	// 获取工作流
	var workflow models.Workflow
	if err := db.First(&workflow, "id = ?", execution.WorkflowID).Error; err != nil {
		return fmt.Errorf("工作流不存在: %w", err)
	}

	// 更新执行状态为运行中
	if err := s.executionService.UpdateExecutionStatus(executionID, models.ExecutionStatusRunning, ""); err != nil {
		return err
	}

	// 构建环境变量映射
	envMap := s.buildEnvMap(workflow.EnvVars, envVars)

	// 执行工作流
	success := true
	var execError error

	// 拓扑排序节点（简化版：按节点ID顺序执行）
	sortedNodes, err := s.topologicalSort(workflow.Nodes, workflow.Edges)
	if err != nil {
		s.executionService.UpdateExecutionStatus(executionID, models.ExecutionStatusFailed, err.Error())
		return err
	}

	// 按顺序执行每个节点
	nodeOutputs := make(map[string]map[string]interface{})
	shouldSkipNext := false // 条件判断控制是否跳过后续节点

	for _, node := range sortedNodes {
		// 检查执行是否被取消
		if err := db.First(&execution, "id = ?", executionID).Error; err == nil {
			if execution.Status == models.ExecutionStatusCancelled {
				log.Info("工作流执行被取消: ExecutionID=%s", executionID)
				break
			}
		}

		// 如果上一个条件节点要求跳过，则跳过当前节点
		if shouldSkipNext {
			nodeLog := &models.NodeExecutionLog{
				NodeID:   node.ID,
				NodeType: node.Type,
				NodeName: s.getNodeName(node),
				Status:   "skipped",
				ToolCode: node.ToolCode,
			}
			startTime := time.Now().Unix()
			nodeLog.StartTime = &startTime
			nodeLog.EndTime = &startTime
			nodeLog.DurationMs = 0
			nodeLog.Output = map[string]interface{}{
				"reason": "条件判断为false，跳过执行",
			}

			if err := s.executionService.AddNodeLog(executionID, *nodeLog); err != nil {
				log.Error("添加节点日志失败: %v", err)
			}
			continue
		}

		// 执行节点
		nodeLog, output, err := s.executeNode(node, envMap, nodeOutputs)
		if err != nil {
			success = false
			execError = err
			nodeLog.Status = models.ExecutionStatusFailed
			nodeLog.Error = err.Error()
		} else {
			nodeLog.Status = models.ExecutionStatusSuccess
			nodeOutputs[node.ID] = output

			// 如果是条件节点，检查条件结果
			if node.Type == "condition" {
				if result, ok := output["result"].(bool); ok && !result {
					shouldSkipNext = true
					log.Info("条件判断为false，将跳过后续节点")
				}
			}
		}

		// 添加节点日志
		if err := s.executionService.AddNodeLog(executionID, *nodeLog); err != nil {
			log.Error("添加节点日志失败: %v", err)
		}

		// 如果节点执行失败，停止后续执行
		if !success {
			break
		}
	}

	// 更新执行结果
	var finalStatus string
	var finalError string

	if success {
		finalStatus = models.ExecutionStatusSuccess
	} else {
		finalStatus = models.ExecutionStatusFailed
		if execError != nil {
			finalError = execError.Error()
		}
	}

	if err := s.executionService.UpdateExecutionStatus(executionID, finalStatus, finalError); err != nil {
		log.Error("更新执行状态失败: %v", err)
	}

	// 更新工作流统计
	if err := s.executionService.UpdateWorkflowStats(workflow.GetID(), success); err != nil {
		log.Error("更新工作流统计失败: %v", err)
	}

	return execError
}

// executeNode 执行单个节点
func (s *EngineService) executeNode(
	node models.WorkflowNode,
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
) (*models.NodeExecutionLog, map[string]interface{}, error) {
	startTime := time.Now().Unix()

	nodeLog := &models.NodeExecutionLog{
		NodeID:     node.ID,
		NodeType:   node.Type,
		NodeName:   s.getNodeName(node),
		Status:     models.ExecutionStatusRunning,
		StartTime:  &startTime,
		RetryCount: 0,
		ToolCode:   node.ToolCode, // 保存工具代码
	}

	// 根据节点类型执行
	var output map[string]interface{}
	var err error

	switch node.Type {
	case "tool":
		output, err = s.executeToolNode(node, envMap, nodeOutputs)
	case "trigger":
		output, err = s.executeTriggerNode(node)
	case "condition":
		output, err = s.executeConditionNode(node, nodeOutputs)
	case "delay":
		output, err = s.executeDelayNode(node)
	case "switch":
		output, err = s.executeSwitchNode(node, nodeOutputs)
	default:
		err = fmt.Errorf("不支持的节点类型: %s", node.Type)
	}

	endTime := time.Now().Unix()
	nodeLog.EndTime = &endTime
	nodeLog.DurationMs = (endTime - startTime) * 1000
	nodeLog.Output = output

	return nodeLog, output, err
}

// executeToolNode 执行工具节点
func (s *EngineService) executeToolNode(
	node models.WorkflowNode,
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
) (map[string]interface{}, error) {
	// 获取工具代码 - 优先从顶级字段读取，兼容旧数据从 Data 中读取
	toolCode := node.ToolCode
	if toolCode == "" {
		if tc, ok := node.Data["tool_code"].(string); ok {
			toolCode = tc
		}
	}
	if toolCode == "" {
		return nil, errors.New("工具代码未配置")
	}

	// 获取工具配置 - 优先从顶级字段读取，兼容旧数据从 Data 中读取
	config := node.Config
	if config == nil || len(config) == 0 {
		if c, ok := node.Data["config"].(map[string]interface{}); ok {
			config = c
		}
	}
	if config == nil {
		return nil, errors.New("工具配置格式错误")
	}

	// 替换环境变量和节点输出引用
	config = s.replaceVariables(config, envMap, nodeOutputs)

	// 获取工具实例
	tool, err := utools.Get(toolCode)
	if err != nil {
		return nil, fmt.Errorf("工具不存在: %s, %w", toolCode, err)
	}

	// 执行工具
	ctx := &utools.ExecutionContext{
		Context:   context.Background(),
		TaskID:    "",
		UserID:    "",
		Variables: make(map[string]interface{}),
		Metadata:  make(map[string]interface{}),
	}

	result, err := tool.Execute(ctx, config)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return result.Output, fmt.Errorf("工具执行失败: %s", result.Message)
	}

	return result.Output, nil
}

// executeTriggerNode 执行触发器节点
func (s *EngineService) executeTriggerNode(node models.WorkflowNode) (map[string]interface{}, error) {
	// 触发器节点不需要执行逻辑，直接返回成功
	return map[string]interface{}{
		"triggered": true,
	}, nil
}

// executeConditionNode 执行条件节点
func (s *EngineService) executeConditionNode(
	node models.WorkflowNode,
	nodeOutputs map[string]map[string]interface{},
) (map[string]interface{}, error) {
	config := node.Config
	if config == nil {
		return nil, errors.New("条件配置为空")
	}

	conditionType, _ := config["conditionType"].(string)

	// 根据条件类型执行
	switch conditionType {
	case "simple":
		return s.evaluateSimpleCondition(config, nodeOutputs)
	case "expression":
		return s.evaluateExpressionCondition(config, nodeOutputs)
	default:
		return map[string]interface{}{
			"result": true,
			"message": "未指定条件类型，默认通过",
		}, nil
	}
}

// evaluateSimpleCondition 执行简单条件判断
func (s *EngineService) evaluateSimpleCondition(
	config map[string]interface{},
	nodeOutputs map[string]map[string]interface{},
) (map[string]interface{}, error) {
	field, _ := config["field"].(string)
	operator, _ := config["operator"].(string)
	expectedValue := config["value"]

	if field == "" || operator == "" {
		return nil, errors.New("条件配置不完整：缺少field或operator")
	}

	// 从上一个节点的输出中获取字段值
	// 支持格式：node_xxx.field_name 或直接 field_name（从最后一个节点获取）
	var actualValue interface{}
	var sourceNode string

	if strings.Contains(field, ".") {
		// 格式：node_xxx.field_name
		parts := strings.SplitN(field, ".", 2)
		nodeID := parts[0]
		fieldName := parts[1]

		if nodeOutput, ok := nodeOutputs[nodeID]; ok {
			actualValue = s.getNestedField(nodeOutput, fieldName)
			sourceNode = nodeID
		} else {
			return nil, fmt.Errorf("找不到节点输出: %s", nodeID)
		}
	} else {
		// 从最后一个节点获取
		for nodeID, output := range nodeOutputs {
			if val := s.getNestedField(output, field); val != nil {
				actualValue = val
				sourceNode = nodeID
			}
		}
		if actualValue == nil {
			return nil, fmt.Errorf("找不到字段: %s", field)
		}
	}

	// 执行比较
	result := s.compareValues(actualValue, operator, expectedValue)

	return map[string]interface{}{
		"result":         result,
		"source_node":    sourceNode,
		"field":          field,
		"operator":       operator,
		"actual_value":   actualValue,
		"expected_value": expectedValue,
		"message":        fmt.Sprintf("%v %s %v = %v", actualValue, operator, expectedValue, result),
	}, nil
}

// evaluateExpressionCondition 执行表达式条件判断
func (s *EngineService) evaluateExpressionCondition(
	config map[string]interface{},
	nodeOutputs map[string]map[string]interface{},
) (map[string]interface{}, error) {
	// TODO: 实现复杂表达式判断
	return map[string]interface{}{
		"result":  true,
		"message": "表达式条件判断暂未实现",
	}, nil
}

// getNestedField 获取嵌套字段值（支持 a.b.c 格式）
func (s *EngineService) getNestedField(data map[string]interface{}, field string) interface{} {
	parts := strings.Split(field, ".")
	var current interface{} = data

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[part]
		default:
			return nil
		}
	}

	return current
}

// compareValues 比较两个值
func (s *EngineService) compareValues(actual interface{}, operator string, expected interface{}) bool {
	// 转换为字符串进行比较
	actualStr := fmt.Sprintf("%v", actual)
	expectedStr := fmt.Sprintf("%v", expected)

	// 尝试转换为数字进行比较
	actualNum, actualIsNum := s.toFloat64(actual)
	expectedNum, expectedIsNum := s.toFloat64(expected)

	switch operator {
	case "equals", "==", "=":
		if actualIsNum && expectedIsNum {
			return actualNum == expectedNum
		}
		return actualStr == expectedStr

	case "not_equals", "!=", "<>":
		if actualIsNum && expectedIsNum {
			return actualNum != expectedNum
		}
		return actualStr != expectedStr

	case "greater_than", ">":
		if actualIsNum && expectedIsNum {
			return actualNum > expectedNum
		}
		return actualStr > expectedStr

	case "greater_than_or_equal", ">=":
		if actualIsNum && expectedIsNum {
			return actualNum >= expectedNum
		}
		return actualStr >= expectedStr

	case "less_than", "<":
		if actualIsNum && expectedIsNum {
			return actualNum < expectedNum
		}
		return actualStr < expectedStr

	case "less_than_or_equal", "<=":
		if actualIsNum && expectedIsNum {
			return actualNum <= expectedNum
		}
		return actualStr <= expectedStr

	case "contains":
		return strings.Contains(actualStr, expectedStr)

	case "not_contains":
		return !strings.Contains(actualStr, expectedStr)

	case "starts_with":
		return strings.HasPrefix(actualStr, expectedStr)

	case "ends_with":
		return strings.HasSuffix(actualStr, expectedStr)

	default:
		return false
	}
}

// toFloat64 尝试将值转换为 float64
func (s *EngineService) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case string:
		if f, err := fmt.Sscanf(v, "%f", new(float64)); err == nil && f == 1 {
			var result float64
			fmt.Sscanf(v, "%f", &result)
			return result, true
		}
	}
	return 0, false
}

// executeDelayNode 执行延迟节点
func (s *EngineService) executeDelayNode(node models.WorkflowNode) (map[string]interface{}, error) {
	// 获取延迟时间（秒） - 优先从 Config 读取，兼容从 Data 读取
	var delaySeconds float64 = 1
	if node.Config != nil {
		if d, ok := node.Config["delay"].(float64); ok {
			delaySeconds = d
		} else if d, ok := node.Config["delaySeconds"].(float64); ok {
			delaySeconds = d
		}
	}
	if delaySeconds == 1 {
		if d, ok := node.Data["delay"].(float64); ok {
			delaySeconds = d
		}
	}

	time.Sleep(time.Duration(delaySeconds) * time.Second)

	return map[string]interface{}{
		"delayed_seconds": delaySeconds,
	}, nil
}

// executeSwitchNode 执行开关节点
func (s *EngineService) executeSwitchNode(
	node models.WorkflowNode,
	nodeOutputs map[string]map[string]interface{},
) (map[string]interface{}, error) {
	// 简化版：返回默认分支
	// TODO: 实现开关逻辑
	return map[string]interface{}{
		"branch": "default",
	}, nil
}

// topologicalSort 拓扑排序（简化版）
func (s *EngineService) topologicalSort(nodes []models.WorkflowNode, edges []models.WorkflowEdge) ([]models.WorkflowNode, error) {
	// 构建入度映射和邻接表
	inDegree := make(map[string]int)
	adjList := make(map[string][]string)
	nodeMap := make(map[string]models.WorkflowNode)

	for _, node := range nodes {
		inDegree[node.ID] = 0
		nodeMap[node.ID] = node
	}

	for _, edge := range edges {
		inDegree[edge.Target]++
		adjList[edge.Source] = append(adjList[edge.Source], edge.Target)
	}

	// 找到所有入度为0的节点（起始节点）
	queue := []string{}
	for nodeID, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, nodeID)
		}
	}

	// BFS 拓扑排序
	result := []models.WorkflowNode{}
	for len(queue) > 0 {
		// 取出队首节点
		nodeID := queue[0]
		queue = queue[1:]

		result = append(result, nodeMap[nodeID])

		// 减少相邻节点的入度
		for _, nextID := range adjList[nodeID] {
			inDegree[nextID]--
			if inDegree[nextID] == 0 {
				queue = append(queue, nextID)
			}
		}
	}

	// 检查是否存在循环依赖
	if len(result) != len(nodes) {
		return nil, errors.New("工作流存在循环依赖")
	}

	return result, nil
}

// buildEnvMap 构建环境变量映射
func (s *EngineService) buildEnvMap(envVars []models.WorkflowEnvVar, runtimeEnvVars map[string]string) map[string]string {
	envMap := make(map[string]string)

	// 添加工作流环境变量
	for _, envVar := range envVars {
		envMap[envVar.Key] = envVar.Value
	}

	// 运行时环境变量覆盖
	for key, value := range runtimeEnvVars {
		envMap[key] = value
	}

	return envMap
}

// replaceVariables 替换变量引用
func (s *EngineService) replaceVariables(
	config map[string]interface{},
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range config {
		switch v := value.(type) {
		case string:
			result[key] = s.replaceStringVariables(v, envMap, nodeOutputs)
		case map[string]interface{}:
			result[key] = s.replaceVariables(v, envMap, nodeOutputs)
		default:
			result[key] = value
		}
	}

	return result
}

// replaceStringVariables 替换字符串中的变量
func (s *EngineService) replaceStringVariables(
	str string,
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
) string {
	// 替换环境变量 {{env.VAR_NAME}}
	for key, value := range envMap {
		str = strings.ReplaceAll(str, fmt.Sprintf("{{env.%s}}", key), value)
	}

	// 替换节点输出 {{nodes.NODE_ID.field}} 或 {{nodes.NODE_ID.output.field}}
	for nodeID, output := range nodeOutputs {
		// 支持 {{nodes.NODE_ID.field}} 格式
		for field, value := range output {
			placeholder := fmt.Sprintf("{{nodes.%s.%s}}", nodeID, field)
			str = strings.ReplaceAll(str, placeholder, fmt.Sprintf("%v", value))
		}

		// 支持 {{nodes.NODE_ID.output.field}} 格式（嵌套output对象）
		if outputObj, ok := output["output"].(map[string]interface{}); ok {
			for field, value := range outputObj {
				placeholder := fmt.Sprintf("{{nodes.%s.output.%s}}", nodeID, field)
				str = strings.ReplaceAll(str, placeholder, fmt.Sprintf("%v", value))
			}
		}
	}

	return str
}

// getNodeName 获取节点名称
func (s *EngineService) getNodeName(node models.WorkflowNode) string {
	// 优先使用顶层 Name 字段
	if node.Name != "" {
		return node.Name
	}

	// 兼容旧数据：从 Data 中读取
	if name, ok := node.Data["label"].(string); ok && name != "" {
		return name
	}
	if name, ok := node.Data["name"].(string); ok && name != "" {
		return name
	}

	// 如果是工具节点，显示工具代码
	if node.Type == "tool" && node.ToolCode != "" {
		return node.ToolCode
	}

	// 最后才返回节点类型
	return node.Type
}

// MarshalJSON 自定义JSON序列化（辅助函数）
func marshalJSON(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
