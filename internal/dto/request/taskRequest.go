package request

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Name          string                 `json:"name" binding:"required"`
	Description   string                 `json:"description"`
	ToolCode      string                 `json:"tool_code" binding:"required"`
	Config        map[string]interface{} `json:"config" binding:"required"`
	ScheduleType  string                 `json:"schedule_type" binding:"required"`
	ScheduleValue string                 `json:"schedule_value" binding:"required"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Name          string                 `json:"name" binding:"required"`
	Description   string                 `json:"description"`
	ToolCode      string                 `json:"tool_code" binding:"required"`
	Config        map[string]interface{} `json:"config" binding:"required"`
	ScheduleType  string                 `json:"schedule_type" binding:"required"`
	ScheduleValue string                 `json:"schedule_value" binding:"required"`
}

// TaskListRequest 任务列表请求
type TaskListRequest struct {
	UserID   string `form:"user_id" binding:"required"`
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}

// TaskExecutionListRequest 任务执行记录列表请求
type TaskExecutionListRequest struct {
	TaskID   string `form:"task_id" binding:"required"`
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}
