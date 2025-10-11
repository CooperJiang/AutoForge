package task

import (
	"strconv"
	"strings"
	"auto-forge/internal/models"
	"auto-forge/internal/repositories/task"
	"auto-forge/pkg/errors"
	"time"
)

// TaskService 任务服务
type TaskService struct {
	taskRepo          *task.TaskRepository
	executionRepo     *task.TaskExecutionRepository
	schedulerCallback func() // 调度器回调，用于通知调度器更新任务
}

var taskService *TaskService

// InitTaskService 初始化任务服务
func InitTaskService() {
	taskService = &TaskService{
		taskRepo:      task.NewTaskRepository(),
		executionRepo: task.NewTaskExecutionRepository(),
	}
}

// GetTaskService 获取任务服务实例
func GetTaskService() *TaskService {
	return taskService
}

// SetSchedulerCallback 设置调度器回调
func (s *TaskService) SetSchedulerCallback(callback func()) {
	s.schedulerCallback = callback
}

// CreateTask 创建任务
func (s *TaskService) CreateTask(userID, name, description, toolCode, config, scheduleType, scheduleValue string) (*models.Task, error) {
	// 验证调度配置
	if err := s.validateSchedule(scheduleType, scheduleValue); err != nil {
		return nil, err
	}

	// 计算下次执行时间
	nextRunTime := s.calculateNextRunTime(scheduleType, scheduleValue)

	task := &models.Task{
		UserID:        userID,
		Name:          name,
		Description:   description,
		ToolCode:      toolCode,
		Config:        config,
		ScheduleType:  scheduleType,
		ScheduleValue: scheduleValue,
		Enabled:       true,
		NextRunTime:   &nextRunTime,
	}

	if err := s.taskRepo.DB.Create(task).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeQueryFailed)
	}

	// 通知调度器更新
	if s.schedulerCallback != nil {
		s.schedulerCallback()
	}

	return task, nil
}

// GetTaskList 获取任务列表
func (s *TaskService) GetTaskList(userID string, page, pageSize int) ([]models.Task, int64, error) {
	return s.taskRepo.FindByUserID(userID, page, pageSize)
}

// GetTaskByID 获取任务详情
func (s *TaskService) GetTaskByID(id, userID string) (*models.Task, error) {
	return s.taskRepo.FindByIDAndUserID(id, userID)
}

// UpdateTask 更新任务
func (s *TaskService) UpdateTask(id, userID, name, description, toolCode, config, scheduleType, scheduleValue string) (*models.Task, error) {
	// 查询任务
	existingTask, err := s.taskRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, errors.New(errors.CodeNotFound, "任务不存在")
	}

	// 验证调度配置
	if err := s.validateSchedule(scheduleType, scheduleValue); err != nil {
		return nil, err
	}

	// 计算下次执行时间
	nextRunTime := s.calculateNextRunTime(scheduleType, scheduleValue)

	// 更新字段
	existingTask.Name = name
	existingTask.Description = description
	existingTask.ToolCode = toolCode
	existingTask.Config = config
	existingTask.ScheduleType = scheduleType
	existingTask.ScheduleValue = scheduleValue
	existingTask.NextRunTime = &nextRunTime

	if err := s.taskRepo.UpdateTask(existingTask); err != nil {
		return nil, errors.Wrap(err, errors.CodeQueryFailed)
	}

	// 通知调度器更新
	if s.schedulerCallback != nil {
		s.schedulerCallback()
	}

	return existingTask, nil
}

// DeleteTask 删除任务
func (s *TaskService) DeleteTask(id, userID string) error {
	if err := s.taskRepo.DeleteByIDAndUserID(id, userID); err != nil {
		return errors.Wrap(err, errors.CodeQueryFailed)
	}

	// 通知调度器更新
	if s.schedulerCallback != nil {
		s.schedulerCallback()
	}

	return nil
}

// EnableTask 启用任务
func (s *TaskService) EnableTask(id, userID string) error {
	if err := s.taskRepo.UpdateEnabled(id, userID, true); err != nil {
		return errors.Wrap(err, errors.CodeQueryFailed)
	}

	// 通知调度器更新
	if s.schedulerCallback != nil {
		s.schedulerCallback()
	}

	return nil
}

// DisableTask 禁用任务
func (s *TaskService) DisableTask(id, userID string) error {
	if err := s.taskRepo.UpdateEnabled(id, userID, false); err != nil {
		return errors.Wrap(err, errors.CodeQueryFailed)
	}

	// 通知调度器更新
	if s.schedulerCallback != nil {
		s.schedulerCallback()
	}

	return nil
}

// GetTaskExecutions 获取任务执行记录
func (s *TaskService) GetTaskExecutions(taskID string, page, pageSize int) ([]models.TaskExecution, int64, error) {
	return s.executionRepo.FindByTaskID(taskID, page, pageSize)
}

// GetExecutionByID 获取执行记录详情
func (s *TaskService) GetExecutionByID(id string) (*models.TaskExecution, error) {
	var execution models.TaskExecution
	if err := s.executionRepo.DB.Where("id = ?", id).First(&execution).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeQueryFailed)
	}
	return &execution, nil
}

// validateSchedule 验证调度配置
func (s *TaskService) validateSchedule(scheduleType, scheduleValue string) error {
	switch scheduleType {
	case "daily":
		// 验证时间格式 HH:MM:SS
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 3 {
			return errors.New(errors.CodeInvalidParameter, "daily调度值格式错误，应为 HH:MM:SS")
		}
		for _, part := range parts {
			if _, err := strconv.Atoi(part); err != nil {
				return errors.New(errors.CodeInvalidParameter, "daily调度值格式错误")
			}
		}
	case "weekly":
		// 验证格式 day1,day2,...:HH:MM:SS
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 4 {
			return errors.New(errors.CodeInvalidParameter, "weekly调度值格式错误，应为 day1,day2:HH:MM:SS")
		}
		// 验证星期几
		days := strings.Split(parts[0], ",")
		if len(days) == 0 {
			return errors.New(errors.CodeInvalidParameter, "weekly至少需要选择一个星期")
		}
		for _, day := range days {
			d, err := strconv.Atoi(day)
			if err != nil || d < 0 || d > 6 {
				return errors.New(errors.CodeInvalidParameter, "星期必须在0-6之间")
			}
		}
		// 验证时间
		for i := 1; i < 4; i++ {
			if _, err := strconv.Atoi(parts[i]); err != nil {
				return errors.New(errors.CodeInvalidParameter, "weekly时间格式错误")
			}
		}
	case "monthly":
		// 验证格式 day:HH:MM:SS
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 4 {
			return errors.New(errors.CodeInvalidParameter, "monthly调度值格式错误，应为 day:HH:MM:SS")
		}
		// 验证日期
		day, err := strconv.Atoi(parts[0])
		if err != nil || day < 1 || day > 31 {
			return errors.New(errors.CodeInvalidParameter, "日期必须在1-31之间")
		}
		// 验证时间
		for i := 1; i < 4; i++ {
			if _, err := strconv.Atoi(parts[i]); err != nil {
				return errors.New(errors.CodeInvalidParameter, "monthly时间格式错误")
			}
		}
	case "hourly":
		// 验证分钟:秒格式 MM:SS，最小间隔5分钟
		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 2 {
			return errors.New(errors.CodeInvalidParameter, "hourly调度值格式错误，应为 MM:SS")
		}
		minute, err := strconv.Atoi(parts[0])
		if err != nil {
			return errors.New(errors.CodeInvalidParameter, "hourly调度值格式错误")
		}
		if minute < 5 {
			return errors.New(errors.CodeInvalidParameter, "每小时执行间隔不能少于5分钟")
		}
		if _, err := strconv.Atoi(parts[1]); err != nil {
			return errors.New(errors.CodeInvalidParameter, "hourly调度值格式错误")
		}
	case "interval":
		// 验证间隔秒数，最小300秒（5分钟）
		seconds, err := strconv.Atoi(scheduleValue)
		if err != nil {
			return errors.New(errors.CodeInvalidParameter, "interval调度值必须是数字（秒）")
		}
		if seconds < 300 {
			return errors.New(errors.CodeInvalidParameter, "间隔执行不能少于300秒（5分钟）")
		}
	case "cron":
		// 这里可以添加 cron 表达式验证
		if scheduleValue == "" {
			return errors.New(errors.CodeInvalidParameter, "cron表达式不能为空")
		}
	default:
		return errors.New(errors.CodeInvalidParameter, "不支持的调度类型")
	}
	return nil
}

// calculateNextRunTime 计算下次执行时间
func (s *TaskService) calculateNextRunTime(scheduleType, scheduleValue string) int64 {
	now := time.Now()

	switch scheduleType {
	case "daily":
		// 解析时间 HH:MM:SS
		parts := strings.Split(scheduleValue, ":")
		hour, _ := strconv.Atoi(parts[0])
		minute, _ := strconv.Atoi(parts[1])
		second, _ := strconv.Atoi(parts[2])

		// 设置为今天的指定时间
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())

		// 如果已经过了今天的时间，设置为明天
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		return next.Unix()

	case "weekly":
		// 解析 day1,day2,...:HH:MM:SS
		parts := strings.Split(scheduleValue, ":")
		dayStrs := strings.Split(parts[0], ",")
		hour, _ := strconv.Atoi(parts[1])
		minute, _ := strconv.Atoi(parts[2])
		second, _ := strconv.Atoi(parts[3])

		// 转换星期几为整数
		var weekdays []int
		for _, dayStr := range dayStrs {
			day, _ := strconv.Atoi(dayStr)
			weekdays = append(weekdays, day)
		}

		// 找到最近的执行日
		currentWeekday := int(now.Weekday())
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())

		// 查找最近的星期几
		minDaysToAdd := 7
		for _, targetWeekday := range weekdays {
			daysToAdd := (targetWeekday - currentWeekday + 7) % 7
			// 如果是今天且时间未过，则为0天
			if daysToAdd == 0 && next.After(now) {
				minDaysToAdd = 0
				break
			}
			// 如果是今天但时间已过，需要等到下周
			if daysToAdd == 0 {
				daysToAdd = 7
			}
			if daysToAdd < minDaysToAdd {
				minDaysToAdd = daysToAdd
			}
		}

		next = next.Add(time.Duration(minDaysToAdd) * 24 * time.Hour)
		return next.Unix()

	case "monthly":
		// 解析 day:HH:MM:SS
		parts := strings.Split(scheduleValue, ":")
		day, _ := strconv.Atoi(parts[0])
		hour, _ := strconv.Atoi(parts[1])
		minute, _ := strconv.Atoi(parts[2])
		second, _ := strconv.Atoi(parts[3])

		// 设置为本月的指定日期和时间
		next := time.Date(now.Year(), now.Month(), day, hour, minute, second, 0, now.Location())

		// 如果已经过了本月的时间，设置为下个月
		if next.Before(now) {
			next = next.AddDate(0, 1, 0)
		}

		// 处理月末不存在的日期（如2月30日）
		for next.Day() != day {
			next = next.AddDate(0, 1, 0)
			next = time.Date(next.Year(), next.Month(), day, hour, minute, second, 0, next.Location())
		}

		return next.Unix()

	case "hourly":
		// 解析 MM:SS
		parts := strings.Split(scheduleValue, ":")
		minute, _ := strconv.Atoi(parts[0])
		second, _ := strconv.Atoi(parts[1])

		// 设置为本小时的指定时间
		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), minute, second, 0, now.Location())

		// 如果已经过了本小时的时间，设置为下一小时
		if next.Before(now) {
			next = next.Add(1 * time.Hour)
		}
		return next.Unix()

	case "interval":
		// 间隔秒数
		seconds, _ := strconv.ParseInt(scheduleValue, 10, 64)
		return now.Add(time.Duration(seconds) * time.Second).Unix()

	case "cron":
		// cron 表达式的处理需要使用 cron 库
		// 这里暂时返回一分钟后
		return now.Add(1 * time.Minute).Unix()

	default:
		return now.Unix()
	}
}

// GetAllEnabledTasks 获取所有启用的任务
func (s *TaskService) GetAllEnabledTasks() ([]models.Task, error) {
	return s.taskRepo.FindEnabledTasks()
}

// RecordExecution 记录任务执行
func (s *TaskService) RecordExecution(execution *models.TaskExecution) error {
	return s.executionRepo.Create(execution)
}

// UpdateNextRunTime 更新任务的下次执行时间
func (s *TaskService) UpdateNextRunTime(taskID string, scheduleType, scheduleValue string) error {
	nextRunTime := s.calculateNextRunTime(scheduleType, scheduleValue)
	return s.taskRepo.UpdateNextRunTime(taskID, nextRunTime)
}

// TriggerTask 手动触发任务
func (s *TaskService) TriggerTask(id, userID string) error {
	// 查询任务
	task, err := s.taskRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return errors.New(errors.CodeNotFound, "任务不存在")
	}

	// 触发任务执行（通过回调通知调度器）
	// 这里先简单返回成功，实际执行由调度器异步完成
	_ = task

	return nil
}

// DeleteExecution 删除执行记录
func (s *TaskService) DeleteExecution(id, userID string) error {
	// 先查询执行记录
	var execution models.TaskExecution
	if err := s.executionRepo.DB.Where("id = ?", id).First(&execution).Error; err != nil {
		return errors.New(errors.CodeNotFound, "执行记录不存在")
	}

	// 验证用户权限
	if execution.UserID != userID {
		return errors.New(errors.CodeForbidden, "无权删除此记录")
	}

	// 删除记录
	if err := s.executionRepo.DB.Where("id = ?", id).Delete(&models.TaskExecution{}).Error; err != nil {
		return errors.Wrap(err, errors.CodeQueryFailed)
	}

	return nil
}

// DeleteAllExecutions 删除任务的所有执行记录
func (s *TaskService) DeleteAllExecutions(taskID, userID string) error {
	// 验证任务是否存在且属于该用户
	task, err := s.taskRepo.FindByIDAndUserID(taskID, userID)
	if err != nil {
		return errors.New(errors.CodeNotFound, "任务不存在")
	}

	// 删除该任务的所有执行记录
	if err := s.executionRepo.DB.Where("task_id = ? AND user_id = ?", task.GetID(), userID).Delete(&models.TaskExecution{}).Error; err != nil {
		return errors.Wrap(err, errors.CodeQueryFailed)
	}

	return nil
}

