package response

import (
	"auto-forge/internal/models"
	"time"
)

type TemplateListResponse struct {
	Items      []TemplateBasicInfo `json:"items"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

type TemplateBasicInfo struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	CoverImage    string    `json:"cover_image"`
	Icon          string    `json:"icon"`
	InstallCount  int       `json:"install_count"`
	ViewCount     int       `json:"view_count"`
	IsOfficial    bool      `json:"is_official"`
	IsFeatured    bool      `json:"is_featured"`
	AuthorName    string    `json:"author_name"`
	RequiredTools []string  `json:"required_tools"`
	CreatedAt     time.Time `json:"created_at"`
}

type TemplateDetailResponse struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	Category      string                `json:"category"`
	CoverImage    string                `json:"cover_image"`
	Icon          string                `json:"icon"`
	InstallCount  int                   `json:"install_count"`
	ViewCount     int                   `json:"view_count"`
	IsOfficial    bool                  `json:"is_official"`
	IsFeatured    bool                  `json:"is_featured"`
	AuthorName    string                `json:"author_name"`
	RequiredTools []string              `json:"required_tools"`
	UsageGuide    string                `json:"usage_guide"`
	TemplateData  models.TemplateData   `json:"template_data"`
	Status        string                `json:"status"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

type InstallTemplateResponse struct {
	WorkflowID   string `json:"workflow_id"`
	WorkflowName string `json:"workflow_name"`
	Message      string `json:"message"`
}
