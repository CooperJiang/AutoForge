package utools

import (
	"context"
	"encoding/json"
)


type Tool interface {

	GetMetadata() *ToolMetadata


	GetSchema() *ConfigSchema


	Validate(config map[string]interface{}) error


	Execute(ctx *ExecutionContext, config map[string]interface{}) (*ExecutionResult, error)
}


type DynamicOutputDescriber interface {
    DescribeOutput(config map[string]interface{}) map[string]OutputFieldDef
}

type ToolMetadata struct {
	Code               string                    `json:"code"`
	Name               string                    `json:"name"`
	Description        string                    `json:"description"`
	Category           string                    `json:"category"`
	Version            string                    `json:"version"`
	Author             string                    `json:"author"`
	AICallable         bool                      `json:"ai_callable"`
	Tags               []string                  `json:"tags"`
	OutputFieldsSchema map[string]OutputFieldDef `json:"output_fields_schema,omitempty"`
}

type OutputFieldDef struct {
	Type     string                    `json:"type"`
	Label    string                    `json:"label"`
	Children map[string]OutputFieldDef `json:"children,omitempty"`
}


type ConfigSchema struct {
	Type       string                    `json:"type"`
	Properties map[string]PropertySchema `json:"properties"`
	Required   []string                  `json:"required"`
}


type PropertySchema struct {
	Type        string        `json:"type"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Default     interface{}   `json:"default,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
	Format      string        `json:"format,omitempty"`
	Pattern     string        `json:"pattern,omitempty"`
	MinLength   *int          `json:"minLength,omitempty"`
	MaxLength   *int          `json:"maxLength,omitempty"`
	Minimum     *float64      `json:"minimum,omitempty"`
	Maximum     *float64      `json:"maximum,omitempty"`
	Items       *PropertySchema `json:"items,omitempty"`
	Properties  map[string]PropertySchema `json:"properties,omitempty"`
	Secret      bool          `json:"secret,omitempty"`
}


type ExecutionContext struct {
	Context   context.Context        `json:"-"`
	TaskID    string                 `json:"task_id"`
	UserID    string                 `json:"user_id"`
	Variables map[string]interface{} `json:"variables"`
	Metadata  map[string]interface{} `json:"metadata"`
}


type ExecutionResult struct {
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
	Output       map[string]interface{} `json:"output"`
	Error        string                 `json:"error"`
	DurationMs   int64                  `json:"duration_ms"`
	StatusCode   int                    `json:"status_code"`
	ResponseBody string                 `json:"response_body"`
	OutputRender *OutputRenderConfig    `json:"output_render,omitempty"`
}



type OutputRenderConfig struct {
	Type    string                 `json:"type"`
	Primary string                 `json:"primary"`
	Fields  map[string]FieldRender `json:"fields"`
}


type FieldRender struct {
	Type    string `json:"type"`
	Label   string `json:"label"`
	Display bool   `json:"display"`
}


type BaseTool struct {
	metadata *ToolMetadata
	schema   *ConfigSchema
}


func NewBaseTool(metadata *ToolMetadata, schema *ConfigSchema) *BaseTool {
	return &BaseTool{
		metadata: metadata,
		schema:   schema,
	}
}


func (bt *BaseTool) GetMetadata() *ToolMetadata {
	return bt.metadata
}


func (bt *BaseTool) GetSchema() *ConfigSchema {
	return bt.schema
}


func (bt *BaseTool) Validate(config map[string]interface{}) error {

	for _, field := range bt.schema.Required {
		if _, ok := config[field]; !ok {
			return &ValidationError{
				Field:   field,
				Message: "required field missing",
			}
		}
	}


	for key, value := range config {
		propSchema, ok := bt.schema.Properties[key]
		if !ok {
			continue
		}

		if err := bt.validateValue(key, value, propSchema); err != nil {
			return err
		}
	}

	return nil
}


func (bt *BaseTool) validateValue(key string, value interface{}, schema PropertySchema) error {

	switch schema.Type {
	case "string":
		str, ok := value.(string)
		if !ok {
			return &ValidationError{Field: key, Message: "must be a string"}
		}

		if schema.MinLength != nil && len(str) < *schema.MinLength {
			return &ValidationError{Field: key, Message: "string too short"}
		}
		if schema.MaxLength != nil && len(str) > *schema.MaxLength {
			return &ValidationError{Field: key, Message: "string too long"}
		}

	case "number":
		num, ok := value.(float64)
		if !ok {
			return &ValidationError{Field: key, Message: "must be a number"}
		}
		if schema.Minimum != nil && num < *schema.Minimum {
			return &ValidationError{Field: key, Message: "number too small"}
		}
		if schema.Maximum != nil && num > *schema.Maximum {
			return &ValidationError{Field: key, Message: "number too large"}
		}

	case "boolean":
		if _, ok := value.(bool); !ok {
			return &ValidationError{Field: key, Message: "must be a boolean"}
		}

	case "array":
		if _, ok := value.([]interface{}); !ok {
			return &ValidationError{Field: key, Message: "must be an array"}
		}

	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return &ValidationError{Field: key, Message: "must be an object"}
		}
	}


	if len(schema.Enum) > 0 {
		found := false
		for _, enumVal := range schema.Enum {
			if value == enumVal {
				found = true
				break
			}
		}
		if !found {
			return &ValidationError{Field: key, Message: "value not in allowed enum"}
		}
	}

	return nil
}


type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}


type ToolError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ToolError) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}
