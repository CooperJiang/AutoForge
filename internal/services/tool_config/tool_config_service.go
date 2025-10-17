package tool_config

import (
	"auto-forge/internal/models"
	"auto-forge/internal/repositories"
	log "auto-forge/pkg/logger"
	"auto-forge/pkg/utils"
	"auto-forge/pkg/utools"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ToolConfigService 工具配置服务接口
type ToolConfigService interface {
	GetAvailableTools() ([]*models.ToolConfig, error)
	GetAllTools() ([]*models.ToolConfig, error)
	GetToolConfig(toolCode string) (*models.ToolConfig, error)
	GetToolConfigDecrypted(toolCode string) (map[string]interface{}, error)
	UpdateToolConfig(toolCode string, configMap map[string]interface{}) error
	UpdateToolSettings(toolCode string, enabled, visible bool, sortOrder int) error
	DeleteTool(id uint) error
	SyncToolsFromRegistry() error
}

type toolConfigService struct {
	repo repositories.ToolConfigRepository
}

// NewToolConfigService 创建工具配置服务实例
func NewToolConfigService() ToolConfigService {
	return &toolConfigService{
		repo: repositories.NewToolConfigRepository(),
	}
}

// GetAvailableTools 获取可用的工具列表（enabled && visible && !is_deprecated）
func (s *toolConfigService) GetAvailableTools() ([]*models.ToolConfig, error) {
	return s.repo.FindAvailable()
}

// GetAllTools 获取所有工具列表
func (s *toolConfigService) GetAllTools() ([]*models.ToolConfig, error) {
	return s.repo.FindAll()
}

// GetToolConfig 获取工具配置（加密的）
func (s *toolConfigService) GetToolConfig(toolCode string) (*models.ToolConfig, error) {
	return s.repo.FindByCode(toolCode)
}

// GetToolConfigDecrypted 获取解密后的工具配置
func (s *toolConfigService) GetToolConfigDecrypted(toolCode string) (map[string]interface{}, error) {
	toolConfig, err := s.repo.FindByCode(toolCode)
	if err != nil {
		return nil, err
	}

	// 解密配置
	return utils.DecryptToolConfig(toolConfig.ConfigJSON)
}

// UpdateToolConfig 更新工具配置
func (s *toolConfigService) UpdateToolConfig(toolCode string, configMap map[string]interface{}) error {
	toolConfig, err := s.repo.FindByCode(toolCode)
	if err != nil {
		return err
	}

	// 加密配置
	encryptedConfig, err := utils.EncryptToolConfig(configMap)
	if err != nil {
		return err
	}

	toolConfig.ConfigJSON = encryptedConfig
	toolConfig.UpdatedAt = time.Now()

	return s.repo.Update(toolConfig)
}

// UpdateToolSettings 更新工具设置（启用/禁用、可见性、排序）
func (s *toolConfigService) UpdateToolSettings(toolCode string, enabled, visible bool, sortOrder int) error {
	toolConfig, err := s.repo.FindByCode(toolCode)
	if err != nil {
		return err
	}

	toolConfig.Enabled = enabled
	toolConfig.Visible = visible
	toolConfig.SortOrder = sortOrder
	toolConfig.UpdatedAt = time.Now()

	return s.repo.Update(toolConfig)
}

// DeleteTool 删除工具配置
func (s *toolConfigService) DeleteTool(id uint) error {
	return s.repo.Delete(id)
}

// SyncToolsFromRegistry 从工具注册表同步工具定义
func (s *toolConfigService) SyncToolsFromRegistry() error {
	// 获取所有已注册的工具
	registry := utools.GetRegistryTools()

	for toolCode, tool := range registry {
		metadata := tool.GetMetadata()
		schema := tool.GetSchema()

		// 检查数据库中是否已存在
		exists, err := s.repo.ExistsByCode(toolCode)
		if err != nil {
			log.Error("检查工具是否存在失败: %v", err)
			continue
		}

		if !exists {
			// 首次注册，创建默认配置
			configSchemaJSON, _ := json.Marshal(schema)
			tags := strings.Join(metadata.Tags, ",")

			toolConfig := &models.ToolConfig{
				ToolCode:     metadata.Code,
				ToolName:     metadata.Name,
				Enabled:      true, // 默认启用
				Visible:      true,
				IsDeprecated: false,
				ConfigJSON:   "", // 空配置
				ConfigSchema: string(configSchemaJSON),
				Description:  metadata.Description,
				Category:     metadata.Category,
				Version:      metadata.Version,
				Author:       metadata.Author,
				Tags:         tags,
				SortOrder:    0,
				LastSyncAt:   time.Now(),
			}

			if err := s.repo.Create(toolConfig); err != nil {
				log.Error("创建工具配置失败: %v", err)
				continue
			}

			log.Info("工具 [%s] 已同步到数据库（默认启用）", metadata.Name)
		} else {
			// 已存在，更新元数据
			toolConfig, err := s.repo.FindByCode(toolCode)
			if err != nil {
				log.Error("获取工具配置失败: %v", err)
				continue
			}

			// 更新元数据（不影响配置和启用状态）
			configSchemaJSON, _ := json.Marshal(schema)
			tags := strings.Join(metadata.Tags, ",")

			toolConfig.ToolName = metadata.Name
			toolConfig.ConfigSchema = string(configSchemaJSON)
			toolConfig.Description = metadata.Description
			toolConfig.Category = metadata.Category
			toolConfig.Version = metadata.Version
			toolConfig.Author = metadata.Author
			toolConfig.Tags = tags
			toolConfig.LastSyncAt = time.Now()

			// 如果工具之前被标记为废弃，现在恢复了
			if toolConfig.IsDeprecated {
				toolConfig.IsDeprecated = false
				log.Info("工具 [%s] 已恢复", metadata.Name)
			}

			if err := s.repo.Update(toolConfig); err != nil {
				log.Error("更新工具配置失败: %v", err)
				continue
			}
		}
	}

	// 检查数据库中有但代码中没有的工具（已删除）
	allToolConfigs, err := s.repo.FindAll()
	if err != nil {
		return err
	}

	for _, toolConfig := range allToolConfigs {
		if _, exists := registry[toolConfig.ToolCode]; !exists {
			// 代码中已删除，标记为废弃
			if !toolConfig.IsDeprecated {
				toolConfig.IsDeprecated = true
				toolConfig.Enabled = false // 自动禁用
				toolConfig.UpdatedAt = time.Now()

				if err := s.repo.Update(toolConfig); err != nil {
					log.Error("标记工具为废弃失败: %v", err)
					continue
				}

				log.Warn("工具 [%s] 在代码中已删除，已标记为废弃并自动禁用", toolConfig.ToolName)
			}
		}
	}

	return nil
}

// GetToolConfigForExecution 获取工具执行所需的配置（供工具内部使用）
func GetToolConfigForExecution(toolCode string) (map[string]interface{}, error) {
	service := NewToolConfigService()

	// 获取工具配置
	toolConfig, err := service.GetToolConfig(toolCode)
	if err != nil {
		return nil, fmt.Errorf("工具配置不存在: %s", toolCode)
	}

	// 检查是否启用
	if !toolConfig.Enabled {
		return nil, fmt.Errorf("工具未启用: %s", toolConfig.ToolName)
	}

	// 检查是否废弃
	if toolConfig.IsDeprecated {
		return nil, fmt.Errorf("工具已废弃: %s", toolConfig.ToolName)
	}

	// 解密并返回配置
	return service.GetToolConfigDecrypted(toolCode)
}
