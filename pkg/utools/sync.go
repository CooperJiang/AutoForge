package utools

// GetRegistry 导出全局注册表（供 tool_config_service 使用）
func GetRegistryTools() map[string]Tool {
	return GetRegistry().GetAllTools()
}
