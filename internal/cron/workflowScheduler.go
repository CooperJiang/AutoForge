package cron

import (
	"auto-forge/internal/models"
	"auto-forge/internal/services/workflow"
	"auto-forge/pkg/database"
	"auto-forge/pkg/logger"
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
)

var (
	workflowScheduler *WorkflowScheduler
)

// WorkflowScheduler 工作流调度器
type WorkflowScheduler struct {
	cron             *cron.Cron
	workflowService  *workflow.WorkflowService
	executionService *workflow.ExecutionService
	engineService    *workflow.EngineService
	workflowIDs      map[string]cron.EntryID // workflowID -> entryID 的映射
}

// InitWorkflowScheduler 初始化工作流调度器
func InitWorkflowScheduler() {
	workflowScheduler = &WorkflowScheduler{
		cron:             cron.New(cron.WithSeconds()),
		workflowService:  workflow.NewWorkflowService(),
		executionService: workflow.NewExecutionService(),
		engineService:    workflow.NewEngineService(),
		workflowIDs:      make(map[string]cron.EntryID),
	}

	// 注册工作流变更回调
	workflow.SetWorkflowChangeCallback(func() {
		if workflowScheduler != nil {
			workflowScheduler.ReloadWorkflows()
		}
	})

	// 加载所有启用的工作流
	workflowScheduler.ReloadWorkflows()

	// 启动调度器
	workflowScheduler.cron.Start()
	logger.Info("工作流调度器已启动")
}

// ReloadWorkflows 重新加载所有工作流
func (ws *WorkflowScheduler) ReloadWorkflows() {
	logger.Info("===== 开始重新加载工作流 =====")

	// 清空现有工作流
	removedCount := len(ws.workflowIDs)
	for workflowID, entryID := range ws.workflowIDs {
		ws.cron.Remove(entryID)
		logger.Info("移除工作流调度: %s", workflowID)
	}
	ws.workflowIDs = make(map[string]cron.EntryID)
	logger.Info("已清空 %d 个现有工作流调度", removedCount)

	// 获取所有启用且有定时调度的工作流
	workflows, err := ws.getAllScheduledWorkflows()
	if err != nil {
		logger.Error("获取工作流列表失败: %v", err)
		return
	}
	logger.Info("从数据库查询到 %d 个待调度工作流", len(workflows))

	// 添加工作流到调度器
	successCount := 0
	for _, wf := range workflows {
		logger.Info("尝试添加工作流: Name=%s, ID=%s, ScheduleType=%s, ScheduleValue=%s, Enabled=%v",
			wf.Name, wf.GetID(), wf.ScheduleType, wf.ScheduleValue, wf.Enabled)
		if err := ws.addWorkflow(&wf); err != nil {
			logger.Error("添加工作流失败 [%s]: %v", wf.Name, err)
		} else {
			successCount++
		}
	}

	logger.Info("===== 工作流调度器加载完成: 成功 %d/%d =====", successCount, len(workflows))
}

// getAllScheduledWorkflows 获取所有需要调度的工作流
func (ws *WorkflowScheduler) getAllScheduledWorkflows() ([]models.Workflow, error) {
	db := database.GetDB()
	var workflows []models.Workflow

	// 查询所有启用且有定时配置的工作流
	logger.Info("查询数据库: enabled = true AND schedule_type != '' AND schedule_type != 'manual'")
	err := db.Where("enabled = ? AND schedule_type != ? AND schedule_type != ?",
		true, "", "manual").Find(&workflows).Error

	if err != nil {
		logger.Error("数据库查询失败: %v", err)
		return nil, err
	}

	logger.Info("数据库查询成功，找到 %d 条记录", len(workflows))
	for i, wf := range workflows {
		logger.Info("  [%d] ID=%s, Name=%s, Type=%s, Value=%s, Enabled=%v",
			i+1, wf.GetID(), wf.Name, wf.ScheduleType, wf.ScheduleValue, wf.Enabled)
	}

	return workflows, err
}

// addWorkflow 添加工作流到调度器
func (ws *WorkflowScheduler) addWorkflow(wf *models.Workflow) error {
	// 如果没有调度配置或者是手动触发，跳过
	if wf.ScheduleType == "" || wf.ScheduleType == "manual" {
		logger.Info("  -> 跳过: 无调度配置或为手动触发")
		return nil
	}

	logger.Info("  -> 构建cron表达式: Type=%s, Value=%s", wf.ScheduleType, wf.ScheduleValue)
	spec, err := ws.buildCronSpec(wf.ScheduleType, wf.ScheduleValue)
	if err != nil {
		logger.Error("  -> 构建cron表达式失败: %v", err)
		return err
	}
	logger.Info("  -> cron表达式: %s", spec)

	entryID, err := ws.cron.AddFunc(spec, func() {
		ws.executeWorkflow(wf)
	})

	if err != nil {
		logger.Error("  -> 添加到cron调度器失败: %v", err)
		return err
	}

	ws.workflowIDs[wf.GetID()] = entryID
	logger.Info("  -> ✓ 工作流已添加到调度器: %s (ID: %s, 调度: %s)", wf.Name, wf.GetID(), spec)

	return nil
}

// buildCronSpec 构建 cron 表达式
func (ws *WorkflowScheduler) buildCronSpec(scheduleType, scheduleValue string) (string, error) {
	switch scheduleType {
	case "daily":
		// scheduleValue: "HH:MM:SS"
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 3 {
			return "", fmt.Errorf("invalid daily format: %s", scheduleValue)
		}
		// cron 格式: 秒 分 时 日 月 周
		return fmt.Sprintf("%s %s %s * * *", parts[2], parts[1], parts[0]), nil

	case "weekly":
		// scheduleValue: "day1,day2,...:HH:MM:SS"
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 4 {
			return "", fmt.Errorf("invalid weekly format: %s", scheduleValue)
		}
		// cron 格式: 秒 分 时 日 月 周
		return fmt.Sprintf("%s %s %s * * %s", parts[3], parts[2], parts[1], parts[0]), nil

	case "monthly":
		// scheduleValue: "day:HH:MM:SS"
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 4 {
			return "", fmt.Errorf("invalid monthly format: %s", scheduleValue)
		}
		// cron 格式: 秒 分 时 日 月 周
		return fmt.Sprintf("%s %s %s %s * *", parts[3], parts[2], parts[1], parts[0]), nil

	case "hourly":
		// scheduleValue: "MM:SS"
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid hourly format: %s", scheduleValue)
		}
		// 每小时的指定分秒执行
		return fmt.Sprintf("%s %s * * * *", parts[1], parts[0]), nil

	case "interval":
		// scheduleValue: 秒数
		// cron 不直接支持间隔，这里使用每N秒执行一次的方式
		return fmt.Sprintf("*/%s * * * * *", scheduleValue), nil

	case "cron":
		// 直接使用 cron 表达式
		return scheduleValue, nil

	default:
		return "", fmt.Errorf("unsupported schedule type: %s", scheduleType)
	}
}

// executeWorkflow 执行工作流
func (ws *WorkflowScheduler) executeWorkflow(wf *models.Workflow) {
	logger.Info("开始执行工作流: %s (ID: %s)", wf.Name, wf.GetID())

	// 创建执行记录
	execution, err := ws.executionService.CreateExecution(wf.GetID(), wf.UserID, "scheduled")
	if err != nil {
		logger.Error("创建执行记录失败: %v", err)
		return
	}

	// 异步执行工作流
	go func() {
		if err := ws.engineService.ExecuteWorkflow(execution.GetID(), nil); err != nil {
			logger.Error("工作流执行失败: WorkflowID=%s, ExecutionID=%s, Error=%v",
				wf.GetID(), execution.GetID(), err)
		} else {
			logger.Info("工作流执行完成: %s (ExecutionID: %s)", wf.Name, execution.GetID())
		}
	}()
}

// StopWorkflowScheduler 停止工作流调度器
func StopWorkflowScheduler() {
	if workflowScheduler != nil && workflowScheduler.cron != nil {
		workflowScheduler.cron.Stop()
		logger.Info("工作流调度器已停止")
	}
}

// GetWorkflowScheduler 获取工作流调度器实例
func GetWorkflowScheduler() *WorkflowScheduler {
	return workflowScheduler
}

// ExecuteWorkflowNow 立即执行工作流（用于手动触发）
func (ws *WorkflowScheduler) ExecuteWorkflowNow(wf *models.Workflow) {
	go ws.executeWorkflow(wf)
}
