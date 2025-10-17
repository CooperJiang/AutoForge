package models

import (
	"time"
)

// ToolConfig 工具配置模型
type ToolConfig struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ToolCode     string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:工具代码" json:"tool_code"`
	ToolName     string    `gorm:"type:varchar(100);not null;comment:工具名称" json:"tool_name"`
	Enabled      bool      `gorm:"default:true;comment:是否启用" json:"enabled"`
	Visible      bool      `gorm:"default:true;comment:是否对外展示" json:"visible"`
	IsDeprecated bool      `gorm:"default:false;comment:是否已废弃" json:"is_deprecated"`
	ConfigJSON   string    `gorm:"type:text;comment:配置JSON(加密)" json:"-"` // 不返回给前端
	ConfigSchema string    `gorm:"type:text;comment:配置结构定义JSON" json:"config_schema"`
	Description  string    `gorm:"type:varchar(500);comment:工具描述" json:"description"`
	Category     string    `gorm:"type:varchar(50);comment:工具分类" json:"category"`
	Version      string    `gorm:"type:varchar(20);comment:工具版本" json:"version"`
	Author       string    `gorm:"type:varchar(100);comment:工具作者" json:"author"`
	Tags         string    `gorm:"type:varchar(500);comment:工具标签(逗号分隔)" json:"tags"`
	SortOrder    int       `gorm:"default:0;comment:排序" json:"sort_order"`
	LastSyncAt   time.Time `gorm:"comment:最后同步时间" json:"last_sync_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (ToolConfig) TableName() string {
	return "tool_config"
}
