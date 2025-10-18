package executor

import (
	"auto-forge/internal/models"
	"auto-forge/pkg/agent/llm"
	"auto-forge/pkg/agent/prompt"
	"auto-forge/pkg/agent/registry"
	"auto-forge/pkg/utools"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// PlanExecutor Plan 执行器
type PlanExecutor struct {
	llmClient    llm.LLMClient
	toolRegistry *registry.ToolRegistry
	temperature  float64
}

// NewPlanExecutor 创建 Plan 执行器
func NewPlanExecutor(
	llmClient llm.LLMClient,
	toolRegistry *registry.ToolRegistry,
	temperature float64,
) *PlanExecutor {
	return &PlanExecutor{
		llmClient:    llmClient,
		toolRegistry: toolRegistry,
		temperature:  temperature,
	}
}

// Execute 执行 Plan 模式
func (e *PlanExecutor) Execute(
	ctx context.Context,
	userMessage string,
	conversationHistory string,
	allowedTools []string,
	maxSteps int,
	streamCallback func(event StreamEvent),
) (*ExecutionResult, error) {
	startTime := time.Now()

	// Step 1: 生成计划
	plan, err := e.generatePlan(ctx, userMessage, conversationHistory, allowedTools)
	if err != nil {
		return nil, fmt.Errorf("生成计划失败: %w", err)
	}

	log.Info(ctx, "生成计划成功，共 %d 步", len(plan.Steps))

	// 发送计划开始事件
	if streamCallback != nil {
		streamCallback(StreamEvent{
			Type: "plan_start",
			Data: map[string]interface{}{
				"plan": plan,
			},
		})
	}

	// Step 2: 执行计划
	trace := &models.AgentTrace{
		Steps:     []models.AgentStep{},
		UsedTools: make(map[string]interface{}),
	}

	for i, planStep := range plan.Steps {
		if i >= maxSteps {
			log.Warn(ctx, "达到最大步骤数 %d", maxSteps)
			break
		}

		// 更新计划步骤状态
		plan.Steps[i].Status = "running"
		if streamCallback != nil {
			streamCallback(StreamEvent{
				Type: "plan_step",
				Data: map[string]interface{}{
					"step_index": i,
					"status":     "running",
				},
			})
		}

		// 执行步骤（传递之前的步骤结果）
		step, err := e.executeStep(ctx, planStep, i+1, userMessage, trace.Steps, streamCallback)
		if err != nil {
			log.Error(ctx, "执行步骤 %d 失败: %v", i+1, err)
			plan.Steps[i].Status = "failed"

			if streamCallback != nil {
				streamCallback(StreamEvent{
					Type: "plan_step",
					Data: map[string]interface{}{
						"step_index": i,
						"status":     "failed",
						"error":      err.Error(),
					},
				})
			}

			step.Error = err.Error()
			trace.Steps = append(trace.Steps, *step)

			// 检查是否是关键步骤失败（工具执行失败）
			// 如果后续步骤依赖这个步骤的输出，应该跳过它们
			if step.ToolOutput == nil || len(step.ToolOutput) == 0 {
				log.Info(ctx, "步骤 %d 失败且无输出，标记后续依赖步骤为跳过", i+1)
				// 标记剩余步骤为 skipped
				for j := i + 1; j < len(plan.Steps); j++ {
					plan.Steps[j].Status = "skipped"
					if streamCallback != nil {
						streamCallback(StreamEvent{
							Type: "plan_step",
							Data: map[string]interface{}{
								"step_index": j,
								"status":     "skipped",
								"error":      fmt.Sprintf("跳过：步骤 %d 失败", i+1),
							},
						})
					}
				}
				break // 终止执行
			}
		} else {
			plan.Steps[i].Status = "completed"
			if streamCallback != nil {
				streamCallback(StreamEvent{
					Type: "plan_step",
					Data: map[string]interface{}{
						"step_index": i,
						"status":     "completed",
					},
				})
			}
			trace.Steps = append(trace.Steps, *step)
		}

		// 更新工具使用统计
		if step.Action != nil && step.Action.Tool != "" {
			if _, ok := trace.UsedTools[step.Action.Tool]; !ok {
				trace.UsedTools[step.Action.Tool] = map[string]interface{}{
					"count":    0,
					"total_ms": int64(0),
				}
			}
			if stats, ok := trace.UsedTools[step.Action.Tool].(map[string]interface{}); ok {
				stats["count"] = stats["count"].(int) + 1
				stats["total_ms"] = stats["total_ms"].(int64) + step.ElapsedMs
				trace.UsedTools[step.Action.Tool] = stats
			}
		}
	}

	// Step 3: 生成最终答案
	finalAnswer, err := e.generateFinalAnswer(ctx, userMessage, trace)
	if err != nil {
		log.Error(ctx, "生成最终答案失败: %v", err)
		finalAnswer = "任务已执行完成，但生成总结失败。"
	}

	trace.FinalAnswer = finalAnswer
	trace.FinishReason = "final"
	trace.TotalMs = time.Since(startTime).Milliseconds()

	// 发送最终事件
	if streamCallback != nil {
		streamCallback(StreamEvent{
			Type: "final",
			Data: map[string]interface{}{
				"answer":        finalAnswer,
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

// generatePlan 生成执行计划
func (e *PlanExecutor) generatePlan(
	ctx context.Context,
	userMessage string,
	conversationHistory string,
	allowedTools []string,
) (*models.AgentPlan, error) {
	// 获取工具定义
	log.Info(ctx, "allowedTools 参数: %v", allowedTools)
	allToolNames := e.toolRegistry.ListTools()
	log.Info(ctx, "注册表中的所有工具: %v", allToolNames)
	toolDefinitions := e.toolRegistry.GetToolDefinitions(allowedTools)
	log.Info(ctx, "可用工具数量: %d", len(toolDefinitions))

	// 转换为 map 格式
	toolDefs := make([]map[string]interface{}, len(toolDefinitions))
	for i, td := range toolDefinitions {
		// 将 FunctionDefinition 结构体转换为 map
		functionMap := map[string]interface{}{
			"name":        td.Function.Name,
			"description": td.Function.Description,
			"parameters":  td.Function.Parameters,
		}

		toolDefs[i] = map[string]interface{}{
			"type":     td.Type,
			"function": functionMap,
		}
	}

	// 构建提示词
	toolDefsStr := prompt.FormatToolDefinitions(toolDefs)

	promptText := prompt.PlanPrompt.Render(map[string]string{
		"tool_definitions": toolDefsStr,
		"user_message":     userMessage,
	})

	// 调用 LLM
	messages := []llm.Message{
		{
			Role:    "system",
			Content: "You are a planning AI. Generate execution plans in JSON format.",
		},
		{
			Role:    "user",
			Content: promptText,
		},
	}

	response, err := e.llmClient.Call(ctx, messages, &llm.CallOptions{
		Temperature:    e.temperature,
		ResponseFormat: "json_object",
	})

	if err != nil {
		log.Error(ctx, "LLM 调用失败: %v", err)
		return nil, fmt.Errorf("LLM 调用失败: %w", err)
	}

	log.Info(ctx, "LLM 返回内容长度: %d, 内容: %s", len(response.Content), response.Content)

	// 解析计划
	var planData struct {
		Steps []struct {
			Step        int    `json:"step"`
			Description string `json:"description"`
			Tool        string `json:"tool"`
			Reasoning   string `json:"reasoning"`
		} `json:"steps"`
	}

	if err := json.Unmarshal([]byte(response.Content), &planData); err != nil {
		log.Error(ctx, "解析计划失败: %v, 原始内容: %s", err, response.Content)
		return nil, fmt.Errorf("解析计划失败: %w", err)
	}

	log.Info(ctx, "解析到 %d 个步骤", len(planData.Steps))

	// 构建计划
	plan := &models.AgentPlan{
		Steps:       make([]models.AgentPlanStep, len(planData.Steps)),
		TotalSteps:  len(planData.Steps),
		CreatedAt:   time.Now().Format(time.RFC3339),
		GeneratedBy: e.llmClient.GetModelInfo().Model,
	}

	for i, s := range planData.Steps {
		plan.Steps[i] = models.AgentPlanStep{
			Step:        s.Step,
			Description: s.Description,
			Tool:        s.Tool,
			Status:      "pending",
		}
	}

	return plan, nil
}

// executeStep 执行单个步骤
func (e *PlanExecutor) executeStep(
	ctx context.Context,
	planStep models.AgentPlanStep,
	stepIndex int,
	userMessage string,
	previousSteps []models.AgentStep,
	streamCallback func(event StreamEvent),
) (*models.AgentStep, error) {
	stepStartTime := time.Now()

	// 发送步骤开始事件
	if streamCallback != nil {
		streamCallback(StreamEvent{
			Type: "step_start",
			Data: map[string]interface{}{
				"step": stepIndex,
				"tool": planStep.Tool,
			},
		})
	}

	// 如果没有指定工具，跳过执行
	if planStep.Tool == "" {
		observation := planStep.Description
		elapsedMs := time.Since(stepStartTime).Milliseconds()

		step := &models.AgentStep{
			Step:        stepIndex,
			Observation: observation,
			ElapsedMs:   elapsedMs,
			Timestamp:   time.Now().Format(time.RFC3339),
		}

		if streamCallback != nil {
			streamCallback(StreamEvent{
				Type: "step_end",
				Data: map[string]interface{}{
					"step":        stepIndex,
					"observation": observation,
					"elapsed_ms":  elapsedMs,
				},
			})
		}

		return step, nil
	}

	// 使用 LLM 生成工具参数（传递之前的步骤结果）
	args, err := e.generateToolArgs(ctx, planStep, userMessage, previousSteps)
	if err != nil {
		log.Error(ctx, "生成工具参数失败: %v", err)
		elapsedMs := time.Since(stepStartTime).Milliseconds()

		step := &models.AgentStep{
			Step: stepIndex,
			Action: &models.AgentAction{
				Type: "action",
				Tool: planStep.Tool,
				Args: nil,
			},
			Observation: fmt.Sprintf("Error: %s", err.Error()),
			Error:       err.Error(),
			ElapsedMs:   elapsedMs,
			Timestamp:   time.Now().Format(time.RFC3339),
		}

		if streamCallback != nil {
			streamCallback(StreamEvent{
				Type: "step_end",
				Data: map[string]interface{}{
					"step":        stepIndex,
					"tool":        planStep.Tool,
					"observation": step.Observation,
					"elapsed_ms":  elapsedMs,
				},
			})
		}

		return step, err
	}

	log.Info(ctx, "执行工具: %s, 参数: %v", planStep.Tool, args)

	// 执行工具
	toolResult, toolErr := e.toolRegistry.Execute(ctx, planStep.Tool, args)

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

	// 构建步骤
	step := &models.AgentStep{
		Step: stepIndex,
		Action: &models.AgentAction{
			Type: "action",
			Tool: planStep.Tool,
			Args: args,
		},
		Observation: observation,
		ElapsedMs:   elapsedMs,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	// 只有在工具执行成功时才设置 ToolOutput
	if toolErr == nil && toolOutput != nil {
		// toolOutput 是 *utools.ExecutionResult 类型
		if result, ok := toolOutput.(*utools.ExecutionResult); ok && result.Output != nil {
			step.ToolOutput = result.Output
		} else if outputMap, ok := toolOutput.(map[string]interface{}); ok {
			step.ToolOutput = outputMap
		}
	}

	if toolErr != nil {
		step.Error = toolErr.Error()
	}

	// 发送步骤结束事件
	if streamCallback != nil {
		streamCallback(StreamEvent{
			Type: "step_end",
			Data: map[string]interface{}{
				"step":        stepIndex,
				"tool":        planStep.Tool,
				"observation": observation,
				"elapsed_ms":  elapsedMs,
			},
		})
	}

	if toolErr != nil {
		return step, toolErr
	}

	return step, nil
}

// generateToolArgs 使用 LLM 生成工具参数
func (e *PlanExecutor) generateToolArgs(
	ctx context.Context,
	planStep models.AgentPlanStep,
	userMessage string,
	previousSteps []models.AgentStep,
) (map[string]interface{}, error) {
	// 获取工具定义
	tool, err := e.toolRegistry.GetTool(planStep.Tool)
	if err != nil {
		return nil, err
	}

	schema := tool.GetSchema()
	schemaJSON, _ := json.Marshal(schema)

	// 构建之前步骤的上下文
	var contextBuilder strings.Builder
	if len(previousSteps) > 0 {
		contextBuilder.WriteString("\n\nPrevious Steps Results:\n")
		for _, prevStep := range previousSteps {
			contextBuilder.WriteString(fmt.Sprintf("\nStep %d (%s):\n", prevStep.Step, prevStep.Action.Tool))
			if prevStep.Error != "" {
				contextBuilder.WriteString(fmt.Sprintf("  Status: FAILED\n  Error: %s\n", prevStep.Error))
			} else {
				contextBuilder.WriteString("  Status: SUCCESS\n")
				if prevStep.ToolOutput != nil {
					outputJSON, _ := json.MarshalIndent(prevStep.ToolOutput, "  ", "  ")
					contextBuilder.WriteString(fmt.Sprintf("  Output: %s\n", string(outputJSON)))
				}
			}
		}
	}

	// 构建提示词
	promptText := fmt.Sprintf(`Generate parameters for the following tool call:

Tool: %s
Description: %s
Task: %s
User's Original Request: %s
%s

CRITICAL INSTRUCTIONS: 
1. Look at "Previous Steps Results" above - each step shows its Output with actual data.
2. Extract values from "Output" field using the access paths shown.
3. For example, if you see:
   Step 1 (openai_image):
     Output: { "url": "https://filesystem.site/cdn/xxx.webp", "revised_prompt": "..." }
   Then extract: "https://filesystem.site/cdn/xxx.webp" from output.url

4. DO NOT use placeholder values like "https://example.com" or empty strings.
5. DO NOT copy the parameter schema definitions into your response.
6. Return ONLY a JSON object with actual parameter values extracted from previous outputs.

Parameter Schema (for reference only):
%s

CORRECT example (using real data from previous step):
{
  "url": "https://filesystem.site/cdn/20251018/bmXOwy6fccuyb9gQmH72TMaeupUKvW.webp"
}

WRONG examples:
{
  "url": ""  ← Empty string
}
{
  "url": "https://example.com/xxx.webp"  ← Placeholder URL
}
{
  "file": {"type": "object", "title": "..."}  ← Schema definition
}`,
		planStep.Tool,
		tool.GetMetadata().Description,
		planStep.Description,
		userMessage,
		contextBuilder.String(),
		string(schemaJSON),
	)

	messages := []llm.Message{
		{
			Role:    "user",
			Content: promptText,
		},
	}

	response, err := e.llmClient.Call(ctx, messages, &llm.CallOptions{
		Temperature:    0.3, // 低温度以获得更确定的结果
		ResponseFormat: "json_object",
	})

	if err != nil {
		return nil, err
	}

	// 解析参数
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(response.Content), &args); err != nil {
		return nil, fmt.Errorf("解析工具参数失败: %w", err)
	}

	return args, nil
}

// generateFinalAnswer 生成最终答案
func (e *PlanExecutor) generateFinalAnswer(
	ctx context.Context,
	userMessage string,
	trace *models.AgentTrace,
) (string, error) {
	// 格式化执行轨迹
	steps := make([]map[string]interface{}, len(trace.Steps))
	for i, step := range trace.Steps {
		steps[i] = map[string]interface{}{
			"step":        step.Step,
			"action":      step.Action,
			"observation": step.Observation,
		}
	}

	traceStr := prompt.FormatExecutionTrace(steps)

	// 构建提示词
	promptText := prompt.SummaryPrompt.Render(map[string]string{
		"user_message":    userMessage,
		"execution_trace": traceStr,
	})

	messages := []llm.Message{
		{
			Role:    "user",
			Content: promptText,
		},
	}

	response, err := e.llmClient.Call(ctx, messages, &llm.CallOptions{
		Temperature: 0.7,
	})

	if err != nil {
		return "", err
	}

	return response.Content, nil
}
