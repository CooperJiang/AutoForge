package request

// CreateAgentConversationRequest 创建对话请求
type CreateAgentConversationRequest struct {
	Title string `json:"title" binding:"required"`
}

// UpdateAgentConversationRequest 更新对话请求
type UpdateAgentConversationRequest struct {
	Title string `json:"title" binding:"required"`
}

// AgentConfig Agent 配置
type AgentConfig struct {
	Mode         string   `json:"mode,omitempty"`          // plan/direct
	Model        string   `json:"model,omitempty"`         // gpt-4o-mini, etc.
	MaxSteps     int      `json:"max_steps,omitempty"`     // 最大步骤数
	Temperature  float64  `json:"temperature,omitempty"`   // 温度
	AllowedTools []string `json:"allowed_tools,omitempty"` // 允许的工具列表
}

// SendAgentMessageRequest 发送消息请求
type SendAgentMessageRequest struct {
	Message string       `json:"message" binding:"required"`
	Config  *AgentConfig `json:"config"`
}
