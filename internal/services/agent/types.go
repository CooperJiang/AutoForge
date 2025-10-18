package agent

import "auto-forge/internal/models"

// AgentStreamEvent Agent 流式事件
type AgentStreamEvent struct {
	Type string      `json:"type"` // plan_start / plan_step / step_start / step_end / final / error
	Data interface{} `json:"data"`
}

// PlanStartEvent 计划开始事件
type PlanStartEvent struct {
	Plan *models.AgentPlan `json:"plan"`
}

// PlanStepEvent 计划步骤更新事件
type PlanStepEvent struct {
	StepIndex int    `json:"step_index"`
	Status    string `json:"status"` // pending / running / completed / failed / skipped
}

// StepStartEvent 步骤开始事件
type StepStartEvent struct {
	Step   int                    `json:"step"`
	Action *models.AgentAction    `json:"action"`
	Tool   string                 `json:"tool"`
	Args   map[string]interface{} `json:"args"`
}

// StepEndEvent 步骤结束事件
type StepEndEvent struct {
	Step        int                    `json:"step"`
	Observation string                 `json:"observation"`
	Output      map[string]interface{} `json:"output,omitempty"`
	ElapsedMs   int64                  `json:"elapsed_ms"`
	Error       string                 `json:"error,omitempty"`
}

// FinalEvent 最终结果事件
type FinalEvent struct {
	Answer       string                 `json:"answer"`
	FinishReason string                 `json:"finish_reason"`
	Trace        *models.AgentTrace     `json:"trace"`
	TokenUsage   *models.TokenUsage     `json:"token_usage,omitempty"`
	UsedTools    map[string]interface{} `json:"used_tools,omitempty"`
}

// ErrorEvent 错误事件
type ErrorEvent struct {
	Error   string `json:"error"`
	Step    int    `json:"step,omitempty"`
	Partial bool   `json:"partial"` // 是否是部分执行后的错误
}


