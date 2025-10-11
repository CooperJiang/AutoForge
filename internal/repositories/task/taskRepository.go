package task

import (
	"auto-forge/internal/models"
	"auto-forge/pkg/database"

	"gorm.io/gorm"
)

// TaskRepository 任务仓库
type TaskRepository struct {
	DB *gorm.DB
}

// NewTaskRepository 创建任务仓库实例
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		DB: database.GetDB(),
	}
}

// FindByUserID 根据用户ID查询任务列表（分页）
func (r *TaskRepository) FindByUserID(userID string, page, pageSize int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.DB.Model(&models.Task{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := r.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// FindByIDAndUserID 根据ID和用户ID查询任务
func (r *TaskRepository) FindByIDAndUserID(id, userID string) (*models.Task, error) {
	var task models.Task
	if err := r.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateEnabled 更新任务启用状态
func (r *TaskRepository) UpdateEnabled(id, userID string, enabled bool) error {
	return r.DB.Model(&models.Task{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("enabled", enabled).Error
}

// UpdateNextRunTime 更新下次执行时间
func (r *TaskRepository) UpdateNextRunTime(id string, nextRunTime int64) error {
	return r.DB.Model(&models.Task{}).
		Where("id = ?", id).
		Update("next_run_time", nextRunTime).Error
}

// FindEnabledTasks 查询所有启用的任务
func (r *TaskRepository) FindEnabledTasks() ([]models.Task, error) {
	var tasks []models.Task
	if err := r.DB.Where("enabled = ?", true).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// DeleteByIDAndUserID 根据ID和用户ID删除任务
func (r *TaskRepository) DeleteByIDAndUserID(id, userID string) error {
	return r.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{}).Error
}

// UpdateTask 更新任务
func (r *TaskRepository) UpdateTask(task *models.Task) error {
	return r.DB.Save(task).Error
}

// Transaction 执行事务
func (r *TaskRepository) Transaction(fn func(*gorm.DB) error) error {
	return r.DB.Transaction(fn)
}
