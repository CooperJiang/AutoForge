package response

import "auto-forge/internal/models"

// ToolConfigDetailResponse 工具配置详情响应
type ToolConfigDetailResponse struct {
	*models.ToolConfig
	DecryptedConfig map[string]interface{} `json:"decrypted_config"`
}
