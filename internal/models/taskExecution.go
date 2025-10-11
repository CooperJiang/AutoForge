package models

import (
	"gorm.io/gorm"
)

// TaskExecution 任务执行记录模型
type TaskExecution struct {
	BaseModel
	TaskID         string        `gorm:"type:char(36);not null;index:idx_task_id_started" json:"task_id"`
	Task           *Task         `gorm:"foreignKey:TaskID;references:ID" json:"task,omitempty"` // 关联任务
	UserID         string        `gorm:"type:char(36);not null;index:idx_user_id_started" json:"user_id"`
	User           *User         `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"` // 关联用户
	Status         string        `gorm:"size:20;not null" json:"status"`                        // success/failed/timeout
	RequestURL     string        `gorm:"type:text" json:"request_url"`
	RequestMethod  string        `gorm:"size:10" json:"request_method"`
	RequestHeaders KeyValueArray `gorm:"type:json" json:"request_headers"`
	RequestParams  KeyValueArray `gorm:"type:json" json:"request_params"`
	RequestBody    string        `gorm:"type:text" json:"request_body"` // 请求体
	ResponseStatus int           `json:"response_status"`                // HTTP 状态码
	ResponseBody   string        `gorm:"type:longtext" json:"response_body"`
	DurationMs     int64         `json:"duration_ms"` // 执行耗时（毫秒）
	ErrorMessage   string        `gorm:"type:longtext" json:"error_message"`
	StartedAt      int64         `gorm:"index:idx_task_id_started;index:idx_user_id_started" json:"started_at"` // Unix timestamp
	CompletedAt    int64         `json:"completed_at"`                                                           // Unix timestamp
}

// TableName 指定表名
func (TaskExecution) TableName() string {
	return "task_execution"
}

// BeforeCreate 创建前的钩子
func (te *TaskExecution) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的BeforeCreate
	if err := te.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 初始化空数组
	if te.RequestHeaders == nil {
		te.RequestHeaders = KeyValueArray{}
	}
	if te.RequestParams == nil {
		te.RequestParams = KeyValueArray{}
	}

	return nil
}
