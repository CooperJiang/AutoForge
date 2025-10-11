package controllers

import (
	"strconv"
	"auto-forge/internal/services"
	"auto-forge/pkg/common"
	"auto-forge/pkg/errors"

	"github.com/gin-gonic/gin"
)

// AdminController 管理员控制器
type AdminController struct {
	adminService *services.AdminService
}

// NewAdminController 创建管理员控制器实例
func NewAdminController() *AdminController {
	return &AdminController{
		adminService: services.NewAdminService(),
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// Login 管理员登录
func (c *AdminController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "请求参数错误")
		return
	}

	token, err := c.adminService.Login(req.Password)
	if err != nil {
		common.Unauthorized(ctx, err.Error())
		return
	}

	common.Success(ctx, gin.H{
		"token":      token,
		"expires_in": 3600,
	}, "登录成功")
}

// Logout 管理员登出
func (c *AdminController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token != "" && len(token) > 7 {
		token = token[7:] // 移除 "Bearer " 前缀
		c.adminService.Logout(token)
	}

	common.SuccessWithMessage(ctx, "登出成功")
}

// TaskQueryParams 任务查询参数
type TaskQueryParams struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	UserID   string `form:"user_id"`
	Status   string `form:"status"` // all, enabled, disabled
	Keyword  string `form:"keyword"`
}

// GetTasks 获取所有任务
func (c *AdminController) GetTasks(ctx *gin.Context) {
	var params TaskQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		common.BadRequest(ctx, "请求参数错误")
		return
	}

	// 设置默认值
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 20
	}

	tasks, total, err := c.adminService.GetAllTasks(
		params.Page,
		params.PageSize,
		params.UserID,
		params.Status,
		params.Keyword,
	)

	if err != nil {
		common.ServerError(ctx, err.Error())
		return
	}

	common.Success(ctx, gin.H{
		"total": total,
		"tasks": tasks,
	}, "获取成功")
}

// UpdateTaskStatusRequest 更新任务状态请求
type UpdateTaskStatusRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateTaskStatus 更新任务状态
func (c *AdminController) UpdateTaskStatus(ctx *gin.Context) {
	taskID := ctx.Param("id")
	if taskID == "" {
		common.BadRequest(ctx, "无效的任务ID")
		return
	}

	var req UpdateTaskStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "请求参数错误")
		return
	}

	if err := c.adminService.UpdateTaskStatus(taskID, req.Enabled); err != nil {
		common.HandleError(ctx, errors.Wrap(err, errors.CodeInternal))
		return
	}

	common.SuccessWithMessage(ctx, "更新成功")
}

// DeleteTask 删除任务
func (c *AdminController) DeleteTask(ctx *gin.Context) {
	taskID := ctx.Param("id")
	if taskID == "" {
		common.BadRequest(ctx, "无效的任务ID")
		return
	}

	if err := c.adminService.DeleteTask(taskID); err != nil {
		common.HandleError(ctx, errors.Wrap(err, errors.CodeInternal))
		return
	}

	common.SuccessWithMessage(ctx, "删除成功")
}

// ExecuteTask 立即执行任务
func (c *AdminController) ExecuteTask(ctx *gin.Context) {
	taskID := ctx.Param("id")
	if taskID == "" {
		common.BadRequest(ctx, "无效的任务ID")
		return
	}

	err := c.adminService.ExecuteTask(taskID)
	if err != nil {
		common.HandleError(ctx, errors.Wrap(err, errors.CodeInternal))
		return
	}

	common.SuccessWithMessage(ctx, "任务已开始执行")
}

// GetStats 获取统计数据
func (c *AdminController) GetStats(ctx *gin.Context) {
	stats, err := c.adminService.GetStats()
	if err != nil {
		common.ServerError(ctx, err.Error())
		return
	}

	common.Success(ctx, stats, "获取成功")
}

// ExecutionQueryParams 执行记录查询参数
type ExecutionQueryParams struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	UserID   string `form:"user_id"`
	TaskID   string `form:"task_id"`
	Status   string `form:"status"` // success, failed, all
}

// GetExecutions 获取所有执行记录
func (c *AdminController) GetExecutions(ctx *gin.Context) {
	var params ExecutionQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		common.BadRequest(ctx, "请求参数错误")
		return
	}

	// 设置默认值
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 20
	}

	executions, total, err := c.adminService.GetAllExecutions(
		params.Page,
		params.PageSize,
		params.UserID,
		params.TaskID,
		params.Status,
	)

	if err != nil {
		common.ServerError(ctx, err.Error())
		return
	}

	common.Success(ctx, gin.H{
		"total":      total,
		"executions": executions,
	}, "获取成功")
}

// GetUsers 获取用户列表
func (c *AdminController) GetUsers(ctx *gin.Context) {
	page := 1
	pageSize := 20
	keyword := ctx.Query("keyword")
	status := 0

	if p := ctx.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if ps := ctx.Query("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err == nil && val > 0 {
			pageSize = val
		}
	}
	if s := ctx.Query("status"); s != "" {
		if val, err := strconv.Atoi(s); err == nil {
			status = val
		}
	}

	users, total, err := c.adminService.GetUsers(page, pageSize, keyword, status)
	if err != nil {
		common.ServerError(ctx, err.Error())
		return
	}

	common.Success(ctx, gin.H{
		"total": total,
		"users": users,
	}, "获取成功")
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	Status int `json:"status" binding:"required"`
}

// UpdateUserStatus 更新用户状态
func (c *AdminController) UpdateUserStatus(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		common.BadRequest(ctx, "用户ID不能为空")
		return
	}

	var req UpdateUserStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "请求参数错误")
		return
	}

	if err := c.adminService.UpdateUserStatus(userID, req.Status); err != nil {
		common.ServerError(ctx, err.Error())
		return
	}

	common.Success(ctx, nil, "更新成功")
}
