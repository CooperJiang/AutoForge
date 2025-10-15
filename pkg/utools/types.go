package utools

import (
	"context"
	"encoding/json"
)

// Tool 工具接口 - 所有工具必须实现此接口
type Tool interface {
	// GetMetadata 获取工具元数据
	GetMetadata() *ToolMetadata

	// GetSchema 获取工具配置的 JSON Schema
	GetSchema() *ConfigSchema

	// Validate 验证工具配置是否有效
	Validate(config map[string]interface{}) error

	// Execute 执行工具任务
	Execute(ctx *ExecutionContext, config map[string]interface{}) (*ExecutionResult, error)
}

type ToolMetadata struct {
	Code               string                    `json:"code"`
	Name               string                    `json:"name"`
	Description        string                    `json:"description"`
	Category           string                    `json:"category"`
	Version            string                    `json:"version"`
	Author             string                    `json:"author"`
	AICallable         bool                      `json:"ai_callable"`
	Tags               []string                  `json:"tags"`
	OutputFieldsSchema map[string]OutputFieldDef `json:"output_fields_schema,omitempty"`
}

type OutputFieldDef struct {
	Type     string                    `json:"type"`
	Label    string                    `json:"label"`
	Children map[string]OutputFieldDef `json:"children,omitempty"`
}

// ConfigSchema 工具配置的 JSON Schema 定义
type ConfigSchema struct {
	Type       string                    `json:"type"`       // 通常是 "object"
	Properties map[string]PropertySchema `json:"properties"` // 配置项定义
	Required   []string                  `json:"required"`   // 必填字段列表
}

// PropertySchema 配置项的 Schema 定义
type PropertySchema struct {
	Type        string        `json:"type"`                  // string, number, boolean, array, object
	Title       string        `json:"title"`                 // 显示标题
	Description string        `json:"description"`           // 描述说明
	Default     interface{}   `json:"default,omitempty"`     // 默认值
	Enum        []interface{} `json:"enum,omitempty"`        // 枚举值（下拉选项）
	Format      string        `json:"format,omitempty"`      // 格式: url, email, date-time, etc.
	Pattern     string        `json:"pattern,omitempty"`     // 正则表达式校验
	MinLength   *int          `json:"minLength,omitempty"`   // 最小长度
	MaxLength   *int          `json:"maxLength,omitempty"`   // 最大长度
	Minimum     *float64      `json:"minimum,omitempty"`     // 最小值
	Maximum     *float64      `json:"maximum,omitempty"`     // 最大值
	Items       *PropertySchema `json:"items,omitempty"`     // 数组项定义
	Properties  map[string]PropertySchema `json:"properties,omitempty"` // 对象属性（嵌套对象）
	Secret      bool          `json:"secret,omitempty"`      // 是否是敏感信息（需加密存储）
}

// ExecutionContext 执行上下文
type ExecutionContext struct {
	Context   context.Context        `json:"-"`           // Go Context (超时控制、取消等)
	TaskID    string                 `json:"task_id"`     // 任务 ID
	UserID    string                 `json:"user_id"`     // 用户 ID
	Variables map[string]interface{} `json:"variables"`   // 上下文变量（用于工作流中传递数据）
	Metadata  map[string]interface{} `json:"metadata"`    // 额外元数据
}

// ExecutionResult 执行结果
type ExecutionResult struct {
	Success      bool                   `json:"success"`                 // 是否成功
	Message      string                 `json:"message"`                 // 消息
	Output       map[string]interface{} `json:"output"`                  // 输出数据
	Error        string                 `json:"error"`                   // 错误信息
	DurationMs   int64                  `json:"duration_ms"`             // 执行耗时(毫秒)
	StatusCode   int                    `json:"status_code"`             // 状态码（如 HTTP 工具返回的状态码）
	ResponseBody string                 `json:"response_body"`           // 响应体（原始内容）
	OutputRender *OutputRenderConfig    `json:"output_render,omitempty"` // 输出渲染配置（可选）
}

// OutputRenderConfig 输出渲染配置
// 用于指示前端如何渲染输出数据
type OutputRenderConfig struct {
	Type    string                 `json:"type"`    // 输出类型：image, video, html, markdown, text, gallery, json
	Primary string                 `json:"primary"` // 主要显示字段的路径，如 "content" 或 "data.url"
	Fields  map[string]FieldRender `json:"fields"`  // 各字段的渲染配置
}

// FieldRender 字段渲染配置
type FieldRender struct {
	Type    string `json:"type"`    // 字段类型：image, video, url, text, json, code, markdown
	Label   string `json:"label"`   // 显示标签
	Display bool   `json:"display"` // 是否显示该字段
}

// BaseTool 工具基础实现 - 提供通用功能
type BaseTool struct {
	metadata *ToolMetadata
	schema   *ConfigSchema
}

// NewBaseTool 创建基础工具
func NewBaseTool(metadata *ToolMetadata, schema *ConfigSchema) *BaseTool {
	return &BaseTool{
		metadata: metadata,
		schema:   schema,
	}
}

// GetMetadata 获取元数据
func (bt *BaseTool) GetMetadata() *ToolMetadata {
	return bt.metadata
}

// GetSchema 获取配置 Schema
func (bt *BaseTool) GetSchema() *ConfigSchema {
	return bt.schema
}

// Validate 验证配置（基础校验）
func (bt *BaseTool) Validate(config map[string]interface{}) error {
	// 检查必填字段
	for _, field := range bt.schema.Required {
		if _, ok := config[field]; !ok {
			return &ValidationError{
				Field:   field,
				Message: "required field missing",
			}
		}
	}

	// 校验每个字段类型
	for key, value := range config {
		propSchema, ok := bt.schema.Properties[key]
		if !ok {
			continue // 未定义的字段跳过
		}

		if err := bt.validateValue(key, value, propSchema); err != nil {
			return err
		}
	}

	return nil
}

// validateValue 验证单个值
func (bt *BaseTool) validateValue(key string, value interface{}, schema PropertySchema) error {
	// 类型校验
	switch schema.Type {
	case "string":
		str, ok := value.(string)
		if !ok {
			return &ValidationError{Field: key, Message: "must be a string"}
		}
		// 长度校验
		if schema.MinLength != nil && len(str) < *schema.MinLength {
			return &ValidationError{Field: key, Message: "string too short"}
		}
		if schema.MaxLength != nil && len(str) > *schema.MaxLength {
			return &ValidationError{Field: key, Message: "string too long"}
		}

	case "number":
		num, ok := value.(float64)
		if !ok {
			return &ValidationError{Field: key, Message: "must be a number"}
		}
		if schema.Minimum != nil && num < *schema.Minimum {
			return &ValidationError{Field: key, Message: "number too small"}
		}
		if schema.Maximum != nil && num > *schema.Maximum {
			return &ValidationError{Field: key, Message: "number too large"}
		}

	case "boolean":
		if _, ok := value.(bool); !ok {
			return &ValidationError{Field: key, Message: "must be a boolean"}
		}

	case "array":
		if _, ok := value.([]interface{}); !ok {
			return &ValidationError{Field: key, Message: "must be an array"}
		}

	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return &ValidationError{Field: key, Message: "must be an object"}
		}
	}

	// 枚举校验
	if len(schema.Enum) > 0 {
		found := false
		for _, enumVal := range schema.Enum {
			if value == enumVal {
				found = true
				break
			}
		}
		if !found {
			return &ValidationError{Field: key, Message: "value not in allowed enum"}
		}
	}

	return nil
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

// ToolError 工具执行错误
type ToolError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ToolError) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}
