package cron

import (
	log "auto-forge/pkg/logger"
	"fmt"
	"os"
	"path/filepath"
)

// CleanupExecutionFiles 清理指定执行记录的临时文件（工作流执行完成后调用）
func CleanupExecutionFiles(executionID string) error {
	baseDir := "/tmp/workflow-files"
	dirPath := filepath.Join(baseDir, executionID)

	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil // 文件不存在，不需要清理
	}

	// 删除整个执行目录
	if err := os.RemoveAll(dirPath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	log.Info("已清理执行记录的临时文件: %s", executionID)
	return nil
}
