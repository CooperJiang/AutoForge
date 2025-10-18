package executor

import (
	"auto-forge/internal/models"
	"auto-forge/pkg/agent/llm"
	"auto-forge/pkg/agent/prompt"
	"auto-forge/pkg/agent/registry"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"auto-forge/pkg/logger"
)

var log = logger.GetLogger()

// ReActExecutor ReAct 执行器
type ReActExecutor struct {
	llmClient    llm.LLMClient
	toolRegistry *registry.ToolRegistry
	maxSteps     int
	temperature  float64
}

// NewReActExecutor 创建 ReAct 执行器
func NewReActExecutor(
	llmClient llm.LLMClient,
	toolRegistry *registry.ToolRegistry,
	maxSteps int,
	temperature float64,
) *ReActExecutor {
	return &ReActExecutor{
		llmClient:    llmClient,
		toolRegistry: toolRegistry,
		maxSteps:     maxSteps,
		temperature:  temperature,
	}
}

// Execute 执行 ReAct 循环
func (e *ReActExecutor) Execute(
	ctx context.Context,
	userMessage string,
	conversationHistory string,
	allowedTools []string,
	streamCallback func(event StreamEvent),
) (*ExecutionResult, error) {
	startTime := time.Now()

	// 构建消息列表
	messages := e.buildMessages(userMessage, conversationHistory, allowedTools)

	// 获取工具定义
	toolDefinitions := e.toolRegistry.GetToolDefinitions(allowedTools)

	// 执行轨迹
	trace := &models.AgentTrace{
		Steps:     []models.AgentStep{},
		UsedTools: make(map[string]interface{}),
	}

	// ReAct 循环
	currentStep := 0
	for currentStep < e.maxSteps {
		currentStep++

		log.Info(ctx, "ReAct 步骤 %d 开始", currentStep)

		// 调用 LLM
		response, err := e.llmClient.Call(ctx, messages, &llm.CallOptions{
			Temperature: e.temperature,
			Tools:       toolDefinitions,
			ToolChoice:  "auto",
		})

		if err != nil {
			return nil, fmt.Errorf("LLM 调用失败: %w", err)
		}

		// 更新 token 使用
		if trace.TokenUsage == nil {
			trace.TokenUsage = &models.TokenUsage{}
		}
		trace.TokenUsage.PromptTokens += response.Usage.PromptTokens
		trace.TokenUsage.CompletionTokens += response.Usage.CompletionTokens
		trace.TokenUsage.TotalTokens += response.Usage.TotalTokens

		// 检查是否有工具调用
		if len(response.ToolCalls) == 0 {
			// 没有工具调用，说明 LLM 给出了最终答案
			log.Info(ctx, "ReAct 完成，原因: %s", response.FinishReason)

			trace.FinalAnswer = response.Content
			trace.FinishReason = "final"
			trace.TotalMs = time.Since(startTime).Milliseconds()

			// 发送最终事件
			if streamCallback != nil {
				streamCallback(StreamEvent{
					Type: "final",
					Data: map[string]interface{}{
						"answer":        response.Content,
						"finish_reason": "final",
						"trace":         trace,
						"token_usage":   trace.TokenUsage,
					},
				})
			}

			return &ExecutionResult{
				Trace:   trace,
				Success: true,
			}, nil
		}

		// 处理工具调用
		assistantMessage := llm.Message{
			Role:      "assistant",
			Content:   response.Content,
			ToolCalls: response.ToolCalls,
		}
		messages = append(messages, assistantMessage)

		// 执行每个工具调用
		for _, toolCall := range response.ToolCalls {
			stepStartTime := time.Now()

			// 发送步骤开始事件
			if streamCallback != nil {
				streamCallback(StreamEvent{
					Type: "step_start",
					Data: map[string]interface{}{
						"step": currentStep,
						"tool": toolCall.Function.Name,
					},
				})
			}

			// 解析参数
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				log.Error(ctx, "解析工具参数失败: %v", err)
				args = make(map[string]interface{})
			}

			log.Info(ctx, "执行工具: %s, 参数: %v", toolCall.Function.Name, args)

			// 执行工具
			toolResult, toolErr := e.toolRegistry.Execute(ctx, toolCall.Function.Name, args)

			// 格式化结果
			var observation string
			var toolOutput interface{}

			if toolErr != nil {
				observation = fmt.Sprintf("Error: %s", toolErr.Error())
				log.Error(ctx, "工具执行失败: %v", toolErr)
			} else {
				toolOutput = toolResult
				observation = registry.FormatToolResult(toolResult)
				log.Info(ctx, "工具执行成功，结果长度: %d", len(observation))
			}

			elapsedMs := time.Since(stepStartTime).Milliseconds()

			// 记录步骤
			step := models.AgentStep{
				Step: currentStep,
				Action: &models.AgentAction{
					Type: "action",
					Tool: toolCall.Function.Name,
					Args: args,
				},
				Observation: observation,
				ToolOutput:  toolOutput.(map[string]interface{}),
				ElapsedMs:   elapsedMs,
				Timestamp:   time.Now().Format(time.RFC3339),
			}

			if toolErr != nil {
				step.Error = toolErr.Error()
			}

			trace.Steps = append(trace.Steps, step)

			// 更新工具使用统计
			if _, ok := trace.UsedTools[toolCall.Function.Name]; !ok {
				trace.UsedTools[toolCall.Function.Name] = map[string]interface{}{
					"count":    0,
					"total_ms": int64(0),
				}
			}
			stats := trace.UsedTools[toolCall.Function.Name].(map[string]interface{})
			stats["count"] = stats["count"].(int) + 1
			stats["total_ms"] = stats["total_ms"].(int64) + elapsedMs
			trace.UsedTools[toolCall.Function.Name] = stats

			// 发送步骤结束事件
			if streamCallback != nil {
				streamCallback(StreamEvent{
					Type: "step_end",
					Data: map[string]interface{}{
						"step":        currentStep,
						"tool":        toolCall.Function.Name,
						"observation": observation,
						"elapsed_ms":  elapsedMs,
					},
				})
			}

			// 将工具结果添加到消息列表
			toolMessage := llm.Message{
				Role:       "tool",
				Content:    observation,
				ToolCallID: toolCall.ID,
				Name:       toolCall.Function.Name,
			}
			messages = append(messages, toolMessage)
		}
	}

	// 达到最大步骤数
	log.Warn(ctx, "达到最大步骤数 %d", e.maxSteps)

	trace.FinalAnswer = "已达到最大步骤数限制，任务可能未完全完成。"
	trace.FinishReason = "max_steps"
	trace.TotalMs = time.Since(startTime).Milliseconds()

	// 发送最终事件（即使未完全完成）
	if streamCallback != nil {
		streamCallback(StreamEvent{
			Type: "final",
			Data: map[string]interface{}{
				"answer":        trace.FinalAnswer,
				"finish_reason": "max_steps",
				"trace":         trace,
				"token_usage":   trace.TokenUsage,
			},
		})
	}

	return &ExecutionResult{
		Trace:   trace,
		Success: false,
		Error:   "达到最大步骤数限制",
	}, nil
}

// buildMessages 构建消息列表
func (e *ReActExecutor) buildMessages(
	userMessage string,
	conversationHistory string,
	allowedTools []string,
) []llm.Message {
	messages := []llm.Message{
		{
			Role:    "system",
			Content: prompt.SystemPrompt.Template,
		},
	}

	// 添加对话历史（如果有）
	if conversationHistory != "" {
		messages = append(messages, llm.Message{
			Role:    "user",
			Content: "Previous conversation context:\n" + conversationHistory,
		})
	}

	// 添加用户消息
	messages = append(messages, llm.Message{
		Role:    "user",
		Content: userMessage,
	})

	return messages
}

// StreamEvent 流式事件
type StreamEvent struct {
	Type string                 // step_start, step_end, final, error
	Data map[string]interface{} // 事件数据
}

// ExecutionResult 执行结果
type ExecutionResult struct {
	Trace   *models.AgentTrace // 执行轨迹
	Success bool               // 是否成功
	Error   string             // 错误信息
}
