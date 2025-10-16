package request

// CreateTemplateCategoryRequest 创建模板分类请求
type CreateTemplateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateTemplateCategoryRequest 更新模板分类请求
type UpdateTemplateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	SortOrder   *int    `json:"sort_order"`
	IsActive    *bool   `json:"is_active"`
}

// TemplateCategoryListQuery 模板分类列表查询参数
type TemplateCategoryListQuery struct {
	Page     int   `form:"page" binding:"omitempty,min=1"`
	PageSize int   `form:"page_size" binding:"omitempty,min=1,max=100"`
	IsActive *bool `form:"is_active"`
}
