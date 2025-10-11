package response

import "auto-forge/internal/models"

// TaskResponse 任务响应
type TaskResponse struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ToolCode      string `json:"tool_code"`
	Config        string `json:"config"`
	ScheduleType  string `json:"schedule_type"`
	ScheduleValue string `json:"schedule_value"`
	Enabled       bool   `json:"enabled"`
	NextRunTime   *int64 `json:"next_run_time"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Items    []TaskResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// TaskExecutionResponse 任务执行记录响应
type TaskExecutionResponse struct {
	ID             string `json:"id"`
	TaskID         string `json:"task_id"`
	UserID         string `json:"user_id"`
	Status         string `json:"status"`
	ResponseStatus int    `json:"response_status"`
	ResponseBody   string `json:"response_body"`
	DurationMs     int64  `json:"duration_ms"`
	ErrorMessage   string `json:"error_message"`
	StartedAt      int64  `json:"started_at"`
	CompletedAt    int64  `json:"completed_at"`
	CreatedAt      string `json:"created_at"`
}

// TaskExecutionListResponse 任务执行记录列表响应
type TaskExecutionListResponse struct {
	Items    []TaskExecutionResponse `json:"items"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
}

// ConvertTaskToResponse 转换任务模型到响应
func ConvertTaskToResponse(task *models.Task) TaskResponse {
	return TaskResponse{
		ID:            task.GetID(),
		UserID:        task.UserID,
		Name:          task.Name,
		Description:   task.Description,
		ToolCode:      task.ToolCode,
		Config:        task.Config,
		ScheduleType:  task.ScheduleType,
		ScheduleValue: task.ScheduleValue,
		Enabled:       task.Enabled,
		NextRunTime:   task.NextRunTime,
		CreatedAt:     task.GetCreatedAt().Format("2006-01-02 15:04:05"),
		UpdatedAt:     task.GetUpdatedAt().Format("2006-01-02 15:04:05"),
	}
}

// ConvertTasksToResponse 转换任务列表到响应
func ConvertTasksToResponse(tasks []models.Task) []TaskResponse {
	responses := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = ConvertTaskToResponse(&task)
	}
	return responses
}

// ConvertTaskExecutionToResponse 转换执行记录到响应
func ConvertTaskExecutionToResponse(execution *models.TaskExecution) TaskExecutionResponse {
	return TaskExecutionResponse{
		ID:             execution.GetID(),
		TaskID:         execution.TaskID,
		UserID:         execution.UserID,
		Status:         execution.Status,
		ResponseStatus: execution.ResponseStatus,
		ResponseBody:   execution.ResponseBody,
		DurationMs:     execution.DurationMs,
		ErrorMessage:   execution.ErrorMessage,
		StartedAt:      execution.StartedAt,
		CompletedAt:    execution.CompletedAt,
		CreatedAt:      execution.GetCreatedAt().Format("2006-01-02 15:04:05"),
	}
}

// ConvertTaskExecutionsToResponse 转换执行记录列表到响应
func ConvertTaskExecutionsToResponse(executions []models.TaskExecution) []TaskExecutionResponse {
	responses := make([]TaskExecutionResponse, len(executions))
	for i, execution := range executions {
		responses[i] = ConvertTaskExecutionToResponse(&execution)
	}
	return responses
}
