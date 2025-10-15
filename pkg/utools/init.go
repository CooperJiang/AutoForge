package utools

import (
	"auto-forge/pkg/logger"
)


func InitTools() {
	count := GetRegistry().Count()
	logger.Info("Tool system initialized with %d tools", count)

	tools := GetRegistry().List()
	for _, tool := range tools {
		logger.Info("  - [%s] %s (v%s) - %s", tool.Code, tool.Name, tool.Version, tool.Category)
	}
}
