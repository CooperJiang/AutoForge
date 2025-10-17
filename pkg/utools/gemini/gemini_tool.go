package gemini

import (
	"auto-forge/pkg/utools"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type GeminiTool struct {
	*utools.BaseTool
}

// GeminiMessage 消息结构
type GeminiMessage struct {
	Role  string                   `json:"role"`
	Parts []map[string]interface{} `json:"parts"`
}

// GeminiRequest API 请求结构
type GeminiRequest struct {
	Contents         []GeminiMessage          `json:"contents"`
	GenerationConfig map[string]interface{}   `json:"generationConfig,omitempty"`
	SafetySettings   []map[string]interface{} `json:"safetySettings,omitempty"`
}

// GeminiResponse API 响应结构
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason  string `json:"finishReason"`
		Index         int    `json:"index"`
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"candidates"`
	PromptFeedback struct {
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"promptFeedback,omitempty"`
}

func NewGeminiTool() *GeminiTool {
	metadata := &utools.ToolMetadata{
		Code:        "gemini_chat",
		Name:        "Gemini AI 对话",
		Description: "使用 Google Gemini AI 模型进行对话交互",
		Category:    utools.CategoryAI,
		Version:     "1.0.0",
		Author:      "Cooper Team",
		AICallable:  true,
		Tags:        []string{"ai", "gemini", "llm", "google", "chat"},
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"response": {
				Type:  "object",
				Label: "完整响应",
				Children: map[string]utools.OutputFieldDef{
					"text": {
						Type:  "string",
						Label: "AI 回复内容",
					},
					"model": {
						Type:  "string",
						Label: "使用的模型",
					},
					"finish_reason": {
						Type:  "string",
						Label: "结束原因",
					},
					"full_response": {
						Type:  "object",
						Label: "完整 API 响应",
					},
				},
			},
			"text": {
				Type:  "string",
				Label: "AI 回复内容（快捷访问）",
			},
		},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"model": {
				Type:        "string",
				Title:       "模型",
				Description: "要使用的 Gemini 模型名称，支持变量",
				Default:     "gemini-pro",
			},
			"prompt": {
				Type:        "string",
				Title:       "提示词",
				Description: "发送给 AI 的内容，支持变量",
			},
			"system_instruction": {
				Type:        "string",
				Title:       "系统指令",
				Description: "设定 AI 的角色和行为（可选）",
				Default:     "",
			},
			"temperature": {
				Type:        "number",
				Title:       "温度",
				Description: "控制输出的随机性，0-2 之间，越高越随机",
				Default:     0.7,
				Minimum:     float64Ptr(0),
				Maximum:     float64Ptr(2),
			},
			"max_tokens": {
				Type:        "integer",
				Title:       "最大 Token 数",
				Description: "生成内容的最大长度，不同模型限制不同",
				Default:     2048,
				Minimum:     float64Ptr(1),
			},
			"top_p": {
				Type:        "number",
				Title:       "Top P",
				Description: "核采样参数，0-1 之间",
				Default:     0.95,
				Minimum:     float64Ptr(0),
				Maximum:     float64Ptr(1),
			},
			"top_k": {
				Type:        "integer",
				Title:       "Top K",
				Description: "采样时考虑的 token 数量",
				Default:     40,
				Minimum:     float64Ptr(1),
			},
			"image": {
				Type:        "object",
				Title:       "图片输入",
				Description: "传入图片文件对象（可选，仅 vision 模型支持），支持变量",
			},
		},
		Required: []string{"prompt"},
	}

	return &GeminiTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

func (t *GeminiTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()

	// 1. 解析配置
	prompt, ok := config["prompt"].(string)
	if !ok || prompt == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "提示词不能为空",
			Error:      "prompt is required",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("提示词不能为空")
	}

	model := "gemini-pro"
	if m, ok := config["model"].(string); ok && m != "" {
		model = m
	}

	systemInstruction := ""
	if si, ok := config["system_instruction"].(string); ok {
		systemInstruction = si
	}

	temperature := 0.7
	if temp, ok := config["temperature"].(float64); ok {
		temperature = temp
	}

	maxTokens := 2048
	if mt, ok := config["max_tokens"]; ok {
		switch v := mt.(type) {
		case float64:
			maxTokens = int(v)
		case int:
			maxTokens = v
		}
	}

	topP := 0.95
	if tp, ok := config["top_p"].(float64); ok {
		topP = tp
	}

	topK := 40
	if tk, ok := config["top_k"]; ok {
		switch v := tk.(type) {
		case float64:
			topK = int(v)
		case int:
			topK = v
		}
	}

	// 2. 从数据库加载 Gemini 配置
	// TODO: 集成工具配置中心后取消注释
	// dbConfig, err := toolConfigService.GetToolConfigForExecution("gemini_chat")
	// if err != nil {
	//     return &utools.ExecutionResult{
	//         Success:    false,
	//         Message:    "Gemini 配置错误",
	//         Error:      err.Error(),
	//         DurationMs: time.Since(startTime).Milliseconds(),
	//     }, err
	// }
	// apiKey, _ := dbConfig["api_key"].(string)

	// 临时方案：从环境变量读取
	apiKey := getEnvOrDefault("GEMINI_API_KEY", "")
	if apiKey == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "Gemini API Key 未配置",
			Error:      "missing gemini api_key",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("gemini api_key not configured")
	}

	// 3. 构建请求
	contents := []GeminiMessage{}

	// 如果有系统指令，作为第一条消息
	if systemInstruction != "" {
		contents = append(contents, GeminiMessage{
			Role: "user",
			Parts: []map[string]interface{}{
				{"text": systemInstruction},
			},
		})
		contents = append(contents, GeminiMessage{
			Role: "model",
			Parts: []map[string]interface{}{
				{"text": "好的，我明白了。"},
			},
		})
	}

	// 用户提示词（支持图片）
	userParts := []map[string]interface{}{
		{"text": prompt},
	}

	// 检查是否有图片输入（支持 file 对象和 base64 字符串）
	if imageParam := config["image"]; imageParam != nil {
		var base64Data string
		var mimeType string

		switch v := imageParam.(type) {
		case map[string]interface{}:
			// 格式 1: File 对象 {path: "...", mime_type: "..."}
			imagePath, _ := v["path"].(string)
			if imagePath != "" {
				// 读取图片文件
				imageData, err := os.ReadFile(imagePath)
				if err != nil {
					return &utools.ExecutionResult{
						Success:    false,
						Message:    "读取图片文件失败",
						Error:      err.Error(),
						DurationMs: time.Since(startTime).Milliseconds(),
					}, fmt.Errorf("读取图片文件失败: %w", err)
				}

				base64Data = base64.StdEncoding.EncodeToString(imageData)

				// 检测 MIME 类型
				if mt, ok := v["mime_type"].(string); ok && mt != "" {
					mimeType = mt
				} else {
					// 根据文件扩展名判断
					ext := strings.ToLower(filepath.Ext(imagePath))
					switch ext {
					case ".png":
						mimeType = "image/png"
					case ".jpg", ".jpeg":
						mimeType = "image/jpeg"
					case ".webp":
						mimeType = "image/webp"
					case ".gif":
						mimeType = "image/gif"
					default:
						mimeType = "image/jpeg"
					}
				}
			}

		case string:
			// 格式 2: Base64 字符串
			inputStr := v

			// 检查是否是 Data URI 格式 (data:image/png;base64,xxx)
			if strings.HasPrefix(inputStr, "data:") {
				// 解析 Data URI
				parts := strings.SplitN(inputStr, ",", 2)
				if len(parts) == 2 {
					base64Data = parts[1]
					// 提取 MIME 类型
					headerParts := strings.Split(parts[0], ";")
					if len(headerParts) > 0 {
						mimeType = strings.TrimPrefix(headerParts[0], "data:")
					}
				}
			} else {
				// 纯 Base64 字符串，默认为 PNG
				base64Data = inputStr
				mimeType = "image/png"
			}
		}

		// 如果成功提取了 base64 数据，添加到请求
		if base64Data != "" {
			if mimeType == "" {
				mimeType = "image/png" // 默认类型
			}

			userParts = append(userParts, map[string]interface{}{
				"inline_data": map[string]interface{}{
					"mime_type": mimeType,
					"data":      base64Data,
				},
			})
		}
	}

	contents = append(contents, GeminiMessage{
		Role:  "user",
		Parts: userParts,
	})

	reqBody := GeminiRequest{
		Contents: contents,
		GenerationConfig: map[string]interface{}{
			"temperature":     temperature,
			"maxOutputTokens": maxTokens,
			"topP":            topP,
			"topK":            topK,
		},
	}

	// 4. 发送请求
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "请求序列化失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("请求序列化失败: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "创建请求失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "请求失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "读取响应失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    fmt.Sprintf("API 返回错误: %d", resp.StatusCode),
			Error:      string(respBody),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("API 返回错误: %d, %s", resp.StatusCode, string(respBody))
	}

	// 5. 解析响应
	var geminiResp GeminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "解析响应失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("解析响应失败: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "API 未返回结果",
			Error:      "no candidates in response",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("API 未返回结果")
	}

	// 提取文本
	candidate := geminiResp.Candidates[0]
	text := ""
	if len(candidate.Content.Parts) > 0 {
		text = candidate.Content.Parts[0].Text
	}

	// 6. 构建输出
	responseData := map[string]interface{}{
		"text":          text,
		"model":         model,
		"finish_reason": candidate.FinishReason,
		"full_response": geminiResp,
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: "Gemini AI 执行成功",
		Output: map[string]interface{}{
			"response": responseData,
			"text":     text,
		},
		DurationMs: time.Since(startTime).Milliseconds(),
	}, nil
}

func float64Ptr(v float64) *float64 {
	return &v
}

func getEnvOrDefault(key, defaultValue string) string {
	// 临时实现，后续集成工具配置中心
	return defaultValue
}

func init() {
	utools.Register(NewGeminiTool())
}
