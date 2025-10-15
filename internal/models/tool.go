package models

import "auto-forge/pkg/utools"

// Tool 工具模型（用于API响应，不存储在数据库）
type Tool struct {
    Code               string                                `json:"code"`
    Name               string                                `json:"name"`
    Description        string                                `json:"description"`
    Category           string                                `json:"category"`
    Version            string                                `json:"version"`
    Author             string                                `json:"author"`
    Icon               string                                `json:"icon"`
    ConfigSchema       string                                `json:"config_schema"`
    AICallable         bool                                  `json:"ai_callable"`
    Enabled            bool                                  `json:"enabled"`
    Tags               string                                `json:"tags"`
    OutputFieldsSchema map[string]utools.OutputFieldDef      `json:"output_fields_schema,omitempty"`
}
