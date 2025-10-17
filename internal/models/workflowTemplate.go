package models

import (
	"auto-forge/pkg/common"
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type WorkflowTemplate struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"type:varchar(50);not null;default:'other';index" json:"category"`

	TemplateData TemplateData `gorm:"type:json;not null" json:"template_data"`

	CoverImage string `gorm:"type:varchar(512)" json:"cover_image"`
	Icon       string `gorm:"type:varchar(100)" json:"icon"`

	RequiredTools StringArray `gorm:"type:json" json:"required_tools"`
	CaseImages    StringArray `gorm:"type:json" json:"case_images"`

	Status     string `gorm:"type:varchar(20);default:'published';index" json:"status"`
	IsOfficial bool   `gorm:"default:true" json:"is_official"`
	IsFeatured bool   `gorm:"default:false;index" json:"is_featured"`

	InstallCount int `gorm:"default:0;index" json:"install_count"`
	ViewCount    int `gorm:"default:0" json:"view_count"`

	AuthorID   string `gorm:"type:varchar(36);index" json:"author_id"`
	AuthorName string `gorm:"type:varchar(100)" json:"author_name"`

	UsageGuide string `gorm:"type:text" json:"usage_guide"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type TemplateData struct {
	Nodes   []WorkflowNode  `json:"nodes"`
	Edges   []WorkflowEdge  `json:"edges"`
	EnvVars []WorkflowEnvVar `json:"env_vars,omitempty"`
}

func (td TemplateData) Value() (driver.Value, error) {
	return json.Marshal(td)
}

func (td *TemplateData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, td)
}

type StringArray []string

func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return "[]", nil
	}
	return json.Marshal(sa)
}

func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		*sa = []string{}
		return nil
	}
	return json.Unmarshal(bytes, sa)
}

func (WorkflowTemplate) TableName() string {
	return "workflow_templates"
}

// BeforeCreate hook to generate ID
func (wt *WorkflowTemplate) BeforeCreate(tx *gorm.DB) error {
	if wt.ID == "" {
		wt.ID = common.NewUUID().String()
	}
	return nil
}
