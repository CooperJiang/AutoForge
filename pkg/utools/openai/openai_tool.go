package openai

import (
	"auto-forge/pkg/config"
	"auto-forge/pkg/utools"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)


type OpenAITool struct {
	*utools.BaseTool
}


func NewOpenAITool() *OpenAITool {
	metadata := &utools.ToolMetadata{
		Code:         "openai_chatgpt",
		Name:         "OpenAI Chat",
		Description:  "使用 OpenAI Chat API 进行对话，支持 GPT-3.5、GPT-4、GPT-4o 等对话模型",
		Category:     "ai",
		Version:      "1.0.0",
		Author:       "AutoForge",
		AICallable:   true,
		Tags:         []string{"openai", "chatgpt", "ai", "llm", "gpt", "chat"},
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"response": {
				Type:  "object",
				Label: "OpenAI 原始响应",
				Children: map[string]utools.OutputFieldDef{
					"id":      {Type: "string", Label: "响应ID"},
					"object":  {Type: "string", Label: "对象类型"},
					"created": {Type: "number", Label: "创建时间戳"},
					"model":   {Type: "string", Label: "使用的模型"},
					"choices": {
						Type:  "array",
						Label: "回复选项数组",
						Children: map[string]utools.OutputFieldDef{
							"0": {
								Type:  "object",
								Label: "第一个回复",
								Children: map[string]utools.OutputFieldDef{
									"index": {Type: "number", Label: "索引"},
									"message": {
										Type:  "object",
										Label: "消息对象",
										Children: map[string]utools.OutputFieldDef{
											"role":    {Type: "string", Label: "角色"},
											"content": {Type: "string", Label: "回复内容"},
										},
									},
									"finish_reason": {Type: "string", Label: "结束原因"},
								},
							},
						},
					},
					"usage": {
						Type:  "object",
						Label: "Token 使用情况",
						Children: map[string]utools.OutputFieldDef{
							"prompt_tokens":     {Type: "number", Label: "提示词 Token 数"},
							"completion_tokens": {Type: "number", Label: "回复 Token 数"},
							"total_tokens":      {Type: "number", Label: "总 Token 数"},
						},
					},
				},
			},
			"content":           {Type: "string", Label: "ChatGPT 回复内容"},
			"content_json":      {Type: "object", Label: "JSON对象（如果回复是JSON，会自动解析）"},
			"model":             {Type: "string", Label: "使用的模型"},
			"finish_reason":     {Type: "string", Label: "结束原因"},
			"prompt_tokens":     {Type: "number", Label: "提示词 Token 数"},
			"completion_tokens": {Type: "number", Label: "回复 Token 数"},
			"total_tokens":      {Type: "number", Label: "总 Token 数"},
			"prompt":            {Type: "string", Label: "原始提示词"},
		},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"messages_json": {
				Type:        "string",
				Title:       "消息数组(JSON)",
				Description: "符合 OpenAI Chat Completions 的 messages 数组 JSON 字符串，存在时优先使用",
			},
			"model": {
				Type:        "string",
				Title:       "模型",
				Description: "使用的对话模型，例如：gpt-3.5-turbo、gpt-4、gpt-4o 等",
				Default:     "gpt-3.5-turbo",
			},
			"prompt": {
				Type:        "string",
				Title:       "提示词",
				Description: "发送给 AI 的问题或指令，支持变量",
			},
			"system_message": {
				Type:        "string",
				Title:       "系统消息",
				Description: "系统角色消息（可选），用于设定 AI 的行为和角色",
			},
			"temperature": {
				Type:        "number",
				Title:       "温度",
				Description: "控制回复的随机性，0-2 之间，越高越随机",
				Default:     0.7,
			},
			"max_tokens": {
				Type:        "number",
				Title:       "最大 Token 数",
				Description: "生成回复的最大 token 数量（可选）",
			},
			"timeout": {
				Type:        "number",
				Title:       "超时时间",
				Description: "API 请求超时时间（秒），默认 300 秒",
				Default:     300,
			},
			"context_config": {
				Type:        "object",
				Title:       "对话上下文配置",
				Description: "启用后自动管理多轮对话的上下文记忆",
				Properties: map[string]utools.PropertySchema{
					"enabled": {
						Type:        "boolean",
						Title:       "启用对话记忆",
						Description: "开启后会自动保存和加载对话历史",
						Default:     false,
					},
					"session_key": {
						Type:        "string",
						Title:       "会话ID",
						Description: "用于区分不同对话，支持变量如 {{params.user_id}}",
						Default:     "{{params.session_id}}",
					},
					"window_size": {
						Type:        "number",
						Title:       "窗口条数",
						Description: "保留最近 N 条消息",
						Default:     10.0,
						Minimum:     func() *float64 { v := 1.0; return &v }(),
					},
					"ttl_seconds": {
						Type:        "number",
						Title:       "过期时间（秒）",
						Description: "对话历史的过期时间，默认 7 天",
						Default:     604800.0,
						Minimum:     func() *float64 { v := 60.0; return &v }(),
					},
				},
			},
		},
		Required: []string{"model", "prompt"},
	}

	return &OpenAITool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}


func (t *OpenAITool) Execute(ctx *utools.ExecutionContext, toolConfig map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()


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
		apiBase = "https://api.openai.com"
	}


	model, _ := toolConfig["model"].(string)
	if model == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "模型不能为空",
			Error:      "missing model",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("model is required")
	}

	prompt, _ := toolConfig["prompt"].(string)
	if prompt == "" {
		if raw, ok := toolConfig["messages_json"].(string); !ok || strings.TrimSpace(raw) == "" {
			return &utools.ExecutionResult{
				Success:    false,
				Message:    "提示词不能为空",
				Error:      "missing prompt",
				DurationMs: time.Since(startTime).Milliseconds(),
			}, fmt.Errorf("prompt is required")
		}
	}

	systemMessage, _ := toolConfig["system_message"].(string)
	temperature := 0.7
	if temp, ok := toolConfig["temperature"].(float64); ok {
		temperature = temp
	}

	timeout := 300
	if t, ok := toolConfig["timeout"].(float64); ok && t > 0 {
		timeout = int(t)
	}

	messages := []map[string]interface{}{}

	if raw, ok := toolConfig["messages_json"].(string); ok {
		trim := bytes.TrimSpace([]byte(raw))
		if len(trim) > 0 {
			var arr []map[string]interface{}
			if err := json.Unmarshal(trim, &arr); err != nil {
				return &utools.ExecutionResult{
					Success:    false,
					Message:    "messages_json 解析失败",
					Error:      err.Error(),
					DurationMs: time.Since(startTime).Milliseconds(),
				}, err
			}
			messages = arr
		}
	}

	if len(messages) == 0 {
		if systemMessage != "" {
			messages = append(messages, map[string]interface{}{
				"role":    "system",
				"content": systemMessage,
			})
		}
		messages = append(messages, map[string]interface{}{
			"role":    "user",
			"content": prompt,
		})
	}

	requestBody := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"temperature": temperature,
	}

	if maxTokens, ok := toolConfig["max_tokens"].(float64); ok && maxTokens > 0 {
		requestBody["max_tokens"] = int(maxTokens)
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

	url := fmt.Sprintf("%s/chat/completions", apiBase)
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


	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "OpenAI API 返回格式异常：未找到 choices 字段",
			Error:      "no choices in response",
			Output: map[string]interface{}{
				"raw_response": result,
				"status_code":  resp.StatusCode,
			},
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("no choices in response, raw: %v", result)
	}


    content := ""
    finishReason := ""
    if len(choices) > 0 {
        if choice, ok := choices[0].(map[string]interface{}); ok {
            if message, ok := choice["message"].(map[string]interface{}); ok {
                if c, ok := message["content"].(string); ok {
                    content = c
                }
            }
            if fr, ok := choice["finish_reason"].(string); ok {
                finishReason = fr
            }
        }
    }


    var promptTokens, completionTokens, totalTokens int
    if usage, ok := result["usage"].(map[string]interface{}); ok {
        if v, ok := usage["prompt_tokens"].(float64); ok { promptTokens = int(v) }
        if v, ok := usage["completion_tokens"].(float64); ok { completionTokens = int(v) }
        if v, ok := usage["total_tokens"].(float64); ok { totalTokens = int(v) }
    }
    modelStr, _ := result["model"].(string)


    output := map[string]interface{}{
        "response":          result,
        "content":           content,
        "model":             modelStr,
        "finish_reason":     finishReason,
        "prompt_tokens":     promptTokens,
        "completion_tokens": completionTokens,
        "total_tokens":      totalTokens,
        "prompt":            prompt,
    }


    if content != "" {
        var contentJSON map[string]interface{}
        if err := json.Unmarshal([]byte(content), &contentJSON); err == nil {
            output["content_json"] = contentJSON


            for k, v := range contentJSON {
                output[k] = v
            }
        }
    }

    return &utools.ExecutionResult{
        Success:    true,
        Message:    "ChatGPT 调用成功",
        Output:     output,
        DurationMs: time.Since(startTime).Milliseconds(),
    }, nil
}



func (t *OpenAITool) DescribeOutput(config map[string]interface{}) map[string]utools.OutputFieldDef {
	return map[string]utools.OutputFieldDef{
		"response":          {Type: "object", Label: "OpenAI 原始响应"},
		"content":           {Type: "string", Label: "ChatGPT 回复内容（字符串）"},
		"content_json":      {Type: "object", Label: "JSON对象（如果回复是JSON，会自动解析）"},
		"model":             {Type: "string", Label: "使用的模型"},
		"finish_reason":     {Type: "string", Label: "结束原因"},
		"prompt_tokens":     {Type: "number", Label: "提示词 Token 数"},
		"completion_tokens": {Type: "number", Label: "回复 Token 数"},
		"total_tokens":      {Type: "number", Label: "总 Token 数"},
		"prompt":            {Type: "string", Label: "原始提示词"},
	}
}


func init() {
	tool := NewOpenAITool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register OpenAI tool: %v", err))
	}
}
