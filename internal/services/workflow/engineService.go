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
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type EngineService struct {
	executionService *ExecutionService
}

func NewEngineService() *EngineService {
	return &EngineService{
		executionService: NewExecutionService(),
	}
}

func (s *EngineService) ExecuteWorkflow(executionID string, envVars map[string]string, externalParams map[string]interface{}) error {
	db := database.GetDB()

	var execution models.WorkflowExecution
	if err := db.First(&execution, "id = ?", executionID).Error; err != nil {
		return fmt.Errorf("执行记录不存在: %w", err)
	}

	var workflow models.Workflow
	if err := db.First(&workflow, "id = ?", execution.WorkflowID).Error; err != nil {
		return fmt.Errorf("工作流不存在: %w", err)
	}

	if err := s.executionService.UpdateExecutionStatus(executionID, models.ExecutionStatusRunning, ""); err != nil {
		return err
	}

	envMap := s.buildEnvMap(workflow.EnvVars, envVars)

	if externalParams != nil {
		for key, value := range externalParams {
			envMap["external."+key] = fmt.Sprintf("%v", value)
		}
	}

	success := true
	var execError error

	sortedNodes, err := s.topologicalSort(workflow.Nodes, workflow.Edges)
	if err != nil {
		s.executionService.UpdateExecutionStatus(executionID, models.ExecutionStatusFailed, err.Error())
		return err
	}

	nodeOutputs := make(map[string]map[string]interface{})
	shouldSkipNext := false

	for _, node := range sortedNodes {

		if err := db.First(&execution, "id = ?", executionID).Error; err == nil {
			if execution.Status == models.ExecutionStatusCancelled {
				log.Info("工作流执行被取消: ExecutionID=%s", executionID)
				break
			}
		}

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

		nodeLog, output, err := s.executeNode(executionID, node, envMap, nodeOutputs, externalParams)
		if err != nil {
			success = false
			execError = err
			nodeLog.Status = models.ExecutionStatusFailed
			nodeLog.Error = err.Error()
		} else {
			nodeLog.Status = models.ExecutionStatusSuccess
			nodeOutputs[node.ID] = output

			if node.Type == "condition" {
				if result, ok := output["result"].(bool); ok && !result {
					shouldSkipNext = true
					log.Info("条件判断为false，将跳过后续节点")
				}
			}
		}

		if err := s.executionService.AddNodeLog(executionID, *nodeLog); err != nil {
			log.Error("添加节点日志失败: %v", err)
		}

		if !success {
			break
		}
	}

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

	if err := s.executionService.UpdateWorkflowStats(workflow.GetID(), success); err != nil {
		log.Error("更新工作流统计失败: %v", err)
	}

	// 清理临时文件（如果有上传的文件）
	s.cleanupExecutionFiles(executionID)

	return execError
}

func (s *EngineService) executeNode(
	executionID string,
	node models.WorkflowNode,
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
	externalParams map[string]interface{},
) (*models.NodeExecutionLog, map[string]interface{}, error) {
	startTime := time.Now().Unix()

	inputData := make(map[string]interface{})

	if node.Config != nil && len(node.Config) > 0 {
		inputData["config"] = node.Config
	}

	externalParamsForLog := make(map[string]interface{})
	if externalParams != nil {
		for key, value := range externalParams {
			externalParamsForLog[key] = value
		}
	}
	if len(externalParamsForLog) > 0 {
		inputData["external_params"] = externalParamsForLog
	}

	nodeLog := &models.NodeExecutionLog{
		NodeID:     node.ID,
		NodeType:   node.Type,
		NodeName:   s.getNodeName(node),
		Status:     models.ExecutionStatusRunning,
		StartTime:  &startTime,
		RetryCount: 0,
		ToolCode:   node.ToolCode,
		Input:      inputData,
	}

	var output map[string]interface{}
	var outputRender *models.OutputRenderConfig
	var err error

	if node.Type == "tool" {
		replacedConfig := s.replaceVariables(node.Config, envMap, nodeOutputs, externalParams)
		if len(replacedConfig) > 0 {
			inputData["resolved_config"] = replacedConfig
		}
	}

	if err := s.executionService.AddNodeLog(executionID, *nodeLog); err != nil {
		log.Error("添加节点开始日志失败: %v", err)
	}

	switch node.Type {
	case "tool":
		output, outputRender, err = s.executeToolNode(node, envMap, nodeOutputs, externalParams)
	case "trigger", "external_trigger":
		output, err = s.executeTriggerNode(node)
	case "condition":
		output, err = s.executeConditionNode(node, nodeOutputs)
	case "delay":
		output, err = s.executeDelayNode(node, envMap, nodeOutputs, externalParams)
	case "switch":
		output, err = s.executeSwitchNode(node, envMap, nodeOutputs, externalParams)
	default:
		err = fmt.Errorf("不支持的节点类型: %s", node.Type)
	}

	endTime := time.Now().Unix()
	nodeLog.EndTime = &endTime
	nodeLog.DurationMs = (endTime - startTime) * 1000
	nodeLog.Output = output
	nodeLog.OutputRender = outputRender

	return nodeLog, output, err
}

func (s *EngineService) executeToolNode(
	node models.WorkflowNode,
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
	externalParams map[string]interface{},
) (map[string]interface{}, *models.OutputRenderConfig, error) {
	toolCode := node.ToolCode
	if toolCode == "" {
		if tc, ok := node.Data["tool_code"].(string); ok {
			toolCode = tc
		}
	}
	if toolCode == "" {
		return nil, nil, errors.New("工具代码未配置")
	}

	config := node.Config
	if config == nil || len(config) == 0 {
		if c, ok := node.Data["config"].(map[string]interface{}); ok {
			config = c
		}
	}
	if config == nil {
		return nil, nil, errors.New("工具配置格式错误")
	}

	config = s.replaceVariables(config, envMap, nodeOutputs, externalParams)

	tool, err := utools.Get(toolCode)
	if err != nil {
		return nil, nil, fmt.Errorf("工具不存在: %s, %w", toolCode, err)
	}

	ctx := &utools.ExecutionContext{
		Context:   context.Background(),
		TaskID:    "",
		UserID:    "",
		Variables: make(map[string]interface{}),
		Metadata:  make(map[string]interface{}),
	}

	nodeVariables := make(map[string]interface{}, len(nodeOutputs))
	for key, val := range nodeOutputs {
		nodeVariables[key] = val
	}
	ctx.Variables["nodes"] = nodeVariables

	envVars := make(map[string]interface{}, len(envMap))
	for key, val := range envMap {
		envVars[key] = val
	}
	ctx.Variables["env"] = envVars

	// 将外部参数作为对象传递，保留原始类型
	if externalParams != nil {
		ctx.Variables["external"] = externalParams
	}

	ctx.Metadata["current"] = map[string]interface{}{
		"nodeId":   node.ID,
		"nodeType": node.Type,
		"nodeName": s.getNodeName(node),
	}
	ctx.Variables["current"] = ctx.Metadata["current"]

	result, err := tool.Execute(ctx, config)
	if err != nil {
		return nil, nil, err
	}

	if !result.Success {
		return result.Output, nil, fmt.Errorf("工具执行失败: %s", result.Message)
	}

	output := result.Output
	if output == nil {
		output = make(map[string]interface{})
	}

	var outputRender *models.OutputRenderConfig
	if result.OutputRender != nil {
		fields := make(map[string]models.FieldRender)
		for k, v := range result.OutputRender.Fields {
			fields[k] = models.FieldRender{
				Type:    v.Type,
				Label:   v.Label,
				Display: v.Display,
			}
		}
		outputRender = &models.OutputRenderConfig{
			Type:    result.OutputRender.Type,
			Primary: result.OutputRender.Primary,
			Fields:  fields,
		}
	}

	return output, outputRender, nil
}

func (s *EngineService) executeTriggerNode(node models.WorkflowNode) (map[string]interface{}, error) {

	return map[string]interface{}{
		"triggered": true,
	}, nil
}

func (s *EngineService) executeConditionNode(
	node models.WorkflowNode,
	nodeOutputs map[string]map[string]interface{},
) (map[string]interface{}, error) {
	config := node.Config
	if config == nil {
		return nil, errors.New("条件配置为空")
	}

	conditionType, _ := config["conditionType"].(string)

	switch conditionType {
	case "simple":
		return s.evaluateSimpleCondition(config, nodeOutputs)
	case "expression":
		return s.evaluateExpressionCondition(config, nodeOutputs)
	default:
		return map[string]interface{}{
			"result":  true,
			"message": "未指定条件类型，默认通过",
		}, nil
	}
}

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

	var actualValue interface{}
	var sourceNode string

	if strings.Contains(field, ".") {

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

func (s *EngineService) evaluateExpressionCondition(
	config map[string]interface{},
	nodeOutputs map[string]map[string]interface{},
) (map[string]interface{}, error) {

	return map[string]interface{}{
		"result":  true,
		"message": "表达式条件判断暂未实现",
	}, nil
}

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

func (s *EngineService) compareValues(actual interface{}, operator string, expected interface{}) bool {

	actualStr := fmt.Sprintf("%v", actual)
	expectedStr := fmt.Sprintf("%v", expected)

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

func (s *EngineService) executeDelayNode(
	node models.WorkflowNode,
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
	externalParams map[string]interface{},
) (map[string]interface{}, error) {

	config := s.replaceVariables(node.Config, envMap, nodeOutputs, externalParams)

	var duration float64 = 5
	unit := "seconds"

	if durationVal, ok := config["duration"]; ok {
		switch v := durationVal.(type) {
		case float64:
			duration = v
		case string:

			if f, err := strconv.ParseFloat(v, 64); err == nil {
				duration = f
			}
		}
	}

	if unitVal, ok := config["unit"].(string); ok {
		unit = unitVal
	}

	var delaySeconds float64
	switch unit {
	case "minutes":
		delaySeconds = duration * 60
	case "hours":
		delaySeconds = duration * 3600
	default:
		delaySeconds = duration
	}

	time.Sleep(time.Duration(delaySeconds) * time.Second)

	return map[string]interface{}{
		"delayed_seconds": delaySeconds,
		"duration":        duration,
		"unit":            unit,
	}, nil
}

func (s *EngineService) executeSwitchNode(
	node models.WorkflowNode,
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
	externalParams map[string]interface{},
) (map[string]interface{}, error) {

	config := s.replaceVariables(node.Config, envMap, nodeOutputs, externalParams)

	fieldValue := ""
	if field, ok := config["field"].(string); ok {
		fieldValue = field
	}

	cases := []map[string]interface{}{}
	if casesVal, ok := config["cases"].([]interface{}); ok {
		for _, c := range casesVal {
			if caseMap, ok := c.(map[string]interface{}); ok {
				cases = append(cases, caseMap)
			}
		}
	}

	matchedBranch := "default"
	matchedValue := ""

	for i, caseItem := range cases {
		caseValue := ""
		if v, ok := caseItem["value"].(string); ok {
			caseValue = v
		}

		if fieldValue == caseValue {
			matchedBranch = fmt.Sprintf("case_%d", i)
			if label, ok := caseItem["label"].(string); ok && label != "" {
				matchedBranch = label
			}
			matchedValue = caseValue
			break
		}
	}

	return map[string]interface{}{
		"branch":        matchedBranch,
		"field_value":   fieldValue,
		"matched_value": matchedValue,
	}, nil
}

func (s *EngineService) topologicalSort(nodes []models.WorkflowNode, edges []models.WorkflowEdge) ([]models.WorkflowNode, error) {

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

	queue := []string{}
	for nodeID, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, nodeID)
		}
	}

	result := []models.WorkflowNode{}
	for len(queue) > 0 {

		nodeID := queue[0]
		queue = queue[1:]

		result = append(result, nodeMap[nodeID])

		for _, nextID := range adjList[nodeID] {
			inDegree[nextID]--
			if inDegree[nextID] == 0 {
				queue = append(queue, nextID)
			}
		}
	}

	if len(result) != len(nodes) {
		return nil, errors.New("工作流存在循环依赖")
	}

	return result, nil
}

func (s *EngineService) buildEnvMap(envVars []models.WorkflowEnvVar, runtimeEnvVars map[string]string) map[string]string {
	envMap := make(map[string]string)

	for _, envVar := range envVars {
		envMap[envVar.Key] = envVar.Value
	}

	for key, value := range runtimeEnvVars {
		envMap[key] = value
	}

	return envMap
}

func (s *EngineService) replaceVariables(
	config map[string]interface{},
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
	externalParams map[string]interface{},
) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range config {
		switch v := value.(type) {
		case string:
			if key == "data_source" {
				result[key] = v
			} else {
				// 检查是否是完整的变量引用（如 "{{external.image}}"）
				if s.isCompleteVariableRef(v) {
					// 直接解析变量，返回对象（不转字符串）
					resolved := s.resolveVariable(v, envMap, nodeOutputs, externalParams)
					result[key] = resolved
				} else {
					// 普通字符串，进行模板替换
					result[key] = s.replaceStringVariables(v, envMap, nodeOutputs, externalParams)
				}
			}
		case map[string]interface{}:
			result[key] = s.replaceVariables(v, envMap, nodeOutputs, externalParams)
		case []interface{}:

			result[key] = s.replaceArray(v, envMap, nodeOutputs, externalParams)
		default:
			result[key] = value
		}
	}

	return result
}

func (s *EngineService) replaceArray(
	arr []interface{},
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
	externalParams map[string]interface{},
) []interface{} {
	result := make([]interface{}, len(arr))
	for i, item := range arr {
		switch v := item.(type) {
		case string:
			result[i] = s.replaceStringVariables(v, envMap, nodeOutputs, externalParams)
		case map[string]interface{}:
			result[i] = s.replaceVariables(v, envMap, nodeOutputs, externalParams)
		case []interface{}:
			result[i] = s.replaceArray(v, envMap, nodeOutputs, externalParams)
		default:
			result[i] = item
		}
	}
	return result
}

func (s *EngineService) replaceStringVariables(
	str string,
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
	externalParams map[string]interface{},
) string {
	// 处理环境变量 {{env.xxx}}
	for key, value := range envMap {
		str = strings.ReplaceAll(str, fmt.Sprintf("{{env.%s}}", key), value)
	}

	// 处理外部参数 {{external.xxx}}，支持嵌套路径
	if externalParams != nil {
		re := regexp.MustCompile(`\{\{external\.([^}]+)\}\}`)
		str = re.ReplaceAllStringFunc(str, func(match string) string {
			path := strings.TrimPrefix(match, "{{external.")
			path = strings.TrimSuffix(path, "}}")

			value := s.getNestedValue(externalParams, path)
			if value == nil {
				return match
			}

			return fmt.Sprintf("%v", value)
		})
	}

	// 处理节点输出 {{nodes.xxx.yyy}}
	re := regexp.MustCompile(`\{\{nodes\.([^}]+)\}\}`)
	str = re.ReplaceAllStringFunc(str, func(match string) string {
		path := strings.TrimPrefix(match, "{{nodes.")
		path = strings.TrimSuffix(path, "}}")

		parts := strings.SplitN(path, ".", 2)
		if len(parts) < 2 {
			return match
		}

		nodeID := parts[0]
		fieldPath := parts[1]

		output, ok := nodeOutputs[nodeID]
		if !ok {
			return match
		}

		value := s.getNestedValue(output, fieldPath)
		if value == nil {
			return match
		}

		return fmt.Sprintf("%v", value)
	})

	return str
}

func (s *EngineService) getNestedValue(data interface{}, path string) interface{} {
	parts := strings.Split(path, ".")

	current := data
	for _, part := range parts {
		if current == nil {
			return nil
		}

		if idx, err := strconv.Atoi(part); err == nil {
			switch arr := current.(type) {
			case []interface{}:
				if idx >= 0 && idx < len(arr) {
					current = arr[idx]
				} else {
					return nil
				}
			case []map[string]interface{}:
				if idx >= 0 && idx < len(arr) {
					current = arr[idx]
				} else {
					return nil
				}
			default:
				return nil
			}
		} else {
			switch obj := current.(type) {
			case map[string]interface{}:
				var ok bool
				current, ok = obj[part]
				if !ok {
					return nil
				}
			default:
				return nil
			}
		}
	}

	return current
}

func (s *EngineService) getNodeName(node models.WorkflowNode) string {
	if node.Name != "" {
		return node.Name
	}

	if name, ok := node.Data["label"].(string); ok && name != "" {
		return name
	}
	if name, ok := node.Data["name"].(string); ok && name != "" {
		return name
	}

	if node.Type == "tool" && node.ToolCode != "" {
		return node.ToolCode
	}

	return node.Type
}

// isCompleteVariableRef 检查字符串是否是完整的变量引用（如 "{{external.image}}"）
func (s *EngineService) isCompleteVariableRef(str string) bool {
	str = strings.TrimSpace(str)
	return strings.HasPrefix(str, "{{") && strings.HasSuffix(str, "}}") && strings.Count(str, "{{") == 1
}

// resolveVariable 解析变量引用，返回实际值（可能是对象、字符串等）
func (s *EngineService) resolveVariable(
	varRef string,
	envMap map[string]string,
	nodeOutputs map[string]map[string]interface{},
	externalParams map[string]interface{},
) interface{} {
	// 去掉 {{ 和 }}
	varRef = strings.TrimSpace(varRef)
	varRef = strings.TrimPrefix(varRef, "{{")
	varRef = strings.TrimSuffix(varRef, "}}")
	varRef = strings.TrimSpace(varRef)

	// 解析变量路径
	if strings.HasPrefix(varRef, "external.") {
		// {{external.image}} 或 {{external.image.path}}
		path := strings.TrimPrefix(varRef, "external.")
		if externalParams != nil {
			return s.getNestedValue(externalParams, path)
		}
	} else if strings.HasPrefix(varRef, "nodes.") {
		// {{nodes.upload.url}}
		path := strings.TrimPrefix(varRef, "nodes.")
		parts := strings.SplitN(path, ".", 2)
		if len(parts) >= 1 {
			nodeID := parts[0]
			if output, ok := nodeOutputs[nodeID]; ok {
				if len(parts) == 2 {
					return s.getNestedValue(output, parts[1])
				}
				return output
			}
		}
	} else if strings.HasPrefix(varRef, "env.") {
		// {{env.api_key}}
		key := strings.TrimPrefix(varRef, "env.")
		if val, ok := envMap[key]; ok {
			return val
		}
	}

	return nil
}

func marshalJSON(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

// cleanupExecutionFiles 清理执行记录的临时文件
func (s *EngineService) cleanupExecutionFiles(executionID string) {
	baseDir := "/tmp/workflow-files"
	dirPath := filepath.Join(baseDir, executionID)

	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return // 没有文件，不需要清理
	}

	// 删除整个执行目录
	if err := os.RemoveAll(dirPath); err != nil {
		log.Error("清理临时文件失败: ExecutionID=%s, Error=%v", executionID, err)
	} else {
		log.Info("已清理临时文件: ExecutionID=%s", executionID)
	}
}
