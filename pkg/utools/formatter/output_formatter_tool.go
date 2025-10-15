package formatter

import (
	"auto-forge/pkg/utools"
	"fmt"
	"time"
)

// OutputFormatterTool 输出格式化工具
// 用于将上��节点的输出数据格式化为特定的展示类型
type OutputFormatterTool struct {
	*utools.BaseTool
}

// NewOutputFormatterTool 创建输出格式化工具实例
func NewOutputFormatterTool() *OutputFormatterTool {
	metadata := &utools.ToolMetadata{
		Code:        "output_formatter",
		Name:        "输出格式化",
		Description: "将数据格式化为指定的展示类型（图片、视频、HTML、Markdown 等），用于控制最终输出的显示方式",
		Category:    "utility",
		Version:     "1.0.0",
		Author:      "AutoForge",
		AICallable:  false, // 通常不被 AI 直接调用
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
					"image",    // 单张图片
					"video",    // 视频
					"html",     // HTML 内容
					"html-url", // HTML URL 预览
					"markdown", // Markdown 文本
					"text",     // 纯文本
					"gallery",  // 图片画廊（多张图片）
					"json",     // JSON 数据（默认）
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
				Description: "根据类型不同，内容含义不同：\n- image: 图片 URL（支持变量，如 {{nodes.xxx.response.data[0].url}}）\n- video: 视频 URL\n- html: HTML 字符串\n- html-url: HTML 页面的 URL 地址\n- markdown: Markdown 文本\n- text: 纯文本\n- gallery: 图片 URL 数组（JSON 字符串）\n- json: JSON 数据",
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

// Execute 执行输出格式化
func (t *OutputFormatterTool) Execute(ctx *utools.ExecutionContext, toolConfig map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()

	// 获取配置 - 注意：这里的 toolConfig 已经是替换过变量的了
	outputType, _ := toolConfig["output_type"].(string)
	if outputType == "" {
		outputType = "json"
	}

	title, _ := toolConfig["title"].(string)
	description, _ := toolConfig["description"].(string)
	content, _ := toolConfig["content"].(string)
	altText, _ := toolConfig["alt_text"].(string)
	thumbnail, _ := toolConfig["thumbnail"].(string)
	metadataStr, _ := toolConfig["metadata"].(string)

	// 验证必填字段
	if content == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "content 字段不能为空",
			Error:      "missing content",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("content is required")
	}

	// Debug: 打印 content 的值
	fmt.Printf("[OutputFormatter] Content value: %s\n", content)

	// 确定实际渲染类型
	renderType := outputType
	if outputType == "html-url" {
		renderType = "url" // html-url 使用 url 渲染器（UrlViewer）
	}

	// 构建输出渲染配置
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

	// 根据输出类型添加相应的字段配置
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
			Display: false, // 不直接显示，用于图片 alt 属性
		}
	}

	if thumbnail != "" {
		outputRender.Fields["thumbnail"] = utools.FieldRender{
			Type:    "image",
			Label:   "缩略图",
			Display: outputType == "video", // 视频类型时显示缩略图
		}
	}

	// 构建输出数据
	output := map[string]interface{}{
		"content": content,
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

	return &utools.ExecutionResult{
		Success:      true,
		Message:      fmt.Sprintf("输出格式化成功，类型：%s", outputType),
		Output:       output,
		OutputRender: outputRender,
		DurationMs:   time.Since(startTime).Milliseconds(),
	}, nil
}

// init 自动注册工具
func init() {
	tool := NewOutputFormatterTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register Output Formatter tool: %v", err))
	}
}
