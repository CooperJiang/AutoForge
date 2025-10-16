package migrations

import (
	"auto-forge/internal/models"
	log "auto-forge/pkg/logger"

	"gorm.io/gorm"
)

// AddDefaultCategories 添加默认模板分类
func AddDefaultCategories(db *gorm.DB) error {
	// 定义默认分类列表
	defaultCategories := []models.TemplateCategory{
		{
			Name:        "自动化工作流",
			Description: "用于自动化日常任务和流程的工作流模板",
			SortOrder:   10,
			IsActive:    true,
		},
		{
			Name:        "数据处理",
			Description: "数据采集、转换、分析等数据处理相关模板",
			SortOrder:   20,
			IsActive:    true,
		},
		{
			Name:        "通知提醒",
			Description: "消息推送、邮件通知、告警提醒等通知类模板",
			SortOrder:   30,
			IsActive:    true,
		},
		{
			Name:        "系统集成",
			Description: "第三方系统集成、API 对接等集成类模板",
			SortOrder:   40,
			IsActive:    true,
		},
		{
			Name:        "监控告警",
			Description: "系统监控、性能检测、异常告警等监控类模板",
			SortOrder:   50,
			IsActive:    true,
		},
		{
			Name:        "定时任务",
			Description: "定期执行的任务，如备份、清理、报表生成等",
			SortOrder:   60,
			IsActive:    true,
		},
		{
			Name:        "DevOps",
			Description: "CI/CD、部署、测试等 DevOps 相关模板",
			SortOrder:   70,
			IsActive:    true,
		},
		{
			Name:        "办公协同",
			Description: "审批流程、文档处理、会议管理等办公协同模板",
			SortOrder:   80,
			IsActive:    true,
		},
		{
			Name:        "其他",
			Description: "其他未分类的工作流模板",
			SortOrder:   999,
			IsActive:    true,
		},
	}

	// 批量插入，如果已存在则跳过
	for _, category := range defaultCategories {
		var count int64
		if err := db.Model(&models.TemplateCategory{}).
			Where("name = ?", category.Name).
			Count(&count).Error; err != nil {
			log.Error("检查分类是否存在失败: %v", err)
			continue
		}

		if count > 0 {
			log.Info("分类 %s 已存在，跳过", category.Name)
			continue
		}

		if err := db.Create(&category).Error; err != nil {
			log.Error("创建默认分类 %s 失败: %v", category.Name, err)
			return err
		}
		log.Info("创建默认分类: %s", category.Name)
	}

	log.Info("默认分类添加完成")
	return nil
}
