package request

// UpdateToolConfigRequest 更新工具配置请求
type UpdateToolConfigRequest struct {
	Config map[string]interface{} `json:"config" binding:"required"`
}

// UpdateToolSettingsRequest 更新工具设置请求
type UpdateToolSettingsRequest struct {
	Enabled   bool `json:"enabled"`
	Visible   bool `json:"visible"`
	SortOrder int  `json:"sort_order"`
}
