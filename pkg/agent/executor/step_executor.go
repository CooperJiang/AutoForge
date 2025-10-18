package executor

import (
	"auto-forge/internal/models"
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

// StepExecutor 步骤执行器
type StepExecutor struct {
	toolRegistry *registry.ToolRegistry
	validator    *tooling.PlanValidator
}

// NewStepExecutor 创建步骤执行器
func NewStepExecutor(toolRegistry *registry.ToolRegistry) *StepExecutor {
	return &StepExecutor{
		toolRegistry: toolRegistry,
		validator:    tooling.NewPlanValidator(toolRegistry),
	}
}

// ExecuteStepRequest 执行步骤请求
type ExecuteStepRequest struct {
	PlanStep       models.AgentPlanStep
	StepIndex      int
	UserMessage    string
	PreviousSteps  []models.AgentStep
	StreamCallback func(event StreamEvent)
}

// ExecuteStepResult 执行步骤结果
type ExecuteStepResult struct {
	Step  *models.AgentStep
	Error error
}

// ExecuteStep 执行单个步骤（模块化版本）
func (e *StepExecutor) ExecuteStep(
	ctx context.Context,
	req *ExecuteStepRequest,
) *ExecuteStepResult {
	stepStartTime := time.Now()

	// 发送步骤开始事件
	if req.StreamCallback != nil {
		req.StreamCallback(StreamEvent{
			Type: "step_start",
			Data: map[string]interface{}{
				"step": req.StepIndex,
				"tool": req.PlanStep.Tool,
			},
		})
	}

	// 如果没有指定工具，跳过执行
	if req.PlanStep.Tool == "" {
		return e.createSkippedStep(req, stepStartTime)
	}

	// 生成工具参数
	args, err := e.generateToolArgs(ctx, req)
	if err != nil {
		return e.createErrorStep(req, stepStartTime, err)
	}

	// 获取工具
	tool, err := e.toolRegistry.GetTool(req.PlanStep.Tool)
	if err != nil {
		return e.createErrorStep(req, stepStartTime, err)
	}

	// 获取工具配置
	toolConfig := e.getToolConfig(tool)

	// 创建工具执行器
	toolExecutor := tooling.NewToolExecutor(toolConfig)

	// 执行工具（带超时和重试）
	logger.Info("执行工具: %s, 参数: %v", req.PlanStep.Tool, args)

	execResult := toolExecutor.ExecuteWithProgress(ctx, tool, args, func(attempt int, message string) {
		// 报告进度
		if req.StreamCallback != nil {
			req.StreamCallback(StreamEvent{
				Type: "tool_progress",
				Data: map[string]interface{}{
					"step":    req.StepIndex,
					"tool":    req.PlanStep.Tool,
					"attempt": attempt,
					"message": message,
				},
			})
		}
	})

	// 处理执行结果
	return e.processExecutionResult(req, stepStartTime, args, execResult)
}

// generateToolArgs 生成工具参数
func (e *StepExecutor) generateToolArgs(
	ctx context.Context,
	req *ExecuteStepRequest,
) (map[string]interface{}, error) {
	// 这里调用原来的 generateToolArgs 逻辑
	// 为了保持模块化，这里只是一个占位符
	// 实际实现会在集成时完成
	return map[string]interface{}{}, nil
}

// processExecutionResult 处理执行结果
func (e *StepExecutor) processExecutionResult(
	req *ExecuteStepRequest,
	startTime time.Time,
	args map[string]interface{},
	execResult *tooling.ExecutionResult,
) *ExecuteStepResult {
	elapsedMs := time.Since(startTime).Milliseconds()

	// 构建步骤
	step := &models.AgentStep{
		Step: req.StepIndex,
		Action: &models.AgentAction{
			Type: "action",
			Tool: req.PlanStep.Tool,
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
	if req.StreamCallback != nil {
		req.StreamCallback(StreamEvent{
			Type: "step_end",
			Data: map[string]interface{}{
				"step":        req.StepIndex,
				"tool":        req.PlanStep.Tool,
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

// createSkippedStep 创建跳过的步骤
func (e *StepExecutor) createSkippedStep(
	req *ExecuteStepRequest,
	startTime time.Time,
) *ExecuteStepResult {
	observation := req.PlanStep.Description
	elapsedMs := time.Since(startTime).Milliseconds()

	step := &models.AgentStep{
		Step:        req.StepIndex,
		Observation: observation,
		ElapsedMs:   elapsedMs,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	if req.StreamCallback != nil {
		req.StreamCallback(StreamEvent{
			Type: "step_end",
			Data: map[string]interface{}{
				"step":        req.StepIndex,
				"observation": observation,
				"elapsed_ms":  elapsedMs,
			},
		})
	}

	return &ExecuteStepResult{
		Step:  step,
		Error: nil,
	}
}

// createErrorStep 创建错误步骤
func (e *StepExecutor) createErrorStep(
	req *ExecuteStepRequest,
	startTime time.Time,
	err error,
) *ExecuteStepResult {
	elapsedMs := time.Since(startTime).Milliseconds()

	step := &models.AgentStep{
		Step: req.StepIndex,
		Action: &models.AgentAction{
			Type: "action",
			Tool: req.PlanStep.Tool,
		},
		Observation: fmt.Sprintf("Error: %s", err.Error()),
		Error:       err.Error(),
		ElapsedMs:   elapsedMs,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	if req.StreamCallback != nil {
		req.StreamCallback(StreamEvent{
			Type: "step_end",
			Data: map[string]interface{}{
				"step":        req.StepIndex,
				"tool":        req.PlanStep.Tool,
				"observation": step.Observation,
				"elapsed_ms":  elapsedMs,
			},
		})
	}

	return &ExecuteStepResult{
		Step:  step,
		Error: err,
	}
}

// getToolConfig 获取工具配置
func (e *StepExecutor) getToolConfig(tool interface{}) *tooling.ExecutionConfig {
	// 尝试从工具获取配置
	if configurable, ok := tool.(tooling.ConfigurableTool); ok {
		return configurable.GetExecutionConfig()
	}

	// 返回默认配置
	return tooling.DefaultExecutionConfig()
}

// BuildPreviousStepsContext 构建之前步骤的上下文
func BuildPreviousStepsContext(previousSteps []models.AgentStep) string {
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

	return contextBuilder.String()
}
