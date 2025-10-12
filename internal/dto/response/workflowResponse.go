package response

import "auto-forge/internal/models"

// WorkflowResponse 工作流响应
type WorkflowResponse struct {
	ID              string                  `json:"id"`
	UserID          string                  `json:"user_id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	Nodes           []models.WorkflowNode   `json:"nodes"`
	Edges           []models.WorkflowEdge   `json:"edges"`
	EnvVars         []models.WorkflowEnvVar `json:"env_vars"`
	ScheduleType    string                  `json:"schedule_type"`
	ScheduleValue   string                  `json:"schedule_value"`
	Enabled         bool                    `json:"enabled"`
	NextRunTime     *int64                  `json:"next_run_time"`
	TotalExecutions int                     `json:"total_executions"`
	SuccessCount    int                     `json:"success_count"`
	FailedCount     int                     `json:"failed_count"`
	LastExecutedAt  *int64                  `json:"last_executed_at"`
	CreatedAt       int64                   `json:"created_at"`
	UpdatedAt       int64                   `json:"updated_at"`
}

// WorkflowListResponse 工作流列表响应
type WorkflowListResponse struct {
	Items    []WorkflowResponse `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// WorkflowExecutionResponse 工作流执行记录响应
type WorkflowExecutionResponse struct {
	ID           string                      `json:"id"`
	WorkflowID   string                      `json:"workflow_id"`
	UserID       string                      `json:"user_id"`
	Status       string                      `json:"status"`
	TriggerType  string                      `json:"trigger_type"`
	StartTime    *int64                      `json:"start_time"`
	EndTime      *int64                      `json:"end_time"`
	DurationMs   int64                       `json:"duration_ms"`
	TotalNodes   int                         `json:"total_nodes"`
	SuccessNodes int                         `json:"success_nodes"`
	FailedNodes  int                         `json:"failed_nodes"`
	SkippedNodes int                         `json:"skipped_nodes"`
	NodeLogs     []models.NodeExecutionLog   `json:"node_logs"`
	Error        string                      `json:"error,omitempty"`
	CreatedAt    int64                       `json:"created_at"`
	UpdatedAt    int64                       `json:"updated_at"`
}

// ExecutionListResponse 执行历史列表响应
type ExecutionListResponse struct {
	Items    []WorkflowExecutionResponse `json:"items"`
	Total    int64                       `json:"total"`
	Page     int                         `json:"page"`
	PageSize int                         `json:"page_size"`
}

// ExecuteWorkflowResponse 执行工作流响应
type ExecuteWorkflowResponse struct {
	ExecutionID string `json:"execution_id"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

// WorkflowStatsResponse 工作流统计响应
type WorkflowStatsResponse struct {
	TotalExecutions int    `json:"total_executions"`
	SuccessCount    int    `json:"success_count"`
	FailedCount     int    `json:"failed_count"`
	AvgDurationMs   int64  `json:"avg_duration_ms"`
	LastExecutionAt *int64 `json:"last_execution_at"`
}

// ValidateWorkflowResponse 验证工作流响应
type ValidateWorkflowResponse struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
}
