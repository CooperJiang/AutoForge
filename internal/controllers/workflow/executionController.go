package workflow

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/dto/response"
	"auto-forge/internal/services/workflow"
	"auto-forge/pkg/errors"
	log "auto-forge/pkg/logger"

	"github.com/gin-gonic/gin"
)

var executionService = workflow.NewExecutionService()
var engineService = workflow.NewEngineService()


func ExecuteWorkflow(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "工作流ID不能为空"))
		return
	}

	var req request.ExecuteWorkflowRequest
	_ = c.ShouldBindJSON(&req)


	execution, err := executionService.CreateExecution(workflowID, userID, "manual")
	if err != nil {
		log.Error("创建执行记录失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "执行失败: "+err.Error()))
		return
	}


	go func() {
		if err := engineService.ExecuteWorkflow(execution.GetID(), req.EnvVars, req.Params); err != nil {
			log.Error("工作流执行失败: ExecutionID=%s, Error=%v", execution.GetID(), err)
		}
	}()

	errors.ResponseSuccess(c, response.ExecuteWorkflowResponse{
		ExecutionID: execution.GetID(),
		Status:      execution.Status,
		Message:     "工作流已开始执行",
	}, "工作流已开始执行")
}


func GetExecutionList(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "工作流ID不能为空"))
		return
	}

	var query request.ExecutionListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	result, err := executionService.GetExecutionList(workflowID, userID, &query)
	if err != nil {
		log.Error("获取执行历史失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "获取执行历史失败"))
		return
	}

	errors.ResponseSuccess(c, result, "获取执行历史成功")
}


func GetExecutionDetail(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	executionID := c.Param("executionId")
	if executionID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "执行ID不能为空"))
		return
	}

	execution, err := executionService.GetExecutionByID(executionID, userID)
	if err != nil {
		log.Error("获取执行详情失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeNotFound, "获取执行详情失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, executionService.ToExecutionResponse(execution), "获取执行详情成功")
}


func StopExecution(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	executionID := c.Param("executionId")
	if executionID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "执行ID不能为空"))
		return
	}


	execution, err := executionService.GetExecutionByID(executionID, userID)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeNotFound, "执行记录不存在"))
		return
	}


	if execution.Status != "running" && execution.Status != "pending" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "只能停止运行中的执行"))
		return
	}


	if err := executionService.UpdateExecutionStatus(executionID, "cancelled", "用户手动取消"); err != nil {
		log.Error("停止执行失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "停止失败"))
		return
	}

	errors.ResponseSuccess(c, nil, "已停止执行")
}


func DeleteExecution(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	executionID := c.Param("executionId")
	if executionID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "执行ID不能为空"))
		return
	}


	if err := executionService.DeleteExecution(executionID, userID); err != nil {
		log.Error("删除执行记录失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, nil, "删除成功")
}
