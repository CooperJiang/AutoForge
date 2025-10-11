package task

import (
	taskService "auto-forge/internal/services/task"
	"auto-forge/pkg/errors"

	"github.com/gin-gonic/gin"
)

// DeleteExecution 删除执行记录
func DeleteExecution(c *gin.Context) {
	id := c.Param("id")
	userID := c.Query("user_id")

	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "user_id不能为空"))
		return
	}

	service := taskService.GetTaskService()
	if err := service.DeleteExecution(id, userID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除执行记录成功")
}

// DeleteAllExecutions 删除任务的所有执行记录
func DeleteAllExecutions(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.Query("user_id")

	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "user_id不能为空"))
		return
	}

	service := taskService.GetTaskService()
	if err := service.DeleteAllExecutions(taskID, userID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "所有执行记录已删除")
}
