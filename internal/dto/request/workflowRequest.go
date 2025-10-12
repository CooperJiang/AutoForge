package request

import "auto-forge/internal/models"

// CreateWorkflowRequest 创建工作流请求
type CreateWorkflowRequest struct {
	Name          string                  `json:"name" binding:"required"`
	Description   string                  `json:"description"`
	Nodes         []models.WorkflowNode   `json:"nodes" binding:"required"`
	Edges         []models.WorkflowEdge   `json:"edges" binding:"required"`
	EnvVars       []models.WorkflowEnvVar `json:"env_vars"`
	ScheduleType  string                  `json:"schedule_type"`
	ScheduleValue string                  `json:"schedule_value"`
	Enabled       bool                    `json:"enabled"`
}

// UpdateWorkflowRequest 更新工作流请求
type UpdateWorkflowRequest struct {
	Name          *string                  `json:"name"`
	Description   *string                  `json:"description"`
	Nodes         *[]models.WorkflowNode   `json:"nodes"`
	Edges         *[]models.WorkflowEdge   `json:"edges"`
	EnvVars       *[]models.WorkflowEnvVar `json:"env_vars"`
	ScheduleType  *string                  `json:"schedule_type"`
	ScheduleValue *string                  `json:"schedule_value"`
	Enabled       *bool                    `json:"enabled"`
}

// ExecuteWorkflowRequest 执行工作流请求
type ExecuteWorkflowRequest struct {
	EnvVars map[string]string `json:"env_vars"` // 临时环境变量
}

// WorkflowListQuery 工作流列表查询参数
type WorkflowListQuery struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword"`
	Enabled  *bool  `form:"enabled"`
}

// ExecutionListQuery 执行历史列表查询参数
type ExecutionListQuery struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status    string `form:"status"`
	StartTime *int64 `form:"start_time"`
	EndTime   *int64 `form:"end_time"`
}

// ToggleEnabledRequest 切换启用状态请求
type ToggleEnabledRequest struct {
	Enabled bool `json:"enabled"`
}

// ImportWorkflowRequest 导入工作流请求
type ImportWorkflowRequest struct {
	Name         string `json:"name" binding:"required"`
	WorkflowJSON string `json:"workflow_json" binding:"required"`
}

// ValidateWorkflowRequest 验证工作流请求
type ValidateWorkflowRequest struct {
	Nodes   []models.WorkflowNode   `json:"nodes" binding:"required"`
	Edges   []models.WorkflowEdge   `json:"edges" binding:"required"`
	EnvVars []models.WorkflowEnvVar `json:"env_vars"`
}
