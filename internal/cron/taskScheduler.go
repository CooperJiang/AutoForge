package cron

import (
	"strings"
	"context"
	"encoding/json"
	"fmt"
	"auto-forge/internal/models"
	taskService "auto-forge/internal/services/task"
	"auto-forge/pkg/logger"
	"auto-forge/pkg/utools"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	taskScheduler *TaskScheduler
)

// TaskScheduler 任务调度器
type TaskScheduler struct {
	cron    *cron.Cron
	service *taskService.TaskService
	taskIDs map[string]cron.EntryID // taskID -> entryID 的映射
}

// InitTaskScheduler 初始化任务调度器
func InitTaskScheduler() {
	taskScheduler = &TaskScheduler{
		cron:    cron.New(cron.WithSeconds()),
		taskIDs: make(map[string]cron.EntryID),
	}

	taskScheduler.service = taskService.GetTaskService()

	// 设置服务层回调
	taskScheduler.service.SetSchedulerCallback(taskScheduler.ReloadTasks)

	// 加载所有任务
	taskScheduler.ReloadTasks()

	// 启动调度器
	taskScheduler.cron.Start()
	logger.Info("任务调度器已启动")
}

// ReloadTasks 重新加载所有任务
func (ts *TaskScheduler) ReloadTasks() {
	logger.Info("重新加载任务...")

	// 清空现有任务
	for _, entryID := range ts.taskIDs {
		ts.cron.Remove(entryID)
	}
	ts.taskIDs = make(map[string]cron.EntryID)

	// 获取所有启用的任务
	tasks, err := ts.service.GetAllEnabledTasks()
	if err != nil {
		logger.Error("获取任务列表失败: %v", err)
		return
	}

	// 添加任务到调度器
	for _, task := range tasks {
		if err := ts.addTask(&task); err != nil {
			logger.Error("添加任务失败 [%s]: %v", task.Name, err)
		}
	}

	logger.Info("成功加载 %d 个任务", len(ts.taskIDs))
}

// addTask 添加任务到调度器
func (ts *TaskScheduler) addTask(task *models.Task) error {
	spec, err := ts.buildCronSpec(task.ScheduleType, task.ScheduleValue)
	if err != nil {
		return err
	}

	entryID, err := ts.cron.AddFunc(spec, func() {
		ts.executeTask(task)
	})

	if err != nil {
		return err
	}

	ts.taskIDs[task.GetID()] = entryID
	logger.Info("任务已添加到调度器: %s (ID: %s, 调度: %s)", task.Name, task.GetID(), spec)

	return nil
}

// buildCronSpec 构建 cron 表达式
func (ts *TaskScheduler) buildCronSpec(scheduleType, scheduleValue string) (string, error) {
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
		// 注意：day1,day2,... 直接用于周字段
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

// executeTask 执行任务 - 直接通过工具
func (ts *TaskScheduler) executeTask(task *models.Task) {
	logger.Info("开始执行任务: %s (ID: %s)", task.Name, task.GetID())

	startTime := time.Now()

	// 解析工具配置
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(task.Config), &config); err != nil {
		logger.Error("解析工具配置失败: %v", err)
		ts.recordFailedExecution(task, startTime, fmt.Sprintf("配置解析失败: %v", err))
		return
	}

	// 获取工具
	registry := utools.GetRegistry()
	tool, err := registry.Get(task.ToolCode)
	if err != nil {
		logger.Error("获取工具失败 [%s]: %v", task.ToolCode, err)
		ts.recordFailedExecution(task, startTime, fmt.Sprintf("工具不存在: %v", err))
		return
	}

	// 构建执行上下文
	ctx := &utools.ExecutionContext{
		Context:   context.Background(),
		TaskID:    task.GetID(),
		UserID:    task.UserID,
		Variables: make(map[string]interface{}),
		Metadata:  make(map[string]interface{}),
	}

	// 执行工具
	result, err := tool.Execute(ctx, config)

	// 记录执行结果
	execution := &models.TaskExecution{
		TaskID:    task.GetID(),
		UserID:    task.UserID,
		StartedAt: startTime.Unix(),
		CompletedAt: time.Now().Unix(),
		DurationMs: time.Since(startTime).Milliseconds(),
	}

	if err != nil || !result.Success {
		execution.Status = "failed"
		if err != nil {
			execution.ErrorMessage = err.Error()
		} else {
			execution.ErrorMessage = result.Error
		}
		execution.ResponseBody = result.ResponseBody
	} else {
		execution.Status = "success"
		execution.ResponseStatus = result.StatusCode
		execution.ResponseBody = result.ResponseBody
	}

	// 保存执行记录
	if err := ts.service.RecordExecution(execution); err != nil {
		logger.Error("保存执行记录失败: %v", err)
	}

	// 更新下次执行时间
	if err := ts.service.UpdateNextRunTime(task.GetID(), task.ScheduleType, task.ScheduleValue); err != nil {
		logger.Error("更新下次执行时间失败: %v", err)
	}

	logger.Info("任务执行完成: %s, 状态: %s, 耗时: %dms", task.Name, execution.Status, execution.DurationMs)
}

// recordFailedExecution 记录失败的执行
func (ts *TaskScheduler) recordFailedExecution(task *models.Task, startTime time.Time, errorMsg string) {
	execution := &models.TaskExecution{
		TaskID:       task.GetID(),
		UserID:       task.UserID,
		Status:       "failed",
		ErrorMessage: errorMsg,
		StartedAt:    startTime.Unix(),
		CompletedAt:  time.Now().Unix(),
		DurationMs:   time.Since(startTime).Milliseconds(),
	}

	if err := ts.service.RecordExecution(execution); err != nil {
		logger.Error("保存执行记录失败: %v", err)
	}
}


// truncateString 截断字符串用于日志预览
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// StopTaskScheduler 停止任务调度器
func StopTaskScheduler() {
	if taskScheduler != nil && taskScheduler.cron != nil {
		taskScheduler.cron.Stop()
		logger.Info("任务调度器已停止")
	}
}

// GetTaskScheduler 获取任务调度器实例
func GetTaskScheduler() *TaskScheduler {
	return taskScheduler
}

// ExecuteTaskNow 立即执行任务（用于手动触发）
func (ts *TaskScheduler) ExecuteTaskNow(task *models.Task) {
	go ts.executeTask(task)
}

