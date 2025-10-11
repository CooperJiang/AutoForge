package task

import (
	"auto-forge/internal/models"
	"auto-forge/pkg/database"

	"gorm.io/gorm"
)

// TaskExecutionRepository 任务执行记录仓库
type TaskExecutionRepository struct {
	DB *gorm.DB
}

// NewTaskExecutionRepository 创建任务执行记录仓库实例
func NewTaskExecutionRepository() *TaskExecutionRepository {
	return &TaskExecutionRepository{
		DB: database.GetDB(),
	}
}

// FindByTaskID 根据任务ID查询执行记录（分页）
func (r *TaskExecutionRepository) FindByTaskID(taskID string, page, pageSize int) ([]models.TaskExecution, int64, error) {
	var executions []models.TaskExecution
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.DB.Model(&models.TaskExecution{}).Where("task_id = ?", taskID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := r.DB.Where("task_id = ?", taskID).
		Order("started_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&executions).Error; err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// FindByUserID 根据用户ID查询执行记录（分页）
func (r *TaskExecutionRepository) FindByUserID(userID string, page, pageSize int) ([]models.TaskExecution, int64, error) {
	var executions []models.TaskExecution
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.DB.Model(&models.TaskExecution{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := r.DB.Where("user_id = ?", userID).
		Order("started_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&executions).Error; err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// Create 创建执行记录
func (r *TaskExecutionRepository) Create(execution *models.TaskExecution) error {
	return r.DB.Create(execution).Error
}
