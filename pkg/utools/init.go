package utools

import (
	"auto-forge/pkg/logger"
)

// InitTools 初始化工具系统
func InitTools() {
	count := GetRegistry().Count()
	logger.Info("Tool system initialized with %d tools", count)

	tools := GetRegistry().List()
	for _, tool := range tools {
		logger.Info("  - [%s] %s (v%s) - %s", tool.Code, tool.Name, tool.Version, tool.Category)
	}
}
