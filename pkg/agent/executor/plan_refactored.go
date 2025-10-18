package executor

import (
	"auto-forge/internal/models"
	"auto-forge/pkg/agent/llm"
	"auto-forge/pkg/agent/prompt"
	"auto-forge/pkg/agent/registry"
	"auto-forge/pkg/agent/tooling"
	"auto-forge/pkg/logger"
	"auto-forge/pkg/utools"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// PlanExecutorV2 Plan 执行器（重构版）
type PlanExecutorV2 struct {
	llmClient     llm.LLMClient
	toolRegistry  *registry.ToolRegistry
	stepExecutor  *StepExecutor
	validator     *tooling.PlanValidator
	temperature   float64
	enableRetry   bool
	enableTimeout bool
}

// NewPlanExecutorV2 创建 Plan 执行器（重构版）
func NewPlanExecutorV2(
	llmClient llm.LLMClient,
	toolRegistry *registry.ToolRegistry,
	temperature float64,
) *PlanExecutorV2 {
	return &PlanExecutorV2{
		llmClient:     llmClient,
		toolRegistry:  toolRegistry,
		stepExecutor:  NewStepExecutor(toolRegistry),
		validator:     tooling.NewPlanValidator(toolRegistry),
		temperature:   temperature,
		enableRetry:   true, // 默认启用重试
		enableTimeout: true, // 默认启用超时
	}
}

// Execute 执行 Plan 模式（重构版）
func (e *PlanExecutorV2) Execute(
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

	logger.Info("生成计划成功，共 %d 步", len(plan.Steps))

	// Step 2: 验证计划
	validationResult := e.validator.Validate(plan)
	if !validationResult.Valid {
		logger.Warn("计划验证失败: %v", validationResult.Errors)
		// 可以选择返回错误或继续执行
		// 这里我们记录警告但继续执行
		for _, warning := range validationResult.Warnings {
			logger.Warn("计划警告: %s", warning)
		}
	}

	// 发送计划开始事件
	if streamCallback != nil {
		streamCallback(StreamEvent{
			Type: "plan_start",
			Data: map[string]interface{}{
				"plan":       plan,
				"validation": validationResult,
			},
		})
	}

	// Step 3: 执行计划
	trace := &models.AgentTrace{
		Steps:     []models.AgentStep{},
		UsedTools: make(map[string]interface{}),
	}

	for i, planStep := range plan.Steps {
		if i >= maxSteps {
			logger.Warn("达到最大步骤数 %d", maxSteps)
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

		// 执行步骤（使用新的模块化执行器）
		stepResult := e.executeStepV2(ctx, planStep, i+1, userMessage, trace.Steps, streamCallback)

		if stepResult.Error != nil {
			logger.Error("执行步骤 %d 失败: %v", i+1, stepResult.Error)
			plan.Steps[i].Status = "failed"

			if streamCallback != nil {
				streamCallback(StreamEvent{
					Type: "plan_step",
					Data: map[string]interface{}{
						"step_index": i,
						"status":     "failed",
						"error":      stepResult.Error.Error(),
					},
				})
			}

			stepResult.Step.Error = stepResult.Error.Error()
			trace.Steps = append(trace.Steps, *stepResult.Step)

			// 智能失败处理：检查是否应该跳过后续步骤
			if e.shouldSkipRemainingSteps(stepResult.Step, plan.Steps[i+1:]) {
				logger.Info("步骤 %d 失败且无输出，标记后续依赖步骤为跳过", i+1)
				e.skipRemainingSteps(plan, i+1, fmt.Sprintf("跳过：步骤 %d 失败", i+1), streamCallback)
				break
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
			trace.Steps = append(trace.Steps, *stepResult.Step)
		}

		// 更新工具使用统计
		e.updateToolStats(trace, stepResult.Step)
	}

	// Step 4: 生成最终答案
	finalAnswer, err := e.generateFinalAnswer(ctx, userMessage, trace)
	if err != nil {
		logger.Error("生成最终答案失败: %v", err)
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

// executeStepV2 执行单个步骤（使用新的模块化执行器）
func (e *PlanExecutorV2) executeStepV2(
	ctx context.Context,
	planStep models.AgentPlanStep,
	stepIndex int,
	userMessage string,
	previousSteps []models.AgentStep,
	streamCallback func(event StreamEvent),
) *ExecuteStepResult {
	// 如果没有指定工具，跳过执行
	if planStep.Tool == "" {
		return e.stepExecutor.createSkippedStep(&ExecuteStepRequest{
			PlanStep:    planStep,
			StepIndex:   stepIndex,
			UserMessage: userMessage,
		}, time.Now())
	}

	// 生成工具参数
	args, err := e.generateToolArgs(ctx, planStep, stepIndex, userMessage, previousSteps)
	if err != nil {
		return e.stepExecutor.createErrorStep(&ExecuteStepRequest{
			PlanStep:    planStep,
			StepIndex:   stepIndex,
			UserMessage: userMessage,
		}, time.Now(), err)
	}

	// 获取工具
	tool, err := e.toolRegistry.GetTool(planStep.Tool)
	if err != nil {
		return e.stepExecutor.createErrorStep(&ExecuteStepRequest{
			PlanStep:    planStep,
			StepIndex:   stepIndex,
			UserMessage: userMessage,
		}, time.Now(), err)
	}

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

	// 获取工具配置
	toolConfig := e.getToolConfig(tool)

	// 创建工具执行器
	toolExecutor := tooling.NewToolExecutor(toolConfig)

	// 执行工具（带超时和重试）
	logger.Info("执行工具: %s, 参数: %v", planStep.Tool, args)

	stepStartTime := time.Now()
	execResult := toolExecutor.ExecuteWithProgress(ctx, tool, args, func(attempt int, message string) {
		// 报告进度
		if streamCallback != nil {
			streamCallback(StreamEvent{
				Type: "tool_progress",
				Data: map[string]interface{}{
					"step":    stepIndex,
					"tool":    planStep.Tool,
					"attempt": attempt,
					"message": message,
				},
			})
		}
	})

	// 处理执行结果
	return e.processExecutionResult(planStep, stepIndex, args, stepStartTime, execResult, streamCallback)
}

// processExecutionResult 处理执行结果
func (e *PlanExecutorV2) processExecutionResult(
	planStep models.AgentPlanStep,
	stepIndex int,
	args map[string]interface{},
	startTime time.Time,
	execResult *tooling.ExecutionResult,
	streamCallback func(event StreamEvent),
) *ExecuteStepResult {
	elapsedMs := time.Since(startTime).Milliseconds()

	// 构建步骤
	step := &models.AgentStep{
		Step: stepIndex,
		Action: &models.AgentAction{
			Type: "action",
			Tool: planStep.Tool,
			Args: args,
		},
		ElapsedMs: elapsedMs,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 处理错误
	if execResult.Error != nil {
		step.Error = execResult.Error.Error()
		step.Observation = fmt.Sprintf("Error: %s", execResult.Error.Error())
		logger.Error("工具执行失败（尝试 %d 次）: %v", execResult.Attempts, execResult.Error)
	} else {
		// 成功
		toolOutput := execResult.Output
		step.Observation = registry.FormatToolResult(toolOutput)

		// 设置 ToolOutput
		if toolOutput != nil {
			if result, ok := toolOutput.(*utools.ExecutionResult); ok && result.Output != nil {
				step.ToolOutput = result.Output
			} else if outputMap, ok := toolOutput.(map[string]interface{}); ok {
				step.ToolOutput = outputMap
			}
		}

		logger.Info("工具执行成功（尝试 %d 次），结果长度: %d",
			execResult.Attempts, len(step.Observation))
	}

	// 发送步骤结束事件
	if streamCallback != nil {
		streamCallback(StreamEvent{
			Type: "step_end",
			Data: map[string]interface{}{
				"step":        stepIndex,
				"tool":        planStep.Tool,
				"observation": step.Observation,
				"elapsed_ms":  elapsedMs,
				"attempts":    execResult.Attempts,
			},
		})
	}

	return &ExecuteStepResult{
		Step:  step,
		Error: execResult.Error,
	}
}

// shouldSkipRemainingSteps 判断是否应该跳过剩余步骤
func (e *PlanExecutorV2) shouldSkipRemainingSteps(
	failedStep *models.AgentStep,
	remainingSteps []models.AgentPlanStep,
) bool {
	// 如果失败的步骤没有输出，且后续还有步骤，则跳过
	return (failedStep.ToolOutput == nil || len(failedStep.ToolOutput) == 0) && len(remainingSteps) > 0
}

// skipRemainingSteps 跳过剩余步骤
func (e *PlanExecutorV2) skipRemainingSteps(
	plan *models.AgentPlan,
	startIndex int,
	reason string,
	streamCallback func(event StreamEvent),
) {
	for j := startIndex; j < len(plan.Steps); j++ {
		plan.Steps[j].Status = "skipped"
		if streamCallback != nil {
			streamCallback(StreamEvent{
				Type: "plan_step",
				Data: map[string]interface{}{
					"step_index": j,
					"status":     "skipped",
					"error":      reason,
				},
			})
		}
	}
}

// updateToolStats 更新工具使用统计
func (e *PlanExecutorV2) updateToolStats(trace *models.AgentTrace, step *models.AgentStep) {
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

// getToolConfig 获取工具配置
func (e *PlanExecutorV2) getToolConfig(tool interface{}) *tooling.ExecutionConfig {
	// 尝试从工具获取配置
	if configurable, ok := tool.(tooling.ConfigurableTool); ok {
		return configurable.GetExecutionConfig()
	}

	// 返回默认配置
	config := tooling.DefaultExecutionConfig()

	// 根据全局设置调整配置
	if !e.enableRetry {
		config.Retry = nil
	}
	if !e.enableTimeout {
		config.TimeoutSeconds = 0
	}

	return config
}

// generatePlan 生成执行计划（保持原有逻辑）
func (e *PlanExecutorV2) generatePlan(
	ctx context.Context,
	userMessage string,
	conversationHistory string,
	allowedTools []string,
) (*models.AgentPlan, error) {
	// 获取工具定义
	logger.Info("allowedTools 参数: %v", allowedTools)
	allToolNames := e.toolRegistry.ListTools()
	logger.Info("注册表中的所有工具: %v", allToolNames)
	toolDefinitions := e.toolRegistry.GetToolDefinitions(allowedTools)
	logger.Info("可用工具数量: %d", len(toolDefinitions))

	// 转换为 map 格式
	toolDefs := make([]map[string]interface{}, len(toolDefinitions))
	for i, td := range toolDefinitions {
		// 将 FunctionDefinition 结构体转换为 map
		functionMap := map[string]interface{}{
			"name":        td.Function.Name,
			"description": td.Function.Description,
			"parameters":  td.Function.Parameters,
			"metadata":    td.Function.Metadata,
		}

		toolDefs[i] = map[string]interface{}{
			"type":     td.Type,
			"function": functionMap,
		}
	}

	// 构建提示词
	toolDefsStr := prompt.FormatToolDefinitions(toolDefs)
	logger.Info("工具定义长度: %d 字符", len(toolDefsStr))

	promptText := prompt.PlanPrompt.Render(map[string]string{
		"tool_definitions": toolDefsStr,
		"user_message":     userMessage,
	})
	logger.Info("完整提示词长度: %d 字符", len(promptText))

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
		logger.Error("LLM 调用失败: %v", err)
		return nil, fmt.Errorf("LLM 调用失败: %w", err)
	}

	logger.Info("LLM 返回内容长度: %d", len(response.Content))

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
		logger.Error("解析计划失败: %v, 原始内容: %s", err, response.Content)
		return nil, fmt.Errorf("解析计划失败: %w", err)
	}

	logger.Info("解析到 %d 个步骤", len(planData.Steps))

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

// generateToolArgs 生成工具参数（保持原有逻辑）
func (e *PlanExecutorV2) generateToolArgs(
	ctx context.Context,
	planStep models.AgentPlanStep,
	stepIndex int,
	userMessage string,
	previousSteps []models.AgentStep,
) (map[string]interface{}, error) {
	// 获取工具
	tool, err := e.toolRegistry.GetTool(planStep.Tool)
	if err != nil {
		return nil, fmt.Errorf("工具不存在: %s", planStep.Tool)
	}

	// 获取工具的参数 schema
	schema := tool.GetSchema()
	schemaJSON, _ := json.MarshalIndent(schema, "", "  ")

	// 构建上下文
	var contextBuilder strings.Builder
	contextBuilder.WriteString(BuildPreviousStepsContext(previousSteps))

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
6. The Parameter Schema shows the TOP-LEVEL parameters. Each parameter is a direct property of the root object.
7. Return ONLY a JSON object with actual parameter values extracted from previous outputs.

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

	logger.Info("生成工具参数的提示词:\n%s", promptText)

	// 调用 LLM
	messages := []llm.Message{
		{
			Role:    "system",
			Content: "You are a parameter generation AI. Generate tool parameters in JSON format based on context.",
		},
		{
			Role:    "user",
			Content: promptText,
		},
	}

	response, err := e.llmClient.Call(ctx, messages, &llm.CallOptions{
		Temperature:    0.1, // 使用较低的温度以获得更确定的结果
		ResponseFormat: "json_object",
	})

	if err != nil {
		return nil, fmt.Errorf("LLM 调用失败: %w", err)
	}

	logger.Info("LLM 生成的参数 JSON:\n%s", response.Content)

	// 解析参数
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(response.Content), &args); err != nil {
		return nil, fmt.Errorf("解析参数失败: %w", err)
	}

	argsJSON, _ := json.MarshalIndent(args, "", "  ")
	logger.Info("解析后的参数:\n%s", string(argsJSON))

	return args, nil
}

// generateFinalAnswer 生成最终答案（保持原有逻辑）
func (e *PlanExecutorV2) generateFinalAnswer(
	ctx context.Context,
	userMessage string,
	trace *models.AgentTrace,
) (string, error) {
	// 构建步骤摘要
	var stepsBuilder strings.Builder
	for _, step := range trace.Steps {
		stepsBuilder.WriteString(fmt.Sprintf("\nStep %d: %s\n", step.Step, step.Action.Tool))
		if step.Error != "" {
			stepsBuilder.WriteString(fmt.Sprintf("  Status: FAILED\n  Error: %s\n", step.Error))
		} else {
			stepsBuilder.WriteString(fmt.Sprintf("  Status: SUCCESS\n  Result: %s\n", step.Observation))
		}
	}

	promptText := prompt.SummaryPrompt.Render(map[string]string{
		"user_message": userMessage,
		"steps":        stepsBuilder.String(),
	})

	messages := []llm.Message{
		{
			Role:    "system",
			Content: "You are a helpful AI assistant. Summarize the execution results.",
		},
		{
			Role:    "user",
			Content: promptText,
		},
	}

	response, err := e.llmClient.Call(ctx, messages, &llm.CallOptions{
		Temperature: e.temperature,
	})

	if err != nil {
		return "", err
	}

	return response.Content, nil
}
