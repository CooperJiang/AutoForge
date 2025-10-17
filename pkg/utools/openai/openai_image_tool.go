package openai

import (
	toolConfigService "auto-forge/internal/services/tool_config"
	"auto-forge/pkg/utools"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAIImageTool struct {
	*utools.BaseTool
}

func NewOpenAIImageTool() *OpenAIImageTool {
	metadata := &utools.ToolMetadata{
		Code:        "openai_image",
		Name:        "OpenAI Image",
		Description: "使用 OpenAI DALL-E Gpt-image-1 等模型生成图片",
		Category:    utools.CategoryAI,
		Version:     "1.0.0",
		Author:      "AutoForge",
		AICallable:  true,
		Tags:        []string{"openai", "dalle", "image", "ai", "picture"},
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"response": {
				Type:  "object",
				Label: "OpenAI 原始响应",
				Children: map[string]utools.OutputFieldDef{
					"created": {Type: "number", Label: "创建时间戳"},
					"data": {
						Type:  "array",
						Label: "图片数组",
						Children: map[string]utools.OutputFieldDef{
							"0": {
								Type:  "object",
								Label: "第一张图片",
								Children: map[string]utools.OutputFieldDef{
									"url":            {Type: "string", Label: "图片 URL（当 response_format=url）"},
									"b64_json":       {Type: "string", Label: "Base64 内容（当 response_format=b64_json）"},
									"revised_prompt": {Type: "string", Label: "修订后的提示词（可选）"},
								},
							},
						},
					},
				},
			},
		},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"model": {
				Type:        "string",
				Title:       "模型",
				Description: "图片生成模型，例如：dall-e-2、dall-e-3、gpt-image-1 等，支持变量",
				Default:     "dall-e-3",
			},
			"prompt": {
				Type:        "string",
				Title:       "提示词",
				Description: "描述你想要生成的图片内容，支持变量",
			},
			"n": {
				Type:        "string",
				Title:       "生成数量",
				Description: "生成图片的数量（dall-e-3 只支持 1 张），支持变量",
				Default:     "1",
			},
			"size": {
				Type:        "string",
				Title:       "图片尺寸",
				Description: "dall-e-2: 256x256, 512x512, 1024x1024; dall-e-3: 1024x1024, 1792x1024, 1024x1792，支持变量",
				Default:     "1024x1024",
			},
			"quality": {
				Type:        "string",
				Title:       "图片质量",
				Description: "standard（标准）或 hd（高清），仅部分模型支持，留空使用默认值，支持变量",
			},
			"response_format": {
				Type:        "string",
				Title:       "返回格式",
				Description: "url（图片链接）或 b64_json（base64 编码），支持变量",
				Default:     "url",
			},
			"timeout": {
				Type:        "string",
				Title:       "超时时间",
				Description: "API 请求超时时间（秒），默认 300 秒，支持变量",
				Default:     "300",
			},
		},
		Required: []string{"model", "prompt"},
	}

	return &OpenAIImageTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

func (t *OpenAIImageTool) Execute(ctx *utools.ExecutionContext, toolConfig map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()

	// 从数据库加载 OpenAI 配置
	dbConfig, err := toolConfigService.GetToolConfigForExecution("openai_image")
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "OpenAI 配置错误",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 解析配置字段
	apiKey, _ := dbConfig["api_key"].(string)
	apiBase, _ := dbConfig["api_base"].(string)

	// 验证配置
	if apiKey == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "OpenAI API Key 未配置",
			Error:      "missing openai api_key in config",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("openai api_key not configured")
	}

	if apiBase == "" {
		apiBase = "https://api.openai.com"
	}

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

	n := 1
	if nVal, ok := toolConfig["n"].(float64); ok && nVal > 0 {
		n = int(nVal)
	} else if nStr, ok := toolConfig["n"].(string); ok && nStr != "" {

		var nParsed int
		if _, err := fmt.Sscanf(nStr, "%d", &nParsed); err == nil && nParsed > 0 {
			n = nParsed
		}
	}

	size, _ := toolConfig["size"].(string)
	if size == "" {
		size = "1024x1024"
	}

	quality, _ := toolConfig["quality"].(string)

	responseFormat, _ := toolConfig["response_format"].(string)
	if responseFormat == "" {
		responseFormat = "url"
	}

	timeout := 300
	if t, ok := toolConfig["timeout"].(float64); ok && t > 0 {
		timeout = int(t)
	} else if tStr, ok := toolConfig["timeout"].(string); ok && tStr != "" {

		var tParsed int
		if _, err := fmt.Sscanf(tStr, "%d", &tParsed); err == nil && tParsed > 0 {
			timeout = tParsed
		}
	}

	requestBody := map[string]interface{}{
		"model":           model,
		"prompt":          prompt,
		"n":               n,
		"size":            size,
		"response_format": responseFormat,
	}

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
			"response": result,
		},
		DurationMs: time.Since(startTime).Milliseconds(),
	}, nil
}

func init() {
	tool := NewOpenAIImageTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register OpenAI Image tool: %v", err))
	}
}
