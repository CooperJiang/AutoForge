package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

// KeyValue 键值对结构
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// KeyValueArray 键值对数组，用于存储 headers 和 params
type KeyValueArray []KeyValue

// Scan 实现 sql.Scanner 接口
func (kva *KeyValueArray) Scan(value interface{}) error {
	if value == nil {
		*kva = KeyValueArray{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, kva)
}

// Value 实现 driver.Valuer 接口
func (kva KeyValueArray) Value() (driver.Value, error) {
	if len(kva) == 0 {
		return "[]", nil
	}
	return json.Marshal(kva)
}

// Task 定时任务模型 - 直接基于工具执行
type Task struct {
	BaseModel
	UserID      string  `gorm:"type:char(36);not null;index:idx_user_id" json:"user_id"` // 关联User表的UUID
	User        *User   `gorm:"-" json:"user,omitempty"` // 关联用户（不创建外键约束）
	Name        string  `gorm:"size:255;not null" json:"name"`                           // 任务名称
	Description string  `gorm:"type:text" json:"description"`                            // 任务描述
	ToolCode    string  `gorm:"size:50;not null;index" json:"tool_code"`                 // 工具代码（如 http_request）
	Config      string  `gorm:"type:text;not null" json:"config"`                        // 工具配置（JSON格式）

	// 调度配置
	ScheduleType  string `gorm:"size:20;not null" json:"schedule_type"`                  // daily/hourly/interval/cron
	ScheduleValue string `gorm:"size:100;not null" json:"schedule_value"`                // 调度值
	Enabled       bool   `gorm:"default:true;index:idx_enabled_next_run" json:"enabled"` // 是否启用
	NextRunTime   *int64 `gorm:"index:idx_enabled_next_run" json:"next_run_time"`        // 下次执行时间(Unix timestamp)
}

// TableName 指定表名
func (Task) TableName() string {
	return "task"
}

// BeforeCreate 创建前的钩子
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的BeforeCreate
	if err := t.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}
