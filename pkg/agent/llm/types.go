package llm

import "context"

// LLMClient LLM 客户端接口
type LLMClient interface {
	// Call 同步调用
	Call(ctx context.Context, messages []Message, options *CallOptions) (*Response, error)

	// Stream 流式调用
	Stream(ctx context.Context, messages []Message, options *CallOptions) (<-chan StreamChunk, error)

	// GetModelInfo 获取模型信息
	GetModelInfo() ModelInfo
}

// Message 消息
type Message struct {
	Role       string     `json:"role"`                   // system, user, assistant, tool
	Content    string     `json:"content"`                // 文本内容
	Name       string     `json:"name,omitempty"`         // 可选：消息发送者名称
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // 可选：工具调用
	ToolCallID string     `json:"tool_call_id,omitempty"` // 可选：工具调用ID（用于tool角色）
}

// CallOptions 调用选项
type CallOptions struct {
	Temperature    float64          `json:"temperature,omitempty"`     // 温度 0-2
	MaxTokens      int              `json:"max_tokens,omitempty"`      // 最大 token 数
	TopP           float64          `json:"top_p,omitempty"`           // 核采样
	Stop           []string         `json:"stop,omitempty"`            // 停止序列
	Tools          []ToolDefinition `json:"tools,omitempty"`           // 可用工具列表
	ToolChoice     interface{}      `json:"tool_choice,omitempty"`     // auto, none, required, {"type": "function", "function": {"name": "xxx"}}
	ResponseFormat string           `json:"response_format,omitempty"` // text, json_object
}

// Response 响应
type Response struct {
	Content      string     `json:"content"`       // 文本内容
	ToolCalls    []ToolCall `json:"tool_calls"`    // 工具调用
	FinishReason string     `json:"finish_reason"` // stop, tool_calls, length, content_filter
	Usage        TokenUsage `json:"usage"`         // Token 使用情况
}

// StreamChunk 流式响应块
type StreamChunk struct {
	Content      string     `json:"content"`       // 增量内容
	ToolCalls    []ToolCall `json:"tool_calls"`    // 工具调用（增量）
	FinishReason string     `json:"finish_reason"` // 终止原因
	Done         bool       `json:"done"`          // 是否完成
	Error        error      `json:"error"`         // 错误
}

// ToolCall 工具调用
type ToolCall struct {
	ID       string       `json:"id"`       // 工具调用ID
	Type     string       `json:"type"`     // function
	Function FunctionCall `json:"function"` // 函数调用
}

// FunctionCall 函数调用
type FunctionCall struct {
	Name      string `json:"name"`      // 函数名
	Arguments string `json:"arguments"` // JSON 字符串参数
}

// ToolDefinition 工具定义
type ToolDefinition struct {
	Type     string             `json:"type"`     // "function"
	Function FunctionDefinition `json:"function"` // 函数定义
}

// FunctionDefinition 函数定义
type FunctionDefinition struct {
	Name        string                 `json:"name"`               // 函数名
	Description string                 `json:"description"`        // 描述
	Parameters  map[string]interface{} `json:"parameters"`         // JSON Schema
	Metadata    map[string]interface{} `json:"metadata,omitempty"` // 元数据（包含输出字段定义等）
}

// TokenUsage Token 使用情况
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`     // 提示词 token 数
	CompletionTokens int `json:"completion_tokens"` // 完成 token 数
	TotalTokens      int `json:"total_tokens"`      // 总 token 数
}

// ModelInfo 模型信息
type ModelInfo struct {
	Provider    string `json:"provider"`     // 提供商：openai, gemini, custom
	Model       string `json:"model"`        // 模型名称
	MaxTokens   int    `json:"max_tokens"`   // 最大 token 数
	SupportTool bool   `json:"support_tool"` // 是否支持工具调用
}
