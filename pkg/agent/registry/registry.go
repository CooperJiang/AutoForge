package registry

import (
	"auto-forge/pkg/agent/llm"
	"auto-forge/pkg/utools"
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// ToolRegistry 工具注册表
type ToolRegistry struct {
	tools map[string]*ToolWrapper
	mu    sync.RWMutex
}

// ToolWrapper 工具包装器
type ToolWrapper struct {
	Tool       utools.Tool
	Definition llm.ToolDefinition
}

// NewToolRegistry 创建工具注册表
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]*ToolWrapper),
	}
}

// RegisterFromUTools 从 UTools 注册所有工具
func (r *ToolRegistry) RegisterFromUTools() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 获取所有已注册的工具
	allTools := utools.GetRegistry().GetAllTools()

	for _, tool := range allTools {
		metadata := tool.GetMetadata()

		// 跳过不适合 Agent 使用的工具
		if shouldSkipTool(metadata.Code) {
			continue
		}

		// 跳过 AICallable = false 的工具
		if !metadata.AICallable {
			continue
		}

		// 生成工具定义
		definition, err := generateToolDefinition(tool)
		if err != nil {
			return fmt.Errorf("生成工具定义失败 [%s]: %w", metadata.Code, err)
		}

		r.tools[metadata.Code] = &ToolWrapper{
			Tool:       tool,
			Definition: definition,
		}
	}

	return nil
}

// GetToolDefinitions 获取工具定义列表（用于 LLM）
func (r *ToolRegistry) GetToolDefinitions(allowedTools []string) []llm.ToolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var definitions []llm.ToolDefinition

	// 如果没有指定工具，返回所有工具
	if len(allowedTools) == 0 {
		for _, wrapper := range r.tools {
			definitions = append(definitions, wrapper.Definition)
		}
		return definitions
	}

	// 返回指定的工具
	for _, toolName := range allowedTools {
		if wrapper, ok := r.tools[toolName]; ok {
			definitions = append(definitions, wrapper.Definition)
		}
	}

	return definitions
}

// Execute 执行工具
func (r *ToolRegistry) Execute(ctx context.Context, toolName string, args map[string]interface{}) (interface{}, error) {
	r.mu.RLock()
	wrapper, ok := r.tools[toolName]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("工具不存在: %s", toolName)
	}

	// 执行工具
	result, err := wrapper.Tool.Execute(&utools.ExecutionContext{Context: ctx}, args)
	if err != nil {
		return nil, fmt.Errorf("工具执行失败: %w", err)
	}

	return result, nil
}

// GetTool 获取工具
func (r *ToolRegistry) GetTool(toolName string) (utools.Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	wrapper, ok := r.tools[toolName]
	if !ok {
		return nil, fmt.Errorf("工具不存在: %s", toolName)
	}

	return wrapper.Tool, nil
}

// ListTools 列出所有工具名称
func (r *ToolRegistry) ListTools() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]string, 0, len(r.tools))
	for name := range r.tools {
		tools = append(tools, name)
	}
	return tools
}

// shouldSkipTool 判断是否应该跳过该工具
func shouldSkipTool(toolCode string) bool {
	// 跳过控制流工具和内部工具
	skipList := []string{
		"condition",        // 条件判断
		"switch",           // 分支
		"delay",            // 延迟
		"loop",             // 循环
		"external_trigger", // 外部触发
	}

	for _, skip := range skipList {
		if toolCode == skip {
			return true
		}
	}

	return false
}

// generateToolDefinition 从 UTool 生成 LLM 工具定义
func generateToolDefinition(tool utools.Tool) (llm.ToolDefinition, error) {
	// 获取工具的配置 Schema
	schema := tool.GetSchema()

	// 构建参数 Schema
	parameters := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
		"required":   []string{},
	}

	properties := parameters["properties"].(map[string]interface{})
	var required []string

	// 解析 Schema
	if schema != nil {
		// 直接遍历 Properties (已经是 map[string]utools.PropertySchema)
		for key, propSchema := range schema.Properties {
			// 构建属性定义
			prop := make(map[string]interface{})

			// 类型
			if propSchema.Type != "" {
				prop["type"] = propSchema.Type
			}

			// 描述
			if propSchema.Description != "" {
				prop["description"] = propSchema.Description
			}

			// 枚举
			if len(propSchema.Enum) > 0 {
				prop["enum"] = propSchema.Enum
			}

			// 默认值
			if propSchema.Default != nil {
				prop["default"] = propSchema.Default
			}

			properties[key] = prop
		}

		// Required 字段 (已经是 []string)
		if len(schema.Required) > 0 {
			required = append(required, schema.Required...)
		}
	}

	if len(required) > 0 {
		parameters["required"] = required
	}

	// 获取工具元数据
	metadata := tool.GetMetadata()

	// 构建工具定义
	definition := llm.ToolDefinition{
		Type: "function",
		Function: llm.FunctionDefinition{
			Name:        metadata.Code,
			Description: metadata.Description,
			Parameters:  parameters,
			Metadata: map[string]interface{}{
				"output_fields_schema": metadata.OutputFieldsSchema,
			},
		},
	}

	return definition, nil
}

// FormatToolResult 格式化工具结果为字符串
func FormatToolResult(result interface{}) string {
	if result == nil {
		return "null"
	}

	// 如果已经是字符串，直接返回
	if str, ok := result.(string); ok {
		return str
	}

	// 尝试 JSON 序列化
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("%v", result)
	}

	return string(data)
}
