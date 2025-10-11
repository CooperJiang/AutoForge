package task

import (
	"encoding/json"
	"strconv"
	"auto-forge/internal/dto/request"
	"auto-forge/internal/dto/response"
	taskService "auto-forge/internal/services/task"
	"auto-forge/pkg/common"
	"auto-forge/pkg/errors"

	"github.com/gin-gonic/gin"
)

// CreateTask 创建任务
func CreateTask(c *gin.Context) {
	// 从认证中间件设置的上下文中获取user_id
	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户未登录"))
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "用户ID无效"))
		return
	}

	req, err := common.ValidateRequest[request.CreateTaskRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 序列化配置
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "配置序列化失败"))
		return
	}

	service := taskService.GetTaskService()
	task, err := service.CreateTask(
		userIDStr,
		req.Name,
		req.Description,
		req.ToolCode,
		string(configJSON),
		req.ScheduleType,
		req.ScheduleValue,
	)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, response.ConvertTaskToResponse(task), "创建任务成功")
}

// GetTaskList 获取任务列表
func GetTaskList(c *gin.Context) {
	// 从认证中间件设置的上下文中获取user_id
	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户未登录"))
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "用户ID无效"))
		return
	}

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	pageInt := 1
	pageSizeInt := 20
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageInt = p
	}
	if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 {
		pageSizeInt = ps
	}

	service := taskService.GetTaskService()
	tasks, total, err := service.GetTaskList(userIDStr, pageInt, pageSizeInt)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	resp := response.TaskListResponse{
		Items:    response.ConvertTasksToResponse(tasks),
		Total:    total,
		Page:     pageInt,
		PageSize: pageSizeInt,
	}

	errors.ResponseSuccess(c, resp, "获取任务列表成功")
}

// GetTask 获取任务详情
func GetTask(c *gin.Context) {
	// 从认证中间件设置的上下文中获取user_id
	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户未登录"))
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "用户ID无效"))
		return
	}

	id := c.Param("id")

	service := taskService.GetTaskService()
	task, err := service.GetTaskByID(id, userIDStr)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, response.ConvertTaskToResponse(task), "获取任务详情成功")
}

// UpdateTask 更新任务
func UpdateTask(c *gin.Context) {
	// 从认证中间件设置的上下文中获取user_id
	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户未登录"))
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "用户ID无效"))
		return
	}

	id := c.Param("id")

	req, err := common.ValidateRequest[request.UpdateTaskRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 序列化配置
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "配置序列化失败"))
		return
	}

	service := taskService.GetTaskService()
	task, err := service.UpdateTask(
		id,
		userIDStr,
		req.Name,
		req.Description,
		req.ToolCode,
		string(configJSON),
		req.ScheduleType,
		req.ScheduleValue,
	)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, response.ConvertTaskToResponse(task), "更新任务成功")
}

// DeleteTask 删除任务
func DeleteTask(c *gin.Context) {
	// 从认证中间件设置的上下文中获取user_id
	userID, exists := c.Get("user_id")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户未登录"))
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "用户ID无效"))
		return
	}

	id := c.Param("id")

	service := taskService.GetTaskService()
	if err := service.DeleteTask(id, userIDStr); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除任务成功")
}

// EnableTask 启用任务
func EnableTask(c *gin.Context) {
	id := c.Param("id")
	userID := c.Query("user_id")

	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "user_id不能为空"))
		return
	}

	service := taskService.GetTaskService()
	if err := service.EnableTask(id, userID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "启用任务成功")
}

// DisableTask 禁用任务
func DisableTask(c *gin.Context) {
	id := c.Param("id")
	userID := c.Query("user_id")

	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "user_id不能为空"))
		return
	}

	service := taskService.GetTaskService()
	if err := service.DisableTask(id, userID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "禁用任务成功")
}

// GetTaskExecutions 获取任务执行记录
func GetTaskExecutions(c *gin.Context) {
	taskID := c.Param("id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	pageInt := 1
	pageSizeInt := 20
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageInt = p
	}
	if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 {
		pageSizeInt = ps
	}

	service := taskService.GetTaskService()
	executions, total, err := service.GetTaskExecutions(taskID, pageInt, pageSizeInt)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	resp := response.TaskExecutionListResponse{
		Items:    response.ConvertTaskExecutionsToResponse(executions),
		Total:    total,
		Page:     pageInt,
		PageSize: pageSizeInt,
	}

	errors.ResponseSuccess(c, resp, "获取执行记录成功")
}

// GetExecution 获取执行记录详情
func GetExecution(c *gin.Context) {
	id := c.Param("id")

	service := taskService.GetTaskService()
	execution, err := service.GetExecutionByID(id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, response.ConvertTaskExecutionToResponse(execution), "获取执行记录详情成功")
}

// TriggerTask 手动触发任务
func TriggerTask(c *gin.Context) {
	id := c.Param("id")
	userID := c.Query("user_id")

	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "user_id不能为空"))
		return
	}

	service := taskService.GetTaskService()
	if err := service.TriggerTask(id, userID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "触发任务成功")
}

