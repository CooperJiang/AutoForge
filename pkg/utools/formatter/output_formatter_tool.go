package formatter

import (
	"auto-forge/pkg/utools"
	"encoding/json"
	"fmt"
	"time"
)



type OutputFormatterTool struct {
    *utools.BaseTool
}


func NewOutputFormatterTool() *OutputFormatterTool {
	metadata := &utools.ToolMetadata{
		Code:        "output_formatter",
		Name:        "输出格式化",
		Description: "将数据格式化为指定的展示类型（图片、视频、HTML、Markdown 等），用于控制最终输出的显示方式",
		Category:    utools.CategoryUtility,
		Version:     "1.0.0",
		Author:      "AutoForge",
		AICallable:  false,
		Tags:        []string{"formatter", "output", "display", "render"},
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"content": {
				Type:  "string",
				Label: "格式化后的内容",
			},
			"type": {
				Type:  "string",
				Label: "输出类型",
			},
			"title": {
				Type:  "string",
				Label: "标题",
			},
			"description": {
				Type:  "string",
				Label: "描述",
			},
			"message": {
				Type:  "string",
				Label: "执行消息",
			},
		},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"output_type": {
				Type:        "string",
				Title:       "输出类型",
				Description: "指定输出的展示类型",
				Default:     "json",
				Enum: []interface{}{
					"image",
					"video",
					"html",
					"html-url",
					"markdown",
					"text",
					"gallery",
					"json",
				},
			},
			"title": {
				Type:        "string",
				Title:       "标题",
				Description: "输出内容的标题（可选）",
			},
			"description": {
				Type:        "string",
				Title:       "描述",
				Description: "输出内容的描述信息（可选）",
			},
			"content": {
				Type:        "string",
				Title:       "主要内容",
				Description: "根据类型不同，内容含义不同：\n- image: 图片 URL（支持变量，如 {{nodes.xxx.response.data[0].url}}）\n- video: 视频 URL\n- html: HTML 字符串\n- html-url: HTML 页面的 URL 地址\n- markdown: Markdown 文本\n- text: 纯文本\n- gallery: 图片 URL 数组（JSON 字符串）\n- json: JSON 对象（支持变量引用和 JSON 字符串）",
			},
			"alt_text": {
				Type:        "string",
				Title:       "替代文本",
				Description: "图片或视频的替代文本描述（可选，支持变量）",
			},
			"thumbnail": {
				Type:        "string",
				Title:       "缩略图",
				Description: "缩略图 URL（可选，用于视频等）",
			},
			"metadata": {
				Type:        "string",
				Title:       "元数据",
				Description: "额外的元数据信息（JSON 字符串）",
			},
		},
		Required: []string{"output_type", "content"},
	}

	return &OutputFormatterTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}


func (t *OutputFormatterTool) Execute(ctx *utools.ExecutionContext, toolConfig map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()


	outputType, _ := toolConfig["output_type"].(string)
	if outputType == "" {
		outputType = "json"
	}

	title, _ := toolConfig["title"].(string)
	description, _ := toolConfig["description"].(string)
	altText, _ := toolConfig["alt_text"].(string)
	thumbnail, _ := toolConfig["thumbnail"].(string)
	metadataStr, _ := toolConfig["metadata"].(string)


	var contentValue interface{}
	if contentStr, ok := toolConfig["content"].(string); ok {

		if contentStr == "" {
			return &utools.ExecutionResult{
				Success:    false,
				Message:    "content 字段不能为空",
				Error:      "missing content",
				DurationMs: time.Since(startTime).Milliseconds(),
			}, fmt.Errorf("content is required")
		}


		if outputType == "json" {
			var jsonObj interface{}
			if err := json.Unmarshal([]byte(contentStr), &jsonObj); err == nil {

				contentValue = jsonObj
			} else {

				contentValue = contentStr
			}
		} else {
			contentValue = contentStr
		}
	} else if contentObj := toolConfig["content"]; contentObj != nil {

		contentValue = contentObj
	} else {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "content 字段不能为空",
			Error:      "missing content",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("content is required")
	}


	renderType := outputType
	if outputType == "html-url" {
		renderType = "url"
	}


	outputRender := &utools.OutputRenderConfig{
		Type:    renderType,
		Primary: "content",
		Fields: map[string]utools.FieldRender{
			"content": {
				Type:    renderType,
				Label:   title,
				Display: true,
			},
		},
	}


	if description != "" {
		outputRender.Fields["description"] = utools.FieldRender{
			Type:    "text",
			Label:   "描述",
			Display: true,
		}
	}

	if altText != "" {
		outputRender.Fields["alt_text"] = utools.FieldRender{
			Type:    "text",
			Label:   "替代文本",
			Display: false,
		}
	}

	if thumbnail != "" {
		outputRender.Fields["thumbnail"] = utools.FieldRender{
			Type:    "image",
			Label:   "缩略图",
			Display: outputType == "video",
		}
	}


	var output map[string]interface{}


	if outputType == "json" {
		if contentObj, ok := contentValue.(map[string]interface{}); ok {

			output = contentObj
		} else {

			output = map[string]interface{}{
				"content": contentValue,
				"type":    outputType,
			}
		}
	} else {

		output = map[string]interface{}{
			"content": contentValue,
			"type":    outputType,
		}

		if title != "" {
			output["title"] = title
		}
		if description != "" {
			output["description"] = description
		}
		if altText != "" {
			output["alt_text"] = altText
		}
		if thumbnail != "" {
			output["thumbnail"] = thumbnail
		}
		if metadataStr != "" {
			output["metadata"] = metadataStr
		}
	}

	return &utools.ExecutionResult{
		Success:      true,
		Message:      fmt.Sprintf("输出格式化成功，类型：%s", outputType),
		Output:       output,
		OutputRender: outputRender,
		DurationMs:   time.Since(startTime).Milliseconds(),
	}, nil
}


func (t *OutputFormatterTool) DescribeOutput(config map[string]interface{}) map[string]utools.OutputFieldDef {

    return map[string]utools.OutputFieldDef{
        "content":     {Type: "string", Label: "格式化后的内容/主内容"},
        "type":        {Type: "string", Label: "输出类型"},
        "title":       {Type: "string", Label: "标题"},
        "description": {Type: "string", Label: "描述"},
    }
}


func init() {
	tool := NewOutputFormatterTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register Output Formatter tool: %v", err))
	}
}
