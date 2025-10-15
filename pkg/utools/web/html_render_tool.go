package web

import (
	"auto-forge/pkg/utools"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// HtmlRenderTool HTML 渲染工具
// 将 HTML 内容保存为文件并生成访问 URL，或者直接渲染传入的 URL
type HtmlRenderTool struct {
	*utools.BaseTool
}

// NewHtmlRenderTool 创建 HTML 渲染工具实例
func NewHtmlRenderTool() *HtmlRenderTool {
	metadata := &utools.ToolMetadata{
		Code:        "html_render",
		Name:        "HTML 内容保存",
		Description: "将 HTML 内容保存为网页文件并生成可访问的 URL，支持分享和预览",
		Category:    "utility",
		Version:     "1.0.0",
		Author:      "AutoForge",
		AICallable:  true,
		Tags:        []string{"html", "save", "web", "preview", "share"},
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"url": {
				Type:  "string",
				Label: "💡 预览 URL 地址",
			},
			"title": {
				Type:  "string",
				Label: "页面标题",
			},
			"content_length": {
				Type:  "number",
				Label: "HTML 内容长度",
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
			"content": {
				Type:        "string",
				Title:       "HTML 内容",
				Description: "要保存的 HTML 内容",
			},
			"title": {
				Type:        "string",
				Title:       "页面标题",
				Description: "网页标题（可选）",
			},
			"expires_hours": {
				Type:        "number",
				Title:       "过期时间（小时）",
				Description: "生成的预览链接过期时间，0 表示永不过期",
				Default:     float64(0),
			},
		},
		Required: []string{"content"},
	}

	return &HtmlRenderTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

// Execute 执行 HTML 内容保存
func (t *HtmlRenderTool) Execute(ctx *utools.ExecutionContext, toolConfig map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()

	// 获取配置
	content, _ := toolConfig["content"].(string)
	if content == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "HTML 内容不能为空",
			Error:      "missing content",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("content is required")
	}

	// 清理内容：移除 AI 思考标签和 markdown 代码块标记
	content = cleanHtmlContent(content)

	title, _ := toolConfig["title"].(string)
	expiresHours := float64(0)
	if h, ok := toolConfig["expires_hours"].(float64); ok {
		expiresHours = h
	}

	// 生成唯一 ID（基于内容的 MD5）
	hash := md5.Sum([]byte(fmt.Sprintf("%s-%d", content, time.Now().UnixNano())))
	fileID := hex.EncodeToString(hash[:])[:16]

	// 确保预览目录存在
	previewDir := "./data/html-preview"
	if err := os.MkdirAll(previewDir, 0755); err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "创建预览目录失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 保存 HTML 文件
	filePath := filepath.Join(previewDir, fileID+".html")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "保存 HTML 文件失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 生成访问 URL
	renderUrl := fmt.Sprintf("/preview/%s", fileID)

	// 计算过期时间
	var expiresAt *int64
	if expiresHours > 0 {
		expires := time.Now().Add(time.Duration(expiresHours) * time.Hour).Unix()
		expiresAt = &expires
	}

	// 构建输出
	output := map[string]interface{}{
		"url":            renderUrl,
		"title":          title,
		"content_length": len(content),
	}

	if expiresAt != nil {
		output["expires_at"] = *expiresAt
	}

	// 构建输出渲染配置
	outputRender := &utools.OutputRenderConfig{
		Type:    "url",
		Primary: "url",
		Fields: map[string]utools.FieldRender{
			"url": {
				Type:    "url",
				Label:   "预览地址",
				Display: true,
			},
		},
	}

	return &utools.ExecutionResult{
		Success:      true,
		Message:      fmt.Sprintf("HTML 内容保存成功，访问地址：%s", renderUrl),
		Output:       output,
		OutputRender: outputRender,
		DurationMs:   time.Since(startTime).Milliseconds(),
	}, nil
}

// cleanHtmlContent 清理 HTML 内容
// 移除 AI 思考标签、markdown 代码块标记等
func cleanHtmlContent(content string) string {
	// 1. 移除 <think>...</think> 标签及其内容
	thinkRegex := regexp.MustCompile(`(?s)<think>.*?</think>`)
	content = thinkRegex.ReplaceAllString(content, "")

	// 2. 移除 markdown 代码块标记（```html 或 ``` 开头/结尾）
	// 移除开头的代码块标记
	content = regexp.MustCompile(`(?m)^` + "```" + `html?\s*\n`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`(?m)^` + "```" + `\s*\n`).ReplaceAllString(content, "")
	// 移除结尾的代码块标记
	content = regexp.MustCompile(`(?m)\n` + "```" + `\s*$`).ReplaceAllString(content, "")

	// 3. 移除开头和结尾的空白字符
	content = strings.TrimSpace(content)

	return content
}

// init 自动注册工具
func init() {
	tool := NewHtmlRenderTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register HTML Render tool: %v", err))
	}
}
