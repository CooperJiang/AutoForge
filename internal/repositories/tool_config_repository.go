package repositories

import (
	"auto-forge/internal/models"
	"auto-forge/pkg/database"

	"gorm.io/gorm"
)

// ToolConfigRepository 工具配置仓储接口
type ToolConfigRepository interface {
	Create(toolConfig *models.ToolConfig) error
	Update(toolConfig *models.ToolConfig) error
	FindByCode(toolCode string) (*models.ToolConfig, error)
	FindAll() ([]*models.ToolConfig, error)
	FindAvailable() ([]*models.ToolConfig, error)
	ExistsByCode(toolCode string) (bool, error)
	Delete(id uint) error
}

type toolConfigRepository struct {
	db *gorm.DB
}

// NewToolConfigRepository 创建工具配置仓储实例
func NewToolConfigRepository() ToolConfigRepository {
	return &toolConfigRepository{db: database.GetDB()}
}

// Create 创建工具配置
func (r *toolConfigRepository) Create(toolConfig *models.ToolConfig) error {
	return r.db.Create(toolConfig).Error
}

// Update 更新工具配置
func (r *toolConfigRepository) Update(toolConfig *models.ToolConfig) error {
	return r.db.Save(toolConfig).Error
}

// FindByCode 根据工具代码查找
func (r *toolConfigRepository) FindByCode(toolCode string) (*models.ToolConfig, error) {
	var toolConfig models.ToolConfig
	err := r.db.Where("tool_code = ?", toolCode).First(&toolConfig).Error
	if err != nil {
		return nil, err
	}
	return &toolConfig, nil
}

// FindAll 查找所有工具配置
func (r *toolConfigRepository) FindAll() ([]*models.ToolConfig, error) {
	var toolConfigs []*models.ToolConfig
	err := r.db.Order("sort_order ASC, created_at ASC").Find(&toolConfigs).Error
	return toolConfigs, err
}

// FindAvailable 查找可用的工具（启用且可见且未废弃）
func (r *toolConfigRepository) FindAvailable() ([]*models.ToolConfig, error) {
	var toolConfigs []*models.ToolConfig
	err := r.db.Where("enabled = ? AND visible = ? AND is_deprecated = ?", true, true, false).
		Order("sort_order ASC, created_at ASC").
		Find(&toolConfigs).Error
	return toolConfigs, err
}

// ExistsByCode 检查工具是否存在
func (r *toolConfigRepository) ExistsByCode(toolCode string) (bool, error) {
	var count int64
	err := r.db.Model(&models.ToolConfig{}).Where("tool_code = ?", toolCode).Count(&count).Error
	return count > 0, err
}

// Delete 删除工具配置
func (r *toolConfigRepository) Delete(id uint) error {
	return r.db.Delete(&models.ToolConfig{}, id).Error
}
