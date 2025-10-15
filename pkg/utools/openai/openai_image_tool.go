package openai

import (
	"auto-forge/pkg/config"
	"auto-forge/pkg/utools"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIImageTool OpenAI 图片生成工具
type OpenAIImageTool struct {
	*utools.BaseTool
}

// NewOpenAIImageTool 创建 OpenAI 图片生成工具实例
func NewOpenAIImageTool() *OpenAIImageTool {
	metadata := &utools.ToolMetadata{
		Code:        "openai_image",
		Name:        "OpenAI Image",
		Description: "使用 OpenAI DALL-E Gpt-image-1 等模型生成图片",
		Category:    "ai",
		Version:     "1.0.0",
		Author:      "AutoForge",
		AICallable:  true,
		Tags:        []string{"openai", "dalle", "image", "ai", "picture"},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"model": {
				Type:        "string",
				Title:       "模型",
				Description: "图片生成模型，例如：dall-e-2、dall-e-3、gpt-image-1 等",
				Default:     "dall-e-3",
			},
			"prompt": {
				Type:        "string",
				Title:       "提示词",
				Description: "描述你想要生成的图片内容，支持变量",
			},
			"n": {
				Type:        "number",
				Title:       "生成数量",
				Description: "生成图片的数量（dall-e-3 只支持 1 张）",
				Default:     1,
			},
			"size": {
				Type:        "string",
				Title:       "图片尺寸",
				Description: "dall-e-2: 256x256, 512x512, 1024x1024; dall-e-3: 1024x1024, 1792x1024, 1024x1792",
				Default:     "1024x1024",
				Enum:        []interface{}{"256x256", "512x512", "1024x1024", "1792x1024", "1024x1792"},
			},
			"quality": {
				Type:        "string",
				Title:       "图片质量",
				Description: "standard（标准）或 hd（高清），仅部分模型支持，留空使用默认值",
				Enum:        []interface{}{"", "standard", "hd"},
			},
			"response_format": {
				Type:        "string",
				Title:       "返回格式",
				Description: "url（图片链接）或 b64_json（base64 编码）",
				Default:     "url",
				Enum:        []interface{}{"url", "b64_json"},
			},
			"timeout": {
				Type:        "number",
				Title:       "超时时间",
				Description: "API 请求超时时间（秒），默认 300 秒",
				Default:     300,
			},
		},
		Required: []string{"model", "prompt"},
	}

	return &OpenAIImageTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

// Execute 执行图片生成
func (t *OpenAIImageTool) Execute(ctx *utools.ExecutionContext, toolConfig map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()

	// 从配置文件读取 API 凭证
	cfg := config.GetConfig()
	apiKey := cfg.OpenAI.APIKey
	if apiKey == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "OpenAI API Key 未配置，请在 config.yaml 中配置 openai.api_key",
			Error:      "missing openai api_key in config",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("openai api_key not configured")
	}

	apiBase := cfg.OpenAI.APIBase
	if apiBase == "" {
		apiBase = "https://api.openai.com/v1"
	}

	// 从用户配置获取参数
	model, _ := toolConfig["model"].(string)
	if model == "" {
		model = "dall-e-3"
	}

	prompt, _ := toolConfig["prompt"].(string)
	if prompt == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "提示词不能为空",
			Error:      "missing prompt",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("prompt is required")
	}

	// 获取其他参数
	n := 1
	if nVal, ok := toolConfig["n"].(float64); ok && nVal > 0 {
		n = int(nVal)
	}

	size, _ := toolConfig["size"].(string)
	if size == "" {
		size = "1024x1024"
	}

	quality, _ := toolConfig["quality"].(string)
	// quality 可以为空，使用模型默认值

	responseFormat, _ := toolConfig["response_format"].(string)
	if responseFormat == "" {
		responseFormat = "url"
	}

	// 超时时间，默认 300 秒
	timeout := 300
	if t, ok := toolConfig["timeout"].(float64); ok && t > 0 {
		timeout = int(t)
	}

	// 构建请求体
	requestBody := map[string]interface{}{
		"model":           model,
		"prompt":          prompt,
		"n":               n,
		"size":            size,
		"response_format": responseFormat,
	}

	// quality 参数（可选）
	if quality != "" {
		requestBody["quality"] = quality
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "构建请求失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 发送请求
	url := fmt.Sprintf("%s/images/generations", apiBase)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "创建请求失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "API 请求失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "读取响应失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "解析响应失败",
			Error:      err.Error(),
			Output:     map[string]interface{}{"raw_response": string(body)},
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != 200 {
		errorMsg := fmt.Sprintf("HTTP %d", resp.StatusCode)
		if errorObj, ok := result["error"].(map[string]interface{}); ok {
			if msg, ok := errorObj["message"].(string); ok {
				errorMsg = msg
			}
		}
		return &utools.ExecutionResult{
			Success:    false,
			Message:    fmt.Sprintf("OpenAI API 错误: %s", errorMsg),
			Error:      errorMsg,
			Output:     result,
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("openai api error: %s", errorMsg)
	}

	// 检查错误
	if errorObj, ok := result["error"].(map[string]interface{}); ok {
		errorMsg, _ := errorObj["message"].(string)
		errorType, _ := errorObj["type"].(string)
		return &utools.ExecutionResult{
			Success:    false,
			Message:    fmt.Sprintf("OpenAI API 错误: %s", errorMsg),
			Error:      errorType,
			Output:     result,
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("openai api error: %s", errorMsg)
	}

	// 提取图片数据
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return &utools.ExecutionResult{
			Success: false,
			Message: "OpenAI API 返回格式异常：未找到 data 字段",
			Error:   "no data in response",
			Output: map[string]interface{}{
				"raw_response": result,
				"status_code":  resp.StatusCode,
			},
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("no data in response, raw: %v", result)
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("图片生成成功，共 %d 张", len(data)),
		Output: map[string]interface{}{
			"response": result, // OpenAI 原始完整响应
		},
		DurationMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// init 自动注册工具
func init() {
	tool := NewOpenAIImageTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register OpenAI Image tool: %v", err))
	}
}
