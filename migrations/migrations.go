package migrations

import (
	log "auto-forge/pkg/logger"
	"time"

	"gorm.io/gorm"
)

// 注册迁移任务
type migrationTask struct {
	name string
	fn   func(*gorm.DB) error
}

// 注册的迁移列表
var registeredMigrations = []migrationTask{
	{"add_default_categories", AddDefaultCategories}, // 添加默认模板分类
}

// RunAllMigrations 执行所有迁移
func RunAllMigrations(db *gorm.DB) error {
	startTime := time.Now()

	// 确保迁移表存在
	if err := EnsureMigrationTable(db); err != nil {
		log.Error("确保迁移表存在失败: %v", err)
		return err
	}

	executedCount := 0
	skippedCount := 0

	// 执行迁移
	for _, task := range registeredMigrations {
		// 检查迁移是否已执行
		applied, err := IsMigrationApplied(db, task.name)
		if err != nil {
			log.Error("检查迁移状态失败: %v", err)
			continue
		}

		if applied {
			log.Info("迁移 %s 已应用，跳过", task.name)
			skippedCount++
			continue
		}

		// 执行迁移
		log.Info("执行迁移: %s", task.name)
		if err := task.fn(db); err != nil {
			log.Error("迁移 %s 执行失败: %v", task.name, err)
			return err
		}

		// 记录迁移成功状态
		if err := RecordMigration(db, task.name); err != nil {
			log.Error("记录迁移状态失败: %s, 错误: %v", task.name, err)
			return err
		}

		log.Info("迁移 %s 执行成功并已记录", task.name)
		executedCount++
	}

	duration := time.Since(startTime)
	log.Info("迁移执行完成，执行: %d，跳过: %d，耗时: %v", executedCount, skippedCount, duration)
	return nil
}

// GetAllMigrationNames 获取所有迁移的名称
func GetAllMigrationNames() []string {
	names := make([]string, len(registeredMigrations))
	for i, task := range registeredMigrations {
		names[i] = task.name
	}
	return names
}
