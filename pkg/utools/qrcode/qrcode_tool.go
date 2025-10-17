package qrcode

import (
	"auto-forge/pkg/utools"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/skip2/go-qrcode"
)

type QRCodeTool struct {
	*utools.BaseTool
}

func NewQRCodeTool() *QRCodeTool {
	metadata := &utools.ToolMetadata{
		Code:        "qrcode_generator",
		Name:        "二维码生成",
		Description: "生成二维码图片，支持自定义尺寸、颜色和错误纠正级别",
		Category:    utools.CategoryUtility,
		Version:     "1.0.0",
		Author:      "Cooper Team",
		AICallable:  true,
		Tags:        []string{"qrcode", "image", "generator", "utility"},
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"response": {
				Type:  "object",
				Label: "完整响应",
				Children: map[string]utools.OutputFieldDef{
					"content": {
						Type:  "string",
						Label: "二维码内容",
					},
					"format": {
						Type:  "string",
						Label: "输出格式",
					},
					"size": {
						Type:  "integer",
						Label: "图片尺寸",
					},
					"level": {
						Type:  "string",
						Label: "错误纠正级别",
					},
					"data": {
						Type:  "string",
						Label: "Base64 数据（仅 base64 模式）",
					},
					"file": {
						Type:  "object",
						Label: "文件对象（仅 file 模式）",
						Children: map[string]utools.OutputFieldDef{
							"path": {
								Type:  "string",
								Label: "文件路径",
							},
							"filename": {
								Type:  "string",
								Label: "文件名",
							},
							"size": {
								Type:  "integer",
								Label: "文件大小",
							},
							"mime_type": {
								Type:  "string",
								Label: "MIME 类型",
							},
						},
					},
					"created_at": {
						Type:  "string",
						Label: "创建时间",
					},
				},
			},
			"data": {
				Type:  "string",
				Label: "Base64 数据（快捷访问，仅 base64 模式）",
			},
			"file": {
				Type:  "object",
				Label: "文件对象（快捷访问，仅 file 模式）",
				Children: map[string]utools.OutputFieldDef{
					"path": {
						Type:  "string",
						Label: "文件路径",
					},
					"filename": {
						Type:  "string",
						Label: "文件名",
					},
					"size": {
						Type:  "integer",
						Label: "文件大小",
					},
					"mime_type": {
						Type:  "string",
						Label: "MIME 类型",
					},
				},
			},
			"content": {
				Type:  "string",
				Label: "原始内容（快捷访问）",
			},
		},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"content": {
				Type:        "string",
				Title:       "二维码内容",
				Description: "要编码的文本内容（URL、文本等），支持变量",
			},
			"size": {
				Type:        "integer",
				Title:       "图片尺寸",
				Description: "二维码图片尺寸（像素），范围：64-2048",
				Default:     256,
				Minimum:     float64Ptr(64),
				Maximum:     float64Ptr(2048),
			},
			"level": {
				Type:        "string",
				Title:       "错误纠正级别",
				Description: "二维码容错能力，越高越能容错但尺寸越大",
				Default:     "Medium",
				Enum:        []interface{}{"Low", "Medium", "High", "Highest"},
			},
			"output_format": {
				Type:        "string",
				Title:       "输出格式",
				Description: "base64: 返回 Base64 编码字符串; file: 返回文件对象",
				Default:     "base64",
				Enum:        []interface{}{"base64", "file"},
			},
		},
		Required: []string{"content"},
	}

	return &QRCodeTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

func (t *QRCodeTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()

	// 1. 解析配置
	content, ok := config["content"].(string)
	if !ok || content == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "二维码内容不能为空",
			Error:      "content is required",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("二维码内容不能为空")
	}

	// 尺寸
	size := 256
	if sizeVal, ok := config["size"]; ok {
		switch v := sizeVal.(type) {
		case float64:
			size = int(v)
		case int:
			size = v
		case string:
			if parsed, err := strconv.Atoi(v); err == nil {
				size = parsed
			}
		}
	}
	if size < 64 {
		size = 64
	}
	if size > 2048 {
		size = 2048
	}

	// 错误纠正级别
	levelStr := "Medium"
	if level, ok := config["level"].(string); ok && level != "" {
		levelStr = level
	}
	recoveryLevel := t.parseRecoveryLevel(levelStr)

	// 输出格式
	outputFormat := "base64"
	if format, ok := config["output_format"].(string); ok && format != "" {
		outputFormat = format
	}

	// 2. 生成二维码
	qr, err := qrcode.New(content, recoveryLevel)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "生成二维码失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("生成二维码失败: %w", err)
	}

	// 生成 PNG 字节数据
	pngBytes, err := qr.PNG(size)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "生成 PNG 图片失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("生成 PNG 图片失败: %w", err)
	}

	// 3. 根据输出格式处理
	var output map[string]interface{}

	if outputFormat == "file" {
		// 生成临时文件
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("qrcode_%s.png", timestamp)

		// 创建临时目录
		tempDir := filepath.Join(os.TempDir(), "autoforge-qrcode")
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			return &utools.ExecutionResult{
				Success:    false,
				Message:    "创建临时目录失败",
				Error:      err.Error(),
				DurationMs: time.Since(startTime).Milliseconds(),
			}, fmt.Errorf("创建临时目录失败: %w", err)
		}

		filePath := filepath.Join(tempDir, filename)

		// 写入文件
		if err := os.WriteFile(filePath, pngBytes, 0644); err != nil {
			return &utools.ExecutionResult{
				Success:    false,
				Message:    "保存文件失败",
				Error:      err.Error(),
				DurationMs: time.Since(startTime).Milliseconds(),
			}, fmt.Errorf("保存文件失败: %w", err)
		}

		// 构建文件对象（与 external_trigger 传递的文件对象格式一致）
		fileObject := map[string]interface{}{
			"path":      filePath,
			"filename":  filename,
			"size":      int64(len(pngBytes)),
			"mime_type": "image/png",
		}

		// 构建响应
		response := map[string]interface{}{
			"content":    content,
			"format":     "file",
			"size":       size,
			"level":      levelStr,
			"file":       fileObject,
			"created_at": time.Now().Format(time.RFC3339),
		}

		output = map[string]interface{}{
			"response": response,
			"file":     fileObject, // 快捷访问：文件对象
			"content":  content,
		}
	} else {
		// Base64 编码
		base64Data := base64.StdEncoding.EncodeToString(pngBytes)

		// 构建响应
		response := map[string]interface{}{
			"content":    content,
			"format":     "base64",
			"size":       size,
			"level":      levelStr,
			"data":       base64Data,
			"created_at": time.Now().Format(time.RFC3339),
		}

		output = map[string]interface{}{
			"response": response,
			"data":     base64Data, // 快捷访问：base64 字符串
			"content":  content,
		}
	}

	return &utools.ExecutionResult{
		Success:    true,
		Message:    fmt.Sprintf("二维码生成成功（%s）", outputFormat),
		Output:     output,
		DurationMs: time.Since(startTime).Milliseconds(),
	}, nil
}

func (t *QRCodeTool) parseRecoveryLevel(level string) qrcode.RecoveryLevel {
	switch level {
	case "Low":
		return qrcode.Low
	case "Medium":
		return qrcode.Medium
	case "High":
		return qrcode.High
	case "Highest":
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

func float64Ptr(v float64) *float64 {
	return &v
}

func init() {
	utools.Register(NewQRCodeTool())
}
