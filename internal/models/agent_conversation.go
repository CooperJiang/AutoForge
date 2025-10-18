package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// AgentConversation Agent 对话
type AgentConversation struct {
	ID        string `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID    string `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Title     string `gorm:"type:varchar(255);not null" json:"title"`
	CreatedAt int64  `gorm:"not null" json:"created_at"`
	UpdatedAt int64  `gorm:"not null" json:"updated_at"`
}

// TableName 指定表名
func (AgentConversation) TableName() string {
	return "agent_conversations"
}

// GetID 获取ID
func (c *AgentConversation) GetID() string {
	return c.ID
}

// BeforeCreate GORM 钩子：创建前
func (c *AgentConversation) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().Unix()
	if c.CreatedAt == 0 {
		c.CreatedAt = now
	}
	if c.UpdatedAt == 0 {
		c.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (c *AgentConversation) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now().Unix()
	return nil
}

// AgentMessage Agent 消息
type AgentMessage struct {
	ID             string       `gorm:"type:varchar(36);primaryKey" json:"id"`
	ConversationID string       `gorm:"type:varchar(36);not null;index" json:"conversation_id"`
	Role           string       `gorm:"type:varchar(20);not null" json:"role"` // user/agent/system
	Content        string       `gorm:"type:text;not null" json:"content"`
	Files          AgentFiles   `gorm:"type:json" json:"files,omitempty"`        // 用户上传的文件
	Trace          *AgentTrace  `gorm:"type:json" json:"trace,omitempty"`        // Agent 执行轨迹
	Plan           *AgentPlan   `gorm:"type:json" json:"plan,omitempty"`         // Agent 执行计划
	Config         *AgentConfig `gorm:"type:json" json:"config,omitempty"`       // Agent 配置
	TokenUsage     *TokenUsage  `gorm:"type:json" json:"token_usage,omitempty"`  // Token 使用情况
	Status         string       `gorm:"type:varchar(20);not null" json:"status"` // pending/running/completed/failed
	Error          string       `gorm:"type:text" json:"error,omitempty"`        // 错误信息
	CreatedAt      int64        `gorm:"not null" json:"created_at"`
}

// TableName 指定表名
func (AgentMessage) TableName() string {
	return "agent_messages"
}

// GetID 获取ID
func (m *AgentMessage) GetID() string {
	return m.ID
}

// BeforeCreate GORM 钩子：创建前
func (m *AgentMessage) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedAt == 0 {
		m.CreatedAt = time.Now().Unix()
	}
	if m.Status == "" {
		m.Status = "pending"
	}
	return nil
}

// AgentConfig Agent 配置
type AgentConfig struct {
	Mode         string   `json:"mode,omitempty"`          // plan/direct
	Model        string   `json:"model,omitempty"`         // gpt-4o-mini, etc.
	MaxSteps     int      `json:"max_steps,omitempty"`     // 最大步骤数
	Temperature  float64  `json:"temperature,omitempty"`   // 温度
	AllowedTools []string `json:"allowed_tools,omitempty"` // 允许的工具列表
}

// Scan 实现 sql.Scanner 接口
func (c *AgentConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, c)
}

// Value 实现 driver.Valuer 接口
func (c AgentConfig) Value() (driver.Value, error) {
	if c.Model == "" && c.Mode == "" {
		return nil, nil
	}
	return json.Marshal(c)
}

// AgentFiles 文件列表
type AgentFiles []AgentFile

// AgentFile 文件信息
type AgentFile struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

// Scan 实现 sql.Scanner 接口
func (f *AgentFiles) Scan(value interface{}) error {
	if value == nil {
		*f = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, f)
}

// Value 实现 driver.Valuer 接口
func (f AgentFiles) Value() (driver.Value, error) {
	if len(f) == 0 {
		return nil, nil
	}
	return json.Marshal(f)
}

// AgentTrace Agent 执行轨迹
type AgentTrace struct {
	Steps        []AgentStep            `json:"steps"`
	FinalAnswer  string                 `json:"final_answer"`
	FinishReason string                 `json:"finish_reason"` // final/max_steps/timeout/error
	UsedTools    map[string]interface{} `json:"used_tools,omitempty"`
	TotalMs      int64                  `json:"total_ms"`
	TokenUsage   *TokenUsage            `json:"token_usage,omitempty"` // Token 使用统计
}

// TokenUsage Token 使用统计
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Scan 实现 sql.Scanner 接口
func (t *AgentTrace) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, t)
}

// Value 实现 driver.Valuer 接口
func (t AgentTrace) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// AgentStep Agent 执行步骤
type AgentStep struct {
	Step        int                    `json:"step"`
	Action      *AgentAction           `json:"action,omitempty"`
	Observation string                 `json:"observation,omitempty"`
	ToolOutput  map[string]interface{} `json:"tool_output,omitempty"`
	ElapsedMs   int64                  `json:"elapsed_ms"`
	Timestamp   string                 `json:"timestamp"`
	Error       string                 `json:"error,omitempty"` // 错误信息
}

// AgentAction Agent 动作
type AgentAction struct {
	Type string                 `json:"type"` // action/final
	Tool string                 `json:"tool,omitempty"`
	Args map[string]interface{} `json:"args,omitempty"`
}

// AgentPlan Agent 执行计划（Plan 模式）
type AgentPlan struct {
	Steps       []AgentPlanStep `json:"steps"`
	TotalSteps  int             `json:"total_steps"`
	CreatedAt   string          `json:"created_at"`
	GeneratedBy string          `json:"generated_by"` // 生成计划的模型
}

// Scan 实现 sql.Scanner 接口
func (p *AgentPlan) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

// Value 实现 driver.Valuer 接口
func (p AgentPlan) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// AgentPlanStep 计划步骤
type AgentPlanStep struct {
	Step        int    `json:"step"`
	Description string `json:"description"`
	Tool        string `json:"tool,omitempty"`
	Status      string `json:"status"` // pending/running/completed/skipped/failed
}

// TokenUsage Token 使用情况

// Scan 实现 sql.Scanner 接口
func (t *TokenUsage) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, t)
}

// Value 实现 driver.Valuer 接口
func (t TokenUsage) Value() (driver.Value, error) {
	return json.Marshal(t)
}
