package services

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"auto-forge/internal/cron"
	"auto-forge/internal/models"
	"auto-forge/pkg/config"
	"auto-forge/pkg/database"
	"auto-forge/pkg/errors"
)


type AdminSession struct {
	Token     string
	ExpiresAt time.Time
}


type AdminService struct {
	sessions sync.Map
	cfg      *config.Config
}

var (
	adminServiceInstance *AdminService
	adminServiceOnce     sync.Once
)


func NewAdminService() *AdminService {
	adminServiceOnce.Do(func() {
		adminServiceInstance = &AdminService{
			cfg: config.GetConfig(),
		}
	})
	return adminServiceInstance
}


func (s *AdminService) Login(password string) (string, error) {

	if password != s.cfg.Admin.Password {
		return "", errors.New(errors.CodeUnauthorized, "密码错误")
	}


	token, err := s.generateToken()
	if err != nil {
		return "", err
	}


	expiresAt := time.Now().Add(time.Duration(s.cfg.Admin.SessionExpires) * time.Second)


	s.sessions.Store(token, AdminSession{
		Token:     token,
		ExpiresAt: expiresAt,
	})

	return token, nil
}


func (s *AdminService) ValidateToken(token string) bool {
	value, ok := s.sessions.Load(token)
	if !ok {
		return false
	}

	session := value.(AdminSession)


	if time.Now().After(session.ExpiresAt) {
		s.sessions.Delete(token)
		return false
	}

	return true
}


func (s *AdminService) Logout(token string) {
	s.sessions.Delete(token)
}


func (s *AdminService) generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}


func (s *AdminService) GetAllTasks(page, pageSize int, userID, status, keyword string) ([]models.Task, int64, error) {
	db := database.GetDB()


	query := db.Model(&models.Task{})


	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status == "enabled" {
		query = query.Where("enabled = ?", true)
	} else if status == "disabled" {
		query = query.Where("enabled = ?", false)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}


	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}


	var tasks []models.Task
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}


func (s *AdminService) UpdateTaskStatus(taskID string, enabled bool) error {
	db := database.GetDB()

	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return errors.New(errors.CodeNotFound, "任务不存在")
	}

	task.Enabled = enabled
	if err := db.Save(&task).Error; err != nil {
		return err
	}


	scheduler := cron.GetTaskScheduler()
	if scheduler != nil {
		scheduler.ReloadTasks()
	}

	return nil
}


func (s *AdminService) DeleteTask(taskID string) error {
	db := database.GetDB()

	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return errors.New(errors.CodeNotFound, "任务不存在")
	}


	if err := db.Where("task_id = ?", taskID).Delete(&models.TaskExecution{}).Error; err != nil {
		return err
	}


	if err := db.Delete(&task).Error; err != nil {
		return err
	}


	scheduler := cron.GetTaskScheduler()
	if scheduler != nil {
		scheduler.ReloadTasks()
	}

	return nil
}


func (s *AdminService) GetStats() (map[string]interface{}, error) {
	db := database.GetDB()

	stats := make(map[string]interface{})


	var totalUsers int64
	db.Model(&models.Task{}).Distinct("user_id").Count(&totalUsers)
	stats["total_users"] = totalUsers


	var totalTasks int64
	db.Model(&models.Task{}).Count(&totalTasks)
	stats["total_tasks"] = totalTasks


	var activeTasks int64
	db.Model(&models.Task{}).Where("enabled = ?", true).Count(&activeTasks)
	stats["active_tasks"] = activeTasks


	todayStart := time.Now().Truncate(24 * time.Hour).Unix()
	var todayExecutions int64
	db.Model(&models.TaskExecution{}).Where("started_at >= ?", todayStart).Count(&todayExecutions)
	stats["today_executions"] = todayExecutions


	var successCount int64
	var totalCount int64
	db.Model(&models.TaskExecution{}).Order("id DESC").Limit(100).Count(&totalCount)
	db.Model(&models.TaskExecution{}).Where("status = ?", "success").Order("id DESC").Limit(100).Count(&successCount)

	successRate := 0.0
	if totalCount > 0 {
		successRate = float64(successCount) / float64(totalCount) * 100
	}
	stats["success_rate"] = successRate

	// 获取工作流总数
	var totalWorkflows int64
	db.Model(&models.Workflow{}).Count(&totalWorkflows)
	stats["total_workflows"] = totalWorkflows

	// 获取工作流模板总数
	var totalTemplates int64
	db.Model(&models.WorkflowTemplate{}).Count(&totalTemplates)
	stats["total_templates"] = totalTemplates

	type UserActivity struct {
		UserID      string
		TaskCount   int64
		LastActive  time.Time
	}

	var recentUsers []UserActivity
	db.Raw(`
		SELECT user_id, COUNT(*) as task_count, MAX(updated_at) as last_active
		FROM tasks
		GROUP BY user_id
		ORDER BY last_active DESC
		LIMIT 5
	`).Scan(&recentUsers)

	stats["recent_users"] = recentUsers

	return stats, nil
}


func (s *AdminService) ExecuteTask(taskID string) error {
	db := database.GetDB()

	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return errors.New(errors.CodeNotFound, "任务不存在")
	}


	scheduler := cron.GetTaskScheduler()
	if scheduler == nil {
		return errors.New(errors.CodeInternal, "任务调度器未初始化")
	}

	scheduler.ExecuteTaskNow(&task)
	return nil
}


func (s *AdminService) GetAllExecutions(page, pageSize int, userID, taskID, status string) ([]models.TaskExecution, int64, error) {
	db := database.GetDB()


	query := db.Model(&models.TaskExecution{})


	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if taskID != "" {
		query = query.Where("task_id = ?", taskID)
	}
	if status == "success" {
		query = query.Where("status = ?", "success")
	} else if status == "failed" {
		query = query.Where("status != ?", "success")
	}


	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}


	var executions []models.TaskExecution
	offset := (page - 1) * pageSize
	if err := query.Preload("Task").Order("started_at DESC").Offset(offset).Limit(pageSize).Find(&executions).Error; err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}


func (s *AdminService) DeleteExecution(executionID string) error {
	db := database.GetDB()


	var execution models.TaskExecution
	if err := db.Where("id = ?", executionID).First(&execution).Error; err != nil {
		return errors.New(errors.CodeNotFound, "执行记录不存在")
	}


	if err := db.Delete(&execution).Error; err != nil {
		return errors.New(errors.CodeInternal, "删除执行记录失败")
	}

	return nil
}


type UserWithStats struct {
	models.User
	TotalTasks   int64 `json:"total_tasks"`
	EnabledTasks int64 `json:"enabled_tasks"`
}


func (s *AdminService) GetUsers(page, pageSize int, keyword string, status int) ([]UserWithStats, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New(errors.CodeDBConnectionFailed, "数据库连接失败")
	}

	query := db.Model(&models.User{})


	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}


	if status > 0 {
		query = query.Where("status = ?", status)
	}


	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}


	var users []models.User
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}


	usersWithStats := make([]UserWithStats, 0, len(users))
	for _, user := range users {
		userStats := UserWithStats{User: user}


		db.Model(&models.Task{}).Where("user_id = ?", user.ID.String()).Count(&userStats.TotalTasks)


		db.Model(&models.Task{}).Where("user_id = ? AND enabled = ?", user.ID.String(), true).Count(&userStats.EnabledTasks)

		usersWithStats = append(usersWithStats, userStats)
	}

	return usersWithStats, total, nil
}


func (s *AdminService) UpdateUserStatus(userID string, status int) error {
	db := database.GetDB()
	if db == nil {
		return errors.New(errors.CodeDBConnectionFailed, "数据库连接失败")
	}


	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}


	if user.Role == 1 && status == 2 {
		return errors.New(errors.CodeForbidden, "不能禁用超级管理员")
	}


	if err := db.Model(&user).Update("status", status).Error; err != nil {
		return errors.New(errors.CodeQueryFailed, "更新用户状态失败")
	}

	return nil
}
