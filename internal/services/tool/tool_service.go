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


type ToolService struct{}


func GetToolService() *ToolService {
	once.Do(func() {
		toolService = &ToolService{}
	})
	return toolService
}


func InitToolService() {
	GetToolService()
}


func (s *ToolService) ListTools() ([]models.Tool, error) {
	registry := utools.GetRegistry()
	metadataList := registry.List()

	tools := make([]models.Tool, 0, len(metadataList))

	for _, metadata := range metadataList {

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
            Icon:         "",
            ConfigSchema: string(schemaJSON),
            AICallable:   metadata.AICallable,
            Enabled:      true,
            Tags:         string(tagsJSON),
            OutputFieldsSchema: metadata.OutputFieldsSchema,
        }

		tools = append(tools, toolModel)
	}

	return tools, nil
}


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
		Icon:         "",
		ConfigSchema: string(schemaJSON),
		AICallable:   metadata.AICallable,
		Enabled:      true,
		Tags:         string(tagsJSON),
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
