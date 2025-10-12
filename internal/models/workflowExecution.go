package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

// ExecutionStatus 执行状态
const (
	ExecutionStatusPending   = "pending"   // 等待中
	ExecutionStatusRunning   = "running"   // 执行中
	ExecutionStatusSuccess   = "success"   // 成功
	ExecutionStatusFailed    = "failed"    // 失败
	ExecutionStatusCancelled = "cancelled" // 已取消
)

// NodeExecutionLog 节点执行日志
type NodeExecutionLog struct {
	NodeID      string                 `json:"node_id"`
	NodeType    string                 `json:"node_type"`
	NodeName    string                 `json:"node_name"`
	Status      string                 `json:"status"` // pending/running/success/failed/skipped
	StartTime   *int64                 `json:"start_time"`
	EndTime     *int64                 `json:"end_time"`
	DurationMs  int64                  `json:"duration_ms"`
	Input       map[string]interface{} `json:"input,omitempty"`
	Output      map[string]interface{} `json:"output,omitempty"`
	Error       string                 `json:"error,omitempty"`
	RetryCount  int                    `json:"retry_count"`
	ToolCode    string                 `json:"tool_code,omitempty"`
	ToolVersion string                 `json:"tool_version,omitempty"`
}

// NodeExecutionLogs 节点执行日志数组
type NodeExecutionLogs []NodeExecutionLog

// Scan 实现 sql.Scanner 接口
func (nel *NodeExecutionLogs) Scan(value interface{}) error {
	if value == nil {
		*nel = NodeExecutionLogs{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, nel)
}

// Value 实现 driver.Valuer 接口
func (nel NodeExecutionLogs) Value() (driver.Value, error) {
	if len(nel) == 0 {
		return "[]", nil
	}
	return json.Marshal(nel)
}

// WorkflowExecution 工作流执行记录
type WorkflowExecution struct {
	BaseModel
	WorkflowID   string            `gorm:"type:char(36);not null;index:idx_workflow_id" json:"workflow_id"`
	Workflow     *Workflow         `gorm:"-" json:"workflow,omitempty"`
	UserID       string            `gorm:"type:char(36);not null;index:idx_user_id" json:"user_id"`
	User         *User             `gorm:"-" json:"user,omitempty"`
	Status       string            `gorm:"size:20;not null;index:idx_status" json:"status"`
	TriggerType  string            `gorm:"size:20" json:"trigger_type"` // manual/schedule/webhook
	StartTime    *int64            `gorm:"index:idx_start_time" json:"start_time"`
	EndTime      *int64            `json:"end_time"`
	DurationMs   int64             `json:"duration_ms"`
	TotalNodes   int               `gorm:"default:0" json:"total_nodes"`
	SuccessNodes int               `gorm:"default:0" json:"success_nodes"`
	FailedNodes  int               `gorm:"default:0" json:"failed_nodes"`
	SkippedNodes int               `gorm:"default:0" json:"skipped_nodes"`
	NodeLogs     NodeExecutionLogs `gorm:"type:json" json:"node_logs"`
	Error        string            `gorm:"type:text" json:"error,omitempty"`
}

// TableName 指定表名
func (WorkflowExecution) TableName() string {
	return "workflow_execution"
}

// BeforeCreate 创建前的钩子
func (we *WorkflowExecution) BeforeCreate(tx *gorm.DB) error {
	if err := we.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}
