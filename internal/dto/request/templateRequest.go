package request

type CreateTemplateRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	Category      string   `json:"category" binding:"required"`
	WorkflowID    string   `json:"workflow_id" binding:"required"`
	CoverImage    string   `json:"cover_image"`
	Icon          string   `json:"icon"`
	RequiredTools []string `json:"required_tools"`
	CaseImages    []string `json:"case_images"`
	UsageGuide    string   `json:"usage_guide"`
	IsFeatured    bool     `json:"is_featured"`
}

type UpdateTemplateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	CoverImage  string   `json:"cover_image"`
	Icon        string   `json:"icon"`
	CaseImages  []string `json:"case_images"`
	UsageGuide  string   `json:"usage_guide"`
	IsFeatured  bool     `json:"is_featured"`
	Status      string   `json:"status"`
}

type ListTemplatesRequest struct {
	Category   string `form:"category"`
	IsFeatured *bool  `form:"is_featured"`
	Search     string `form:"search"`
	Status     string `form:"status"`      // 状态筛选 (active/inactive/published等)
	ShowAll    bool   `form:"show_all"`    // 是否显示所有状态（管理后台使用）
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=100"`
}

type InstallTemplateRequest struct {
	TemplateID   string `json:"template_id" binding:"required"`
	WorkflowName string `json:"workflow_name"`
	EnvVars      map[string]string `json:"env_vars"`
}
