package task

import (
	"strconv"
	"strings"
	"auto-forge/internal/models"
	"auto-forge/internal/repositories/task"
	"auto-forge/pkg/errors"
	"time"
)


type TaskService struct {
	taskRepo          *task.TaskRepository
	executionRepo     *task.TaskExecutionRepository
	schedulerCallback func()
}

var taskService *TaskService


func InitTaskService() {
	taskService = &TaskService{
		taskRepo:      task.NewTaskRepository(),
		executionRepo: task.NewTaskExecutionRepository(),
	}
}


func GetTaskService() *TaskService {
	return taskService
}


func (s *TaskService) SetSchedulerCallback(callback func()) {
	s.schedulerCallback = callback
}


func (s *TaskService) CreateTask(userID, name, description, toolCode, config, scheduleType, scheduleValue string) (*models.Task, error) {

	if err := s.validateSchedule(scheduleType, scheduleValue); err != nil {
		return nil, err
	}


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


	if s.schedulerCallback != nil {
		s.schedulerCallback()
	}

	return task, nil
}


func (s *TaskService) GetTaskList(userID string, page, pageSize int) ([]models.Task, int64, error) {
	return s.taskRepo.FindByUserID(userID, page, pageSize)
}


func (s *TaskService) GetTaskByID(id, userID string) (*models.Task, error) {
	return s.taskRepo.FindByIDAndUserID(id, userID)
}


func (s *TaskService) UpdateTask(id, userID, name, description, toolCode, config, scheduleType, scheduleValue string) (*models.Task, error) {

	existingTask, err := s.taskRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, errors.New(errors.CodeNotFound, "任务不存在")
	}


	if err := s.validateSchedule(scheduleType, scheduleValue); err != nil {
		return nil, err
	}


	nextRunTime := s.calculateNextRunTime(scheduleType, scheduleValue)


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


	if s.schedulerCallback != nil {
		s.schedulerCallback()
	}

	return existingTask, nil
}


func (s *TaskService) DeleteTask(id, userID string) error {
	if err := s.taskRepo.DeleteByIDAndUserID(id, userID); err != nil {
		return errors.Wrap(err, errors.CodeQueryFailed)
	}


	if s.schedulerCallback != nil {
		s.schedulerCallback()
	}

	return nil
}


func (s *TaskService) EnableTask(id, userID string) error {
	if err := s.taskRepo.UpdateEnabled(id, userID, true); err != nil {
		return errors.Wrap(err, errors.CodeQueryFailed)
	}


	if s.schedulerCallback != nil {
		s.schedulerCallback()
	}

	return nil
}


func (s *TaskService) DisableTask(id, userID string) error {
	if err := s.taskRepo.UpdateEnabled(id, userID, false); err != nil {
		return errors.Wrap(err, errors.CodeQueryFailed)
	}


	if s.schedulerCallback != nil {
		s.schedulerCallback()
	}

	return nil
}


func (s *TaskService) GetTaskExecutions(taskID string, page, pageSize int) ([]models.TaskExecution, int64, error) {
	return s.executionRepo.FindByTaskID(taskID, page, pageSize)
}


func (s *TaskService) GetExecutionByID(id string) (*models.TaskExecution, error) {
	var execution models.TaskExecution
	if err := s.executionRepo.DB.Where("id = ?", id).First(&execution).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeQueryFailed)
	}
	return &execution, nil
}


func (s *TaskService) validateSchedule(scheduleType, scheduleValue string) error {
	switch scheduleType {
	case "daily":

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

		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 4 {
			return errors.New(errors.CodeInvalidParameter, "weekly调度值格式错误，应为 day1,day2:HH:MM:SS")
		}

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

		for i := 1; i < 4; i++ {
			if _, err := strconv.Atoi(parts[i]); err != nil {
				return errors.New(errors.CodeInvalidParameter, "weekly时间格式错误")
			}
		}
	case "monthly":

		parts := strings.Split(scheduleValue, ":")
		if len(parts) != 4 {
			return errors.New(errors.CodeInvalidParameter, "monthly调度值格式错误，应为 day:HH:MM:SS")
		}

		day, err := strconv.Atoi(parts[0])
		if err != nil || day < 1 || day > 31 {
			return errors.New(errors.CodeInvalidParameter, "日期必须在1-31之间")
		}

		for i := 1; i < 4; i++ {
			if _, err := strconv.Atoi(parts[i]); err != nil {
				return errors.New(errors.CodeInvalidParameter, "monthly时间格式错误")
			}
		}
	case "hourly":

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

		seconds, err := strconv.Atoi(scheduleValue)
		if err != nil {
			return errors.New(errors.CodeInvalidParameter, "interval调度值必须是数字（秒）")
		}
		if seconds < 300 {
			return errors.New(errors.CodeInvalidParameter, "间隔执行不能少于300秒（5分钟）")
		}
	case "cron":

		if scheduleValue == "" {
			return errors.New(errors.CodeInvalidParameter, "cron表达式不能为空")
		}
	default:
		return errors.New(errors.CodeInvalidParameter, "不支持的调度类型")
	}
	return nil
}


func (s *TaskService) calculateNextRunTime(scheduleType, scheduleValue string) int64 {
	now := time.Now()

	switch scheduleType {
	case "daily":

		parts := strings.Split(scheduleValue, ":")
		hour, _ := strconv.Atoi(parts[0])
		minute, _ := strconv.Atoi(parts[1])
		second, _ := strconv.Atoi(parts[2])


		next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())


		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		return next.Unix()

	case "weekly":

		parts := strings.Split(scheduleValue, ":")
		dayStrs := strings.Split(parts[0], ",")
		hour, _ := strconv.Atoi(parts[1])
		minute, _ := strconv.Atoi(parts[2])
		second, _ := strconv.Atoi(parts[3])


		var weekdays []int
		for _, dayStr := range dayStrs {
			day, _ := strconv.Atoi(dayStr)
			weekdays = append(weekdays, day)
		}


		currentWeekday := int(now.Weekday())
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())


		minDaysToAdd := 7
		for _, targetWeekday := range weekdays {
			daysToAdd := (targetWeekday - currentWeekday + 7) % 7

			if daysToAdd == 0 && next.After(now) {
				minDaysToAdd = 0
				break
			}

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

		parts := strings.Split(scheduleValue, ":")
		day, _ := strconv.Atoi(parts[0])
		hour, _ := strconv.Atoi(parts[1])
		minute, _ := strconv.Atoi(parts[2])
		second, _ := strconv.Atoi(parts[3])


		next := time.Date(now.Year(), now.Month(), day, hour, minute, second, 0, now.Location())


		if next.Before(now) {
			next = next.AddDate(0, 1, 0)
		}


		for next.Day() != day {
			next = next.AddDate(0, 1, 0)
			next = time.Date(next.Year(), next.Month(), day, hour, minute, second, 0, next.Location())
		}

		return next.Unix()

	case "hourly":

		parts := strings.Split(scheduleValue, ":")
		minute, _ := strconv.Atoi(parts[0])
		second, _ := strconv.Atoi(parts[1])


		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), minute, second, 0, now.Location())


		if next.Before(now) {
			next = next.Add(1 * time.Hour)
		}
		return next.Unix()

	case "interval":

		seconds, _ := strconv.ParseInt(scheduleValue, 10, 64)
		return now.Add(time.Duration(seconds) * time.Second).Unix()

	case "cron":


		return now.Add(1 * time.Minute).Unix()

	default:
		return now.Unix()
	}
}


func (s *TaskService) GetAllEnabledTasks() ([]models.Task, error) {
	return s.taskRepo.FindEnabledTasks()
}


func (s *TaskService) RecordExecution(execution *models.TaskExecution) error {
	return s.executionRepo.Create(execution)
}


func (s *TaskService) UpdateNextRunTime(taskID string, scheduleType, scheduleValue string) error {
	nextRunTime := s.calculateNextRunTime(scheduleType, scheduleValue)
	return s.taskRepo.UpdateNextRunTime(taskID, nextRunTime)
}


func (s *TaskService) TriggerTask(id, userID string) error {

	task, err := s.taskRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return errors.New(errors.CodeNotFound, "任务不存在")
	}



	_ = task

	return nil
}


func (s *TaskService) DeleteExecution(id, userID string) error {

	var execution models.TaskExecution
	if err := s.executionRepo.DB.Where("id = ?", id).First(&execution).Error; err != nil {
		return errors.New(errors.CodeNotFound, "执行记录不存在")
	}


	if execution.UserID != userID {
		return errors.New(errors.CodeForbidden, "无权删除此记录")
	}


	if err := s.executionRepo.DB.Where("id = ?", id).Delete(&models.TaskExecution{}).Error; err != nil {
		return errors.Wrap(err, errors.CodeQueryFailed)
	}

	return nil
}


func (s *TaskService) DeleteAllExecutions(taskID, userID string) error {

	task, err := s.taskRepo.FindByIDAndUserID(taskID, userID)
	if err != nil {
		return errors.New(errors.CodeNotFound, "任务不存在")
	}


	if err := s.executionRepo.DB.Where("task_id = ? AND user_id = ?", task.GetID(), userID).Delete(&models.TaskExecution{}).Error; err != nil {
		return errors.Wrap(err, errors.CodeQueryFailed)
	}

	return nil
}

