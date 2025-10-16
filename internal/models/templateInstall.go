package models

import (
	"auto-forge/pkg/common"
	"time"

	"gorm.io/gorm"
)

type TemplateInstall struct {
	ID         string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	TemplateID string    `gorm:"type:varchar(36);not null;index" json:"template_id"`
	UserID     string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	WorkflowID string    `gorm:"type:varchar(36);index" json:"workflow_id"`
	InstalledAt time.Time `gorm:"autoCreateTime" json:"installed_at"`
}

func (TemplateInstall) TableName() string {
	return "template_installs"
}

// BeforeCreate hook to generate ID
func (ti *TemplateInstall) BeforeCreate(tx *gorm.DB) error {
	if ti.ID == "" {
		ti.ID = common.NewUUID().String()
	}
	return nil
}
