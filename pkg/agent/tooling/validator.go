package tooling

import (
	"auto-forge/internal/models"
	"auto-forge/pkg/agent/registry"
	"fmt"
	"strings"
)

// PlanValidator 计划验证器
type PlanValidator struct {
	toolRegistry *registry.ToolRegistry
}

// NewPlanValidator 创建计划验证器
func NewPlanValidator(toolRegistry *registry.ToolRegistry) *PlanValidator {
	return &PlanValidator{
		toolRegistry: toolRegistry,
	}
}

// ValidationResult 验证结果
type ValidationResult struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

// Validate 验证计划
func (v *PlanValidator) Validate(plan *models.AgentPlan) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	if plan == nil || len(plan.Steps) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "计划为空")
		return result
	}

	// 收集每个步骤提供的类型
	providedTypes := make(map[string]bool)
	usedTools := make(map[string]bool)

	for i, step := range plan.Steps {
		stepNum := i + 1

		// 检查工具是否存在
		tool, err := v.toolRegistry.GetTool(step.Tool)
		if err != nil {
			result.Valid = false
			result.Errors = append(result.Errors,
				fmt.Sprintf("步骤 %d: 工具不存在 (%s)", stepNum, step.Tool))
			continue
		}

		// 获取工具配置
		config := v.getToolConfig(tool)
		if config == nil || config.Dependencies == nil {
			continue
		}

		deps := config.Dependencies

		// 检查依赖是否满足
		for _, required := range deps.Requires {
			if !providedTypes[required] {
				result.Valid = false
				result.Errors = append(result.Errors,
					fmt.Sprintf("步骤 %d (%s): 需要类型 '%s'，但前面的步骤未提供",
						stepNum, step.Tool, required))
			}
		}

		// 检查冲突
		for _, conflict := range deps.ConflictsWith {
			if usedTools[conflict] {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("步骤 %d (%s): 与已使用的工具 '%s' 冲突",
						stepNum, step.Tool, conflict))
			}
		}

		// 检查建议的前置工具
		if len(deps.SuggestedPredecessors) > 0 {
			hasPredecessor := false
			for _, pred := range deps.SuggestedPredecessors {
				if usedTools[pred] {
					hasPredecessor = true
					break
				}
			}
			if !hasPredecessor {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("步骤 %d (%s): 建议先使用 %s",
						stepNum, step.Tool, strings.Join(deps.SuggestedPredecessors, " 或 ")))
			}
		}

		// 记录此步骤提供的类型
		for _, provided := range deps.Provides {
			providedTypes[provided] = true
		}

		// 记录已使用的工具
		usedTools[step.Tool] = true
	}

	return result
}

// ValidateAndFix 验证并尝试修复计划
func (v *PlanValidator) ValidateAndFix(plan *models.AgentPlan) (*models.AgentPlan, *ValidationResult) {
	result := v.Validate(plan)

	// 如果验证通过或只有警告，直接返回
	if result.Valid {
		return plan, result
	}

	// TODO: 实现自动修复逻辑
	// 例如：自动插入缺失的前置步骤

	return plan, result
}

// getToolConfig 获取工具的执行配置
func (v *PlanValidator) getToolConfig(tool interface{}) *ExecutionConfig {
	// 尝试从工具获取配置
	// 这里需要工具实现一个接口来提供配置
	// 暂时返回 nil，后续扩展

	// 如果工具实现了 ConfigurableTool 接口
	if configurable, ok := tool.(ConfigurableTool); ok {
		return configurable.GetExecutionConfig()
	}

	return nil
}

// ConfigurableTool 可配置工具接口
type ConfigurableTool interface {
	GetExecutionConfig() *ExecutionConfig
}
