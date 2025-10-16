package models

import (
	"gorm.io/gorm"
)

// TemplateCategory 模板分类
type TemplateCategory struct {
	BaseModel
	Name        string `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	SortOrder   int    `gorm:"default:100;index" json:"sort_order"`
	IsActive    bool   `gorm:"default:true;index" json:"is_active"`
}

// TableName 指定表名
func (TemplateCategory) TableName() string {
	return "template_category"
}

// BeforeCreate 创建前的钩子
func (tc *TemplateCategory) BeforeCreate(tx *gorm.DB) error {
	if err := tc.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}
