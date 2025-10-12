package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

// WorkflowNode 工作流节点
type WorkflowNode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"` // tool/trigger/condition/delay/switch
	ToolCode string                 `json:"toolCode"` // 工具代码（仅当type=tool时）
	Name     string                 `json:"name"`
	Config   map[string]interface{} `json:"config"`
	Retry    *NodeRetryConfig       `json:"retry,omitempty"`
	Position map[string]float64     `json:"position"`
	Data     map[string]interface{} `json:"data,omitempty"` // 保留兼容性
}

// NodeRetryConfig 节点重试配置
type NodeRetryConfig struct {
	Enabled            bool `json:"enabled"`
	MaxRetries         int  `json:"maxRetries"`
	RetryInterval      int  `json:"retryInterval"` // 毫秒
	ExponentialBackoff bool `json:"exponentialBackoff"`
}

// WorkflowEdge 工作流连接线
type WorkflowEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label,omitempty"`
}

// WorkflowEnvVar 工作流环境变量
type WorkflowEnvVar struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Encrypted   bool   `json:"encrypted"`
}

// WorkflowNodes 节点数组类型
type WorkflowNodes []WorkflowNode

// Scan 实现 sql.Scanner 接口
func (wn *WorkflowNodes) Scan(value interface{}) error {
	if value == nil {
		*wn = WorkflowNodes{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, wn)
}

// Value 实现 driver.Valuer 接口
func (wn WorkflowNodes) Value() (driver.Value, error) {
	if len(wn) == 0 {
		return "[]", nil
	}
	return json.Marshal(wn)
}

// WorkflowEdges 连接线数组类型
type WorkflowEdges []WorkflowEdge

// Scan 实现 sql.Scanner 接口
func (we *WorkflowEdges) Scan(value interface{}) error {
	if value == nil {
		*we = WorkflowEdges{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, we)
}

// Value 实现 driver.Valuer 接口
func (we WorkflowEdges) Value() (driver.Value, error) {
	if len(we) == 0 {
		return "[]", nil
	}
	return json.Marshal(we)
}

// WorkflowEnvVars 环境变量数组类型
type WorkflowEnvVars []WorkflowEnvVar

// Scan 实现 sql.Scanner 接口
func (wev *WorkflowEnvVars) Scan(value interface{}) error {
	if value == nil {
		*wev = WorkflowEnvVars{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, wev)
}

// Value 实现 driver.Valuer 接口
func (wev WorkflowEnvVars) Value() (driver.Value, error) {
	if len(wev) == 0 {
		return "[]", nil
	}
	return json.Marshal(wev)
}

// Workflow 工作流模型
type Workflow struct {
	BaseModel
	UserID      string          `gorm:"type:char(36);not null;index:idx_user_id" json:"user_id"`
	User        *User           `gorm:"-" json:"user,omitempty"`
	Name        string          `gorm:"size:255;not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	Nodes       WorkflowNodes   `gorm:"type:json;not null" json:"nodes"`
	Edges       WorkflowEdges   `gorm:"type:json;not null" json:"edges"`
	EnvVars     WorkflowEnvVars `gorm:"type:json" json:"env_vars"`

	// 调度配置
	ScheduleType  string `gorm:"size:20" json:"schedule_type"`
	ScheduleValue string `gorm:"size:100" json:"schedule_value"`
	Enabled       bool   `gorm:"default:false;index:idx_enabled_next_run" json:"enabled"`
	NextRunTime   *int64 `gorm:"index:idx_enabled_next_run" json:"next_run_time"`

	// 统计信息
	TotalExecutions int    `gorm:"default:0" json:"total_executions"`
	SuccessCount    int    `gorm:"default:0" json:"success_count"`
	FailedCount     int    `gorm:"default:0" json:"failed_count"`
	LastExecutedAt  *int64 `gorm:"index" json:"last_executed_at"`
}

// TableName 指定表名
func (Workflow) TableName() string {
	return "workflow"
}

// BeforeCreate 创建前的钩子
func (w *Workflow) BeforeCreate(tx *gorm.DB) error {
	if err := w.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}
