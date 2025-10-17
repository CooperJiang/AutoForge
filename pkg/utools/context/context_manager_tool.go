package context

import (
    "auto-forge/pkg/cache"
    "auto-forge/pkg/utools"
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "time"
)



type ContextManagerTool struct {
    *utools.BaseTool
}


type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}


func NewContextManagerTool() *ContextManagerTool {
	metadata := &utools.ToolMetadata{
		Code:        "context_manager",
		Name:        "对话上下文管理器",
		Description: "管理多轮对话的上下文历史，支持准备消息（prepare）和保存消息（persist）两种模式",
		Category:    utools.CategoryData,
		Version:     "1.0.0",
		Author:      "AutoForge",
		AICallable:  false,
		Tags:        []string{"context", "chat", "memory", "conversation"},
        OutputFieldsSchema: map[string]utools.OutputFieldDef{
            "messages_json": {
                Type:  "string",
                Label: "标准消息数组（JSON字符串）",
            },
            "openai_messages_json": {
                Type:  "string",
                Label: "OpenAI 消息数组（JSON字符串）",
            },
            "gemini_contents_json": {
                Type:  "string",
                Label: "Gemini contents（JSON字符串）",
            },
            "anthropic_messages_json": {
                Type:  "string",
                Label: "Anthropic 消息数组（JSON字符串）",
            },
            "anthropic_system": {
                Type:  "string",
                Label: "Anthropic system（字符串）",
            },
            "messages": {
                Type:  "array",
                Label: "消息数组（对象）",
            },
			"session_key": {
				Type:  "string",
				Label: "会话标识",
			},
			"message_count": {
				Type:  "number",
				Label: "消息总数",
			},
			"preview": {
				Type:  "string",
				Label: "对话预览（文本）",
			},
			"updated": {
				Type:  "boolean",
				Label: "是否已更新（persist模式）",
			},
		},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
        Properties: map[string]utools.PropertySchema{
            "mode": {
				Type:        "string",
				Title:       "工作模式",
				Description: "prepare=准备消息（读取历史+拼接），persist=保存消息（追加到历史）",
				Default:     "prepare",
				Enum:        []interface{}{"prepare", "persist"},
            },
            "session_key": {
                Type:        "string",
                Title:       "会话标识",
                Description: "用于区分不同用户的对话，支持变量如 {{external.session_id}}",
                MinLength:   func() *int { v := 1; return &v }(),
            },
            "scope": {
                Type:        "string",
                Title:       "作用域",
                Description: "node=每节点独立，workflow=同工作流共享",
                Default:     "node",
                Enum:        []interface{}{"node", "workflow"},
            },
            "workflow_id": {
                Type:        "string",
                Title:       "工作流ID（可选）",
                Description: "用于命名空间键；不填则省略",
            },
            "user_input": {
                Type:        "string",
                Title:       "用户输入（prepare模式）",
                Description: "当前用户发送的消息，支持变量",
            },
			"assistant_output": {
				Type:        "string",
				Title:       "AI回复（persist模式）",
				Description: "AI助手的回复内容，通常引用上游节点输出",
			},
			"system_message": {
				Type:        "string",
				Title:       "系统消息（可选）",
				Description: "设定AI角色和行为的系统提示",
			},
			"window_size": {
				Type:        "number",
				Title:       "窗口大小",
				Description: "保留最近N条消息（包括system/user/assistant）",
				Default:     10.0,
				Minimum:     func() *float64 { v := 1.0; return &v }(),
			},
			"ttl_seconds": {
				Type:        "number",
				Title:       "过期时间（秒）",
				Description: "对话历史的存储时长，默认7天",
				Default:     604800.0,
				Minimum:     func() *float64 { v := 60.0; return &v }(),
			},
			"clear_history": {
				Type:        "boolean",
				Title:       "清空历史",
				Description: "执行前清空当前会话的历史记录",
				Default:     false,
			},
		},
		Required: []string{"mode", "session_key"},
	}

	return &ContextManagerTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}


func (t *ContextManagerTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()


	_ = cache.GetCache()


	mode, _ := config["mode"].(string)
	sessionKey, _ := config["session_key"].(string)

	if sessionKey == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "会话标识不能为空",
			Error:      "session_key is required",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("session_key is required")
	}


    cacheKey := fmt.Sprintf("chat:context:%s", sessionKey)


	clearHistory, _ := config["clear_history"].(bool)
	if clearHistory {
		cache.Del(cacheKey)
	}


	switch mode {
	case "prepare":
		return t.executePrepare(config, cacheKey, startTime)
	case "persist":
		return t.executePersist(config, cacheKey, startTime)
	default:
		return &utools.ExecutionResult{
			Success:    false,
			Message:    fmt.Sprintf("不支持的模式: %s", mode),
			Error:      "unsupported mode",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("unsupported mode: %s", mode)
	}
}


func (t *ContextManagerTool) executePrepare(config map[string]interface{}, cacheKey string, startTime time.Time) (*utools.ExecutionResult, error) {
	sessionKey, _ := config["session_key"].(string)
	userInput, _ := config["user_input"].(string)
	systemMessage, _ := config["system_message"].(string)
	windowSize := t.getWindowSize(config)


	var messages []ChatMessage
	historyJSON, err := cache.Get(cacheKey)
	if err == nil && historyJSON != "" {
		if err := json.Unmarshal([]byte(historyJSON), &messages); err != nil {

			messages = []ChatMessage{}
		}
	}


	if systemMessage != "" {
		hasSystem := false
		for _, msg := range messages {
			if msg.Role == "system" {
				hasSystem = true
				break
			}
		}
		if !hasSystem {
			messages = append([]ChatMessage{{
				Role:    "system",
				Content: systemMessage,
			}}, messages...)
		}
	}


	if userInput != "" {
		messages = append(messages, ChatMessage{
			Role:    "user",
			Content: userInput,
		})
	}


	messages = t.trimMessages(messages, windowSize)


    messagesJSON, _ := json.Marshal(messages)
    preview := t.generatePreview(messages)


    openaiJSON := string(messagesJSON)
    geminiJSON := t.toGeminiContentsJSON(messages)
    anthropicMsgsJSON, anthropicSystem := t.toAnthropicJSON(messages)

	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("已准备 %d 条消息", len(messages)),
        Output: map[string]interface{}{
            "messages_json": string(messagesJSON),
            "openai_messages_json": openaiJSON,
            "gemini_contents_json": geminiJSON,
            "anthropic_messages_json": anthropicMsgsJSON,
            "anthropic_system": anthropicSystem,
            "messages":      messages,
            "session_key":   sessionKey,
            "message_count": len(messages),
            "preview":       preview,
        },
		DurationMs: time.Since(startTime).Milliseconds(),
	}, nil
}


func (t *ContextManagerTool) executePersist(config map[string]interface{}, cacheKey string, startTime time.Time) (*utools.ExecutionResult, error) {
	sessionKey, _ := config["session_key"].(string)
	userInput, _ := config["user_input"].(string)
	assistantOutput, _ := config["assistant_output"].(string)
	systemMessage, _ := config["system_message"].(string)
	windowSize := t.getWindowSize(config)
	ttl := t.getTTL(config)


    unlock := t.tryLock(cacheKey)
    defer unlock()


    var messages []ChatMessage
    historyJSON, err := cache.Get(cacheKey)
	if err == nil && historyJSON != "" {
		json.Unmarshal([]byte(historyJSON), &messages)
	}


	if systemMessage != "" {
		hasSystem := false
		for _, msg := range messages {
			if msg.Role == "system" {
				hasSystem = true
				break
			}
		}
		if !hasSystem {
			messages = append([]ChatMessage{{
				Role:    "system",
				Content: systemMessage,
			}}, messages...)
		}
	}


	if userInput != "" {
		messages = append(messages, ChatMessage{
			Role:    "user",
			Content: userInput,
		})
	}


	if assistantOutput != "" {
		messages = append(messages, ChatMessage{
			Role:    "assistant",
			Content: assistantOutput,
		})
	}


	messages = t.trimMessages(messages, windowSize)


    newMessagesJSON, _ := json.Marshal(messages)
    if err := cache.Set(cacheKey, string(newMessagesJSON), ttl); err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "保存失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("已保存 %d 条消息", len(messages)),
		Output: map[string]interface{}{
			"updated":       true,
			"session_key":   sessionKey,
			"message_count": len(messages),
			"ttl_seconds":   ttl.Seconds(),
		},
		DurationMs: time.Since(startTime).Milliseconds(),
	}, nil
}


func (t *ContextManagerTool) trimMessages(messages []ChatMessage, windowSize int) []ChatMessage {
	if len(messages) <= windowSize {
		return messages
	}


	systemMsg := []ChatMessage{}
	otherMsgs := messages

	if len(messages) > 0 && messages[0].Role == "system" {
		systemMsg = []ChatMessage{messages[0]}
		otherMsgs = messages[1:]
	}


	if len(otherMsgs) > windowSize-len(systemMsg) {
		otherMsgs = otherMsgs[len(otherMsgs)-(windowSize-len(systemMsg)):]
	}

	return append(systemMsg, otherMsgs...)
}


func (t *ContextManagerTool) generatePreview(messages []ChatMessage) string {
	var lines []string
	for i, msg := range messages {
		prefix := ""
		switch msg.Role {
		case "system":
			prefix = "[系统]"
		case "user":
			prefix = "[用户]"
		case "assistant":
			prefix = "[AI]"
		}


		content := msg.Content
		if len(content) > 50 {
			content = content[:50] + "..."
		}

		lines = append(lines, fmt.Sprintf("%d. %s %s", i+1, prefix, content))
	}
	return strings.Join(lines, "\n")
}


func (t *ContextManagerTool) getWindowSize(config map[string]interface{}) int {
	windowSize, ok := config["window_size"].(float64)
	if !ok || windowSize <= 0 {
		return 10
	}
	return int(windowSize)
}


func (t *ContextManagerTool) getTTL(config map[string]interface{}) time.Duration {
	ttlSeconds, ok := config["ttl_seconds"].(float64)
	if !ok || ttlSeconds <= 0 {
		return 7 * 24 * time.Hour
	}
	return time.Duration(ttlSeconds) * time.Second
}


func (t *ContextManagerTool) DescribeOutput(config map[string]interface{}) map[string]utools.OutputFieldDef {
    mode, _ := config["mode"].(string)
    if mode == "persist" {
        return map[string]utools.OutputFieldDef{
            "updated":       {Type: "boolean", Label: "是否已更新"},
            "session_key":   {Type: "string", Label: "会话标识"},
            "message_count": {Type: "number", Label: "消息总数"},
            "ttl_seconds":   {Type: "number", Label: "有效期（秒）"},
        }
    }

    return map[string]utools.OutputFieldDef{
        "messages_json":          {Type: "string", Label: "中立消息 JSON"},
        "openai_messages_json":   {Type: "string", Label: "OpenAI 消息 JSON"},
        "gemini_contents_json":   {Type: "string", Label: "Gemini contents JSON"},
        "anthropic_messages_json": {Type: "string", Label: "Anthropic messages JSON"},
        "anthropic_system":       {Type: "string", Label: "Anthropic system"},
        "message_count":          {Type: "number", Label: "消息总数"},
        "session_key":            {Type: "string", Label: "会话标识"},
        "preview":                {Type: "string", Label: "预览文本"},
    }
}


func (t *ContextManagerTool) tryLock(cacheKey string) func() {
    client := cache.GetRedisClient()
    if client == nil {
        return func() {}
    }
    ctx := cache.GetRedisContext()
    lockKey := "lock:" + cacheKey

    for i := 0; i < 5; i++ {
        ok, _ := client.SetNX(ctx, lockKey, "1", 2*time.Second).Result()
        if ok {
            return func() {
                client.Del(context.Background(), lockKey)
            }
        }
        time.Sleep(50 * time.Millisecond)
    }

    return func() {}
}


func (t *ContextManagerTool) toGeminiContentsJSON(messages []ChatMessage) string {
    type part struct{ Text string `json:"text"` }
    type content struct{
        Role string  `json:"role"`
        Parts []part `json:"parts"`
    }
    contents := make([]content, 0, len(messages))
    for _, m := range messages {
        role := m.Role
        if role == "assistant" { role = "model" }
        if role == "system" { role = "user" }
        contents = append(contents, content{Role: role, Parts: []part{{Text: m.Content}}})
    }
    b, _ := json.Marshal(contents)
    return string(b)
}


func (t *ContextManagerTool) toAnthropicJSON(messages []ChatMessage) (string, string) {
    type amsg struct{
        Role string `json:"role"`
        Content string `json:"content"`
    }
    system := ""
    msgs := make([]amsg, 0, len(messages))
    for _, m := range messages {
        if m.Role == "system" {
            if system == "" { system = m.Content }
            continue
        }
        role := m.Role
        if role != "user" && role != "assistant" { role = "user" }
        msgs = append(msgs, amsg{Role: role, Content: m.Content})
    }
    b, _ := json.Marshal(msgs)
    return string(b), system
}


func init() {
	tool := NewContextManagerTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register Context Manager tool: %v", err))
	}
}
