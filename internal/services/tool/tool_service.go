package tool

import (
	"auto-forge/internal/models"
	"auto-forge/internal/repositories"
	"auto-forge/pkg/utools"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

var (
	toolService *ToolService
	once        sync.Once
)

type ToolService struct {
	toolConfigRepo repositories.ToolConfigRepository
}

func GetToolService() *ToolService {
	once.Do(func() {
		toolService = &ToolService{
			toolConfigRepo: repositories.NewToolConfigRepository(),
		}
	})
	return toolService
}

// parseTags 将逗号分隔的 tags 字符串解析为数组
func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}
	tags := strings.Split(tagsStr, ",")
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if trimmed := strings.TrimSpace(tag); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (s *ToolService) ListTools() ([]models.Tool, error) {
	// 从数据库获取已启用且可见的工具配置
	toolConfigs, err := s.toolConfigRepo.FindAvailable()
	if err != nil {
		return nil, fmt.Errorf("获取工具配置失败: %w", err)
	}

	registry := utools.GetRegistry()
	tools := make([]models.Tool, 0, len(toolConfigs))

	for _, config := range toolConfigs {
		// 从注册表获取工具实例
		tool, err := registry.Get(config.ToolCode)
		if err != nil {
			continue
		}

		metadata := tool.GetMetadata()
		schema := tool.GetSchema()
		schemaJSON, _ := json.Marshal(schema)

		toolModel := models.Tool{
			Code:               config.ToolCode,
			Name:               config.ToolName,
			Description:        config.Description,
			Category:           config.Category,
			Version:            config.Version,
			Author:             config.Author,
			Icon:               "",
			ConfigSchema:       string(schemaJSON),
			AICallable:         metadata.AICallable,
			Enabled:            config.Enabled,
			Tags:               parseTags(config.Tags),
			OutputFieldsSchema: metadata.OutputFieldsSchema,
		}

		tools = append(tools, toolModel)
	}

	return tools, nil
}

func (s *ToolService) GetToolByCode(code string) (*models.Tool, error) {
	// 从数据库获取工具配置
	config, err := s.toolConfigRepo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("工具不存在: %s", code)
	}

	// 检查工具是否启用
	if !config.Enabled || config.IsDeprecated {
		return nil, fmt.Errorf("工具不可用: %s", code)
	}

	// 从注册表获取工具实例
	registry := utools.GetRegistry()
	tool, err := registry.Get(code)
	if err != nil {
		return nil, fmt.Errorf("工具未注册: %s", code)
	}

	metadata := tool.GetMetadata()
	schema := tool.GetSchema()
	schemaJSON, _ := json.Marshal(schema)

	toolModel := &models.Tool{
		Code:               config.ToolCode,
		Name:               config.ToolName,
		Description:        config.Description,
		Category:           config.Category,
		Version:            config.Version,
		Author:             config.Author,
		Icon:               "",
		ConfigSchema:       string(schemaJSON),
		AICallable:         metadata.AICallable,
		Enabled:            config.Enabled,
		Tags:               parseTags(config.Tags),
		OutputFieldsSchema: metadata.OutputFieldsSchema,
	}

	return toolModel, nil
}

func (s *ToolService) DescribeOutput(code string, cfg map[string]interface{}) (map[string]utools.OutputFieldDef, error) {
	tool, err := utools.Get(code)
	if err != nil {
		return nil, fmt.Errorf("工具不存在: %s", code)
	}
	if d, ok := tool.(utools.DynamicOutputDescriber); ok {
		schema := d.DescribeOutput(cfg)
		if schema != nil {
			return schema, nil
		}
	}

	md := tool.GetMetadata()
	if md != nil && md.OutputFieldsSchema != nil {
		return md.OutputFieldsSchema, nil
	}
	return map[string]utools.OutputFieldDef{}, nil
}

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

// ToolCategory 工具分类（带工具列表）
type ToolCategory struct {
	Code        string        `json:"code"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	Tools       []models.Tool `json:"tools"`
}

// GetToolsGroupedByCategory 获取按分类分组的工具列表（内置工具在最前面）
func (s *ToolService) GetToolsGroupedByCategory() ([]ToolCategory, error) {
	// 获取所有分类
	categories := utools.GetToolCategories()

	// 获取所有可用工具
	allTools, err := s.ListTools()
	if err != nil {
		return nil, err
	}

	// 创建分类映射
	categoryMap := make(map[string]*ToolCategory)
	result := make([]ToolCategory, 0)

	// 初始化分类
	for _, cat := range categories {
		categoryMap[cat.Code] = &ToolCategory{
			Code:        cat.Code,
			Name:        cat.Name,
			Description: cat.Description,
			Icon:        cat.Icon,
			Tools:       make([]models.Tool, 0),
		}
	}

	// 分配工具到分类
	for _, tool := range allTools {
		if category, exists := categoryMap[tool.Category]; exists {
			category.Tools = append(category.Tools, tool)
		}
	}

	// 按照分类定义的顺序输出，只包含有工具的分类
	for _, cat := range categories {
		if category, exists := categoryMap[cat.Code]; exists && len(category.Tools) > 0 {
			result = append(result, *category)
		}
	}

	return result, nil
}
