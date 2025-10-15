package tool

import (
	"auto-forge/internal/models"
	"auto-forge/pkg/utools"
	"encoding/json"
	"fmt"
	"sync"
)

var (
	toolService *ToolService
	once        sync.Once
)

// ToolService 工具服务
type ToolService struct{}

// GetToolService 获取工具服务单例
func GetToolService() *ToolService {
	once.Do(func() {
		toolService = &ToolService{}
	})
	return toolService
}

// InitToolService 初始化工具服务
func InitToolService() {
	GetToolService()
}

// ListTools 获取所有工具列表（从代码注册表读取）
func (s *ToolService) ListTools() ([]models.Tool, error) {
	registry := utools.GetRegistry()
	metadataList := registry.List()

	tools := make([]models.Tool, 0, len(metadataList))

	for _, metadata := range metadataList {
		// 获取工具的 Schema
		tool, err := registry.Get(metadata.Code)
		if err != nil {
			continue
		}

		schema := tool.GetSchema()
		schemaJSON, _ := json.Marshal(schema)
		tagsJSON, _ := json.Marshal(metadata.Tags)

		toolModel := models.Tool{
			Code:         metadata.Code,
			Name:         metadata.Name,
			Description:  metadata.Description,
			Category:     metadata.Category,
			Version:      metadata.Version,
			Author:       metadata.Author,
			Icon:         "", // Icon is now configured in frontend
			ConfigSchema: string(schemaJSON),
			AICallable:   metadata.AICallable,
			Enabled:      true,
			Tags:         string(tagsJSON),
		}

		tools = append(tools, toolModel)
	}

	return tools, nil
}

// GetToolByCode 根据 code 获取工具详情（从代码注册表读取）
func (s *ToolService) GetToolByCode(code string) (*models.Tool, error) {
	registry := utools.GetRegistry()
	tool, err := registry.Get(code)
	if err != nil {
		return nil, fmt.Errorf("工具不存在: %s", code)
	}

	metadata := tool.GetMetadata()
	schema := tool.GetSchema()
	schemaJSON, _ := json.Marshal(schema)
	tagsJSON, _ := json.Marshal(metadata.Tags)

	toolModel := &models.Tool{
		Code:         metadata.Code,
		Name:         metadata.Name,
		Description:  metadata.Description,
		Category:     metadata.Category,
		Version:      metadata.Version,
		Author:       metadata.Author,
		Icon:         "", // Icon is now configured in frontend
		ConfigSchema: string(schemaJSON),
		AICallable:   metadata.AICallable,
		Enabled:      true,
		Tags:         string(tagsJSON),
	}

	return toolModel, nil
}

// GetToolsByCategory 根据分类获取工具（从代码注册表读取）
func (s *ToolService) GetToolsByCategory(category string) ([]models.Tool, error) {
	allTools, err := s.ListTools()
	if err != nil {
		return nil, err
	}

	tools := make([]models.Tool, 0)
	for _, tool := range allTools {
		if tool.Category == category {
			tools = append(tools, tool)
		}
	}

	return tools, nil
}

// GetAICallableTools 获取所有可被 AI 调用的工具（从代码注册表读取）
func (s *ToolService) GetAICallableTools() ([]models.Tool, error) {
	allTools, err := s.ListTools()
	if err != nil {
		return nil, err
	}

	tools := make([]models.Tool, 0)
	for _, tool := range allTools {
		if tool.AICallable {
			tools = append(tools, tool)
		}
	}

	return tools, nil
}
