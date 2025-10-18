package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAIClient OpenAI 客户端
type OpenAIClient struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// NewOpenAIClient 创建 OpenAI 客户端
func NewOpenAIClient(model, apiKey, baseURL string) *OpenAIClient {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &OpenAIClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Call 同步调用
func (c *OpenAIClient) Call(ctx context.Context, messages []Message, options *CallOptions) (*Response, error) {
	// 构建请求
	reqBody := c.buildRequest(messages, options, false)

	// 发送请求
	respBody, err := c.doRequest(ctx, "/chat/completions", reqBody)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var apiResp openAIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("响应中没有选项")
	}

	choice := apiResp.Choices[0]
	resp := &Response{
		Content:      choice.Message.Content,
		ToolCalls:    convertToolCalls(choice.Message.ToolCalls),
		FinishReason: choice.FinishReason,
		Usage: TokenUsage{
			PromptTokens:     apiResp.Usage.PromptTokens,
			CompletionTokens: apiResp.Usage.CompletionTokens,
			TotalTokens:      apiResp.Usage.TotalTokens,
		},
	}

	return resp, nil
}

// Stream 流式调用
func (c *OpenAIClient) Stream(ctx context.Context, messages []Message, options *CallOptions) (<-chan StreamChunk, error) {
	chunkChan := make(chan StreamChunk, 10)

	// 构建请求
	reqBody := c.buildRequest(messages, options, true)

	go func() {
		defer close(chunkChan)

		// 发送流式请求
		resp, err := c.doStreamRequest(ctx, "/chat/completions", reqBody)
		if err != nil {
			chunkChan <- StreamChunk{Error: err, Done: true}
			return
		}
		defer resp.Body.Close()

		// 读取 SSE 流
		reader := NewSSEReader(resp.Body)
		var fullContent strings.Builder
		var toolCalls []ToolCall

		for {
			event, err := reader.Read()
			if err != nil {
				if err != io.EOF {
					chunkChan <- StreamChunk{Error: err, Done: true}
				}
				break
			}

			// 跳过空事件
			if event.Data == "" || event.Data == "[DONE]" {
				continue
			}

			// 解析事件
			var chunk openAIStreamChunk
			if err := json.Unmarshal([]byte(event.Data), &chunk); err != nil {
				continue
			}

			if len(chunk.Choices) == 0 {
				continue
			}

			delta := chunk.Choices[0].Delta
			finishReason := chunk.Choices[0].FinishReason

			// 处理内容增量
			if delta.Content != "" {
				fullContent.WriteString(delta.Content)
				chunkChan <- StreamChunk{
					Content: delta.Content,
					Done:    false,
				}
			}

			// 处理工具调用增量
			if len(delta.ToolCalls) > 0 {
				for _, tc := range delta.ToolCalls {
					// 确保 toolCalls 数组足够大
					for len(toolCalls) <= tc.Index {
						toolCalls = append(toolCalls, ToolCall{})
					}

					// 更新工具调用
					if tc.ID != "" {
						toolCalls[tc.Index].ID = tc.ID
						toolCalls[tc.Index].Type = tc.Type
					}
					if tc.Function.Name != "" {
						toolCalls[tc.Index].Function.Name = tc.Function.Name
					}
					if tc.Function.Arguments != "" {
						toolCalls[tc.Index].Function.Arguments += tc.Function.Arguments
					}
				}
			}

			// 检查是否完成
			if finishReason != "" {
				chunkChan <- StreamChunk{
					Content:      fullContent.String(),
					ToolCalls:    toolCalls,
					FinishReason: finishReason,
					Done:         true,
				}
				break
			}
		}
	}()

	return chunkChan, nil
}

// GetModelInfo 获取模型信息
func (c *OpenAIClient) GetModelInfo() ModelInfo {
	maxTokens := 4096
	if strings.Contains(c.model, "gpt-4") {
		maxTokens = 8192
	}
	if strings.Contains(c.model, "gpt-4-turbo") || strings.Contains(c.model, "gpt-4o") {
		maxTokens = 128000
	}

	return ModelInfo{
		Provider:    "openai",
		Model:       c.model,
		MaxTokens:   maxTokens,
		SupportTool: true,
	}
}

// buildRequest 构建请求体
func (c *OpenAIClient) buildRequest(messages []Message, options *CallOptions, stream bool) map[string]interface{} {
	req := map[string]interface{}{
		"model":    c.model,
		"messages": convertMessages(messages),
		"stream":   stream,
	}

	if options != nil {
		if options.Temperature > 0 {
			req["temperature"] = options.Temperature
		}
		if options.MaxTokens > 0 {
			req["max_tokens"] = options.MaxTokens
		}
		if options.TopP > 0 {
			req["top_p"] = options.TopP
		}
		if len(options.Stop) > 0 {
			req["stop"] = options.Stop
		}
		if len(options.Tools) > 0 {
			req["tools"] = options.Tools
		}
		if options.ToolChoice != nil {
			req["tool_choice"] = options.ToolChoice
		}
		if options.ResponseFormat == "json_object" {
			req["response_format"] = map[string]string{"type": "json_object"}
		}
	}

	return req
}

// doRequest 发送 HTTP 请求
func (c *OpenAIClient) doRequest(ctx context.Context, path string, body map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 错误 [%d]: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// doStreamRequest 发送流式 HTTP 请求
func (c *OpenAIClient) doStreamRequest(ctx context.Context, path string, body map[string]interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 错误 [%d]: %s", resp.StatusCode, string(respBody))
	}

	return resp, nil
}

// convertMessages 转换消息格式
func convertMessages(messages []Message) []map[string]interface{} {
	result := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		m := map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		}
		if msg.Name != "" {
			m["name"] = msg.Name
		}
		if len(msg.ToolCalls) > 0 {
			m["tool_calls"] = msg.ToolCalls
		}
		if msg.ToolCallID != "" {
			m["tool_call_id"] = msg.ToolCallID
		}
		result[i] = m
	}
	return result
}

// convertToolCalls 转换工具调用格式
func convertToolCalls(apiToolCalls []openAIToolCall) []ToolCall {
	if len(apiToolCalls) == 0 {
		return nil
	}

	result := make([]ToolCall, len(apiToolCalls))
	for i, tc := range apiToolCalls {
		result[i] = ToolCall{
			ID:   tc.ID,
			Type: tc.Type,
			Function: FunctionCall{
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			},
		}
	}
	return result
}

// OpenAI API 响应结构
type openAIResponse struct {
	Choices []struct {
		Message struct {
			Role      string           `json:"role"`
			Content   string           `json:"content"`
			ToolCalls []openAIToolCall `json:"tool_calls"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type openAIToolCall struct {
	Index    int    `json:"index"`
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type openAIStreamChunk struct {
	Choices []struct {
		Delta struct {
			Role      string           `json:"role"`
			Content   string           `json:"content"`
			ToolCalls []openAIToolCall `json:"tool_calls"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}
